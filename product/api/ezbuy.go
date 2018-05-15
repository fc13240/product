package api

import (
	"product/ezbuy"
	"helper/net/igin"
	"helper/configs"
	"order"
	"product"
	"product/alibaba"
	"fmt"
	"github.com/gin-gonic/gin"
	"helper/dbs/mongodb"
	"log"
	"helper/account"
)

type EzBuy int

func (ez *EzBuy)Orders(c *gin.Context){
	gh:=igin.H(c)

	param:=struct {
		Offset int `json:"offset"`
		Limit int `json:"limit"`
		OrderBy string `json:"orderBy"`
	}{}

	if err:=c.BindJSON(&param);err!=nil{
		gh.Fail(err.Error())
		return
	}

	author:=gh.Account()
	items,total:=ezbuy.OrderListing(gh.Account(),param.Offset,param.Limit,param.OrderBy)
	for _,ord:=range items{
		if ord.Itemid<1{
			ez.bindAliSupplier(ord,author)
			if label,ok:=order.FirstLabelLog(ord.OrderNum);ok{
				ord.LabelLog=label
			}
		}
	}

	gh.Succ(gin.H{"items":items,"total":total})
}

func (*EzBuy)bindAliSupplier(ord *ezbuy.Order,author *account.Account){
	if pItem,err:=product.GetBySku(ord.Items[0].SellerSkuId,author);err==nil && pItem.Id>0{
		if aliItem,err:=alibaba.Get(pItem.Sku,author);err==nil{
			ord.AliItem=configs.M{
				"company_name":aliItem.CompanyName,
				"buyAddr":aliItem.BuyAddr,
				"other_source":aliItem.CountOtherSource(),
				"itemid":pItem.Id,
			}
		}
	}
}

func (ez *EzBuy)RefreshOrder(c *gin.Context){
	gh:=igin.H(c)
	author:=gh.Account()
	ord,err:=ezbuy.GetOrder(c.Param("ordernum"),author)
	if err!=nil{
		gh.Fail(err.Error())
		return
	}

	if ord.Itemid<1{
		ez.bindAliSupplier(ord,author)
		if label,ok:=order.FirstLabelLog(ord.OrderNum);ok{
			ord.LabelLog=label
		}
	}
	gh.Succ(gin.H{"item":ord})
	return
}

func (*EzBuy)Listing(c *gin.Context){
	gh:=igin.H(c)

	param:=struct {
		*igin.ParamPage

		OrderBy string `json:"orderBy"`
		IsHot bool `json:"ishot"`
		IsSale bool `json:"issale"`
	}{}

	if err:=c.BindJSON(&param);err!=nil{
		gh.Fail(err.Error())
		return
	}

	author:=gh.Account()

	filter:=configs.M{
		"authorid":author.Uid,
	}

	if param.IsHot{
		filter["ishot"]=true
	}

	if param.IsSale{
		filter["id"]=configs.M{"$gt":0}
	}

	items,total:=ezbuy.Listing(filter,param.Offset,param.Limit,param.OrderBy)
	gh.Succ(gin.H{"items":items,"total":total})
}

func (*EzBuy)SaveItems(c *gin.Context){
	gh:=igin.H(c)
	param:=struct {
		Total int `json:"total"`
		Items []ezbuy.Item `json:"products"`
	}{}

	if err:=c.BindJSON(&param);err!=nil{
		gh.Fail(err.Error())
		return
	}

	if len(param.Items) == 0{
		gh.Fail("商品数为0")
		return
	}

	ezbuy.SaveItems(gh.Account(),param.Items...)
	gh.Succ(nil)
}

func (*EzBuy)SetItemField(c *gin.Context){

	gh:=igin.H(c)

	param:=struct {
		Id int `json:"id"`
		Field string `json:"field"`
		Value string `json:"value"`
	}{}

	if err:=c.BindJSON(&param);err!=nil{
		gh.Fail(err.Error())
		return
	}

	ezbuy.SetItemField(gh.Account(),param.Id,param.Field,param.Value)
	gh.Succ(nil)
}

func(*EzBuy) SaveSetting(c *gin.Context){
	setting:=ezbuy.Setting{}
	gh:=igin.H(c)
	if err:=c.Bind(&setting);err!=nil{
		gh.Fail(err.Error())
		return
	}

	ezbuy.SaveSetting(setting,gh.Account())
	gh.Succ(nil)
}

func (*EzBuy)GetSetting(c *gin.Context){
	gh:=igin.H(c)
	store:=ezbuy.GetSetting(gh.Account())
	store.SecrectKey=gh.User().Token.Value
	gh.Succ(gin.H{"item":store})
}

func (*EzBuy)OrderDetal(c *gin.Context){

}

func (*EzBuy)OrderLabels(c *gin.Context){
	gh:=igin.H(c)
	gh.Succ(gin.H{"items":order.Labels()})
}

func (*EzBuy)OrderLabelLogs(c *gin.Context){
	gh:=igin.H(c)
	gh.Succ(gin.H{"items":order.OrderLabelLogs(c.Param("ordernum"))})
}

func (*EzBuy)AddLabelLog(c *gin.Context){
	gh:=igin.H(c)

	param:=struct {
		LabelId int `json:"labelid"`
		OrderNum []string `json:"ordernum"`
		Remarks string `json:"remarks"`
	}{}

	if err:=c.BindJSON(&param);err!=nil{
		gh.Fail(err.Error())
		return
	}

	if label,err:=order.GetLabel(param.LabelId);err==nil{
		label.AddLog(param.Remarks,param.OrderNum...)
		gh.Succ(nil)
	}else{
		gh.Fail(err.Error())
		return
	}
}

