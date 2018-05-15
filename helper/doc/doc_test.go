package doc

import (
	_"ainit"
	"testing"
	"helper/configs"
	"fmt"
)

func TestFavListing(t *testing.T) {
	rows,total:=FavListing(configs.M{"author_id":21},0,5)
fmt.Println(rows)
fmt.Println(total)
}
