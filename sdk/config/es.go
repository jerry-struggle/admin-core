package config

import (
	"github.com/jerry-struggle/admin-core/storage"
	storageEs "github.com/jerry-struggle/admin-core/storage/es"
	// 引入sms
)

type Es struct {
	Host  string `yaml:"host" json:"host"`
	Index string `yaml:"index" json:"index"`
	Type  string `yaml:"type" json:"type"`
}

var EsConfig = new(EsOption)

type EsOption struct {
	*Es
}

// Empty 空设置
func (e EsOption) Empty() bool {
	return e.Es == nil
}

func (e EsOption) Setup() (storage.AdapterEs, error) {
	if e.Es != nil {
		opt := _cfg.Settings.Es
		esClient, err := storageEs.NewEs(nil, &storageEs.EsOption{
			Host:  opt.Host,
			Index: opt.Index,
			Type:  opt.Type,
		})
		return esClient, err
	}
	return nil, nil
}
