package rpi

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/websocket"

	"helper/configs"
	conn "helper/redisCli"
	"helper/dbs/mongodb"
)

var r = conn.Conn()

var (
	col *mongodb.Collection
)

//树莓派配置表
func Col() *mongodb.Collection{
	if col == nil{
		col=mongodb.Conn().C("rasp.gpio")
	}
	return col
}

const (
	Reverse = "Reverse"
	Parking = "Parking"
	Driving = "Driving"
)

type Sayer interface {
	Say(text string)
	Send(cmd string, s interface{})
}

type Tank struct {
	customers map[int]*Customer
	Name      string `json:"name"`
	speed     int    //1-100速度
	Link      *websocket.Conn
	Gear      string `json:"gear"`
	In1       string `json:"in1"`
	In2       string `json:"in2"`
	In3       string `json:"in3"`
	In4       string `json:"in4"`
	En1       string `json:"en1"`
	En2       string `json:"en2"`
	Icon      string `json:"icon"`
}

func (tank *Tank) Speed(i int) {
	tank.speed = i
	tank.Send("speed", configs.M{"lspeed": i, "rspeed": i})
}

func (tank *Tank) Parking() {
	tank.Gear = Parking
	tank.Speed(0)
}

func (tank *Tank) Register() (err error) {
	if exist, _ := redis.Bool(r.Do("exists", fmt.Sprint("tank:", tank.Name))); exist {
		return errors.New("error : this name is exist.")
	}
	return tank.Save()
}

func (tank *Tank) Save() (err error) {
	_, err = r.Do("hmset", fmt.Sprint("tank:", tank.Name), "in1", tank.In1, "in2", tank.In2, "in3", tank.In3, "in4", tank.In4, "en1", tank.En1, "en2", tank.En2,"name",tank.Name)
	return err
}


func (tank *Tank) Login(name string) (err error) {
	if ok, _ := redis.Bool(r.Do("exists", fmt.Sprint("tank:", name))); ok {
		if info, e := redis.StringMap(r.Do("hgetall", fmt.Sprint("tank:", name))); err == nil {
			tank.In1 = info["in1"]
			tank.In2 = info["in2"]
			tank.In3 = info["in3"]
			tank.In4 = info["in4"]
			tank.En1 = info["en1"]
			tank.En2 = info["en2"]
			tank.Name = name
			tank.Gear = Parking
			tank.customers = map[int]*Customer{}
			return nil
		} else {
			err = e
		}
	} else {
		tank.SendText("tank not exist ,register ....")
		r.Send("HMSET",fmt.Sprint("tank:",name),
			"name",name,
			"in1",0,
			"in2",0,
			"in3",0,
			"in4",0,
			"en1",0,
			"en2",0,
		)
		r.Flush()
		return nil
	}
	return err
}

func (tank *Tank) Exit() {

	tank.SendText("bye...")
	tank.Send("byebye", configs.M{"tank_name": tank.Name})
}

func (tank *Tank) AddCustomer(customer *Customer) {
	tank.customers[customer.Uid] = customer
}


func (tank *Tank) Send(cmd string, args configs.M) {
	for _, custome := range tank.customers {
		custome.Input(cmd,args)
	}
}

func (tank *Tank)Input(cmd string,args interface{}){
	data := configs.M{"cmmand": cmd, "args": args}
	tank.Link.WriteJSON(data)
}

func (tank *Tank) SendText(content string) {
	tank.Send("text",configs.M{"content":content})
}

