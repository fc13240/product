package shopee
import(
	"helper/dbs/mongodb"
	"product"
	"helper/configs"
	"log"
)
const(
	PlatformName="shopee"
)

type Variation struct{
	Id int `json:"variation_id,omitempty"`
	Name string `json:"name,omitempty"`
	Sku string `json:"variation_sku,omitempty"`
	Price float32 `json:"price,omitempty"`
	Stock int `json:"stock,omitempty"`
}

type Attr struct{
	Id int `json:"attributes_id"`
	Value string `json:"value"`

}

type Logistics struct{
	 Logistic_id	int `json:"logistic_id"`
 	 Enabled bool `json:"enabled"`
}

type Item struct{
	Sku string `json:"item_sku" bson:"sku"`
	ItemId int `json:"item_id" bson:"item_id" `
	Cid int `json:"category_id" bson:"category_id"`
	LocalCategoryId int `json:"local_category_id"`  //本地分类ID
	Name string `json:"name" bson:"name"`
	Cnname string `json:"cnname" bson:"cnname"`
	Stock int `json:"stock" bson:"stock"`
	Variations []Variation `json:"variations" bson:"variations"`
	Desc string `json:"description" bson:"description"`
	Weight float32 `json:"weight" bson:"weight"`
	Price float32 `json:"price" bson:"price"`
	Images []string `json:"images" bson:"images"`
	Attributes []Attr `json:"attributes" bson:"attributes"`
	Logistics []Logistics `json:"logistics"`
	AuthorId int `json:"author_id"`
	Skus []product.Sku `json:"skus"`
	Sizes []string `json:"sizes"`
	Colors []string `json:"colors"`
	Currency string `json:"currency"`
}

//获取商品属性
func (item *Item)GetAttribute()(shopee_attrs []Attr){
	shopee_attrs=[]Attr{}
	attrs,_:=product.GetCategoryAttrs(PlatformName,item.LocalCategoryId,false)
	
	for _,att:=range attrs{
		if opt:=product.GetOneAttrSelectedOption(item.Sku,att.Id);len(opt)>0{
			shopee_attrs=append(shopee_attrs,Attr{att.Config.GetInt("shopee_attribute_id"),opt[0].Value})
		}
	}
	return shopee_attrs
}

func (item *Item)SetItemId(item_id int){
	set:=configs.M{
		"$set":configs.M{"item_id":item_id},
	}
	ItemCol().Update(configs.M{"sku":item.Sku},set)
}

func (item *Item)SetVariations(variations []Variation){
	log.Println(variations,item)

	set:=configs.M{
		"$set":configs.M{"variations":variations},
	}
	ItemCol().Update(configs.M{"sku":item.Sku},set)
}
func (item *Item)GetDefaultLogistics()([]Logistics){
	return []Logistics{{28016,true},{28008,true}}
}

func Paste(sku string)(item *Item,err error){
	pro,err:=product.Get(sku)
	item=&Item{}
	if err!=nil{
		return 
	}
	item.Name=pro.EnName
	item.Cnname=pro.Name
	item.Sku=pro.Sku
	item.Weight=pro.Weight
	item.Price=pro.Price
	item.Images=product.GetImages(pro.Sku)
	return 
}

func Get(sku string)(item *Item,err error){
	if Exist(sku) == false{
		
		item,err=Paste(sku)
		if err!=nil{
			log.Println(err)
		}
		return 
	}
	item=&Item{}
	ItemCol().Find(configs.M{"sku":sku}).One(item)
	item.Currency="RM"
	return item,nil
}

func Exist(sku string)bool{
	n,_:=ItemCol().Find(configs.M{"sku":sku}).Count()
	if n>0{
		return true
	}
	return false
}

func ItemCol()*mongodb.Collection{
	return mongodb.Conn().C("shopee.item")
}

func (item *Item)Save(){
	if Exist(item.Sku) == false {
		ItemCol().Insert(item)
	}else{
		ItemCol().Update(configs.M{"sku":item.Sku},configs.M{"$set":item})
	}
}