package storage

import (
	"github.com/jerry-struggle/admin-core/storage/es"
	"time"

	"github.com/bsm/redislock"
	sts "github.com/tencentyun/qcloud-cos-sts-sdk/go"
)

const (
	PrefixKey = "__host"
)

type AdapterCache interface {
	String() string
	Get(key string) (string, error)
	Set(key string, val interface{}, expire int) error
	Del(key string) error
	HashGet(hk, key string) (string, error)
	HashDel(hk, key string) error
	Increase(key string) error
	Decrease(key string) error
	Expire(key string, dur time.Duration) error
	HashSet(key string, values ...interface{}) error
	HashMSet(key string, values ...interface{}) error
	Exists(key string) (int64, error)
	HashGetAll(key string) (map[string]string, error)
	LPush(key string, values ...interface{}) (int64, error)
	RPush(key string, values ...interface{}) (int64, error)
	LRange(key string, start, stop int64) ([]string, error)
}

type AdapterQueue interface {
	String() string
	Append(message Messager) error
	Register(name string, f ConsumerFunc)
	Run()
	Shutdown()
}

type Messager interface {
	SetID(string)
	SetStream(string)
	SetValues(map[string]interface{})
	GetID() string
	GetStream() string
	GetValues() map[string]interface{}
	GetPrefix() string
	SetPrefix(string)
	SetErrorCount(count int)
	GetErrorCount() int
}

type ConsumerFunc func(Messager) error

type AdapterLocker interface {
	String() string
	Lock(key string, ttl int64, options *redislock.Options) (*redislock.Lock, error)
}

type AdapterSms interface {
	SendSms(string, string) error
}

type AdapterOss interface {
	UpLoad(objectName, localFile string) (string, error)
	GetCredential() (*sts.CredentialResult, error)
	GetPresignedURL(objectName string) (string, error)
	GetSignedFileUrl(dir string, isSign bool) ([]string, []string, error)
	DownloadFiles(fileObject []string, localDir string) error
}

type AdapterEs interface {
	AddRecord(string, int, string, string, string, string) (string, error)
	GetRecord(string, int) (*es.Knowledge, error)
	UpdateRecord(string, int, string, string, string, string) error
	DeleteRecord(string, int) error
	PageRecord(string, int, int, string) (int64, []int, error)
	RecordList(string, string) ([]int, error)
}
