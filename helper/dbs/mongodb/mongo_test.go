package mongodb

import (
	_ "ainit"
	"fmt"
	"helper/configs"
	"testing"
)

func TestConn(t *testing.T) {
	col := Conn().C("test")
	col.Insert(&configs.M{"title": "lxg", "q": 2})
	total := 0
	col.ItemsCall(configs.M{}, func(result *Result) {
		for result.Next() {
			total++
			var item configs.M
			result.Scan(&item)
			result.UpSet(configs.M{"test": "oksssssssss"})

		}
	}, 1, 30)

	fmt.Println(total)
}
