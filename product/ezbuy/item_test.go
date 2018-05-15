package ezbuy

import (
	_ "ainit"
	"testing"
	"log"
	"helper/account"
)

func TestExport(t *testing.T) {
	items:=[]*Item{}
	ItemsIds:=[]string{"test"}

	for _,sku:=range ItemsIds{
		if item:=Get(sku);item!=nil{
			items=append(items,item)
		}
	}
	if len(items)!=len(ItemsIds){
		t.Errorf("ids length %d not equal items %d",len(ItemsIds),len(items))
	}
	author:=account.Find(21)

	if len(items)>0{
		filename,err:=Export(author,1,items...)
		if err!=nil {
			t.Error(err)
		}else {
			log.Println("succ",filename)
		}
	}else{

	}
}

func TestItem_Save(t *testing.T) {
	item:=Get("sku1")
	if item == nil{
		t.Error("is empty")
		return
	}

	item.Name="bbbb"
	item.Save()

	item1:=Get("sku2")
	if item1.Name!=item.Name{
		t.Errorf("edit name fail . old:%s new:%s",item.Name,item1.Name)
	}else{
		log.Println("succ ",item.Name,item1.Name)
	}
}