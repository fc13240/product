package ser

import (
	"helper/rpi"
	"fmt"
	"net/http"
	"log"
)


func Rasp(w http.ResponseWriter, r *http.Request){
	conn, err := upgrader.Upgrade(w, r, nil)
	fmt.Println("conn..")
	if err != nil {
		log.Println("conn error:", err)
		return
	}

	rasp := &rpi.Rasp{Link: conn}

	defer func() { //end

		rpi.Rasps.Del(rasp)
		conn.Close()
	}()

	re := &Receive{link: conn}

	if err = conn.ReadJSON(&re); err != nil || re.Cmmand != "auth" {
		log.Println("auth fail:",err.Error())
		return
	}

	args := re.Args

	if args.Get("name") == "" {
		re.Say("error : name can not empty ")
		return
	}else{
		rasp.Name=args.Get("name")
	}

	rpi.Rasps.Add(rasp)

	rasp.SendText( "hello,i am coming up.")

	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("internal error: %v", p)
			fmt.Println("error:",err)
			rasp.Exit()
		}
	}()

	for {
		re = &Receive{link: conn}
		if err = conn.ReadJSON(re); err != nil {
			fmt.Println(err)
			return
		}

		args = re.Args
		rasp.Send(re.Cmmand,re.Args)
	}
}