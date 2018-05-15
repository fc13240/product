package action 
import (
	"helper/configs"
	"product"
	"github.com/gin-gonic/gin"
	"helper/net/igin"
)

func Add(c *gin.Context){
	cate:=&product.Category{}
	ig:=igin.H(c)

	if err:=c.BindJSON(cate);err!=nil{
		ig.Fail(err.Error)
		return
	}
	product.AddCategory(cate)
	ig.Succ(nil)
}

func Get(c *gin.Context){
	ig:=igin.H(c)
	if cid,ok:=c.GetQuery("cid");ok{
		if item,err:=product.GetCategory(configs.Int(cid));err==nil{
			ig.Succ(gin.H{"item":item})
		}else{
			ig.Fail(err.Error)
		}
	}
}

func Categorys(c *gin.Context){
	ig:=igin.H(c)
	items:=product.Categorys(ig.GetInt(c.Param("pid")))

	if sku,ok:=c.GetQuery("sku");ok{
		item,_:=product.GetBySku(sku,ig.Account())
		if item.CategoryId>0 {
			ig.Succ(gin.H{"items":items,"selected":item.CategoryId})
			return
		}
	}
	ig.Succ(gin.H{"items":items})
}

func Listing(c *gin.Context){
	items,count:=product.CategoryListing(nil,0,100)
	ig:=igin.H(c)
	ig.Succ(gin.H{"items":items,"count":count})
}

func GetAttr(c *gin.Context){
	attrid:=configs.Int(c.Param("attid"))
	item:=product.GetAttr(attrid)
	ig:=igin.H(c)
	ig.Succ(gin.H{"item":item})
}

func GetAttrIds(c *gin.Context){
	ig:=igin.H(c)
	cid:=configs.Int(c.Param("cid"))
	platform:=c.Param("platform")
	
	data:=product.GetCategoryAttrIds(platform,cid)
	ig.Succ(gin.H{"data":data})
}


func SetAttr(c *gin.Context){
	ig:=igin.H(c)
	cid:=configs.Int(c.Param("cid"))
	platform:=c.Param("platform")
	param:=&struct{
		Selected map[int]int `json:"selected"`
	}{}

	if err:=c.BindJSON(param);err!=nil{
		ig.Fail(err.Error)
		return
	}
	if cat,err:=product.GetCategory(cid);err==nil{
		
		cat.Platform=platform


		if err=cat.SetAttr(param.Selected);err!=nil{
			ig.Fail(err.Error())
		}else{
			ig.Succ(nil)
		}
	}else{
		ig.Fail(err.Error())
	}
}