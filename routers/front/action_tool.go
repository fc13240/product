package front
import (
	"helper/net/igin"
	"github.com/gin-gonic/gin"
	"tools"
	"factory"
	"helper/configs"
	"helper/util"
	"fmt"
)
type ToolAction int
var toolAction ToolAction

func (*ToolAction)GetAlbabaProductImages(c *gin.Context) {

	param:=struct{
		Url string `json:"url"`
	}{}
	gh:=igin.H(c)

	if err:=c.BindJSON(&param);err!=nil{
		gh.Fail(err.Error())
		return
	}
	items:=tools.GetAlbabaTupian(param.Url)

	gh.Succ(gin.H{"data":items})

}

func (*ToolAction)Md5(c *gin.Context){
	param:=struct{
		Value string `json:"value"`
	}{}

	gh:=igin.H(c)

	if err:=c.BindJSON(&param);err!=nil{
		gh.Fail(err.Error())
		return
	}
	gh.Succ(gin.H{"data":util.Md5(param.Value)})
}


func (*ToolAction)JsonFormat(c *gin.Context){
	param:=struct{
		Value string `json:"value"`
	}{}

	gh:=igin.H(c)

	if err:=c.BindJSON(&param);err!=nil{
		gh.Fail(err.Error())
		return
	}

	fmt.Println(param.Value)

	gh.Succ(gin.H{"data":tools.StringToJson(param.Value)})
}

func (*ToolAction)Apps(c *gin.Context){
	gh:=igin.H(c)
	apps,total:=factory.Search()

	items:=[]configs.M{}
	for _,app:=range apps{
		items=append(items,configs.M{"name":app.Name,"code":app.Code,"seo":app.Seo})
	}
	gh.Succ(gin.H{"items":items,"total":total})
}

func (*ToolAction)Get(c *gin.Context){
	gh:=igin.H(c)
	app,err:=factory.GetByCode(c.Param("code"))

	if err!=nil{
		gh.Fail(err.Error())
		return
	}

	for _,element:=range app.Screen.Elements{
		 factory.ParseElementValue(&element.Value)
	}
	gh.Succ(gin.H{"item":app})
}

func (*ToolAction)Req(c *gin.Context){
	var appid=c.Param("id")
	gh:=igin.H(c)

	var params=configs.M{}

	if e:=c.BindJSON(&params);e!=nil{
		gh.Fail(e.Error())
		return
	}else{
		res,err:=factory.Req(appid,params)
		if err!=nil{
			gh.Fail(err.Error())
			return
		}
		rr:=configs.M{}
		if err:=res.BindJSON(&rr);err==nil{
			c.JSON(200,rr)
		}else{
			gh.Fail(err.Error())
		}
	}
}