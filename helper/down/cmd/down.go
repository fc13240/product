package main

import (
	"fmt"
	"helper/configs"
	"helper/dbs"
	"helper/dbs/mongodb"
	"helper/down"
)

func init() {
	configs.Ini("D:/config.ini")
	mongodb.Conn()
	dbs.Conn()
}

func main() {
	fmt.Println("Start")
	down.RunDown(20)
}
