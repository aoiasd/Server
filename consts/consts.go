package consts

import . "github.com/CyDrive/types"

const (
	HttpListenPort              = 6454
	HttpListenPortStr           = ":6454"
	RpcListenPort               = 6455
	RpcListenPortStr            = ":6455"
	FileTransferorListenPort    = 6456
	FileTransferorListenPortStr = ":6456"
	UserDataDir                 = "user_data"

	MemAccountStoreJsonPath = "accounts.json"

	// The size of file must be not greater than 1GB
	FileSizeLimit int64 = 1 << 30

	// A file with not small than 100MB size should be compressed
	CompressBaseline int64 = 100 << 20

	AccountStoreTypeMem   AccountStoreType = "mem"
	AccountStoreTypeMySQL AccountStoreType = "mysql"

	EnvTypeLocal  = "local"
	EnvTypeRemote = "remote"

	TimeFormat = "2006-01-02 15.04.05"

	DataTaskType_Download DataTaskType = iota
	DataTaskType_Upload
	DataTaskExpireTime int64 = 30 * 60
)
