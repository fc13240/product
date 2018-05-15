package main

import (
	_ "ainit"
	"product/ezbuy"
	"time"


)
func main(){

	go func(){

		for {

			ezbuy.RefreshBy(7807720,7907474,7816659,7782742,7782800,7867702,7797404)
			time.Sleep(time.Second*240)
		}
	}()

	ezbuy.RefreshAll()
}
