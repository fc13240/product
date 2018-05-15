package crypto

import (
	"crypto/md5"
	"encoding/hex"
	"crypto/hmac"
	"crypto/sha256"
)

func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

type HHmac struct{
	Content string 
	Key string
}

func Hmac(c,k string) *HHmac{
	return &HHmac{c,k}
}


func (h *HHmac) Sha256()string{
	mac := hmac.New(sha256.New,[]byte(h.Key))
	mac.Write([]byte( h.Content))
	return hex.EncodeToString(mac.Sum(nil))
}

func (h *HHmac) Md5()string{
	mac := hmac.New(md5.New,[]byte(h.Key))
	mac.Write([]byte( h.Content))
	return hex.EncodeToString(mac.Sum(nil))
}
//sha1


