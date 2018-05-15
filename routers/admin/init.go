package admin

import (
	"helper/configs"
	"helper/net/igin"

	"helper/auth"
	"github.com/gin-gonic/gin"
	"log"
	"regexp"
)

func Index(c *gin.Context){
	c.HTML(200,"index.html",gin.H{})
}

func init(){
	opt:=configs.GetSection("admin")
	dir:=opt["dir"]
	r:=igin.R

	r.Static("/down",dir+"/down")

	r.OPTIONS("/api/:act",func(c *gin.Context){
		c.Header("Access-Control-Allow-Origin","*")
		c.Header("Access-Control-Allow-Headers","Origin, X-Requested-With, Content-Type, Accept, Connection, User-Agent, Cookie,Authorization")
		c.Status(200)
	})

	r.OPTIONS("/api/:act/:act1",func(c *gin.Context){
		c.Header("Access-Control-Allow-Origin","*")
		c.Header("Access-Control-Allow-Headers","Origin, X-Requested-With, Content-Type, Accept, Connection, User-Agent, Cookie,Authorization")
		c.Status(200)
	})

	r.OPTIONS("/api/:act/:act1/:act2",func(c *gin.Context){
		c.Header("Access-Control-Allow-Origin","*")
		c.Header("Access-Control-Allow-Headers","Origin, X-Requested-With, Content-Type, Accept, Connection, User-Agent, Cookie ,Authorization")
		c.Status(200)
	})

	r.OPTIONS("/api/:act/:act1/:act2/:act2",func(c *gin.Context){
		c.Header("Access-Control-Allow-Origin","*")
		c.Header("Access-Control-Allow-Headers","Origin, X-Requested-With, Content-Type, Accept, Connection, User-Agent, Cookie ,Authorization")
		c.Status(200)
	})

	r.OPTIONS("common/api/:act",func(c *gin.Context){
		c.Header("Access-Control-Allow-Origin","*")
		c.Header("Access-Control-Allow-Headers","Origin, X-Requested-With, Content-Type, Accept, Connection, User-Agent, Cookie ,Authorization")
		c.Status(200)
	})


	r.OPTIONS("common/api/:act/:act1",func(c *gin.Context){
		c.Header("Access-Control-Allow-Origin","*")
		c.Header("Access-Control-Allow-Headers","Origin, X-Requested-With, Content-Type, Accept, Connection, User-Agent, Cookie ,Authorization")
		c.Status(200)
	})

	r.OPTIONS("common/api/:act/:act1/:act2",func(c *gin.Context){
		c.Header("Access-Control-Allow-Origin","*")
		c.Header("Access-Control-Allow-Headers","Origin, X-Requested-With, Content-Type, Accept, Connection, User-Agent, Cookie ,Authorization")
		c.Status(200)
	})
	r.Use(CheckIn())
}

func CheckIn() gin.HandlerFunc{
	return func(c *gin.Context){
		defer func(){
			if r:=recover();r!=nil{
				log.Println("异常错误",r)

			}
		}()
		c.Header("Access-Control-Allow-Origin","*")
		c.Header("Access-Control-Allow-Headers","Origin, X-Requested-With, Content-Type, Accept, Connection, User-Agent, Cookie,Authorization")
	}
}

func Auth() gin.HandlerFunc{
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Println("异常错误", r)
			}
		}()
		ok,_:=regexp.MatchString("/api/",c.Request.URL.Path)
		if ok{
			token:=c.Request.Header.Get("Authorization")
			customer:=auth.Customer{}
			customer.VerifyToken(token)

			if !customer.IsLogin(){
				switch c.Request.URL.Path {
				case "/api/account.create":
				case "/api/account.login":
				case "/api/account.sendregistercode":
				case "/api/item.margecontent":
				default:
					c.JSON(200,gin.H{"isSucc":false,"error_msg":"not login","error_code":10000})
					c.Abort()
				}
			}
			c.Set("user",customer)
		}else{
			c.String(400,"error")
			c.Abort()
		}
	}
}