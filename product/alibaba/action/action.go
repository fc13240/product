package action
import (
	"helper/account"
	"github.com/gin-gonic/gin"
	"helper/net/igin"
	"product/alibaba"
	"helper/configs"
	"fmt"
	"product"
	"strings"
	"errors"
	"net/url"
	"log"
	"time"
	"product/gallery"
)
//通过一个url下载商品信息
func  down(sourcer_url string,author *account.Account)(item *alibaba.Item ,err  error ){

	urlInfo,err:=url.ParseRequestURI(sourcer_url)
	if err!=nil{
		return nil,errors.New("错误的网址:"+sourcer_url)
	}
	switch urlInfo.Host{
	case "item.taobao.com":
		taobao:=alibaba.Taobao{}
		return taobao.Down(sourcer_url,author)
	case "detail.1688.com":
		return alibaba.Down(sourcer_url,author)
	}
	return item,errors.New("没有找到啊")
}


func CheckUrlExist(c *gin.Context){
	hg:=igin.H(c)

    param:=struct{
		BuyAddr string `json:"buy_addr`
	}{}

	if err:=c.BindJSON(&param);err!=nil{
		hg.Fail(err.Error())
		return
	}

	id:=alibaba.GetUrlId(param.BuyAddr)

	if id == 0{
		hg.Fail("商品ID为0")
		return
	}
	author:=hg.Account()
	if exist_item,err:=alibaba.GetById(id,author) ;err==nil{ //获取已经存在的商品信息
		hg.Fail(fmt.Sprint("这个商品已经存在,SKU是:",exist_item.Sku),gin.H{"pid":exist_item.Id,"item_id":exist_item.ItemId,"sku":exist_item.Sku})
		return
	}else{
		hg.Succ(nil)
	}
}

func Add(c *gin.Context){
	hg:=igin.H(c)
	var item alibaba.Item

	if err:=c.BindJSON(&item);err!=nil{
		hg.Fail(err.Error())
		return
	}

	item.Id=alibaba.GetUrlId(item.BuyAddr)

	if item.Id == 0{
		hg.Fail("商品ID为0")
		return
	}

	author:=hg.Account()

	item.Authorid=author.Uid
	item.Addtime=time.Now()

	item.Source="alibaba"

	if item.Exist(){
		exist_item,_:=alibaba.GetById(item.Id,author) //获取已经存在的商品信息
		hg.Fail(fmt.Sprint("这个商品已经存在,SKU是:",exist_item.Sku),gin.H{"pid":exist_item.Id,"item_id":exist_item.ItemId,"sku":item.Sku})
		return
	}
	sku,err:=product.NewSku()
	if err==nil{
		item.Sku=sku.Name
	}else{
		hg.Fail(err.Error())
		return 
	}
	
	seller:=&alibaba.Seller{CompanyName:item.CompanyName,CompanyUrl:item.CompanyUrl}
	item.SellerId=seller.Save()


	if len(item.Images) == 0{
		item.DownAllImage()
	}

	item.Save()

	sku.AddBuySite(item.BuyAddr)

	if err:=copyToItem(&item);err!=nil{
		hg.Fail(err.Error())
	}
	hg.Succ(gin.H{"item_id":item.ItemId,"pid":item.Id,"sku":item.Sku})
}

//获取商品属性信息
func GetAttrs(c *gin.Context){
	item_id:=configs.Int(c.Param("item_id"))
	attrs:=alibaba.GetAttrs(item_id)
	hg:=igin.H(c)
	hg.Succ(gin.H{"attrs":attrs})
}

//复制为商品主表
func copyToItem(item *alibaba.Item) error{
	pItem:=product.Item{
		Name:item.Title,
		Id:item.ItemId,
		Remarks:item.CompanyName,
		Desc:item.Desc,
		Authorid:item.Authorid,
		Sku:item.Sku,
		Length:30,
		Height:5,
		Width:20,
		Weight:0.3,
	}

	if item.BuyPrice >0{
		pItem.BuyingPrice=item.BuyPrice
		pItem.Price=item.BuyPrice
	}

	if err:=pItem.Save();err==nil{
		if item.ItemId==0 || item.ItemId!=pItem.Id {
			item.ItemId=pItem.Id
			item.Set(gin.H{"item_id":item.ItemId,"sku":item.Sku})
			if len(item.Images)>0{
				gallery.AddImage(pItem.Sku,item.Images...)
			}
		}
		return nil
	}else{
		return err
	}
}

