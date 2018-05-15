package title

import (
_ "ainit"
"testing"
	"fmt"
	"strings"
)

func TestNewExamples(t *testing.T) {
	NewExamples(
		"Women's Beach Crochet Backless Bohemian Halter Maxi Long Dress",
		"女子沙滩钩针露背长裙波西米亚露背长裙",
		[]string{"https://img.alicdn.com/imgextra/i4/926500743/TB2r9G4lpXXXXavXXXXXXXXXXXX_!!926500743.jpg"},
	)


}

func TestNewLoabel(t *testing.T) {
	NewLabel("Long Dress","长礼服,长裙",100,0)
	NewLabel("Backless Bohemian","露背放荡不羁",100,0)
	NewLabel("Women's Beach","女子沙滩",100,0)
	NewLabel("Women's Slimming Shirt","女式瘦身衬衫",100,0)
	NewLabel("Backless Dress","露背裙",100,0)
	NewLabel("Women's Slimming 3/4 Sleeve Dress","女士减肥3/4袖连衣裙",100,0)
}


func TestSearchTitle(t *testing.T){
	fmt.Println(SearchTitle("dress") )
}

func TestSearchCnTitle(t *testing.T){
	fmt.Println(SearchCnTitle("减肥") )
}

func TestImport(t *testing.T){
body:=`Dots Blouse
圆点 女衬衫

O Neck Crop Tops
圆领上衣

Irregular Cool T-Shirts
不规则 凉爽 T恤

Dots Blouse,Han Shi Fashion  Women O Neck Crop Tops 3/4 Sleeve  Irregular Cool T-Shirts
点女衬衫，韩时时尚女装O领裁缝上衣3/4袖不规则酷T恤`

	ss:=strings.Split(body,"\n\n")
	for _,s:=range ss{
		row:=strings.Split(s,"\n")
		NewLabel(row[0],row[1],1,0)
	}
}

func TestNewLoabelCate(t *testing.T){
	NewLabelCate(100,"袖子")
}

func TestLabels(t *testing.T){
	fmt.Println(GetLabels(100))
}
