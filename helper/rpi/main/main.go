package main

import (
	_ "ainit"
	"flag"
	"fmt"
	"helper/rpi/ser"
	"net/http"

)

var addr = flag.String("localhost", ":1811", "websock")

func main() {
	flag.Parse()
	http.HandleFunc("/rasp", ser.Rasp)
	http.HandleFunc("/manage", ser.Customer)
	http.HandleFunc("/console",ser.Console)
	fmt.Println("Start OK")
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		fmt.Println(err)
	}
}
