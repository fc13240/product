package wss

import (
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
	"log"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}


type Call struct{
	funs map[string]func(m *Message)
}

func (c *Call)Add(name string ,call func(m *Message)){
	if c.funs==nil{
		c.funs=map[string]func(m *Message){}
	}
	c.funs[name]=call
}

func (c *Call)Funs()map[string]func(m *Message){
	return c.funs
}


func EzStart(w http.ResponseWriter, r *http.Request){
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println("conn error:", err)
		return
	}

	fmt.Println("ez client succ")
	ezMessage = &Message{Client: conn,IsOpen: true,Routes:ezcall.Funs()}

	for {
		if cmd, err := ezMessage.Wait(); err == nil {
			cmd.Run(ezMessage)
		} else {
			fmt.Println(err.Error())
			break
		}
	}

	defer func(){
		if err:=recover();err!=nil{
			log.Println(err)
			ezMessage=nil
		}
	}()
}


func Alibaba(w http.ResponseWriter, r *http.Request){
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("conn error:", err)
		return
	}

	pid:=clients.Add(conn)
	aliMessage = &Message{Client: conn,IsOpen:true,Routes:albabaCall.Funs()}

	for {
		if cmd, err := aliMessage.Wait();err == nil{
			fmt.Println("call:",aliMessage.Cmd)
			cmd.Run(aliMessage)
		} else {
			clients.Del(pid)
			fmt.Println(err.Error())
			break
		}
	}

	defer func(){
		if err:=recover();err!=nil{
			log.Println(err)
			aliMessage=nil
		}
	}()
}

func call(w http.ResponseWriter, r *http.Request){
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("conn error:", err)
		return
	}

	pid:=clients.Add(conn)
	aliMessage = &Message{Client: conn,IsOpen:true,Routes:albabaCall.Funs()}
	
	for {
		if cmd, err := aliMessage.Wait();err == nil{
			cmd.Run(aliMessage)
		} else {
			clients.Del(pid)
			fmt.Println(err.Error())
			break
		}
	}

	defer func(){
		if err:=recover();err!=nil{
			log.Println(err)
			aliMessage=nil
		}
	}()
}

func Start(w http.ResponseWriter, r *http.Request){
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println("conn error:", err)
		return
	}

	pid:=clients.Add(conn)

	message := &Message{Client: conn, IsOpen: true,Routes:hcall.Funs()}
	client_num++
	for {
		if cmd, err := message.Wait(); err == nil {
			cmd.Run(message)
		} else {
			clients.Del(pid)
			break
		}
	}


	defer func(){
		if err:=recover();err!=nil{

		}

	}()
}