//下载
func Down(c *gin.Context){
	param:= struct {
		Url string `json:"down_url"`
	}{}
	hg:=igin.H(c)
	if err:=c.BindJSON(&param);err!=nil{
		hg.Fail(err.Error())
		return
	}
	if item,err:=down(param.Url,hg.Account());err!=nil{
		hg.Fail(err.Error())
	}else{
		if sku,err:=product.NewSku();err==nil{
			item.Sku=sku.Name
			sku.AddBuySite(param.Url)
		}else{
			hg.Fail(err.Error())
			return 
		}
		
		if err:=copyToItem(item);err==nil {
			hg.Succ(gin.H{"item_id":item.ItemId})
		}
	}
}

//下载所有源
func DownSource(c *gin.Context){
	param:= struct {
		Source string `json:"source"`
		ItemId int `json:"itemid"`
		Sku    string `json:"sku"`
	}{}

	hg:=igin.H(c)
	author:=hg.Account()
	if err:=c.BindJSON(&param);err!=nil{
		hg.Fail(err.Error())
		return
	}

	if item,err:=down(param.Source,author);err!=nil{
		hg.Fail(err.Error())
		return
	}else{
		if param.Sku!="" && param.ItemId ==0 {
			item.Sku=param.Sku
			if item_id,isExist:=product.SkuExist(param.Sku,author);isExist == false{
				if err:=copyToItem(item);err!=nil{
					hg.Fail(err.Error())
					return
				}
				param.ItemId=item.ItemId
			}else{
				param.ItemId=item_id
			}
		}
		item.Set(configs.M{"item_id":param.ItemId,"sku":param.Sku})

		item.Sku=param.Sku

		hg.Succ(gin.H{"item":gin.H{
			"company_name":item.CompanyName,
			"company_url":item.CompanyUrl,
			"buyaddr":item.BuyAddr,
			"id":item.Id,
		}})
		return
	}
}

//通过ID 获得信息
func BySource(c *gin.Context){
	gh:=igin.H(c)
	if id,ok:=c.GetQuery("id");ok{
		if item,err:=alibaba.GetById(configs.Int64(id),gh.Account());err==nil{
			gh.Succ(gin.H{"item":item})
			return
		}else{
			gh.Fail(err.Error())
		}
	}else if sku,ok:=c.GetQuery("sku");ok{
		if item,err:=alibaba.Get(sku,gh.Account());err==nil{
			gh.Succ(gin.H{"item":item})
			return
		}else{
			gh.Fail(err.Error())
		}
	}
	gh.Fail("id no exist")
}

func Sources(c *gin.Context){

	sku,_:=c.GetQuery("sku")
	hg:=igin.H(c)

	author:=hg.Account()

	items:=[]struct{
		CompanyName string `json:"company_name" bson:"company_name"`
		CompanyUrl string `json:"company_url" bson:"company_url"`
		SellerId  int `json:"seller_id" bson:"seller_id"`
		OrderBy int `json:"orderby" bson:"orderby"`
		Buyaddr string `json:"buyaddr" bson:"buyaddr"`
		Id int `json:"id" bson:"id"`
	}{}

	alibaba.Col().Find(configs.M{"sku":sku,"authorid":author.Uid}).Select(configs.M{"id":true,"company_name":true,"company_url":true,"seller_id":true,"orderBy":true,"buyaddr":true}).All(&items)

	hg.Succ(gin.H{"items":items})
}

func DelSource(c *gin.Context){}

