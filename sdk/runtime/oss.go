package runtime

import (
	"github.com/jerry-struggle/admin-core/storage"
	sts "github.com/tencentyun/qcloud-cos-sts-sdk/go"
)

// NewSms 创建对应上下文短信对象
func NewOss(oss storage.AdapterOss) storage.AdapterOss {
	return &Oss{
		oss: oss,
	}
}

type Oss struct {
	oss storage.AdapterOss
}

func (o *Oss) UpLoad(objectName, localFile string) (string, error) {
	return o.oss.UpLoad(objectName, localFile)
}

func (o *Oss) GetCredential() (*sts.CredentialResult, error) {
	return o.oss.GetCredential()
}

func (o *Oss) GetPresignedURL(objectName string) (string, error) {
	return o.oss.GetPresignedURL(objectName)
}

func (o *Oss) GetSignedFileUrl(dir string, isSign bool) ([]string, []string, error) {
	return o.oss.GetSignedFileUrl(dir, isSign)
}

func (o *Oss) DownloadFiles(fileObject []string, localDir string) error {
	return o.oss.DownloadFiles(fileObject, localDir)
}
