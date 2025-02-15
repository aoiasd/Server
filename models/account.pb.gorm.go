package models

import (
	context "context"
	fmt "fmt"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	gorm1 "github.com/infobloxopen/atlas-app-toolkit/gorm"
	errors "github.com/infobloxopen/protoc-gen-gorm/errors"
	gorm "github.com/jinzhu/gorm"
	field_mask "google.golang.org/genproto/protobuf/field_mask"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	strings "strings"
	time "time"
)

type AccountORM struct {
	Cap       int64 `gorm:"default:0"`
	CreatedAt *time.Time
	Email     string `gorm:"unique;not null"`
	Id        int32  `gorm:"primary_key"`
	Name      string `gorm:"not null"`
	Password  string `gorm:"not null"`
	UpdatedAt *time.Time
	Usage     int64 `gorm:"default:0"`
}

// TableName overrides the default tablename generated by GORM
func (AccountORM) TableName() string {
	return "accounts"
}

// ToORM runs the BeforeToORM hook if present, converts the fields of this
// object to ORM format, runs the AfterToORM hook, then returns the ORM object
func (m *Account) ToORM(ctx context.Context) (AccountORM, error) {
	to := AccountORM{}
	var err error
	if prehook, ok := interface{}(m).(AccountWithBeforeToORM); ok {
		if err = prehook.BeforeToORM(ctx, &to); err != nil {
			return to, err
		}
	}
	to.Id = m.Id
	to.Email = m.Email
	to.Name = m.Name
	to.Usage = m.Usage
	to.Cap = m.Cap
	to.Password = m.Password
	if m.CreatedAt != nil {
		t := m.CreatedAt.AsTime()
		to.CreatedAt = &t
	}
	if m.UpdatedAt != nil {
		t := m.UpdatedAt.AsTime()
		to.UpdatedAt = &t
	}
	if posthook, ok := interface{}(m).(AccountWithAfterToORM); ok {
		err = posthook.AfterToORM(ctx, &to)
	}
	return to, err
}

// ToPB runs the BeforeToPB hook if present, converts the fields of this
// object to PB format, runs the AfterToPB hook, then returns the PB object
func (m *AccountORM) ToPB(ctx context.Context) (Account, error) {
	to := Account{}
	var err error
	if prehook, ok := interface{}(m).(AccountWithBeforeToPB); ok {
		if err = prehook.BeforeToPB(ctx, &to); err != nil {
			return to, err
		}
	}
	to.Id = m.Id
	to.Email = m.Email
	to.Name = m.Name
	to.Usage = m.Usage
	to.Cap = m.Cap
	to.Password = m.Password
	if m.CreatedAt != nil {
		to.CreatedAt = timestamppb.New(*m.CreatedAt)
	}
	if m.UpdatedAt != nil {
		to.UpdatedAt = timestamppb.New(*m.UpdatedAt)
	}
	if posthook, ok := interface{}(m).(AccountWithAfterToPB); ok {
		err = posthook.AfterToPB(ctx, &to)
	}
	return to, err
}

// The following are interfaces you can implement for special behavior during ORM/PB conversions
// of type Account the arg will be the target, the caller the one being converted from

// AccountBeforeToORM called before default ToORM code
type AccountWithBeforeToORM interface {
	BeforeToORM(context.Context, *AccountORM) error
}

// AccountAfterToORM called after default ToORM code
type AccountWithAfterToORM interface {
	AfterToORM(context.Context, *AccountORM) error
}

// AccountBeforeToPB called before default ToPB code
type AccountWithBeforeToPB interface {
	BeforeToPB(context.Context, *Account) error
}

// AccountAfterToPB called after default ToPB code
type AccountWithAfterToPB interface {
	AfterToPB(context.Context, *Account) error
}

