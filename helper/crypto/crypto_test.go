package crypto

import (
	"testing"
	"fmt"
)

func TestAes_Encrypt(t *testing.T) {

	aes:=Aes{
		Iv:"0123456789abcdef",
		Key:"0123456789abcdef",
		Mode:"CBC",
		Padding:"PKCS5",
	}

	b,_:=aes.Encrypt([]byte("123456789456asdfasdddddddddddddddddddddddddddddddddd"))
	fmt.Println(b)
	content,_:=aes.Decrypt(b)

	fmt.Println(string(content))
	//bb:=aes.EncryptTobase64([]byte("123"))
	//fmt.Println(bb)
	//c,_:=aes.Decrypt(b)

	//fmt.Println(c)
}
