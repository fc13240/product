package ezbuy
import (
	_ "ainit"
	"testing"
	"log"
	"helper/account"
)

func TestGetCategorys(t *testing.T) {
	author:=account.Find(21)
	store:=GetSetting(author,1)

	categorys:=GetCategorys(store.StoreCateId)
	if len(categorys)==0{
		t.Error(store.StoreCateId,"erro empty")
	}else{
		log.Println("succ len:",len(categorys))
	}
}