// DefaultCreateAccount executes a basic gorm create call
func DefaultCreateAccount(ctx context.Context, in *Account, db *gorm.DB) (*Account, error) {
	if in == nil {
		return nil, errors.NilArgumentError
	}
	ormObj, err := in.ToORM(ctx)
	if err != nil {
		return nil, err
	}
	if hook, ok := interface{}(&ormObj).(AccountORMWithBeforeCreate_); ok {
		if db, err = hook.BeforeCreate_(ctx, db); err != nil {
			return nil, err
		}
	}
	if err = db.Create(&ormObj).Error; err != nil {
		return nil, err
	}
	if hook, ok := interface{}(&ormObj).(AccountORMWithAfterCreate_); ok {
		if err = hook.AfterCreate_(ctx, db); err != nil {
			return nil, err
		}
	}
	pbResponse, err := ormObj.ToPB(ctx)
	return &pbResponse, err
}

type AccountORMWithBeforeCreate_ interface {
	BeforeCreate_(context.Context, *gorm.DB) (*gorm.DB, error)
}
type AccountORMWithAfterCreate_ interface {
	AfterCreate_(context.Context, *gorm.DB) error
}

func DefaultReadAccount(ctx context.Context, in *Account, db *gorm.DB) (*Account, error) {
	if in == nil {
		return nil, errors.NilArgumentError
	}
	ormObj, err := in.ToORM(ctx)
	if err != nil {
		return nil, err
	}
	if ormObj.Id == 0 {
		return nil, errors.EmptyIdError
	}
	if hook, ok := interface{}(&ormObj).(AccountORMWithBeforeReadApplyQuery); ok {
		if db, err = hook.BeforeReadApplyQuery(ctx, db); err != nil {
			return nil, err
		}
	}
	if db, err = gorm1.ApplyFieldSelection(ctx, db, nil, &AccountORM{}); err != nil {
		return nil, err
	}
	if hook, ok := interface{}(&ormObj).(AccountORMWithBeforeReadFind); ok {
		if db, err = hook.BeforeReadFind(ctx, db); err != nil {
			return nil, err
		}
	}
	ormResponse := AccountORM{}
	if err = db.Where(&ormObj).First(&ormResponse).Error; err != nil {
		return nil, err
	}
	if hook, ok := interface{}(&ormResponse).(AccountORMWithAfterReadFind); ok {
		if err = hook.AfterReadFind(ctx, db); err != nil {
			return nil, err
		}
	}
	pbResponse, err := ormResponse.ToPB(ctx)
	return &pbResponse, err
}

type AccountORMWithBeforeReadApplyQuery interface {
	BeforeReadApplyQuery(context.Context, *gorm.DB) (*gorm.DB, error)
}
type AccountORMWithBeforeReadFind interface {
	BeforeReadFind(context.Context, *gorm.DB) (*gorm.DB, error)
}
type AccountORMWithAfterReadFind interface {
	AfterReadFind(context.Context, *gorm.DB) error
}

func DefaultDeleteAccount(ctx context.Context, in *Account, db *gorm.DB) error {
	if in == nil {
		return errors.NilArgumentError
	}
	ormObj, err := in.ToORM(ctx)
	if err != nil {
		return err
	}
	if ormObj.Id == 0 {
		return errors.EmptyIdError
	}
	if hook, ok := interface{}(&ormObj).(AccountORMWithBeforeDelete_); ok {
		if db, err = hook.BeforeDelete_(ctx, db); err != nil {
			return err
		}
	}
	err = db.Where(&ormObj).Delete(&AccountORM{}).Error
	if err != nil {
		return err
	}
	if hook, ok := interface{}(&ormObj).(AccountORMWithAfterDelete_); ok {
		err = hook.AfterDelete_(ctx, db)
	}
	return err
}

type AccountORMWithBeforeDelete_ interface {
	BeforeDelete_(context.Context, *gorm.DB) (*gorm.DB, error)
}
type AccountORMWithAfterDelete_ interface {
	AfterDelete_(context.Context, *gorm.DB) error
}