//填坑
func Tiankeng(c *gin.Context){
	param:= struct {
		Flag string `json:"flag"`
		Rows []string `json:"rows"`
	}{}

	hg:=igin.H(c)
	if err:=c.BindJSON(&param);err!=nil{
		hg.Fail(err.Error())
		return
	}

	for _,row:=range param.Rows{

		var data=configs.M{}
		strings.Replace(row," ","",-1)
		line:=strings.Split(row,";;")

		for _,v:=range line{
			curr:=strings.Split(v,"::")
			if len(curr)==2{
				field:=curr[0]
				value:=curr[1]
				data[field]=value
			}
		}

		if param.Flag=="buyaddr"{
			id:=alibaba.GetUrlId(data.Get("buyaddr"))

			if id==0{
				continue
			}

			delete(data,"buyaddr")
			if len(data) == 0{
				log.Println("data is empty")
				continue
			}
			if item,err:=alibaba.GetById(id,hg.Account());err==nil{
				item.Set(data)
				if data.Get("sku")!="" {

					pItem, err := product.Get(item.Sku)

					if  err == nil {
						pItem.SetSku(data.Get("sku"))
					}
				}
			}else{
				fmt.Println("get error ",err.Error())
			}
		}
	}
	hg.Succ(nil)
}

func DownAll(c *gin.Context){
	param:= struct {
		Urls []string `json:"urls"`
	}{}

	hg:=igin.H(c)

	if err:=c.BindJSON(&param);err!=nil{
		hg.Fail(err.Error())
		return
	}

	var succ = 0
	
	for _,addr:=range param.Urls{
		if addr==""{
			continue
		}
		if item,err:=down(addr,hg.Account());err!=nil{
			fmt.Println(err.Error())
		}else{
			copyToItem(item)
			succ++
		}
	}
	hg.Fail(fmt.Sprint("处理:",succ))
}

func Save(c *gin.Context){
	item:=alibaba.Item{}
	gh:=igin.H(c)
	if err:=c.BindJSON(&item);err==nil{
		item.Save()
		gh.Succ(nil)
	}else{
		gh.Fail(err.Error())
	}
}

func Set(c *gin.Context){
	param:= struct {
		Id int64 `json:"id"`
		Field string `json:"field"`
		Value string `json:"value"`
	}{}

	gh:=igin.H(c)
	if err:=c.BindJSON(&param);err==nil{

		if item,err:=alibaba.GetById(param.Id,gh.Account());err==nil{
			item.Set(gin.H{param.Field:param.Value})
			gh.Succ(nil)
		}else{
			gh.Fail(err.Error())
		}

	}else{
		gh.Fail(err.Error())
	}
}

func Get(c *gin.Context){
	gh:=igin.H(c)
	if sku,ok:=c.GetQuery("sku");ok {
		if item,err:=alibaba.Get(sku,gh.Account());err==nil{
			gh.Succ(gin.H{"item":item})
			return
		}
	}
	gh.Fail("")
}

func Sellers(c *gin.Context){
	gh:=igin.H(c)
	param:= struct {
		Offset int `json:"offset"`
		Limit int `json:"limit"`
	}{}
	if err:=c.BindJSON(&param);err==nil {
		author_id:=gh.Account().Uid
		if gh.Account().ISAdmin() {
			author_id=0
		}
		rows,total:=alibaba.SellerListing(author_id,param.Offset,param.Limit)
		gh.Succ(gin.H{"items":rows,"total":total})
	}else{
		gh.Fail(err.Error())
	}
}

func Search(c *gin.Context){
	gh:=igin.H(c)
	param:= struct {
		igin.ParamPage
		igin.ParamFilter
	}{}
	if err:=c.BindJSON(&param);err==nil {
		rows,total:=alibaba.Search(param.Filter,param.Offset,param.Limit)
		gh.Succ(gin.H{"items":rows,"total":total})
	}else{
		gh.Fail(err.Error())
	}
}

func Delete(c *gin.Context){
	gh:=igin.H(c)
	if err:=alibaba.Del(configs.Int64(c.Param("id")));err==nil {
		gh.Succ(nil)
	}else{
		gh.Fail(err.Error())
	}
}

func GetColorsAndSizes(c *gin.Context){
	gh:=igin.H(c)
	gh.Succ(gin.H{"data":alibaba.GetColorsAndSizes(c.Param("sku"))})
}