package config

import (
	"github.com/jerry-struggle/admin-core/storage"
	"github.com/jerry-struggle/admin-core/storage/oss"
)

type Oss struct {
	SecretId   string
	SecretKey  string
	AppId      string
	Bucket     string
	Region     string
	ExpireTime int64
}

var OssConfig = new(OssOption)

type OssOption struct {
	*Oss
}

// Empty 空设置
func (e OssOption) Empty() bool {
	return e.Oss == nil
}

func (e OssOption) Setup() (storage.AdapterOss, error) {
	if e.Oss != nil {
		opt := _cfg.Settings.Oss
		ossClient := oss.NewOss(nil, &oss.OssOption{
			SecretId:   opt.SecretId,
			SecretKey:  opt.SecretKey,
			AppId:      opt.AppId,
			Bucket:     opt.Bucket,
			Region:     opt.Region,
			ExpireTime: opt.ExpireTime,
		})
		return ossClient, nil
	}
	return nil, nil
}
