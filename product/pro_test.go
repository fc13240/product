package product

import (
	_ "ainit"
	"helper/configs"
	"testing"
)

func TestEditItem(t *testing.T) {
	if item, err := Get("BBBB"); err == nil {
		item.Quant = 100
		item.Save()

	}
}



func TestItems(t *testing.T) {
	items, total :=Search(configs.M{}, 0, 10,"")

	t.Log("total:", total)

	if len(items) != 10 {
		t.Errorf("items len（%d） != %d", len(items), 10)
	}
}

func TestDetail(t *testing.T) {
	item_id := 59
	if item, err :=Get("BLABC"); err == nil {
		if item.Id != item_id {
			t.Errorf("item.id(%d)!=item_id(%d)", item.Id, item_id)
		} else {
			t.Log(err)
		}

	}

}
