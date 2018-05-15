package ezbuy

import (
	"bytes"


	"helper/account"

	"qiniupkg.com/api.v7/kodocli"
	"os"
)



func  UploadExec(author *account.Account,items ...*Item) (string, error) {

	filename,err:=Export(author,1,items...)

	token:="LljRgfC0Wlx7ScBAhP_g78WXTYOiuGViAiCMw76V:ymm1fxMxHNEA0Bj3ZVOZPwBrQlg=:eyJzY29wZSI6ImRhb2dvdSIsImRlYWRsaW5lIjoxNTA5OTY5MDI2LCJ1cGhvc3RzIjpbImh0dHA6Ly91cC5xaW5pdS5jb20iLCJodHRwOi8vdXBsb2FkLnFpbml1LmNvbSIsIi1IIHVwLnFpbml1LmNvbSBodHRwOi8vMTgzLjEzMS43LjE4Il19"
	//构建一个uploader
	zone := 0
	file,_:=os.Open(filename)

	uploader := kodocli.NewUploader(zone, nil)

	var ret kodocli.PutRet

	buff := &bytes.Buffer{}

	filesize, _ := buff.ReadFrom(file)

	err = uploader.PutWithoutKey(nil, &ret,token, buff, filesize, nil)

	if err != nil {
		return "", err
	}

	return ret.Hash, nil
}