func DefaultDeleteAccountSet(ctx context.Context, in []*Account, db *gorm.DB) error {
	if in == nil {
		return errors.NilArgumentError
	}
	var err error
	keys := []int32{}
	for _, obj := range in {
		ormObj, err := obj.ToORM(ctx)
		if err != nil {
			return err
		}
		if ormObj.Id == 0 {
			return errors.EmptyIdError
		}
		keys = append(keys, ormObj.Id)
	}
	if hook, ok := (interface{}(&AccountORM{})).(AccountORMWithBeforeDeleteSet); ok {
		if db, err = hook.BeforeDeleteSet(ctx, in, db); err != nil {
			return err
		}
	}
	err = db.Where("id in (?)", keys).Delete(&AccountORM{}).Error
	if err != nil {
		return err
	}
	if hook, ok := (interface{}(&AccountORM{})).(AccountORMWithAfterDeleteSet); ok {
		err = hook.AfterDeleteSet(ctx, in, db)
	}
	return err
}

type AccountORMWithBeforeDeleteSet interface {
	BeforeDeleteSet(context.Context, []*Account, *gorm.DB) (*gorm.DB, error)
}
type AccountORMWithAfterDeleteSet interface {
	AfterDeleteSet(context.Context, []*Account, *gorm.DB) error
}

// DefaultStrictUpdateAccount clears / replaces / appends first level 1:many children and then executes a gorm update call
func DefaultStrictUpdateAccount(ctx context.Context, in *Account, db *gorm.DB) (*Account, error) {
	if in == nil {
		return nil, fmt.Errorf("Nil argument to DefaultStrictUpdateAccount")
	}
	ormObj, err := in.ToORM(ctx)
	if err != nil {
		return nil, err
	}
	lockedRow := &AccountORM{}
	db.Model(&ormObj).Set("gorm:query_option", "FOR UPDATE").Where("id=?", ormObj.Id).First(lockedRow)
	if hook, ok := interface{}(&ormObj).(AccountORMWithBeforeStrictUpdateCleanup); ok {
		if db, err = hook.BeforeStrictUpdateCleanup(ctx, db); err != nil {
			return nil, err
		}
	}
	if hook, ok := interface{}(&ormObj).(AccountORMWithBeforeStrictUpdateSave); ok {
		if db, err = hook.BeforeStrictUpdateSave(ctx, db); err != nil {
			return nil, err
		}
	}
	if err = db.Save(&ormObj).Error; err != nil {
		return nil, err
	}
	if hook, ok := interface{}(&ormObj).(AccountORMWithAfterStrictUpdateSave); ok {
		if err = hook.AfterStrictUpdateSave(ctx, db); err != nil {
			return nil, err
		}
	}
	pbResponse, err := ormObj.ToPB(ctx)
	if err != nil {
		return nil, err
	}
	return &pbResponse, err
}

type AccountORMWithBeforeStrictUpdateCleanup interface {
	BeforeStrictUpdateCleanup(context.Context, *gorm.DB) (*gorm.DB, error)
}
type AccountORMWithBeforeStrictUpdateSave interface {
	BeforeStrictUpdateSave(context.Context, *gorm.DB) (*gorm.DB, error)
}
type AccountORMWithAfterStrictUpdateSave interface {
	AfterStrictUpdateSave(context.Context, *gorm.DB) error
}

