package ezbuy

import (
	"fmt"
	"helper/dbs/mongodb"
	"helper/configs"

	"helper/account"
	"log"
	"order"
	"helper/mail"
	"time"

)

var orderCol *mongodb.Collection

func OrderCol() *mongodb.Collection{
	if orderCol==nil{
		orderCol=mongodb.Conn().C("ezbuy.orders")
	}
	return orderCol
}

type OrderItem struct {
	OrderItemId int 	`json:"orderItemId"`
	Quantity    int		`json:"quantity"`
	SellerSkuId string	`json:"sellerSkuId"`
	SkuName     string	`json:"skuName"`
	Status      int		`json:"status"`
	TotalAmount float32	`json:"totalAmount"`
	Track OrderTrack	`json:"track"`
	Type int 		`json:"type"`
	UnitCashOff float32	`json:"unitCashOff"`
	UnitDiscount float32	 `json:"unitDiscount"`
	unitOriginalPrice float32 `json:"unitOriginalPrice"`
	UnitPrice float32 	`json:"unitPrice"`
}

type OrderTrack struct {
	Provider	string
	TrackingNum	string
}

type Order struct {
	CreatedAt          int64   `json:"createdAt"`
	DispatchDelayed    bool `json:"dispatchDelayed"`
	ExpectedDispatchAt int `json:"expectedDispatchAt"`
	FinishedAt         int `json:"finishedAt"`
	Items              []OrderItem `json:"Items"`
	OrderNum           string `json:"orderNum"`
	PrimaryImage       string `json:"primaryImage"`
	ProductName        string `json:"productName"`
	ProductUrl         string `json:"productUrl"`
	Remarks            []configs.M `json:"remarks"`
	SellType           int `json:"sellType"`
	ShippingFee        float32 `json:"shippingFee"`
	Status             int `json:"status"`
	TotalAmount        float32 `json:"totalAmount"`
	Track              OrderTrack `json:"track"`
	Warehouse          int `json:"warehouse"`
	AliItem      	   configs.M `json:"aliinfo"`
	CreateTime         string `json:"createtime"`
	SellerId 	     int `json:"sellerid"`
	AuthorId            int `json:"authorid"`
	Itemid int 	   `json:"item_id" bson:"item_id"`
	LabelLog *order.LabelLog `json:"labellog"`
}

func (ord *Order)SetItemId(item_id int)error{
	ord.Itemid=item_id
	return OrderCol().Update(configs.M{"ordernum":ord.OrderNum},configs.M{"$set":configs.M{"item_id":item_id}})
}

func GetOrder(order_num string,author *account.Account)(item *Order,err error ){
	item =&Order{}
	err=OrderCol().Find(configs.M{"ordernum":order_num,"authorid":author.Uid}).One(item)
	if err!=nil{
		log.Println(err)
	}
	return item,err
}

func OrderListing(author *account.Account,offset,limit int,orderBy string)(items []*Order,total int ){
	 items =[]*Order{}
	 err:=OrderCol().Find(configs.M{"authorid":author.Uid}).Sort(orderBy).Skip(offset).Limit(limit).All(&items)
	 if err!=nil{
		log.Println(err)
	 }
	 total,_=OrderCol().Find(configs.M{"authorid":author.Uid}).Count()
	 return items,total
}

func Listing(filter configs.M,offset,limit int,orderBy string)(items []*Item,total int){
	items =[]*Item{}

	ItemCol().Find(filter).Sort(orderBy).Skip(offset).Limit(limit).All(&items)
	total,_=ItemCol().Find(filter).Count()
	return items,total
}

func SaveItems(author *account.Account,items ...Item){
	for _,item:=range items{
		item.AuthorId=author.Uid
		n,_:=ItemCol().Find(configs.M{"id":item.Id,"authorid":author.Uid}).Count()
		if n>0{
			ItemCol().Update(configs.M{"id":item.Id,"authorid":author.Uid},configs.M{"$set":configs.M{
				"name":item.Name,
				"soldCount":item.SoldCount,
			}})
		}else{
			ItemCol().Insert(item)
		}
	}
}

func SaveOrders(orders []Order,author *account.Account){
	col:=OrderCol()
	for _,order:=range  orders{
		n,_:=col.Find(configs.M{"ordernum":order.OrderNum}).Count()
		order.AuthorId=author.Uid
		order.CreateTime=time.Unix(order.CreatedAt,0).Format("01/02 15:04")
		if n>0{
			col.Update(configs.M{"ordernum":order.OrderNum},configs.M{"$set":order})
		}else{
			col.Insert(order)
		}
	}
}

func CheckNewOrders(orders []Order,author *account.Account)[]Order{
	new_orders:=[]Order{}
	col:=OrderCol()

	for _,order:=range  orders{
		n,_:=col.Find(configs.M{"ordernum":order.OrderNum}).Count()
		order.AuthorId=author.Uid
		order.CreateTime=time.Unix(order.CreatedAt,0).Format("01/02 15:04")
		if n>0{
			col.Update(configs.M{"ordernum":order.OrderNum},configs.M{"$set":order})
		}else{
			new_orders=append(new_orders,order)
			col.Insert(order)
		}
	}

	if len(new_orders)>0{
		body := `<html>
			<body>
			   <h3>您有（%d）个新订单 <a href="http://main.51helper.com/ezbuy/orders">点击查看</a></h3>
			 %s
			</body>
			</html>
		`
		tab:="<ul>"
		for _,o:=range new_orders{
			tab+=fmt.Sprintf(`<li>%s</li>`,o.ProductName)
		}
		body=fmt.Sprintf(body,len(new_orders),tab+"</ul>")
		mail.Send(author.User,"您ezbuy.com上有新的订单", body)
	}
	return new_orders
}

