package action

import (

	"github.com/gin-gonic/gin"
	"helper/net/igin"
	"product/lazada"
)

func SaveUploadImageInfo(c *gin.Context){
	ig:=igin.H(c)
	param:= struct {
		Sku string `json:"sku"`
		Rel string `json:"rel"`
	}{}

	c.BindJSON(&param)
	
	if img,err:=lazada.SaveUploadImageInfo(param.Sku,param.Rel);err==nil{
		ig.Succ(gin.H{"image":img})
	}else{
		ig.Fail(err.Error())
	}
}


//上传一个sku的所有图片到lazada
func UploadImages(c *gin.Context){
	ig:=igin.H(c)
	sku:=c.Param("sku")
	ig.Succ(gin.H{"images":lazada.Images(sku)})
}

//保存产品
func SaveProduct(c *gin.Context){
	ig:=igin.H(c)
	item:=&lazada.Item{}

	if err:=c.BindJSON(item);err==nil{
		item.AuthorId=ig.Account().Uid
		lazada.SaveProduct(item)
		ig.Succ(nil)
	}else{
		ig.Fail(err.Error())
	}
}

func GetProduct(c *gin.Context){
	ig:=igin.H(c)
	sku:=c.Param("sku")

	item:=lazada.GetProduct(ig.Account().Uid,sku)
	ig.Succ(gin.H{"item":item})
}

func Create(c *gin.Context){
	ig:=igin.H(c)
	sku:=c.Param("sku")
	err:=lazada.GetProduct(ig.Account().Uid,sku).Create()
	if err==nil{
		ig.Succ(nil)
	}else{
		ig.Fail(err.Error())
	}
}
