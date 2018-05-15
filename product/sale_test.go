package product

import (
	_ "ainit"
	"log"
	"testing"
)

func TestNews(t *testing.T) {

	sale := NewArrivals()
	sale.Add(94, 1)
	sale.Add(77, 2)
	sale.Add(76, 3)
	sale.Rem(76)
	if _, total := sale.Listing(0, 10); total != 2 {
		t.Errorf("sale %d != 2", total)
	} else {
		log.Println("TestNews:OK")
	}

}