// DefaultPatchAccount executes a basic gorm update call with patch behavior
func DefaultPatchAccount(ctx context.Context, in *Account, updateMask *field_mask.FieldMask, db *gorm.DB) (*Account, error) {
	if in == nil {
		return nil, errors.NilArgumentError
	}
	var pbObj Account
	var err error
	if hook, ok := interface{}(&pbObj).(AccountWithBeforePatchRead); ok {
		if db, err = hook.BeforePatchRead(ctx, in, updateMask, db); err != nil {
			return nil, err
		}
	}
	pbReadRes, err := DefaultReadAccount(ctx, &Account{Id: in.GetId()}, db)
	if err != nil {
		return nil, err
	}
	pbObj = *pbReadRes
	if hook, ok := interface{}(&pbObj).(AccountWithBeforePatchApplyFieldMask); ok {
		if db, err = hook.BeforePatchApplyFieldMask(ctx, in, updateMask, db); err != nil {
			return nil, err
		}
	}
	if _, err := DefaultApplyFieldMaskAccount(ctx, &pbObj, in, updateMask, "", db); err != nil {
		return nil, err
	}
	if hook, ok := interface{}(&pbObj).(AccountWithBeforePatchSave); ok {
		if db, err = hook.BeforePatchSave(ctx, in, updateMask, db); err != nil {
			return nil, err
		}
	}
	pbResponse, err := DefaultStrictUpdateAccount(ctx, &pbObj, db)
	if err != nil {
		return nil, err
	}
	if hook, ok := interface{}(pbResponse).(AccountWithAfterPatchSave); ok {
		if err = hook.AfterPatchSave(ctx, in, updateMask, db); err != nil {
			return nil, err
		}
	}
	return pbResponse, nil
}

type AccountWithBeforePatchRead interface {
	BeforePatchRead(context.Context, *Account, *field_mask.FieldMask, *gorm.DB) (*gorm.DB, error)
}
type AccountWithBeforePatchApplyFieldMask interface {
	BeforePatchApplyFieldMask(context.Context, *Account, *field_mask.FieldMask, *gorm.DB) (*gorm.DB, error)
}
type AccountWithBeforePatchSave interface {
	BeforePatchSave(context.Context, *Account, *field_mask.FieldMask, *gorm.DB) (*gorm.DB, error)
}
type AccountWithAfterPatchSave interface {
	AfterPatchSave(context.Context, *Account, *field_mask.FieldMask, *gorm.DB) error
}

// DefaultPatchSetAccount executes a bulk gorm update call with patch behavior
func DefaultPatchSetAccount(ctx context.Context, objects []*Account, updateMasks []*field_mask.FieldMask, db *gorm.DB) ([]*Account, error) {
	if len(objects) != len(updateMasks) {
		return nil, fmt.Errorf(errors.BadRepeatedFieldMaskTpl, len(updateMasks), len(objects))
	}

	results := make([]*Account, 0, len(objects))
	for i, patcher := range objects {
		pbResponse, err := DefaultPatchAccount(ctx, patcher, updateMasks[i], db)
		if err != nil {
			return nil, err
		}

		results = append(results, pbResponse)
	}

	return results, nil
}

