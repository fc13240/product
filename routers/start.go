package routers

import (
	"helper/configs"
	"helper/net/igin"

	"github.com/gin-gonic/gin"


)

func Index(c *gin.Context){
	c.HTML(200,"index.html",gin.H{})
}

func start(){
	opt:=configs.GetSection("admin")
	dir:=opt["dir"]
	r:=igin.R
	
	r.Static("/down",dir+"/down")

	r.OPTIONS("/api/*action",func(c *gin.Context){
		igin.HttpRequestCommon(c)
		c.Status(200)
	})
	
	r.Use(CheckIn())
}

func CheckIn() gin.HandlerFunc{
	return func(c *gin.Context){
		igin.HttpRequestCommon(c)
	}
}
