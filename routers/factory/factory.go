package factory
import (
	"helper/net/igin"
	"github.com/gin-gonic/gin"
	"factory"
	"helper/configs"

)
func CreateApp(c *gin.Context){
	gh:=igin.H(c)
	app:=factory.App{}
	if err:=c.BindJSON(&app);err!=nil{
		gh.Fail(err.Error())
		return
	}

	app.AuthorId=gh.Account().Uid
	id:=app.Save()
	gh.Succ(gin.H{"data":id})
}

func Apps(c *gin.Context){
	gh:=igin.H(c)
	apps,total:=factory.Listing(gh.Account())
	gh.Succ(gin.H{"items":apps,"total":total})
}

func Get(c *gin.Context){
	gh:=igin.H(c)
	app,err:=factory.Get(c.Param("id"))

	if err!=nil{
		gh.Fail(err.Error())
		return
	}
	gh.Succ(gin.H{"item":app})
}



func Req(c *gin.Context){
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