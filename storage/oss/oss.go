package oss

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"
	sts "github.com/tencentyun/qcloud-cos-sts-sdk/go"
)

type OssOption struct {
	SecretId   string
	SecretKey  string
	AppId      string
	Bucket     string
	Region     string
	ExpireTime int64
}

func NewOss(ossCli *cos.Client, options *OssOption) *OSS {
	// 默认链接有效期5年
	if options.ExpireTime <= 0 {
		options.ExpireTime = 5 * 12 * 30 * 24
	}
	oss := &OSS{Client: ossCli, option: options}
	u, _ := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", options.Bucket, options.Region))
	b := &cos.BaseURL{BucketURL: u}
	oss.Client = cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  options.SecretId,
			SecretKey: options.SecretKey,
		},
	})

	return oss
}

type OSS struct {
	Client *cos.Client
	option *OssOption
}

// UpLoad 文件上传
func (e *OSS) UpLoad(objectName, localFile string) (string, error) {
	_, err := e.Client.Object.PutFromFile(context.Background(), objectName, localFile, nil)
	if err != nil {
		log.Println("Error:", err)
		return "", err
	}
	presignedURL, err := e.Client.Object.GetPresignedURL(
		context.Background(), http.MethodGet, objectName, e.option.SecretId, e.option.SecretKey, time.Duration(e.option.ExpireTime)*time.Hour, nil)
	if err != nil {
		log.Println("Error:", err)
		return "", err
	}
	return presignedURL.String(), nil
}

// 获取临时上传密钥
func (e *OSS) GetCredential() (*sts.CredentialResult, error) {

	c := sts.NewClient(
		e.option.SecretId,
		e.option.SecretKey,
		nil,
		// sts.Host("sts.tencentcloudapi.com"), // 设置域名, 默认域名sts.tencentcloudapi.com
		// sts.Scheme("http"),      // 设置协议, 默认为https，公有云sts获取临时密钥不允许走http，特殊场景才需要设置http
	)
	opt := &sts.CredentialOptions{
		DurationSeconds: int64(time.Hour.Seconds()),
		Region:          e.option.Region,
		Policy: &sts.CredentialPolicy{
			Statement: []sts.CredentialPolicyStatement{
				{
					Action: []string{
						"name/cos:PostObject",
						"name/cos:PutObject",
					},
					Effect: "allow",
					Resource: []string{
						//这里改成允许的路径前缀，可以根据自己网站的用户登录态判断允许上传的具体路径，例子： a.jpg 或者 a/* 或者 * (使用通配符*存在重大安全风险, 请谨慎评估使用)
						"qcs::cos:" + e.option.Region + ":uid/" + e.option.AppId + ":" + e.option.Bucket + "/*",
					},
				},
			},
		},
	}
	return c.GetCredential(opt)
}

// 获取签名链接
func (e *OSS) GetPresignedURL(objectName string) (string, error) {
	presignedURL, err := e.Client.Object.GetPresignedURL(
		context.Background(), http.MethodGet, objectName, e.option.SecretId, e.option.SecretKey, time.Duration(e.option.ExpireTime)*time.Hour, nil)
	if err != nil {
		log.Println("Error:", err)
		return "", err
	}
	return presignedURL.String(), nil
}

// 下载文件
// fileObject - cos文件对象
// localDir - 本地存储路径
func (e *OSS) DownloadFiles(fileObject []string, localDir string) error {

	if len(fileObject) <= 0 {
		return errors.New("无效的文件对象")
	}
	if len(localDir) <= 0 {
		return errors.New("无效的存储路径")
	}

	opt := &cos.MultiDownloadOptions{
		ThreadPoolSize: 5,
	}

	for _, fileUrl := range fileObject {
		filePath := localDir + path.Base(fileUrl)
		log.Println("info fileUrl=", fileUrl)
		log.Println("info filePath=", filePath)
		_, err := e.Client.Object.Download(context.Background(), fileUrl, filePath, opt)
		if err != nil {
			return err
		}
	}
	return nil
}

// 查询目录下的文件
func (e *OSS) GetSignedFileUrl(dir string, isSign bool) ([]string, []string, error) {
	var (
		results     []string
		signResults []string
	)
	opt := &cos.BucketGetOptions{
		Prefix: dir,
		// Delimiter: "/",
		MaxKeys: 1000,
	}
	result, _, err := e.Client.Bucket.Get(context.Background(), opt)
	if err != nil {
		log.Println("Error:", err)
		return nil, nil, err
	}
	for _, content := range result.Contents {
		log.Println("info:", content)
		presignedURL, err := e.Client.Object.GetPresignedURL(
			context.Background(), http.MethodGet, content.Key, e.option.SecretId, e.option.SecretKey, time.Duration(e.option.ExpireTime)*time.Hour, nil)
		if err != nil {
			log.Println("Error:", err)
			return nil, nil, err
		}
		results = append(results, presignedURL.String())
		signResults = append(signResults, content.Key)
	}
	return results, signResults, nil
}
