package product

import (
	"github.com/gin-gonic/gin"
	"helper/net/igin"
	"product/alibaba"
	"product"
	"helper/configs"
	"fmt"
	"helper/dbs"
	"product/ebay"
	product_template "product/template"
	
)

func Search(c *gin.Context) {
	rows := []configs.M{}

	var list []*product.Item
	var total int

	ig:=igin.H(c)

	author:=ig.Account()

	param:=struct{
		*igin.ParamFilter
		*igin.ParamPage
		*igin.ParamSort
	}{}

	if err:=c.BindJSON(&param);err!=nil{
		ig.Fail(err.Error())
		return
	}

	tag:=param.Filter.Get("tag")

	if (tag == "is_new" ){
		list, total = product.NewArrivals().Listing(param.Offset, param.Limit)
	} else if tag == "is_hot" {
		list, total = product.BestSelling().Listing(param.Offset, param.Limit)
	}else {
		
		if author.ISAdmin() == false{
			param.Filter.Set("author_id",author.Uid)
		}

		if tag == "is_under" {
			param.Filter.Set("is_under",1)
		}

		if seller_id:=param.Filter.Int("seller_id");seller_id>0{
			param.Filter.Set("in",alibaba.BySellerItemIds(seller_id) )
		}

	
		list, total = product.Search(param.Filter, param.Offset, param.Limit,param.Sort)

	}
	for _, item := range list {
		row:=configs.M{
			"id":         item.Id,
			"name":       item.Name,
			"quant":      item.Quant,
			"desc":       item.Desc,
			"sku":	      item.Sku,
			"author_id":  item.Authorid,
			"add_time":   item.Addtime.Format("01/02 15:04"),
			"headimg":    item.Headimg,
			"price":      item.Price,
			"old_price":  item.OldPrice,
			"channel_id": item.Channelid,
			"tags":       item.Tags(),
			"labellog"   :item.LabelLog,
			"shopes"	:product.GetSkuUseStatus(author,item.Sku),
		}
		rows = append(rows, row)
	}

	//关联阿里巴巴

	for _,row:=range rows{
		row["aliinfo"]=alibaba.GetAliSupplier(row.Get("sku"),author)
	}

	igin.Succ(c,gin.H{"items": rows, "total": total})
}


func byTag(now_tag string) *product.Sales {

	if now_tag == "is_new" {
		return product.NewArrivals()
	}

	if now_tag == "is_hot" {
		return product.BestSelling()
	}

	return nil
}

func AddToTag(c *gin.Context) {

	if tagname,ok:=c.GetQuery("tag"); ok{
		tag:=byTag(tagname)

		if tag==nil{
			igin.Fail(c,"failing")
			return
		}

		param:=struct{
			Item_ids []int `json:"item_ids"`
		}{}

		c.BindJSON(&param)
		for _, item_id := range param.Item_ids {
			tag.Add(item_id, 255)
		}

		igin.Succ(c,nil)

	}
}

func RemToTag(c *gin.Context) {
	if tagname,ok:=c.GetQuery("tag"); ok {
		tag := byTag(tagname)

		if tag == nil {
			igin.Fail(c, "failing")
			return
		}

		param:=struct{
			Item_ids []int `json:"item_ids"`
		}{}

		c.BindJSON(&param)

		tag.Rem(param.Item_ids...)

	}
}

//设置下架
func SetUnder(c *gin.Context) {

	param:=struct{
		Item_ids []int `json:"item_ids"`
	}{}

	c.BindJSON(&param)
	if len(param.Item_ids)>0 {

		for _, item_id := range param.Item_ids {
			if item, err := product.IdGet(item_id); err == nil {
				item.SetUnder()
			}
		}
		igin.Succ(c,nil)
	} else {
		igin.Fail(c,"parse josn:")
	}
}

//设置上架
func SetUpper(c *gin.Context) {
	param:=struct{
		Item_ids []int `json:"item_ids"`
	}{}
	c.BindJSON(&param)
	if len(param.Item_ids)>0 {
		for _, item_id := range param.Item_ids{
			if item, err := product.IdGet(item_id); err == nil {
				item.SetUpper()
			}
		}
		igin.Succ(c,nil)
	} else {
		igin.Fail(c,"parse josn:")
	}
}