func (*EzBuy)SaveOrders(c *gin.Context){
	gh:=igin.H(c)
	data:=struct {
		Total int `json:"total"`
		Orders []ezbuy.Order `json:"data"`
	}{}

	if err:=c.BindJSON(&data);err!=nil{
		gh.Fail(err.Error())
		return
	}

	ezbuy.SaveOrders(data.Orders,gh.Account())
	gh.Succ(nil)
}

func (*EzBuy)CheckNewOrders(c *gin.Context){
	gh:=igin.H(c)
	data:=struct {
		Total int `json:"total"`
		Orders []ezbuy.Order `json:"data"`
	}{}

	if err:=c.BindJSON(&data);err!=nil{
		gh.Fail(err.Error())
		return
	}
	new_ordes:=ezbuy.CheckNewOrders(data.Orders,gh.Account())
	gh.Succ(gin.H{"newOrderNum":len(new_ordes)})
}

func (*EzBuy)CleanItems(c *gin.Context){
	gh:=igin.H(c)
	if err:=ezbuy.CleanItems(gh.Account());err==nil{
		gh.Succ(nil)
	}else{
		gh.Fail(err.Error())
	}
}

func (*EzBuy)Get(c *gin.Context){
	gh:=igin.H(c)
	sku:=c.Param("sku")
	//是否检查是否已经存在，如果不存在，从主商品中拷贝
	if isCheckPaste,ok:=c.GetQuery("check_paste");ok && isCheckPaste == "1"{
		if ezbuy.SkuExist(sku) == false{
			if err:=ezbuy.Paste(sku,gh.Account());err!=nil{
				gh.Fail(err.Error())
				return
			}
		}
	}
	if item:=ezbuy.Get(sku);item==nil{
		gh.Fail("ezbuy中不存在")
	}else{
		gh.Succ(gin.H{"item":item})
	}
}

func (*EzBuy)Paste(c *gin.Context){
	gh:=igin.H(c)
	sku:=c.Param("sku")
	if err:=ezbuy.Paste(sku,gh.Account());err==nil{
		gh.Succ(nil)
	}else{
		gh.Fail(err.Error())
	}
}

func (*EzBuy)Export(c *gin.Context){
	gh:=igin.H(c)

	defer func(){
		if r:=recover();r!=nil{
			gh.Fail(fmt.Sprint(r))
			c.Abort()
		}
	}()
	data:=struct{
		Skus []string `json:"skus"`
	}{}

	if err:=c.BindJSON(&data);err!=nil{
		gh.Fail(err.Error())
		return
	}

	if len(data.Skus) == 0{
		gh.Fail("请至少选择一个商品")
		return
	}

	items:=[]*ezbuy.Item{}
	for _,sku:=range data.Skus{
		if item:=ezbuy.Get(sku);item!=nil{
			items=append(items,item)
		}else{
			gh.Fail(fmt.Sprint(sku,"：不是EZBUY 商品"))
			return
		}
	}

	if len(items)>0{
		filename,err:=ezbuy.Export(gh.Account(),items...)
		if err==nil{
			gh.Succ(gin.H{"data":filename})
		}else{
			gh.Fail(err.Error())
		}
	}else{
		gh.Fail("not items")
	}
}

func (*EzBuy)Save(c *gin.Context){
	gh:=igin.H(c)
	var item *ezbuy.Item
	if err:=c.BindJSON(&item);err==nil{
		item.Save()
		gh.Succ(nil)
	}else{
		gh.Fail(err.Error())
	}
}

func (*EzBuy)Colors(c *gin.Context){
	gh:=igin.H(c)
	gh.Succ(gin.H{"items":ezbuy.GetColors()})
}

func (*EzBuy)Sizes(c *gin.Context){
	gh:=igin.H(c)
	gh.Succ(gin.H{"items":ezbuy.GetSizes()})
}

func (*EzBuy)MyCategorys(c *gin.Context){
	gh:=igin.H(c)
	gh.Succ(gin.H{"items":ezbuy.GetMyCategorys(gh.Account())})
}

func (*EzBuy)Categorys(c *gin.Context){
	gh:=igin.H(c)
	store:=ezbuy.GetSetting(gh.Account())
	gh.Succ(gin.H{"items":ezbuy.GetCategorys(store.StoreCateId)})
}

func (*EzBuy)AddMyCategory(c *gin.Context){
	gh:=igin.H(c)

	data:=ezbuy.Category{}
	if err:=c.BindJSON(&data);err!=nil{
		gh.Fail(err.Error())
		return
	}
	ezbuy.AddToMyCategory(gh.Account().Uid,data.Id,data.Name)
	gh.Succ(nil)
}

func (*EzBuy)SaveUserProductsFromSource(c *gin.Context){
	param:=struct{
		Results []configs.M `json:"results"`
	}{}

	gh:=igin.H(c)
	err:=c.BindJSON(&param)
	if err!=nil{
		log.Println(err)
		return
	}

	col:=mongodb.Conn().C("ezbuy.fromsource")
	col.Remove(configs.M{"pid":configs.M{"$exists":true}})
	for _,c:=range param.Results{
		col.Insert(c)
	}
	gh.Succ(nil)
}

func (*EzBuy)UserProductsFromSource(c *gin.Context){
	gh:=igin.H(c)
	col:=mongodb.Conn().C("ezbuy.fromsource")
	results:=[]configs.M{}
	col.Find(configs.M{}).All(&results)
	gh.Succ(gin.H{"data":results})
}

func (*EzBuy)GetAttrs(c *gin.Context){
	gh:=igin.H(c)
	cate_id:=configs.Int(c.Param("cid"))

	if cate_id == 93 || cate_id==84 || cate_id == 94 || cate_id==91{
		gh.Succ(gin.H{"styles":ezbuy.Styles})
	}else{
		gh.Fail("not")
	}
}
