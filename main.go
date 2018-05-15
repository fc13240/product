package main
import (
	_ "ainit"
	_ "routers"
	"helper/net/igin"
	"helper/configs"
	
	"fmt"
	"log"
	_ "helper/mongodb/api"
)
func main(){
	if err:=recover();err!=nil{
		log.Println("internal error:", err)
	}
	var r=igin.R
	opt:=configs.GetSection("admin")
	port:=opt["port"] 
	r.Run(fmt.Sprint(":",port))
}