// DefaultApplyFieldMaskAccount patches an pbObject with patcher according to a field mask.
func DefaultApplyFieldMaskAccount(ctx context.Context, patchee *Account, patcher *Account, updateMask *field_mask.FieldMask, prefix string, db *gorm.DB) (*Account, error) {
	if patcher == nil {
		return nil, nil
	} else if patchee == nil {
		return nil, errors.NilArgumentError
	}
	var err error
	var updatedCreatedAt bool
	var updatedUpdatedAt bool
	for i, f := range updateMask.Paths {
		if f == prefix+"Id" {
			patchee.Id = patcher.Id
			continue
		}
		if f == prefix+"Email" {
			patchee.Email = patcher.Email
			continue
		}
		if f == prefix+"Name" {
			patchee.Name = patcher.Name
			continue
		}
		if f == prefix+"Usage" {
			patchee.Usage = patcher.Usage
			continue
		}
		if f == prefix+"Cap" {
			patchee.Cap = patcher.Cap
			continue
		}
		if f == prefix+"Password" {
			patchee.Password = patcher.Password
			continue
		}
		if !updatedCreatedAt && strings.HasPrefix(f, prefix+"CreatedAt.") {
			if patcher.CreatedAt == nil {
				patchee.CreatedAt = nil
				continue
			}
			if patchee.CreatedAt == nil {
				patchee.CreatedAt = &timestamp.Timestamp{}
			}
			childMask := &field_mask.FieldMask{}
			for j := i; j < len(updateMask.Paths); j++ {
				if trimPath := strings.TrimPrefix(updateMask.Paths[j], prefix+"CreatedAt."); trimPath != updateMask.Paths[j] {
					childMask.Paths = append(childMask.Paths, trimPath)
				}
			}
			if err := gorm1.MergeWithMask(patcher.CreatedAt, patchee.CreatedAt, childMask); err != nil {
				return nil, nil
			}
		}
		if f == prefix+"CreatedAt" {
			updatedCreatedAt = true
			patchee.CreatedAt = patcher.CreatedAt
			continue
		}
		if !updatedUpdatedAt && strings.HasPrefix(f, prefix+"UpdatedAt.") {
			if patcher.UpdatedAt == nil {
				patchee.UpdatedAt = nil
				continue
			}
			if patchee.UpdatedAt == nil {
				patchee.UpdatedAt = &timestamp.Timestamp{}
			}
			childMask := &field_mask.FieldMask{}
			for j := i; j < len(updateMask.Paths); j++ {
				if trimPath := strings.TrimPrefix(updateMask.Paths[j], prefix+"UpdatedAt."); trimPath != updateMask.Paths[j] {
					childMask.Paths = append(childMask.Paths, trimPath)
				}
			}
			if err := gorm1.MergeWithMask(patcher.UpdatedAt, patchee.UpdatedAt, childMask); err != nil {
				return nil, nil
			}
		}
		if f == prefix+"UpdatedAt" {
			updatedUpdatedAt = true
			patchee.UpdatedAt = patcher.UpdatedAt
			continue
		}
	}
	if err != nil {
		return nil, err
	}
	return patchee, nil
}

// DefaultListAccount executes a gorm list call
func DefaultListAccount(ctx context.Context, db *gorm.DB) ([]*Account, error) {
	in := Account{}
	ormObj, err := in.ToORM(ctx)
	if err != nil {
		return nil, err
	}
	if hook, ok := interface{}(&ormObj).(AccountORMWithBeforeListApplyQuery); ok {
		if db, err = hook.BeforeListApplyQuery(ctx, db); err != nil {
			return nil, err
		}
	}
	db, err = gorm1.ApplyCollectionOperators(ctx, db, &AccountORM{}, &Account{}, nil, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	if hook, ok := interface{}(&ormObj).(AccountORMWithBeforeListFind); ok {
		if db, err = hook.BeforeListFind(ctx, db); err != nil {
			return nil, err
		}
	}
	db = db.Where(&ormObj)
	db = db.Order("id")
	ormResponse := []AccountORM{}
	if err := db.Find(&ormResponse).Error; err != nil {
		return nil, err
	}
	if hook, ok := interface{}(&ormObj).(AccountORMWithAfterListFind); ok {
		if err = hook.AfterListFind(ctx, db, &ormResponse); err != nil {
			return nil, err
		}
	}
	pbResponse := []*Account{}
	for _, responseEntry := range ormResponse {
		temp, err := responseEntry.ToPB(ctx)
		if err != nil {
			return nil, err
		}
		pbResponse = append(pbResponse, &temp)
	}
	return pbResponse, nil
}

type AccountORMWithBeforeListApplyQuery interface {
	BeforeListApplyQuery(context.Context, *gorm.DB) (*gorm.DB, error)
}
type AccountORMWithBeforeListFind interface {
	BeforeListFind(context.Context, *gorm.DB) (*gorm.DB, error)
}
type AccountORMWithAfterListFind interface {
	AfterListFind(context.Context, *gorm.DB, *[]AccountORM) error
}
