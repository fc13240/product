package api


import (

	"github.com/gin-gonic/gin"
	"helper/net/igin"
	"helper/mongodb"

	"helper/configs"
	"fmt"
)
func init()  {
	m:=igin.R.Group("/api/manage/mongodb")

	m.GET("/cols",func (c *gin.Context){
		gh:=igin.H(c)
		gh.Succ(gin.H{"items":mongodb.Cols()})
	})

	m.POST("/listing", func(c *gin.Context) {
		gh:=igin.H(c)
		param:=struct{
			Name string  `name`
			Find configs.M `find`
			Limit int `json:"limit"`
			Offset int `json:"offset"`
			Sort []string `json:"sort"`
		}{}

		if err:=c.BindJSON(&param);err!=nil{
			gh.Fail(fmt.Sprint("无效的参数:",err.Error()))
			return
		}
		fmt.Println(param)
		list,total:=mongodb.Listing(param.Name,param.Find,param.Offset,param.Limit,param.Sort...)
		gh.Succ(gin.H{"items":list,"total":total})
	})

	m.POST("/commendrun",func(c *gin.Context){
		gh:=igin.H(c)
		param:=struct{
			Commend string `json:"commend"`
		}{}

		if err:=c.BindJSON(&param);err!=nil{
			gh.Fail(fmt.Sprint("无效的参数:",err.Error()))
			return
		}

		result:=mongodb.CommendRun(param.Commend)
		gh.Succ(gin.H{"data":result})
	})

}
