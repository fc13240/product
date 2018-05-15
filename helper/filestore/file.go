package filestore

import (
	"helper/filestore/api/qiniu"
	"io"
)

var store qiniu.Store

type FileStoreer interface {
	Save(file io.Reader) (string, error)
}

func Save(file io.Reader) (string, error) {
	return store.Save(file)
}