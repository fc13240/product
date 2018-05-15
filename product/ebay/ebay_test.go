package ebay

import (
	_ "ainit"
	"product"
	"testing"
)

func TestBindEbay(t *testing.T) {
	item_id := 393
	if item, err := product.Get(item_id); err == nil {
		if err := item.BindEbay(123456); err != nil {
			t.Error(err)
		} else {
			t.Log("Yes OK very good")

		}
	}
}
