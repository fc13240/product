package ezbuy
import(
	"fmt"
	"helper/configs"
	"helper/util"
)

func NewFission(sku string)*Item{
	old_item:=Get(sku)
	n,_:=ItemCol().Find(configs.M{"parent_sku":old_item.SKU}).Count()
	new_item:=old_item
	new_item.SKU=fmt.Sprint(old_item.SKU,"-",n+1)
	new_item.IsChild=true
	new_item.ParentSku=old_item.SKU
	new_item.CreateDate=util.Datetime()
	ItemCol().Insert(new_item)
	return new_item
}

func FissionItemByColor(sku string){
	old_item:=Get(sku)
	n,_:=ItemCol().Find(configs.M{"parent_sku":old_item.SKU}).Count()
	for _,color:=range old_item.Colors{
		new_item:=old_item.Copy()
		new_item.Colors=[]Color{color}
		n++
		new_item.SKU=fmt.Sprint(old_item.SKU,"-",n)
		new_item.IsChild=true
		new_item.ParentSku=old_item.SKU
		new_item.CreateDate=util.Datetime()
		ItemCol().Insert(new_item)
	}
}

func GetItemChilds(sku string)(items []Item){
	ItemCol().Find(configs.M{"parent_sku":sku}).Select(configs.M{"parent_sku":true,"sku":true,"cname":true}).All(&items)
	return 
}

func DelItemChild(sku string){
	ItemCol().Remove(configs.M{"sku":sku,"ischild":true})
}