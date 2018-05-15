package alibaba

import (
	"fmt"
	_ "ainit"
	"testing"
	//"fmt"
	"helper/account"
	"helper/configs"
	"product/gallery"
)

func TestGet(t *testing.T){
	acc:=&account.Account{Uid:1}
	item,err:=Down("https://detail.1688.com/offer/547051334633.html?spm=a2604.8275163.ix056udb.3.awA7VG",acc)

	if err !=nil{
		t.Error(err)
		t.Fail()
	}

	if item.Id<1{
		t.Fail()
	}

	if item.Desc==""{
		t.Error("desc is empty")
	}
}

func TestGetAttrs(t *testing.T) {
	t.Log(GetAttrs(7665))
}


func TestSearch(t *testing.T){
	offset:=0
	limit:=10
	res,total:=Search(configs.M{"sku":"BL0"},offset,limit)
	fmt.Println("total:",total)
for _,row:=range res{
	fmt.Println(row.Sku,row.ItemId,row.Id)
}
}

func TestImprotImage(t *testing.T){
	fmt.Println("start")
	list,_:=Search(configs.M{},0,2000)
	for _,row:=range list{
		for _,image:=range row.Images{
			gallery.AddImage(row.Sku,image)
		}

	}
}