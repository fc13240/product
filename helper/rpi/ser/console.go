package ser

import (
	"helper/rpi"
	"fmt"
	"net/http"
	"time"
	"log"
	"github.com/gorilla/websocket"
)


func Console(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println("conn error:", err)
		return
	}

	log.Println("console ....")


	defer func() {
		conn.Close()
		if p := recover(); p != nil {
			err = fmt.Errorf("internal error: %v", p)
			fmt.Println("error:",err)
		}
	}()

	for {
		time.Sleep(time.Second *1)
		customers:=len(rpi.Customers)
		tanks:=len(rpi.Rasps)
		body:=fmt.Sprint("控制端数量:",customers,"坦克数量:",tanks)


		err:=conn.WriteMessage(websocket.TextMessage,[]byte(body))
		if err!=nil{
			log.Println(err)
			return
		}
	}
}
