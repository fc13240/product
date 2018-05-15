package gallery

import(
	_ "ainit"
	"fmt"
	"testing"
	"helper/configs"
)

func TestList(t *testing.T){
	rows,_:=Listing(configs.M{"sku":"BL18040107"},0,10)
	fmt.Println(rows)
}

func TestDelSkuImage(t *testing.T){
	if err:=DelSkuImage("BL18042953",[]int{9324});err!=nil{
		t.Error(err)
	}
}

func TestSetFlag(t *testing.T){
	SetFlag([]int{1604},CoverImageFlage)
}