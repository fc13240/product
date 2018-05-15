package action

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"helper/net/igin"
	"product/shopee"
	"product/shopee/api"
)

func Get(c *gin.Context){
	ig:=igin.H(c)
	if item,err:=shopee.Get(c.Param("sku"));err!=nil{
		ig.Fail(err)
	}else{
		ig.Succ(gin.H{"item":item})
	}
}

func SaveUploadImageInfo(c *gin.Context){
	
}

func UploadImages(c *gin.Context){
	
}

func DelRemoteItem(c *gin.Context){
	ig:=igin.H(c)
	res:=api.DelItem(c.Param("sku"))
	ig.Succ(gin.H{"data":res})
}

func AddRemoteItem(c *gin.Context){
	ig:=igin.H(c)
	if err:=api.AddItem(c.Param("sku"));err==nil{
		ig.Succ(nil)
	}else{
		ig.Fail(err.Error())
	}
}

func SaveProduct(c *gin.Context){
	ig:=igin.H(c)
	item:=shopee.Item{}
	if err:=c.BindJSON(&item);err==nil{
		item.AuthorId=ig.Account().Uid
		fmt.Println(item.GetAttribute())
		item.Attributes=item.GetAttribute()
		item.Logistics=item.GetDefaultLogistics()
		item.Save()
		ig.Succ(nil)
	}else{
		ig.Fail(err.Error())
	}

}
