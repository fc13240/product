package main

import (
	_ "ainit"

	"product/ezbuy"
	"fmt"

)

func main() {
	defer func(){
		if r:=recover();r!=nil{
			fmt.Println(r)
		}
	}()
	ezbuy.GetOrders(80)
}