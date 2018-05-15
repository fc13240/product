package wss


import (
	"fmt"
	"errors"
	"github.com/gorilla/websocket"
	"helper/configs"
)

type Clients map[int]*websocket.Conn
var clients =Clients{}

var client_num =0

func(data Clients)Add(client *websocket.Conn)int {
	pid:=client_num
	data[pid]=client
	client_num++
	return pid
}

func (data Clients)Del(index int){
	delete(data,index)
}

type Message struct {
	Client *websocket.Conn
	Cmd   *Command `json:"cmd"`
	IsOpen bool
	Routes map[string]func(m *Message)
}

type Command struct {
	Name string   `json:"cmd"`
	Args configs.M `json:"args"`
}

func (cmd *Command) Get(k string) string {
	return cmd.Args.Get(k)
}

func (cmd *Command) Run(m *Message) {
	if _, ok := m.Routes[cmd.Name]; ok {
		m.Routes[cmd.Name](m)
	} else {
		m.SendTxt(fmt.Sprint("无效的命令:", cmd.Name))
	}
}

func (m *Message) Wait() (cmd *Command, err error) {

	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("internal error: %v", p)
		}
	}()

	m.Cmd = &Command{}

	if m.IsOpen == false {
		return m.Cmd, errors.New("连接已断开")
	}

	err = m.Client.ReadJSON(m.Cmd)
	if err == nil {

	}
	return m.Cmd, err
}

func (m *Message) Send(ss interface{}) {
	data := configs.M{ "args": ss,"cmd":m.Cmd.Name}
	for _,client:=range clients{
		client.WriteJSON(data)
	}
}

func (m *Message) SendTxt(s ...interface{}) {
	m.Send(configs.M{"txt":fmt.Sprint(s...)})
}

func (m *Message) SendEnd(s ...interface{}) {
	m.Send(configs.M{"end":fmt.Sprint(s...)})
}

func (m *Message) Fail(s string) {
	m.Send(configs.M{"isFail":true,"error_msg":s})
}

func (m *Message) Succ(dd interface{}){
	m.Send(dd)
}
