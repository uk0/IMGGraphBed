package main

import (
	"context"
	"fmt"
	"github.com/qiniu/api.v7/auth/qbox"
	"testing"

	"github.com/qiniu/api.v7/storage"
)

var (
	accessKey = "9Zctl1ro_QMLg3TaQPSCF8"
	secretKey = "4syYqkVMgvzeSaUThH28A0V"
	bucket    = "bucket-images-firsh"
)

func TestPutQiniu(t *testing.T) {

	localFile := "/Users/zhangjianxin/Pictures/WechatIMG415.png"
	key := "1111111-222222-333333-44444.png"
	putPolicy := storage.PutPolicy{
		Scope:               bucket,
	}
	mac := qbox.NewMac(accessKey, secretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = &storage.ZoneHuabei
	fmt.Println(cfg.Zone)
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	// 构建表单上传的对象
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	// 可选配置
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": "github logo",
		},
	}
	err := formUploader.PutFile(context.Background(), &ret, upToken, key, localFile, &putExtra)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ret.Key,ret.Hash)
}