package store

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/CyDrive/config"
	"github.com/CyDrive/consts"
	"github.com/CyDrive/model"
	"github.com/CyDrive/utils"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
)

type AccountStore interface {
	AddAccount(account *model.Account) error
	GetAccountByEmail(email string) (*model.Account, error)

	AddUsage(email string, usage int64) error
	ExpandCap(email string, newCap int64) error
}

// Store users in memory
// load from a json file
type MemStore struct {
	accountEmailMap map[string]*model.Account
	rwMutex         *sync.RWMutex // guard for accountEmailMap

	accountJsonPath string
	updatedFlag     int32
}

func NewMemStore(accountJsonPath string) *MemStore {
	store := MemStore{
		accountEmailMap: make(map[string]*model.Account),
		rwMutex:         &sync.RWMutex{},

		accountJsonPath: accountJsonPath,
		updatedFlag:     0,
	}

	data, err := ioutil.ReadFile(accountJsonPath)
	if err != nil {
		if !os.IsNotExist(err) {
			panic(err)
		}
	}

	accountArray := make([]*model.Account, 1)
	json.Unmarshal(data, &accountArray)
	for _, user := range accountArray {
		// Get the storage usage
		user.Usage, _ = utils.DirSize(user.DataDir)

		store.accountEmailMap[user.UserName] = user
	}

	return &store
}

func (store *MemStore) AddAccount(account *model.Account) error {
	store.rwMutex.RLock()
	_, ok := store.accountEmailMap[account.Email]
	if ok {
		store.rwMutex.RUnlock()
		return fmt.Errorf("email %v has been registered", account.Email)
	}
	store.rwMutex.RUnlock()

	store.rwMutex.Lock()
	defer store.rwMutex.Unlock()

	store.accountEmailMap[account.Email] = account
	store.updatedFlag++

	return nil
}

func (store *MemStore) GetAccountByEmail(email string) (*model.Account, error) {
	store.rwMutex.RLock()
	defer store.rwMutex.RUnlock()

	return store.accountEmailMap[email], nil
}

func (store *MemStore) AddUsage(email string, usage int64) error {
	store.rwMutex.RLock()
	defer store.rwMutex.RUnlock()

	account, err := store.GetAccountByEmail(email)
	if err != nil {
		return err
	}

	if usage != 0 {
		atomic.AddInt32(&store.updatedFlag, 1)
	}
	atomic.AddInt64(&account.Usage, usage)

	return nil
}

func (store *MemStore) ExpandCap(email string, newCap int64) error {
	store.rwMutex.RLock()
	defer store.rwMutex.RUnlock()

	account, err := store.GetAccountByEmail(email)
	if err != nil {
		return err
	}

	for {
		old := atomic.LoadInt64(&account.Cap)
		if atomic.CompareAndSwapInt64(&account.Cap, old, newCap) {
			if old != newCap {
				atomic.AddInt32(&store.updatedFlag, 1)
				break
			}
		}
	}

	return nil
}

func (store *MemStore) persistThread() {
	for {
		store.rwMutex.RLock()
		updatedNum := atomic.LoadInt32(&store.updatedFlag)
		if updatedNum > 0 {
			store.save()
			atomic.AddInt32(&store.updatedFlag, -updatedNum)
		}
		store.rwMutex.RUnlock()
	}
}

func (store *MemStore) save() {
	accountList := make([]*model.Account, 0, len(store.accountEmailMap))
	for _, account := range store.accountEmailMap {
		accountList = append(accountList, account)
	}

	accountListBytes, err := json.Marshal(accountList)
	if err != nil {
		log.Error(err)
		return
	}

	err = ioutil.WriteFile(store.accountJsonPath, accountListBytes, 0666)
	if err != nil {
		log.Error(err)
		return
	}
}

// Store users in a relational db
type RdbStore struct {
	db *gorm.DB
}

func NewRdbStore(config config.Config) *RdbStore {
	store := RdbStore{}
	store.db, _ = gorm.Open("mysql", config.PackDSN())
	return &store
}

func (store *RdbStore) GetAccountByEmail(name string) *model.Account {
	var user model.Account

	if store.db.First(user, "username = ?", name).RecordNotFound() {
		return nil
	}

	user.DataDir = filepath.Join(consts.UserDataDir, fmt.Sprint(user.Id))
	return &user
}
