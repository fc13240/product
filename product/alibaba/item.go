package alibaba

import (
	"fmt"
	"helper/account"
	"helper/configs"
	"helper/crypto"
	"strings"
	"helper/dbs/mongodb"
	"time"
	"errors"
	"product"
	"encoding/json"
)


func Col() *mongodb.Collection{
	return mongodb.Conn().C("alibaba")
}

type Item struct {
	Id 	   int64 	    `json:"id" bson:"id"`
	ItemId     int 		    	`json:"item_id" bson:"item_id"`
	Authorid   int  	    	`json:"authorid"`
	BuyPrice  float32	    	`json:"buyprice"`
	Title     string            `json:"title"`
	BuyAddr   string            `json:"buyaddr" bson:"buyaddr"`
	SaveAddr  string            `json:"saveAddr"`
	BaseAddr  string            `json:"baseAddr"`
	Addtime   time.Time         `json:"addtime"`
	Uptime    time.Time         `json:"update"`
	SukImages []SkuImage 	    `bson:"sukimages" json:"sukImages"`
	Images    []string          `json:"images"`
	Desc      string 	    	`json:"desc"`
	Descurl   string            `json:"desurl"`
	Colors    []Color	    	`bson:"colors" json:"colors"`
	Sizes     []Size	    	`bson:"sizes"  json:"sizes"`
	SellPrice float32  	    	`json:"sellPrice"`
	ShippingPrice float32       `json:"shippingPrice"`
	CompanyName string 	    	`bson:"company_name" json:"company_name"`
	CompanyUrl string 	    	`bson:"company_url"  json:"company_url"`
	Remark    string 	    	`bson:"remark"  json:"remark"`
	SellerId int 		    	`bson:"seller_id" json:"seller_id"`
	Source   string 	     	`bson:"source" josn:"source"`
	Sku    string  				`bson:"sku" json:"sku"`
	Orderby int 				`bson:"orderby" json:"orderby"`
	Meats []Meat 				`json:"meats"`
	Attrs []Attr  				`json:"attrs"`
}

type Meat struct{
	Name string `json:"name"`
	Content string `json:"content"`
}

type Attr struct{
	Name string `json:"name"`
	Content string `json:"content"`
}

//统计还有其他店有货源
func (item *Item)CountOtherSource() int {
	n,_:=Col().Find(configs.M{"item_id":item.ItemId,"authorid":item.Authorid,"id":configs.M{"$ne":item.Id}}).Count()
	return n
}

type SkuImage struct {
	Name          string `json:"name"`
	Original      string `json:"original"`
	Preview       string `json:"preview"`
	LocalOriginal string `json:"localOriginal"`
	LocalPreview  string `json:"localPreview"`
}

type Color struct {
	Name string `json:"name"`
}

type Size struct {
	Name string `json:"name"`
}

func GetByUrl(url string) (item *Item,err  error){
	Col().Find(configs.M{"md5":crypto.Md5(url)}).One(&item)
	return item, nil
}

func GetById(id int64,author *account.Account) (item *Item,err  error) {
	err=Col().Find(configs.M{"id":id,"authorid":author.Uid}).One(&item)

	if item!=nil && item.Id>0{
		return item, nil
	}else{
		return item,errors.New("没有采购信息")
	}
}

func GetBySku(sku string) (item *Item,err  error) {
	sku=strings.Trim(sku," ")
	if err:=Col().Find(configs.M{"sku":sku}).One(&item);err==nil{

		return item, nil
	}else{
		return item,errors.New("没有找到相同在宝贝")
	}
}

//通过 itemID 获取商品
func Get(sku string,author *account.Account) (item *Item,err  error){
	if sku==""{
		return item,errors.New("SKU is empty")
	}
	Col().Find(configs.M{"sku":sku,"authorid":author.Uid}).Sort("orderby").One(&item)
	if item!=nil && item.Id>0{
		return item, nil
	}else{
		return item,errors.New("没有采购信息")
	}
}

