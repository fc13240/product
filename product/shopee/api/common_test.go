package api

import(
	"testing"
	"fmt"
	_ "ainit"
	"helper/dbs"
	"helper/configs"
)

func TestA(t *testing.T){
	Listing()
}

func TestAdd(t *testing.T){
	Add()
}

func TestGetCategorys(t *testing.T){
	cates:=GetCategorys()
	db:=dbs.Def()
	for _,cate:=range cates{
		db.Insert("shopee_category",configs.M{"parent_id":cate.ParentId,"has_children":cate.HasChildren,"category_id":cate.Id,"category_name":cate.Name})
	}
}

func TestGetAttributes(t *testing.T){
	GetAttributes(6588)
}

func TestGetLogistics(t *testing.T){
	GetLogistics()
}

func TestAddItem(t *testing.T){
	fmt.Println(AddItem("BL18040107"))
}

func TestGetItem(t *testing.T){
	GetItem(1073714213)
}

func TestDel(t *testing.T){
	DelItem("Bl2000162")
}