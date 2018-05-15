package crypto
import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"log"
	"bytes"
	"errors"

)

type Aes struct {
	Iv string
	Key string
	BlockSize int
	Mode string
	Padding string
}

//加密字符串
func (a *Aes) Encrypt(content []byte) ([]byte, error) {
	key := []byte(a.Key)

	var iv = []byte(key)[:aes.BlockSize]

	c, err := aes.NewCipher(key)

	if err != nil {
		err=errors.New(err.Error()+"new cipher error")
		return nil, err
	}


	switch a.Padding{
		case "PKCS5":
			content = PKCS5Padding(content, c.BlockSize())
		case "ZERO":
			content = ZeroPadding(content, c.BlockSize())
	}

	iv=[]byte(a.Iv)

	ciphertext := make([]byte,len(content))

	switch a.Mode{
	case "CBC":
		enc:= cipher.NewCBCEncrypter(c, iv)
		enc.CryptBlocks( ciphertext, content)
	case "CFB":
		enc:= cipher.NewCFBEncrypter(c, iv)
		enc.XORKeyStream( ciphertext, content)
	case "CTR":
		enc:=cipher.NewCTR(c, content)
		enc.XORKeyStream(ciphertext, content)
	case "OFB":
		ofb:=cipher.NewOFB(c, content)
		ofb.XORKeyStream(ciphertext, content)
	default:
		err=errors.New("not mode")
	}
	return ciphertext, err
}

func(a *Aes)EncryptTobase64(strMesg []byte)string{
	s,e:=a.Encrypt(strMesg)
	if e!=nil{
		log.Println(e)
		return ""
	}
	return base64.StdEncoding.EncodeToString(s)
}

func Base64DecodeString(s string )([]byte,error){
	return base64.StdEncoding.DecodeString(s)
}

func Base64EncodeToString(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

//解密字符串
func (a *Aes) Decrypt(content []byte) (strDesc []byte, err error) {
	defer func() {
		//错误处理
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	key:=[]byte(a.Key)
	iv:=[]byte(a.Iv)

	decrypted := make([]byte, len(content))
	var aesBlockDecrypter cipher.Block
	aesBlockDecrypter, err = aes.NewCipher([]byte(key))

	if err != nil {
		return decrypted, err
	}

	switch a.Mode {
		case "CBC":
			aesDecrypter := cipher.NewCBCDecrypter(aesBlockDecrypter, iv)
			aesDecrypter.CryptBlocks(decrypted, content)
		case "CFB":
			aesDecrypter := cipher.NewCFBDecrypter(aesBlockDecrypter, iv)
			aesDecrypter.XORKeyStream(decrypted, content)
	}

	switch a.Padding{
		case "PKCS5":
			decrypted = PKCS5UnPadding(decrypted)
	case "ZERO":
			decrypted = ZeroUnPadding(decrypted)
	}
	return decrypted, nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {

	padding := blockSize - len(ciphertext)%blockSize//需要padding的数目
	//只要少于256就能放到一个byte中，默认的blockSize=16(即采用16*8=128, AES-128长的密钥)
	//最少填充1个byte，如果原文刚好是blocksize的整数倍，则再填充一个blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)//生成填充的文本
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)//用0去填充
	return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
	return bytes.TrimFunc(origData,
		func(r rune) bool {
			return r == rune(0)
	})
}