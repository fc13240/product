package api
import(
	
	"fmt"
	"helper/configs"
	"product/shopee"
	"errors"
	"helper/webtest"
	"product"
	
)

type Category struct{
	ParentId int `json:"parent_id"`
	HasChildren bool `json:"has_children"`
	Name string `json:"category_name"`
	Id int `json:"category_id"`
}
func Listing(){
	parma:=configs.M{
		"pagination_offset":1,
		"pagination_entries_per_page":20,
		"more":true,
	}
	res,_:=POST("/items/get",parma)
	fmt.Println(res.String())
}

func Add(){
	parma:=configs.M{
		"pagination_offset":1,
		"pagination_entries_per_page":20,
		"more":true,
	}
	res,_:=POST("/item/add",parma)
fmt.Println(res.String())

}

func GetCategorys()[]Category{
	param:=configs.M{

	}
	res,_:=POST("/item/categories/get",param)
	data:=struct{
		Categories []Category `json:"categories"`
	}{}
	if err:=res.BindJSON(&data);err!=nil{
		fmt.Println(err)
	}
	return data.Categories
	
}

type Attribute struct{
	Name string `json:"attribute_name"`
	InputType string `json:"input_type"`
	Id int `json:"attribute_id"`
	Type string `json:"attribute_type"`
	IsMandatory bool `json:"is_mandatory"`
	Options []string `json:"options"`
}

func GetAttributes(category_id int){
	param:=configs.M{
		"category_id":category_id,
	}
	res,_:=POST("/item/attributes/get",param)
	data:=struct{
		Attributes []Attribute `attributes`
	}{}
	res.BindJSON(&data)
	for _,attr:=range data.Attributes{
		fmt.Println(attr.Name,attr.IsMandatory,attr.Id)
	}
}

type Logistic struct {
	Id int `json:"logistic_id"`
	Name string `json:"logistic_name"`
	Enabled bool `json:"enabled"`
}

func GetLogistics(){
	param:=configs.M{
		
	}
	res,_:=POST("/logistics/channel/get",param)

	data:=struct{
		Logistics []Logistic `logistics`
	}{}
	res.BindJSON(&data)
	for _,attr:=range data.Logistics{
		fmt.Println(attr.Id,attr.Name,attr.Enabled)
	}
}
func GetItem(item_id int)*shopee.Item{
	
	res,_:=POST("/item/get",configs.M{"item_id":item_id})

	data:=struct{
		Item shopee.Item `json:"item"`
	}{}

	res.BindJSON(&data)
	return &data.Item
}

func DelItem(sku string)configs.M{
	item,_:=shopee.Get(sku)
	res,_:=POST("/item/delete",configs.M{"item_id":item.ItemId})
	item.SetItemId(0)
	body:=configs.M{}
	res.BindJSON(&body)
	return body
}

func AddItem(sku string)error{
	item,_:=shopee.Get(sku)
	
	varations:=[]shopee.Variation{}
	
	if item.ItemId == 0{
		skuId:=1
		for _,color:=range item.Colors{
		for _,size:=range item.Sizes{
			v:=shopee.Variation{}
			v.Name=fmt.Sprint(color," ",size)
			v.Price=item.Price
			v.Stock=99
			v.Sku=fmt.Sprint(item.Sku,"-",skuId)
			skuId++
			varations=append(varations,v)
			item.Stock+=v.Stock
		}
	}
	}else{
		for _,old:=range item.Variations{
			v:=shopee.Variation{}
			v.Id=old.Id
			v.Name=old.Name
			v.Sku=old.Sku
			varations=append(varations,v)
		}
	}
	
	attributes:=item.GetAttribute()
	param:=configs.M{
		"category_id":item.Cid,
		"name":item.Name,
		"description":item.Desc,
		"price":item.Price,
		"stock":item.Stock,
		"item_sku":item.Sku,
		"variations":varations,
		"logistics":item.GetDefaultLogistics(),
		"weight":item.Weight,
		"attributes":attributes,
	}
	var res *webtest.Result

	if item.ItemId == 0{ //新增
		images:=[]configs.M{}
		for _,src:=range product.GetImages(item.Sku){
			images=append(images,configs.M{"url":src})
		}
		param.Add("images",images)
		res,_=POST("/item/add",param)
	}else{//更新
		param.Add("item_id",item.ItemId)
		res,_=POST("/item/update",param)
	}

	re:=struct{
		Msg string `json:"msg"`
		Error string `json:"error"`
		ItemId int `json:"item_id"`
	}{}

	if err:=res.BindJSON(&re);err==nil{
		if re.ItemId>0 { 
			if item.ItemId==0{//新增的
				item.SetItemId(re.ItemId)
				new_item:=GetItem(re.ItemId)
				new_item.SetVariations(new_item.Variations)
			}
			return nil
		}else{
			return errors.New(re.Msg)
		}
	}else{
		return err
	}
	
}
