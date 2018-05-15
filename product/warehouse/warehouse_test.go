package warehouse

import (
	_ "ainit"
	"product"
	"testing"
)

func TestDown(t *testing.T) {
	product.DownToWarehouse("https://detail.1688.com/offer/549556556070.html?spm=a2604.8109329.204.3.p5H639")
}

func TestDetail(t *testing.T) {
	pro, e := product.GetWarehouseProductInfo(1990)
	if e != nil {
		t.Error(e)
	} else {
		if len(pro.Images) == 0 {
			t.Error("没有图片")
		}

		if len(pro.SukImages) == 0 {
			t.Error("没有SKU图片")
		}
	}

}
