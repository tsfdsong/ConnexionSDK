package kodo

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/logger"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/sirupsen/logrus"
)

func KodoUpload(accessKey string, secretKey string, bucket string, publisherID uint64, gameID uint64, email string, cover bool, d []byte) (string, error) {
	originKey := fmt.Sprintf("%s-%s-%s", publisherID, gameID, email)
	data := []byte(originKey)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	scope := bucket
	if cover {
		scope = fmt.Sprintf("%s:%s", bucket, md5str)
	}

	putPolicy := storage.PutPolicy{
		Scope: scope,
	}
	mac := qbox.NewMac(accessKey, secretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	cfg.Zone = &storage.ZoneHuadong
	cfg.UseHTTPS = false
	cfg.UseCdnDomains = false
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"PublisherName": fmt.Sprintf("%s", publisherID),
			"GameID":        fmt.Sprintf("%s", gameID),
			"Email":         email,
		},
	}

	dataLen := int64(len(d))
	err := formUploader.Put(context.Background(), &ret, upToken, md5str, bytes.NewReader(d), dataLen, &putExtra)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error(), "PublisherID": publisherID, "GameID": gameID, "Email": email}).Error("upload kodo failed")
		return "", err
	}
	return ret.Key, nil
}

func KodoPublicGet(domain string, key string) string {
	accessURL := storage.MakePublicURL(domain, key)
	return accessURL
}
