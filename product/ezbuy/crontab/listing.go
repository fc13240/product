package crontab
import(
	"product/ezbuy/api"
	"product/ezbuy"
	"helper/account"
	"helper/dbs/mongodb"
	"helper/configs"
	"helper/util"
	"fmt"
	"errors"
)
type ProductBase struct {
	IsOnSale bool `json:"isOnSale" bson:"isonsale"`
	IsStocked bool
	CategoryId int `json:"categoryId" bson:"categoryid"`
	Name	string  `json:"name"`
	EnName string `json:"enname"`
	OriginCode string 
	Pid int   `json:"pid"`
	Price float32  `json:"price"`
	PrimaryImage string  `bson:"primaryimage" json:"primaryImage"`
	SellType int `json:"sellType"`
	Shipment int 
	SoldCount int `json:"soldCount" bson:"soldcount"`
	IsDownDetail bool `bson:"isdowndetail"`
	Base api.ProductDetail `json:"base"`
	Skus []api.Sku `json:"skus"`
	Sku string `json:"sku,omitempty"`
	AddDate string `json:"add_date,omitempty"`
	IsCopy bool `json:"isCopy" bson:"iscopy"`
	InFor bool `json:"infor" bson:"infor"`
}

func Listing(filter configs.M ,offset,limit int)(rows []ProductBase,total int){
	mdb:=mongodb.Conn()
	total,_=mdb.C("ezbuy.listing").Find(filter).Count()
	rows=[]ProductBase{}
	sekecd:=configs.M{"pid":true,"name":true,"enname":true,"sku":true,"adddate":true,"categoryid":true,"isdowndetail":true,"infor":true,"primaryimage":true,"isonsale":true,"soldcount":true}
	mdb.C("ezbuy.listing").Find(filter).Limit(limit).Select(sekecd).Skip(offset).Sort("pid").All(&rows)
	return 
}

func (p *ProductBase)Exist()bool{
	mdb:=mongodb.Conn()
	count,_:=mdb.C("ezbuy.listing").Find(configs.M{"pid":p.Pid}).Count()
	if count>0{
		return true
	}
	return false 
}

func (p *ProductBase)Save()error {
	col:=mongodb.Conn().C("ezbuy.listing")
	if p.Exist(){
		return col.Update(configs.M{"pid":p.Pid},configs.M{"$set":p})
	}
	return col.Insert(p)
}

func (p *ProductBase)Set(data configs.M) error {
	mdb:=mongodb.Conn()
	return mdb.C("ezbuy.listing").Update(configs.M{"pid":p.Pid},configs.M{"$set":data})
	
}
var store_id=1684
func SynEzbuyListing(page int )([]ProductBase,error){
	ez:=NewEzApi()
	data:=struct{
		Products []ProductBase `products`
	}{}
	rowCount:=40
	res,err:=ez.UserProductList((page-1)*rowCount,rowCount)
	if err==nil{
		if err:=res.BindJSON(&data);err!=nil{
			return nil,err
		}
	}else{
		fmt.Println(err)
	}
	return data.Products,nil
}

func NewEzApi() *api.EzbuyeApi{
	author:=account.Find(21)
	setting:=ezbuy.GetSetting(author,store_id)
	return &api.EzbuyeApi{
		Cookie:setting.Cookie,
		ShopName:setting.StoreName,
		SkuFirst:setting.SkuFirst,
	}
}

func GetDetail(productid int){
	ezitem,err:=NewEzApi().UserProductDetail(productid)
	item:=ezbuy.Get(ezitem.Skus[0].Name)
	if item.SKU == ezitem.Skus[0].Name{
		item.Set(configs.M{"pid":ezitem.Base.Pid,"bind_date":util.Datetime() })
	}
	if err==nil{
		fmt.Println(ezitem)
	}else{
		fmt.Println(err)
	}
}

func Down(pid int)(*api.EzItem,error){
	ez:=NewEzApi()
	item,err:=ez.UserProductDetail(pid)
	if err!=nil{
		return nil,err 
	}
	if len(item.Skus)== 0{
		return nil,errors.New("skus is empty") 
	}
	sku:=item.Skus[0].SellerSkuId
	pro:=ProductBase{Pid:pid}
	pro.Set(configs.M{"sku":sku,"update_date":util.Datetime(),"skus":item.Skus,"base":item.Base,"isdowndetail":true,"enname":item.Base.EnName})
	return item,nil
	 
}

func SynToRemote(pid int)error{
	
	itemDetail:=ProductBase{}
	mdb:=mongodb.Conn()

	err:=mdb.C("ezbuy.listing").Find(configs.M{"pid":pid}).One(&itemDetail)
	if err!=nil{
		return err
	}
	if itemDetail.IsDownDetail == false{
		return errors.New("not down detail")
	}
	item:=&api.EzItem{Base:itemDetail.Base,Skus:itemDetail.Skus}
	item.Base.Source=1
	item.Base.ForceOffSale=true
	item.Base.Pid=0
	item.Base.Name=itemDetail.Name

	for i:=range item.Skus {  //set skuid  0
		item.Skus[i].SkuId=0
	}

	ez:=NewEzApi()
	new_item,err:=ez.UserProductUpdate(item)
	if err!=nil{
		return err
	}
	if len(new_item.Skus)>0{
		new_item.Sku=item.Skus[0].SellerSkuId
		base:=ProductBase{
			Sku:new_item.Sku,
			Name:new_item.Base.Name,
			CategoryId:new_item.Base.CategoryId,
			Pid:new_item.Base.Pid,
			PrimaryImage:new_item.Base.PrimaryImage,
			AddDate:util.Datetime(),
			IsDownDetail:true,
			IsCopy:true,
		}
		base.Save()
	}
	return nil
}

func SetInfor(pid int,infor bool)error{
	pro:=&ProductBase{Pid:pid}
	return pro.Set(configs.M{"infor":infor})

}