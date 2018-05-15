package ser

import (
	"fmt"
	"net/http"
	"helper/configs"
	"helper/rpi"
	"log"
)

func Customer(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println("conn error:", err)
		return
	}

	re := Receive{link: conn, Args: configs.M{}}

	if err = conn.ReadJSON(&re); err != nil || re.Cmmand != "auth" {
		log.Println("验证 Fail")
		re.Say("验证失败")
		return
	}else{
		log.Println("验证 SUCC")
	}

	args := re.Args

	customer := &rpi.Customer{Link: conn}

	token := args.Get("token")

	if customer.Login(token); err != nil {
		log.Println("login fail :",err)
		re.Say(err.Error())
		return
	}else{
		log.Println("login succ.")
	}

	if len(rpi.Rasps)>0{ //如果有树莓派在线，就和这个用户绑定起来
		for rasp_name,_:=range rpi.Rasps{
			fmt.Println("conn:",rasp_name,rpi.Rasps)
			customer.ConnectRasp(rasp_name)
			break
		}
	}

	defer func(){
		if p := recover(); p != nil {
			err = fmt.Errorf("internal error: %v", p)
		}
		fmt.Println("exit:", err)
		customer.Exit()
	}()

	for {
		if err = conn.ReadJSON(&re); err != nil {
			log.Println("exit:", err)
			return
		}

		if err=customer.Send(re.Cmmand,re.Args);err!=nil{
			customer.Input(re.Cmmand,configs.M{"isFail":true,"error_msg":err.Error()})
		}

	}
}
