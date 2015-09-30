package qiniu

import (
	"golang.org/x/net/context"
	"heroku.com/tg-bot/common"
	"io"
	"qiniupkg.com/api.v7/kodo"
)

var config common.Config

func init() {
	config = common.GetConfig()
}

func UploadFile(fileName string, fileSize int, r io.Reader) (ok bool) {

	ok = false
	kodo.SetMac(config.QiniuAccessKey, config.QiniuSecretKey)
	zone := 0
	c := kodo.New(zone, nil)
	bucket := c.Bucket(config.QiniuBucketName)
	ctx := context.Background()

	fileSizeIn64 := int64(fileSize)
	err := bucket.Put(ctx, nil, fileName, r, fileSizeIn64, nil)
	if err == nil {
		ok = true
	}
	return
}
