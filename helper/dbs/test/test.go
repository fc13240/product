package main

import (
	"fmt"
	"helper/dbs"
)

func main() {
	var v string
	err := dbs.One("select name from wsns_groups").Scan(&v)
}
