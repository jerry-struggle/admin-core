package runtime

import (
	"github.com/jerry-struggle/admin-core/storage"
)

// NewSms 创建对应上下文短信对象
func NewSms(sms storage.AdapterSms) storage.AdapterSms {
	return &Sms{
		sms: sms,
	}
}

type Sms struct {
	sms storage.AdapterSms
}

// SendSms 发送短信
func (s *Sms) SendSms(phoneNum, smsCode string) error {
	return s.sms.SendSms(phoneNum, smsCode)
}
