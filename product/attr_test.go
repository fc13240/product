package product
import (
	_ "ainit"
	"testing"
	"fmt"
)

func TestNewAtt(t *testing.T){
	if i,err:=NewAtt("服装尺寸","clothing_size",INPUT_CHECKBOX);err!=nil{
		t.Error(err.Error())
	}else{
		t.Error("NewId:",i)
	}
}

func TestAttrList(t *testing.T){
	att:=GetAttr(1)
	fmt.Println(att.GetOptions())
}

func TestCategoryAttrs(t *testing.T){
	fmt.Println(GetCategoryAttrs("local",100,false))
}

func TestItemAttr(t *testing.T){
	if item,err:=Get("MBL1");err==nil{
		fmt.Println(item.GetOptVal(Color))
	}
}

func TestGetSkuSelectedOptions(t *testing.T){
	opts:=GetSkuSelectedOptions("Bl20002394",nil)
	fmt.Println(opts)
}

func TestS(t *testing.T){
	attr:=GetAttr(2)
	opts:=attr.GetOptions()
	attr.FlagCheckedOption("BL18041784")
	for _,opt:=range opts{
		fmt.Println(opt.Value,opt.IsChecked)
	}
}
