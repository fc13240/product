package rpi

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/websocket"
	"helper/account"
	"helper/configs"
)

type Customer struct {
	Tank *Tank
	Link *websocket.Conn
	Uid  int
	Nick string
	Rasp *Rasp
}

func (customer *Customer) Login(token string) (err error) {
	if token==""{
		err = errors.New("token is empty")
		return
	}
	token = fmt.Sprint("token:", token)

	if ok, err := redis.Bool(r.Do("exists", token)); ok && err==nil {
		account_id, _ := redis.Int(r.Do("get", token))
		user := account.Find(account_id)
		customer.Nick = user.Nick
		customer.Uid = user.Uid
		Customers[customer.Uid] = customer

		return nil
	} else {
		err = errors.New("not login")
	}
	return err
}

func (customer *Customer) Tanks() []interface{} {
	tanks := []interface{}{}
	for _, tank := range Rasps {
		tanks = append(tanks, tank)
	}
	return tanks
}

//我发送到其他
func (customer *Customer) Send(cmd string, args configs.M) error {

	fmt.Println("customer send to rasp:", cmd,args)

	if customer.Rasp == nil {
		return errors.New("not connect tank")
	}

	customer.parseCmmand(cmd,args)

	customer.Rasp.Input(cmd,args)
	return nil
}

func (customer *Customer)parseCmmand(cmd string, args configs.M){
	switch cmd{
	case "new":
		tank:=Tank{}
		tank.In1=args.Get("in1")
		tank.In2=args.Get("in2")
		tank.In3=args.Get("in3")
		tank.In4=args.Get("in4")
		tank.En1=args.Get("en1")
		tank.En2=args.Get("en2")
		tank.Save()

	case "edit":
		tank := customer.Tank
		tank.In1=args.Get("in1")
		tank.In2=args.Get("in2")
		tank.In3=args.Get("in3")
		tank.In4=args.Get("in4")
		tank.En1=args.Get("en1")
		tank.En2=args.Get("en2")
		tank.Save()
	}
}

func (customer *Customer) Exit() {
	if customer.Tank != nil {
		//customer.Tank.DelCustomer(customer)
	}
	delete(Customers, customer.Uid)
}

func (customer *Customer) ConnectRasp(rasp_name string) (err error) {

	if rasp, ok := Rasps[rasp_name]; ok {
		customer.Rasp = rasp
		//customer.Tank.To(customer, fmt.Sprint("hi ", customer.Nick))
		//customer.Tank.AddCustomer(customer)
		return nil
	} else {
		return errors.New(fmt.Sprint(rasp_name, "this tank  already left"))
	}
}

//其他发送到我
func (customer *Customer)Input(cmd string, args interface{}){
	data := configs.M{"cmmand": cmd, "args": args}
	customer.Link.WriteJSON(data)
}

//其他发送到我
func (customer *Customer)input(content string){
	customer.Input("text",configs.M{"content":content})
}