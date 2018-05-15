package rpi

import (
	"github.com/gorilla/websocket"
	"helper/configs"
)
const(
	IN =1
	OUT =0
	SPI = 41
	I2C = 42
	HARD_PWM = 43
	SERIAL = 40
	UNKNOWN = -1
)

type Rasp struct{
	Link      *websocket.Conn
	Name string
}

func (r *Rasp) Send(cmd string, args configs.M) {
	for _, custome := range Customers {
		custome.Input(cmd,args)
	}
}

func (r *Rasp)Input(cmd string,args interface{}){
	data := configs.M{"cmmand": cmd, "args": args}
	r.Link.WriteJSON(data)
}

func (r *Rasp) SendText(content string) {
	r.Send("text",configs.M{"content":content})
}

func(r *Rasp)NewGpio(pin int )*Gpio{
	return &Gpio{R:r,Pin:pin}
}

func (r *Rasp) Exit() {
	delete(Rasps,r.Name)
	r.SendText("bye...")
}

type Gpio struct {
	R *Rasp
	Pin int
}

func (g *Gpio)Setup(f int){
	g.R.Send("tank.setup",configs.M{"pin":g.Pin,"val":f})
}

func (g *Gpio)Output(t bool){
	g.R.Send("tank.output",configs.M{"pin":g.Pin,"val":t})
}