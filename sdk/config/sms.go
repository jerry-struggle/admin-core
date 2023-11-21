package config

import (
	"github.com/jerry-struggle/admin-core/storage"
	storageSms "github.com/jerry-struggle/admin-core/storage/sms"
	// 引入sms
)

type Sms struct {
	SecretId   string `yaml:"secretId" json:"secretId"`
	SecretKey  string `yaml:"secretKey" json:"secretKey"`
	Region     string `yaml:"region" json:"region"`
	SdkAppId   string `yaml:"sdkAppId" json:"sdkAppId"`
	SignName   string `yaml:"signName" json:"signName"`
	TemplateId string `yaml:"templateId" json:"templateId"`
}

var SmsConfig = new(SmsOption)

type SmsOption struct {
	*Sms
}

// Empty 空设置
func (e SmsOption) Empty() bool {
	return e.Sms == nil
}

func (e SmsOption) Setup() (storage.AdapterSms, error) {
	if e.Sms != nil {
		opt := _cfg.Settings.Sms
		smsClient, err := storageSms.NewSms(nil, &storageSms.SmsOption{
			SecretId:   opt.SecretId,
			SecretKey:  opt.SecretKey,
			Region:     opt.Region,
			SdkAppId:   opt.SdkAppId,
			SignName:   opt.SignName,
			TemplateId: opt.TemplateId,
		})
		return smsClient, err
	}
	return nil, nil
}
