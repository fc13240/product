package product

import (
	"time"
	"helper/dbs"
	"helper/configs"
	"fmt"
	"log"
	"helper/util"
)

//关注
type Wish struct{
	ItemId int
	AddTime time.Time
	Sku string
	BuyerId int
}

//添加关注
func AddToWish(spu string,buyerid int)(wish *Wish,err error){

	item,err:=Get(spu)
	if err!= nil{
		return 
	}
	wish=&Wish{Sku:item.Sku,ItemId:item.Id,BuyerId:buyerid}

	if IsAddWish(item.Sku,buyerid){
		return wish,nil
	}
	db:=dbs.Def()
	_,err=db.Insert("product_wishlist",configs.M{"sku":wish.Sku,"item_id":wish.ItemId,"buyer_id":wish.BuyerId,"addtime":util.Datetime()})
	if err!=nil{
		log.Println("增加关注失败",err.Error())
	}
	return wish,err
}

func IsAddWish(sku string,buyerid int)bool{
	var totl int
	db:=dbs.Def()
	db.One("SELECT COUNT(*) FROM product_wishlist WHERE sku=? AND buyer_id=?",sku,buyerid).Scan(&totl)
	if totl>0{
		return true
	}else{
		return false
	}
}

//关注列表
func Wishlish(buyerid int)(items []Item){
	sql:="SELECT pro.id,pro.name,pro.price,pro.oldprice,pro.quant,pro.headimg FROM  product_wishlist AS wish LEFT JOIN product pro ON(pro.id=wish.itemid) WHERE wish.buyerid=%d"
	db:=dbs.Def()
	rows:=db.Rows(fmt.Sprintf(sql,buyerid))
	items =[]Item{}
	for rows.Next(){
		var item Item
		rows.Scan(&item.Id,&item.Name,&item.Price,&item.OldPrice,&item.Quant,&item.Headimg)
		items=append(items,item)
	}
	return items
}