package product

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"helper/configs"
	"helper/redisCli"

	"log"
)

type Cart struct {
	BuyerId     int
	Token       string
	Totalamount float32
}

func OpenCart(token string) *Cart {
	return &Cart{Token: token}
}

type CartItem struct {
	cart     *Cart
	Sku   		string  `redis:"sku"`
	Pid      int     `redis:"product_id"`
	Name     string  `redis:"name"`
	Headimg  string  `redis:"headimg"`
	Price    float32 `redis:"price"`
	Oldprice float32 `redis:"old_price"`
	Color  string     `redis:"color"`
	Size   string     `redis:"size"`
	Quant    int     `redis:"quant"`
}

func (item *CartItem) SetQuant(quant int) {
	r := redisCli.Conn()
	item.Quant = quant
	r.Do("HSET", item.key(), "quant", quant)
}


func (item *CartItem) key() string {
	return fmt.Sprint("cart:", item.cart.Token, ":", item.Sku)
}

func (item *CartItem) Incr(quant int) {

	if item.Quant+quant < 6 {
		r := redisCli.Conn()
		item.Quant += quant
		r.Do("HINCRBY", item.key(), "quant", quant)
	}
}

func (item *CartItem) Minus(quant int) {
	if item.Quant-quant > 0 {
		r := redisCli.Conn()
		item.Quant -= quant
		r.Do("HINCRBY", item.key(), "quant", -quant)
	}
}

func (item *CartItem) Remove() {
	r := redisCli.Conn()
	r.Do("SREM", fmt.Sprint("cart:", item.cart.Token), item.Sku)
	r.Do("DEL", item.key())
}

func (item *CartItem) TotalPrice() float32 {
	return float32(item.Quant) * item.Price
}

func (cart *Cart) Item(sku string)(item *CartItem){

	item=&CartItem{}
	r := redisCli.Conn()

	key:=fmt.Sprint("cart:", cart.Token, ":", sku)

	if vv,err:=redis.Values(r.Do("HGETALL",key));err==nil{
		redis.ScanStruct(vv,item)
		item.cart=cart
	}else{
		log.Println("cart error",err.Error())
	}
	return item
}

func (cart *Cart) Add(sku string) (*CartItem, error) {
	r := redisCli.Conn()

	key := fmt.Sprint("cart:", cart.Token, ":", sku)

	if count, _ := redis.Int(r.Do("EXISTS", key)); count == 1 {
		if vv, err := redis.Values(r.Do("HGETALL", key)); err == nil {
			var item CartItem
			if err := redis.ScanStruct(vv, &item); err == nil {
				item.cart = cart
				item.Sku = sku
				return &item, nil
			}
		} else {
			fmt.Println("read case error ", err)
		}
	}else if item, err := Get(sku); err == nil {

		caritem := &CartItem{
			Sku:      sku,
			Pid:      item.Id,
			Name:     item.Name,
			Price:    item.Price,
			Oldprice: item.OldPrice,
			Headimg:  item.Headimg,
			Quant:    1,
		}

		_, err := r.Do("HMSET",key,
			"product_id", caritem.Pid,
			"name", caritem.Name,
			"price", caritem.Price,
			"old_price", caritem.Oldprice,
			"quant", caritem.Quant,
			"color", caritem.Color,
			"size", caritem.Size,
			"headimg", caritem.Headimg,
			"sku",caritem.Sku,
		)

		if err != nil {
			return nil, err
		}

		if _, err := r.Do("SADD", fmt.Sprint("cart:", cart.Token), sku); err != nil {
			return nil, err
		}

		caritem.cart = cart
		return caritem, nil
	}
	return nil, errors.New("not exist on product")
}

func (cart *Cart) TotalAmount() float32 {
	var total float32
	for _, item := range cart.Items() {
		total += item.TotalPrice()
	}
	return total
}

func (cart *Cart) Items() []*CartItem {
	r := redisCli.Conn()
	var items []*CartItem
	values, _ := redis.Strings(r.Do("SDIFF", fmt.Sprint("cart:", cart.Token)))

	cart.Totalamount = 0

	if len(values) == 0 {
		return items
	}

	for _, sku := range values {
		key := fmt.Sprint("cart:", cart.Token, ":", sku)

		if count, _ := redis.Int(r.Do("EXISTS", key)); count < 1 {
			r.Do("SREM", fmt.Sprint("cart:", cart.Token), sku)
			continue
		}

		if vv, err := redis.Values(r.Do("HGETALL", key)); err == nil {
			var item CartItem

			if err := redis.ScanStruct(vv, &item); err == nil {
				item.Sku = sku
				items = append(items, &item)
			}
		}
	}
	return items
}

func (cart *Cart) Listing() []interface{} {
	items := []interface{}{}
	for _, item := range cart.Items() {
		v := map[string]interface{}{
			"id":          item.Pid,
			"sku":         item.Sku,
			"quant":       item.Quant,
			"price":       item.Price,
			"total_price": configs.Price(item.TotalPrice()),
			"title":       item.Name,
			"headimg":     item.Headimg,
			"color":       item.Color,
			"size":        item.Size,
		}
		items = append(items, v)
	}
	return items
}

func (cart *Cart) Del(sku string){
	r := redisCli.Conn()
	key:=fmt.Sprint("cart:", cart.Token, ":", sku)
	r.Do("DEL", key)
}

func (cart *Cart) Clean() {
	r := redisCli.Conn()
	var indexs []interface{}

	indexs = append(indexs, fmt.Sprint("cart:", cart.Token))
	for _, item := range cart.Items() {
		indexs = append(indexs, fmt.Sprint("cart:", cart.Token, ":", item.Sku))
	}
	r.Do("DEL", indexs...)
}
