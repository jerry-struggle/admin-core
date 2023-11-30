package config

type WxPay struct {
	AppId                      string
	MchId                      string
	MchCertificateSerialNumber string
	PrivateKeyPath             string
	NotifyUrl                  string
	MchAPIv3Key                string
}

var WxPayConfig = new(WxPay)
