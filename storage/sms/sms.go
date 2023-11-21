package sms

import (
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111" // 引入sms
)

type SmsOption struct {
	SecretId   string
	SecretKey  string
	Region     string
	SdkAppId   string
	SignName   string
	TemplateId string
}

func NewSms(cliSms *sms.Client, options *SmsOption) (*Sms, error) {
	var err error
	s := &Sms{Client: cliSms, options: options}
	credential := common.NewCredential(
		options.SecretId,
		options.SecretKey,
	)
	cpf := profile.NewClientProfile()
	s.Client, err = sms.NewClient(credential, options.Region, cpf)
	if err != nil {
		return nil, err
	}
	return s, nil
}

type Sms struct {
	Client  *sms.Client
	options *SmsOption
}

// 发送短信
func (s *Sms) SendSms(phoneNum, smsCode string) error {

	request := sms.NewSendSmsRequest()
	request.SmsSdkAppId = common.StringPtr(s.options.SdkAppId)
	request.SignName = common.StringPtr(s.options.SignName)
	request.TemplateId = common.StringPtr(s.options.TemplateId)
	request.TemplateParamSet = common.StringPtrs([]string{smsCode})
	request.PhoneNumberSet = common.StringPtrs([]string{phoneNum})
	_, err := s.Client.SendSms(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return err
	}
	return nil
}
