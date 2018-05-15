package igin

import (
	"github.com/gin-gonic/gin"
	"helper/auth"
	"fmt"
	"helper/configs"
	"helper/account"
)

var R *gin.Engine

func init(){
	R=gin.Default()
}

func NewDef() *gin.Engine{
	return gin.Default()
}

func HttpRequestCommon(c *gin.Context){
	c.Header("Access-Control-Allow-Origin","*")
	c.Header("Access-Control-Allow-Methods","POST,DELETE,PUT")
	c.Header("Access-Control-Allow-Headers","Origin, X-Requested-With, Content-Type, Accept, Connection, User-Agent, Cookie,Authorization")
}

type Helper struct{
	c *gin.Context
}

func H(c *gin.Context)*Helper{
	return &Helper{c}
}

func (h *Helper)Succ(v gin.H) {
	if v==nil{
		v=gin.H{}
	}
	v["isSucc"]=true
	h.c.JSON(200,v)
}

func (h *Helper)Token()(token auth.Token){
	if v,ok:=h.c.GetQuery("token");ok{
		token=auth.Token(v)
	}
	return token
}

func Succ(c *gin.Context,v gin.H){
	H(c).Succ(v)
}

func Fail(c *gin.Context,error_msg string){
	H(c).Fail(error_msg)
}

func  (h *Helper)Fail(error_msg interface{},other ...gin.H) {
	out:=gin.H{"isSucc":false,"error_msg":error_msg}

	if len(other)>0{
		for _,data:=range other{
			for k,v:=range data{
				out[k]=v
			}
		}
	}
	h.c.JSON(200,out)
}


func (h *Helper)GetInt(k string)int {
	if i,ok:=h.c.GetQuery(k);ok{
		return configs.Int(i)
	}
	return 0
}

func (h *Helper)Account() *account.Account{

	defer func() {
		if e:=recover();e!=nil{
			fmt.Println(e)
		}
	}()
	if v,ok:=h.c.Get("token");ok{
		token:=v.(auth.Token)
		uid:=token.GetInt("uid")
		return account.Find(uid)
	}
	
	panic("token失败")

}

func (h *Helper)User()*auth.Customer{

	defer func() {
		if e:=recover();e!=nil{
			fmt.Println(e)
		}
	}()

	if user,ok:=h.c.Get("user");ok{
		userr:=user.(auth.Customer)
		return &userr
	}

	panic("token失败")
}



type ParamPage struct{
	Offset int `json:"offset"`
	Limit int `json:"limit"`
}

type ParamFilter struct{
	Filter configs.M `json:"filter"`
}

type ParamSort struct{
	Sort string `json:"sort"`
}