//设置模板
func SetTemp(c *gin.Context) {
	param:=struct{
		Item_ids []int `json:"item_ids"`
		TempId int `json:"temp_id"`
		Desc string `json:"desc"`
	}{}
	c.BindJSON(&param)
	if len(param.Item_ids)>0 {
		for _, item_id := range param.Item_ids{
			if item, err := product.IdGet(item_id); err == nil {
				item.SetTemp(param.TempId, param.Desc)
			} else {
				igin.Fail(c,err.Error())
			}
		}
		igin.Succ(c,nil)
	} else {
		igin.Fail(c,"fail")
	}
}


func Detail(c *gin.Context) {
	sku:=c.Param("sku")
	item, _ := product.Get(sku)
	gallery := item.GetSmallImages()
	data := configs.M{
		"id":          item.Id,
		"sku":         item.Sku,
		"name":        item.Name,
		"en_name":     item.EnName,
		"quant":       item.Quant,
		"desc":        item.Desc,
		"author_id":   item.Authorid,
		"add_time":    item.Addtime.Format("01月02 15点04分"),
		"headimg":     item.Headimg,
		"price":       item.Price,
		"old_price":   item.OldPrice,
		"gallery":     gallery,
		"buying_price":item.BuyingPrice,
		"channel_id":  item.Channelid,
		"ebay_itemid": item.EbayItemid,
		"temp_id":     item.Tempid,
		"length":      item.Length,
		"width":       item.Width,
		"height":      item.Height,
		"weight":      item.Weight,
		"category_id": item.CategoryId,
	}
	igin.Succ(c,gin.H{"item": data})
}

func GetChannelList(c *gin.Context) {
	igin.Succ(c,gin.H{"items": product.GetChannelList()})
}

func EditBase(c *gin.Context) {

	param:=struct{
		Itemid int `json:"id"`
		Name string `json:"name"`
		EnName string `json:"en_name"`
		Quant int `json:"quant"`
		Desc string `json:"desc"`
		Price float32 `json:"price"`

	}{}

	c.BindJSON(&param)

	item := &product.Item{
		Id:     param.Itemid,
		Name:   param.Name,
		Quant: param.Quant,
		Desc:   param.Desc,
		Price:  param.Price,
		EnName:param.EnName,
	}
	if err := item.Save(); err == nil {
		igin.Succ(c,nil)
	} else {
		igin.Fail(c,err.Error())
	}
}

func Save(c *gin.Context) {
	item := &product.Item{}
	gh:=igin.H(c)
	if err:=c.BindJSON(item);err!=nil{
		igin.Fail(c,err.Error())
		return
	}

	item.Authorid=gh.Account().Uid
	if err := item.Save(); err == nil {
		igin.Succ(c,gin.H{"item_id": item.Id})
	} else {
		igin.Fail(c,err.Error())
	}
}

func DelItem(c *gin.Context){
	gh:=igin.H(c)
	param:=struct{
		Ids []int	`json:"ids"`
	}{}
	defer func(){
		if r:=recover();r!=nil{
			fmt.Println(r)
		}
	}()
	if err:=c.BindJSON(&param);err==nil{
		if err:=product.DelItem(gh.Account(),param.Ids...);err==nil{
			gh.Succ(nil)
			return
		}else{
			gh.Fail(err.Error())
			return
		}
	}else{
		gh.Fail(err.Error())
		return
	}
}



func SaveToEbay(c *gin.Context) {
	hc:=igin.H(c)
	if item, err := product.IdGet(hc.GetInt("item_id")); err == nil {
		content, _ := product_template.MergeItem(item)
		if res := ebay.UpdateDescript(item.EbayItemid,content); res.IsSuccess() {
			hc.Succ(nil)
		} else {
			hc.Fail(res.Log)
		}
	} else {
		hc.Fail(err.Error())
	}
}




