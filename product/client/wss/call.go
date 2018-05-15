package wss

import (
	"time"
	"fmt"
	"log"
)

var option map[string]string
func CheckNewOrders(){

	if ! setting.CheckNewOrders || runSta.CheckNewordersRunSta == On {
		log.Println("已经运行刷新")
		return
	}

	defer  func(){
		if r:=recover();r!=nil{
			runSta.CheckNewordersRunSta=Off
			log.Println("更新订单失败",r)
		}
	}()
	runSta.CheckNewordersRunSta=On
	for{
		ez.CheckNewOrders()
		time.Sleep(5*time.Minute)
	}
}

func Refresh(){
	if ! setting.OnRefresh || runSta.RefreshRunSta == On {
		log.Println("已经运行刷新")
		return
	}

	defer func(){
		if r:=recover();r!=nil{
			isRefresh=false
			log.Println(r)
		}
	}()
	runSta.RefreshRunSta=On
	ez.RefreshAll(setting)
}

func PrintRunStat(){
	fmt.Println("check New Orders :",runSta.CheckNewordersRunSta)
	fmt.Println("Refresh Items:",runSta.RefreshRunSta)
}