//商品属性
func GetAttrs(item_id int) (attrs []Attr){
	attrs=[]Attr{}

	dd:=struct{
		Attrs []Attr `json:"attrs"`
	}{}
	Col().Find(configs.M{"item_id":item_id}).Select(configs.M{"attrs":true}).One(&dd)
	return dd.Attrs

}

func Search(filter configs.M,offset ,limit int)(res []Item,total int){
	res=[]Item{};
	where:=configs.M{}

	if filter.Get("sku") !="" {
		sku:=mongodb.Rexp(fmt.Sprintf("^%s.*",filter.Get("sku")),"i")
		where["sku"]=sku
	}
	fmt.Println(where)
	Col().Find(where).Skip(offset).Limit(limit).Sort("orderby").All(&res)
	total,_=Col().Find(where).Count()

	return 
}

func Del(id int64)error {
	return Col().Remove(configs.M{"id":id})
}

func (item *Item)Exist() (ok bool){
	fmt.Println(item.Id)
	if n,_:=Col().Find(configs.M{"id":item.Id,"source":item.Source,"authorid":item.Authorid,"item_id":configs.M{"$gt":0} }).Count();n>0{
		return true
	}
	return false
}

func (p *Item) Save(){
	if c,_:=Col().Find(configs.M{"id":p.Id}).Count();c>0{
		p.Uptime=time.Now()
		Col().Update(configs.M{"id":p.Id},configs.M{"$set":p})
	}else{
		Col().Insert(p)
	}
}

func (p *Item) Set(v interface{}) {
	Col().Update(configs.M{"id":p.Id,"authorid":p.Authorid},configs.M{"$set":v})
}

func (item *Item)Encode()(body []byte){
	body,_=json.Marshal(item)
	return 
}

func (item *Item)Unmarshal(body []byte)error{
	return json.Unmarshal(body,item)
}

func Update(selector,date configs.M)error{
	return Col().Update(selector,configs.M{"$set":date})
}

func Join(author *account.Account, items ...*product.Item)(rows map[int]*Item){
	item_ids:=[]int{}

	for _,item:=range items{
		item_ids=append(item_ids,item.Id)
	}

	data:=[]*Item{}
	Col().Find(configs.M{"item_id":configs.M{"$in":item_ids},"author_id":author.Uid}).All(&data)
	rows=map[int]*Item{}
	for _,row:=range data{
		rows[row.ItemId]=row
	}
	return rows
}

//商品的供应商列表
func GetAliSupplier(sku string,author *account.Account)configs.M{
	if aliItem,err:=Get(sku,author);err==nil{
		headimg:="/images/timg.gif"

		if len(aliItem.Images)>0{
			headimg=aliItem.Images[0]
		}
		return configs.M{
			"company_name":aliItem.CompanyName,
			"company_url":aliItem.CompanyUrl,
			"seller_id":aliItem.SellerId,
			"buy_addr":aliItem.BuyAddr,
			"buy_price":aliItem.BuyPrice,
			"remark":aliItem.Remark,
			"other_source":aliItem.CountOtherSource(),
			"headimg":headimg,
		}
	}
	return nil
}

type ColorSize struct{
	Colors    []Color	    	`bson:"colors" json:"colors"`
	Sizes     []Size	    	`bson:"sizes"  json:"sizes"`
}

//获取颜色和尺寸
func GetColorsAndSizes(sku string)(colorsize *ColorSize){
	Col().Find(configs.M{"sku":sku}).Select(configs.M{"sizes":true,"colors":true}).One(&colorsize)
	return colorsize
}


//某个店铺的商品列表
func BySellerItemIds(seller_id int)(item_ids []int){

	ids:=[]struct{
		ItemId int `bson:"item_id"`
	}{}
	Col().Find(configs.M{"seller_id":seller_id}).Select(configs.M{"item_id":1}).All(&ids)
	item_ids=[]int{}
	for _,id:=range ids{
		item_ids=append(item_ids,id.ItemId)
	}
	if len(item_ids) == 0{
		return []int{0}
	}
	return item_ids
}

func filertName(title string) string {
	title = strings.Replace(title, "#", "", -1)
	title = strings.Replace(title, " ", "", -1)
	return title
}


