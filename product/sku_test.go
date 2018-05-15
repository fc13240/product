package product
import (
	"testing"
	_ "ainit"
	"helper/account"
	"fmt"
)

func TestGetSkuUseStatus(t *testing.T){
	store:=NewStore()
	store.StoreId=1
	store.StoreName="Hello"
	store.StoreSite="EZBUY"
	if err:=store.Save() ;err!=nil{
		t.Error(err)
	}
	author:=account.Find(21)
	//AddSkuUploadLog(author,"BL100",1)
	fmt.Println(GetSkuUseStatus(author,"BL100"))
}

func TestNewSku(t *testing.T){
	sku,_:=NewSku()
	sku.AddBuySite("http://baidu.com")
}