/*

func SetQuant(c *gin.Context) {
	cart := product.OpenCart(act.Auth.Token())
	if item, err := cart.Open(act.Get("cart_id")); err == nil {
		item.SetQuant(act.GetInt("quant"))
		act.Succ(configs.M{
			"item_total_price": configs.Price(item.TotalPrice()),
			"total_price":      configs.Price(cart.TotalAmount()),
		})
	} else {
		igin.Fail(c,err.Error())
	}
}

func Orders(c *gin.Context) {

	orders := product.MyOrders(act.Auth.Account, configs.M{
		"unpaid": act.Get("unpaid"),
		"paid":   act.Get("paid"),
	})

	var data []configs.M

	for _, order := range orders {
		item := configs.M{
			"order_id":     order.Id,
			"order_sn":     order.OrderNo,
			"order_amount": configs.Price(order.OrderAmount),
			"buyer_id":     order.BuyerId,
			"seller_id":    order.SellerId,
			"is_payment":   order.IsPayment(),
			"order_status": order.OrderStatus,
			"products":     order.GetProducts(),
		}
		data = append(data, item)
	}
	act.Succ(configs.M{"items": data})
}

//创建订单
func CreateOrder(c *gin.Context) {
	data, _ := act.ParseJson()
	shipping, err := account.GetShippingInfo(data.Int("shipping_id"))
	payment_method := data.Get("payment_method")
	order_remark := data.Get("order_remark")
	if err != nil {
		igin.Fail(c,err.Error())
		return
	}
	cart := product.OpenCart(act.Auth)
	if len(cart.Items()) == 0 {
		act.Fail("Your Shopping Cart is empty.")
		return
	}
	order, tx, err := product.CreateOrder(act.Auth.Account, cart.TotalAmount(), payment_method, order_remark)
	defer tx.Rollback()

	if err == nil {
		if err := order.AddItems(cart.Items()); err != nil {
			igin.Fail(c,err.Error())
			return
		}

		if err := order.AddShipping(shipping); err != nil {
			igin.Fail(c,err.Error())
			return
		} else {
			tx.Commit()
			cart.Clean()
			act.Succ(configs.M{"order_id": order.Id})
		}
	} else {
		igin.Fail(c,err.Error())
	}
	//出异常回滚
}

//创建付款
func CreatePayment(c *gin.Context) {

	if order, err := product.OpenOrder(act.GetInt("order_id")); err == nil {
		if url, err := order.CreatePayment(); err == nil {
			act.Succ(configs.M{"payment_url": url})
		} else {
			igin.Fail(c,err.Error())
		}
	} else {
		act.Fail("failing 2")
	}
}

//确认付款
func ApprovedPayment(c *gin.Context) {
	err := product.ApprovedPayment(act.Get("paymentId"), act.Get("PayerID"), act.Get("token"))
	if err == nil {
		igin.Succ(c,nil)
	} else {
		igin.Fail(c,err.Error())
	}
}
*/
//国家
func Country(c *gin.Context) {
	rows := dbs.Rows("SELECT `code`,`name` FROM country WHERE is_disable=?", 0)
	data := []configs.M{}

	defer rows.Close()

	for rows.Next() {
		var code string
		var name string
		rows.Scan(&code, &name)
		data = append(data, configs.M{"id": code, "name": name})
	}
	igin.Succ(c,gin.H{"items": data})
}

//缩略图
/*
func ResizeAllByProduct(c *gin.Context) {

	hc:=igin.H(c)
	if item_id := hc.GetInt("item_id"); item_id > 0 {
		product.ResizeAllByProduct(item_id)
		hc.Succ(nil)
	} else {
		param:=struct{
			Item_ids []int `json:"item_ids"`
		}{}
		c.BindJSON(&param)
		for _, item_id := range param.Item_ids {
			product.ResizeAllByProduct(item_id)
		}
		hc.Succ(nil)
	}
}


func RemoveImage(c *gin.Context) {
	hc:=igin.H(c)
	if image_id := hc.GetInt("image_id"); image_id > 0 {
	//	gallery.(image_id)
		hc.Succ(nil)
	} else {
		hc.Fail("error")
	}
}

*/

func MargeContent(c *gin.Context) {

	c.Header("Content-Type","text/html")

	var content string
	var err error

	item_id,_:=c.GetQuery("item_id")

	tmp_id,_:=c.GetQuery("tmp_id")

	if ispreview,ok:=c.Get("ispreview") ;ok && ispreview == "1" {
		body := c.PostForm("body")
		content, err = product_template.MergeContent(configs.Int(item_id), body)
	} else {
		itme,_:=product.IdGet(configs.Int(item_id))
		itme.Tempid=configs.Int(tmp_id)
		content, err = product_template.MergeItem(itme)
	}
	if err != nil {
		c.String(200,err.Error())
	} else {
		c.String(200,content)
	}
}

