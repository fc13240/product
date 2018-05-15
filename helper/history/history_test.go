package history

import (
	_ "ainit"
	"testing"
	"helper/account"
)

func TestAdd(t *testing.T){
	author:=account.Find(21)
	var his Flag ="article"
	t.Log(his.Add(author,"http://baidu.com","baidu"))
}