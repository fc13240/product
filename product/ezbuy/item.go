package ezbuy

import (
	"helper/util"

	"product/alibaba"
	"helper/account"
	"helper/dbs/mongodb"
	"helper/configs"
	"product"
	"errors"
	"encoding/json"
	"reflect"

)

var itemCol *mongodb.Collection

func ItemCol() *mongodb.Collection{
	if itemCol==nil{
		itemCol=mongodb.Conn().C("ezbuy.items")
	}
	return itemCol
}

type Item struct{
	Id     int `json:"pid" bson:"id"`
	ItemId  int `json:"itemid"`
	AuthorId int `json:"authorid"`
	StoreID  int  `xlsx:"0" json:"storeid"`
	StoreName string  `xlsx:"1" json:"storename"`
	CID int `json:"cid" xlsx:"2"`
	CName string `json:"cname" xlsx:"3"`
	Name  string `json:"name" xlsx:"4"`
	CNName  string `json:"cnname" xlsx:"5"`
	Desc string `json:"desc" xlsx:"6"`
	Images []string `json:"images" xlsx:"7"`
	N string `xlsx:"8" json:"n"`
	ISONE string `xlsx:"9" json:"isone"`
	SKU  string `json:"sku"  bson:"sku" xlsx:"10"`
	BuyingPrice float32 `json:"buying_price"`
	OldPrice float32 `json:"oldprice" xlsx:"11"`
	Price float32 ` json:"price" xlsx:"12"`
	Quant int `json:"quant" xlsx:"13"`
	Weight float32 `json:"weight" xlsx:"14"`
	Length float32 `json:"length" xlsx:"15"`
	Width float32 `json:"width" xlsx:"16"`
	Height float32 `json:"height" xlsx:"17"`
	SkuImages []string `json:"skuimages" xlsx:"18"`
	Update string `json:"update"`
	AddDate string `json:"adddate"`
	Colors []Color `json:"colors"`
	Sizes []Size `json:"sizes"`
	SoldCount int `json:"soldcount"`
	Materials []int `json:"materials"`
	Styles Style `json:"styles"`
	CreateDate string `json:"create_date" bson:"create_date"`
	ParentSku string `json:"parent_sku" bson:"parent_sku"` //parent sku
	IsChild bool `json:"ischild" bson:"ischild"`
	Pid int `json:"pid" bson:"pid"`  //对应ez平台的多个产品，重复上传的可能
}

//颜色
type Color struct{
	ID int  `json:"id"`
	Name string `json:"name"`
	Image string `json:"image"`
}

//尺寸
type Size struct{
	ID int  `json:"id"`
	Name string  `json:"name"`
}

//材质
type Material struct{
	ID int  `json:"id"`
	Name string  `json:"name"`
}

//获取一个ez商品
func Get(sku string )(item *Item){
	ItemCol().Find(configs.M{"sku":sku}).One(&item)
	return item
}

//从商品和阿里巴巴中，组成一个商品
func Paste(sku string,author *account.Account)error{
	pro,err:=product.Get(sku)
	if err!=nil{
		return err
	}
	if pro.Sku==""{
		return errors.New("sku is empty")
	}
	if info, err := alibaba.Get(sku,author); err == nil {
		detal := &Item{
			AuthorId:author.Uid,
			CNName:pro.Name,
			ItemId:pro.Id,
			Length:pro.Length,
			Width:pro.Width,
			Weight:pro.Weight,
			Height:pro.Height,
			Price:pro.Price,
			OldPrice:pro.OldPrice,
			SKU:pro.Sku,
			Desc: pro.Desc,
			Images:info.Images,
			Quant:99,
			BuyingPrice:pro.BuyingPrice, //采购价格
		}

		aliItem,_:=alibaba.Get(sku,author)
		if len(aliItem.SukImages)>0{
			for id, alicolor := range aliItem.SukImages {
				color:=Color{id,alicolor.Name,alicolor.Original}
				detal.Colors = append(detal.Colors,color)
			}
		}else{
			for _, color := range info.Colors {
				detal.Colors = append(detal.Colors, GetColor(color.Name))
			}
		}
		for _, size := range info.Sizes {
			detal.Sizes = append(detal.Sizes, GetSize(size.Name))
		}
		if SkuExist(sku){
			detal.Update=util.Datetime()
			ItemCol().Update(configs.M{"sku":sku},configs.M{"$set":detal})
		}else{
			detal.AddDate=util.Datetime()
			ItemCol().Insert(detal)
		}
		return nil
	}else{
		return err
	}
}

//是否存在
func Exist(item_id int) bool{
	count,_:=ItemCol().Find(configs.M{"itemid":item_id}).Count()
	if count>0{
		return true
	}
	return false
}

//SKU是否存在
func SkuExist(sku string) bool{
	count,_:=ItemCol().Find(configs.M{"sku":sku}).Count()
	if count>0{
		return true
	}
	return false
}

//更新商品信息
func (item *Item) Save() error{
	item.Update=util.Datetime()
	return ItemCol().Update(configs.M{"sku":item.SKU},configs.M{"$set":item})
}

func (item *Item)Set(data configs.M)error {
	return ItemCol().Update(configs.M{"sku":item.SKU},configs.M{"$set":data})
}

//添加
func Add(item *Item)error{
	return ItemCol().Insert(item)
}

func (item *Item)Encode()(body []byte){
	body,_=json.Marshal(item)
	return 
}

func (item *Item)Unmarshal(body []byte)error{
	return json.Unmarshal(body,item)
}

func (item *Item)GetColorImage(color_name string)string{
	for _,color:=range item.Colors{
		if color.Name == color_name{
			return color.Image
		}
	}
	return  ""
}

//通过ID，改变修改一个商品属性
func SetItemField(author *account.Account,id int,field,value string ){
	ItemCol().Update(configs.M{"authorid":author.Uid,"id":id},configs.M{"$set":configs.M{field:value}})
}

//清空所有商品
func CleanItems(author *account.Account) error{
	return errors.New("这是一个很危险的动作，不能执行")
	_,err:= ItemCol().RemoveAll(configs.M{"authorid":author.Uid})
	return err
}

func (item *Item)Copy() *Item{
	new_item:=&Item{}
    sval:= reflect.ValueOf(item).Elem()
    dval:= reflect.ValueOf(new_item).Elem()

	
    for i := 0; i < sval.NumField(); i++ {
        value := sval.Field(i)
        name := sval.Type().Field(i).Name

        dvalue := dval.FieldByName(name)
        if dvalue.IsValid() == false {
            continue
	}
        dvalue.Set(value) //这里默认共同成员的类型一样，否则这个地方可能导致 panic，需要简单修改一下。
	}
	return new_item
}