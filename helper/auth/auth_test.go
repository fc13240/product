package auth

import (
	"testing"
	_"ainit"
	"fmt"
	"time"
)
func TestNewToken(t *testing.T){
	i:=0
	for {
	token:=NewToken()
	fmt.Println(token)
	i++
		token.Set("abcd","hello")
		fmt.Println(token.Get("abcd"))
		time.Sleep(time.Microsecond*100)
		fmt.Println(i)
	}
}

func TestLogin(t *testing.T){

	token:=NewToken()
	token.Set("uid",10)
	if false == token.IsLogin(){
		t.Error("fail")
	}

	token=NewToken()
	if false == token.IsLogin(){
		t.Error("is ok")
	}
}
