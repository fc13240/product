package qiniu

import (
	"bytes"
	"io"
	"qiniupkg.com/api.v7/conf"
	"qiniupkg.com/api.v7/kodo"
	"qiniupkg.com/api.v7/kodocli"

)

const (
	AccessKey = "siA16n8i4wWgUoOjRHh0hsudB2BLn0qoocJ0FT25"
	SecretKey = "WsplJnYBc8ZpWOx_Zph0i5xlxYkKnSHO8l-B7jkI"
)

var (
	//设置上传到的空间
	bucket = "product-images"
	token  string
)

type Store struct {
	File     io.Reader
	FileSize int64
}

func (store *Store) genToken() string {
	if token != "" {
		return token
	}
	//初始化AK，SK
	conf.ACCESS_KEY = AccessKey
	conf.SECRET_KEY = SecretKey

	//创建一个Client
	c := kodo.New(0, nil)

	//设置上传的策略
	policy := &kodo.PutPolicy{
		Scope: bucket,
		//设置Token过期时间
		Expires: 31536000,
	}
	//生成一个上传token
	token = c.MakeUptoken(policy)

	return token
}

func (store *Store) Save(file io.Reader) (string, error) {
	//构建一个uploader
	zone := 0
	uploader := kodocli.NewUploader(zone, nil)

	var ret kodocli.PutRet

	buff := &bytes.Buffer{}
	filesize, _ := buff.ReadFrom(file)

	err := uploader.PutWithoutKey(nil, &ret, store.genToken(), buff, filesize, nil)

	if err != nil {
		return "", err
	}

	return ret.Hash, nil
}


