package action
import (
	"github.com/gin-gonic/gin"
	"helper/net/igin"
	"product"
	"helper/configs"
)

//获取分类下的属性
func GetCategoryAttrs(c *gin.Context){
	
	param:=struct{
		Platform string `json:"platform"`
		Sku string  `json:"sku"`
		Type string `json:"type"`
	}{}

	ig:=igin.H(c)
	if err:=c.BindJSON(&param);err!=nil{
		ig.Fail(err.Error())
		return
	}

	pro,err:=product.Get(param.Sku)

	if err!=nil{
		ig.Fail(err.Error())
		return
	}

	if pro.CategoryId==0{
		ig.Fail("没有设置分类")
		return
	}
	var attrs []product.Attr
	if param.Type=="sku"{
		attrs,err=product.GetCategorySkuAttrs(param.Platform,pro.CategoryId,true)
	}else{
		attrs,err=product.GetCategoryAttrs(param.Platform,pro.CategoryId,true)
	}

	for _,att:=range attrs{
		att.FlagCheckedOption(param.Sku)
	}
	if err==nil{
		ig.Succ(gin.H{"items":attrs})
	}else{
		ig.Fail(err.Error())
	}
}

//保存产品属性
func SaveAttrVal(c *gin.Context){
	params:=struct{
		Sku string `json:"sku"`
		Attid int `json:"attid"`
		Selected []product.AttrOptions `json:"selected"`
		StoreGroup string `json:"store_group"`
		Type string `json:"att_type"`
	}{}
	ig:=igin.H(c)
	if err:=c.BindJSON(&params);err==nil{
		product.SaveAttrVal(params.Sku,params.Attid,params.Selected)
		ig.Succ(nil)
	}else{
		ig.Fail(err.Error)
	}
}

//获得所有属性
func GetAttrs(c *gin.Context){
	ig:=igin.H(c)
	attrs:=product.GetAttrs()
	ig.Succ(gin.H{"data":attrs})
}

//添加属性
func AddAttr(c *gin.Context){
	ig:=igin.H(c)
	param:=&struct{
		Lable string `json:"label"`
		Name string `json:"name"`
		InputType string `json:"input_type"`
		Type string `json:"type"`
	}{}
	
	if err:=c.BindJSON(param);err!=nil{
		ig.Fail(err)
		return
	}

	if _,err:=product.NewAtt(param.Lable,param.Name,param.InputType);err==nil{
		ig.Fail(err.Error())
	}else{
		ig.Succ(nil)
	}
}

func AddAttrOpt(c *gin.Context){
	ig:=igin.H(c)
	param:=&struct{
		Attid int `json:"attid"`
		Value string `json:"value"`
		CnValue string `json:"cnvalue"`
	}{}
	
	if err:=c.BindJSON(param);err!=nil{
		ig.Fail(err)
		return
	}
	if _,err:=product.NewAttOption(param.Attid,param.Value,param.CnValue);err==nil{
		ig.Fail(err.Error())
	}else{
		ig.Succ(nil)
	}
}

func GetAttOptions(c *gin.Context){
	attrid:=configs.Int(c.Param("attid"))
	items:=product.GetAttr(attrid).GetOptions()
	ig:=igin.H(c)
	ig.Succ(gin.H{"items":items})
}

//获得商品选中的属性
func GetAttrVal(c *gin.Context){
	sku:=c.Param("sku")
	selected:=product.GetAttrVal(sku)
	ig:=igin.H(c)
	ig.Succ(gin.H{"data":selected})
}

//
func GetOneAttrSelectedOption(c *gin.Context){
	sku:=c.Param("sku")
	attid:=configs.Int(c.Param("attid"))
	items:=product.GetOneAttrSelectedOption(sku,attid)
	ig:=igin.H(c)
	ig.Succ(gin.H{"items":items})
}

//包裹长宽高，重量填充
func PackageFill(c *gin.Context){
	hg:=igin.H(c)
	data:= []struct{
		length  int `json:"length"`
		Width int `json:"width"`
		Height int `json:"height"`
		Weight float32 `json:"weight"`
	}{{30,20,5,0.3},{30,20,5,0.2}}
	hg.Succ(gin.H{"data":data})
}