package rpi

var(
	Rasps     = RaspStorage{}
	Customers = CustomerStorage{}
)

type RaspStorage map[string]*Rasp


func (storage RaspStorage)Add(rasp *Rasp){
	storage[rasp.Name]=rasp
}

func (storage RaspStorage)Del(rasp *Rasp){
	delete(storage,rasp.Name)
}

type CustomerStorage map[int]*Customer

func (storage CustomerStorage)Add(c *Customer){
	storage[c.Uid]=c
}

func (storage CustomerStorage)Del(c *Customer){
	delete(storage,c.Uid)
}