/*
//下载到仓库
func DownToWarehouse(c *gin.Context) {
	data, err := act.ParseJson()
	if err != nil {
		igin.Fail(c,err.Error())
		return
	}
	if err := product.DownToWarehouse(data.Get("down_url"), act.Auth.Account); err == nil {
		igin.Succ(c,nil)
	} else {
		igin.Fail(c,err.Error())
	}
}

//查看仓库商品
func GetWarehouseProductInfo(c *gin.Context) {
	if detail, err := product.GetWarehouseProductInfo(act.GetInt("item_id")); err == nil {
		act.Succ(configs.M{"item": detail})
	} else {
		igin.Fail(c,err.Error())
	}
}
*/
func TemplateOpts(c *gin.Context) {
	list := product_template.GetTemplateTags()
	igin.Succ(c,gin.H{"items": list})
}

func Labels(c *gin.Context){
	gh:=igin.H(c)
	gh.Succ(gin.H{"items":product.Labels()})
}

func LabelLogs(c *gin.Context){
	gh:=igin.H(c)
	gh.Succ(gin.H{"items":product.ItemLabelLogs(configs.Int(c.Param("ordernum")))})
}

func AddLabelLog(c *gin.Context){
	gh:=igin.H(c)

	param:=struct {
		LabelId int `json:"labelid"`
		Itemids []int `json:"itemids"`
		Remarks string `json:"remarks"`
	}{}

	if err:=c.BindJSON(&param);err!=nil{
		gh.Fail(err.Error())
		return
	}

	if label,err:=product.GetLabel(param.LabelId);err==nil{
		label.AddLog(param.Remarks,param.Itemids...)
		gh.Succ(nil)
	}else{
		gh.Fail(err.Error())
		return
	}
}




//上传记录
func AddSkuUploadLog(c *gin.Context){
	ig:=igin.H(c)
	parmas:=struct{
		Sku string `json:"sku"`
		Storeid int `json:"store_id"`
	}{}
	if err:=c.BindJSON(&parmas);err!=nil{
		return 
	}
	
	product.AddSkuUploadLog(ig.Account(),parmas.Sku,parmas.Storeid)
}

func GetSkuUseStatus(c *gin.Context){
	sku:=c.Param("sku")
	ig:=igin.H(c)
	items:=product.GetSkuUseStatus(ig.Account(),sku)
	ig.Succ(gin.H{"items":items})
}

func SaveField(c *gin.Context){
	sku:=c.Param("sku")
	ig:=igin.H(c)
	if item,err:=product.Get(sku);err==nil{
		param:=struct{
			Field string `json:"field"`
			Value interface{} `json:"value"`
		}{}
		if err:=	c.BindJSON(&param);err==nil{
		item.SaveField(param.Field,param.Value)
		ig.Succ(nil)
	}else{
		ig.Fail(err.Error())
	}
	}else{
		ig.Fail(err.Error())
	}
}

func GetSkuSelectedOptions(c *gin.Context){
	ig:=igin.H(c)
	opts:=product.GetSkuSelectedOptions(c.Param("sku"),nil)
	ig.Succ(gin.H{"items":opts})
}

//获取平台分类信息
func GetPlatformCategorys(c *gin.Context){
	cid:=configs.Int(c.Param("cid"))
	platform:=c.Param("platform")
	igin.Succ(c,gin.H{
		"data":product.GetPlatformCategorys(cid,platform),
	})
}

//获取已选择的分类
func GetPlatformCategorySelected(c *gin.Context){
	cid:=configs.Int(c.Param("cid"))
	platform:=c.Param("platform")
	igin.Succ(c,gin.H{
		"data":product.GetPlatformCategorySelected(cid,platform),
	})
}

//产品的分类信息
func GetItemCategoryInfo(c *gin.Context){
	sku:=c.Param("sku")
	if item,err:=product.Get(sku);err==nil{
		if item.CategoryId>0{
			if cate,err:=product.GetCategory(item.CategoryId);err==nil{
				igin.Succ(c,gin.H{"data":cate})
				return 
			}
		}else{
			igin.Fail(c,"not set category")
		}
	}else{
		igin.Fail(c,err.Error())
	}
}