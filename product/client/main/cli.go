package main

import (
	"product/ezbuy"
	"log"
	"github.com/gin-gonic/gin"
	"helper/configs"
	"product/alibaba"
	"helper/images"
	"helper/util"
	"time"
	"os"
	"fmt"
	"path/filepath"
	ezapi "product/ezbuy/api"
	ezclient "product/ezbuy/client"

)

var setting *ezbuy.Setting
var client *ezclient.HClient
var ez *ezapi.EzbuyeApi

var isRefresh =false

func Start() gin.HandlerFunc{
	return func(c *gin.Context){

		c.Header("Access-Control-Allow-Origin","*")
		c.Header("Access-Control-Allow-Headers","Origin, X-Requested-With, Content-Type, Accept, Connection, User-Agent, Cookie ,Authorization")

		token:=c.Request.Header.Get("Authorization")

		if token==""{
			c.JSON(200,gin.H{"isSucc":false,"error_msg":"token不存在"})
			c.Abort()
			return
		}

		client=&ezclient.HClient{Token:token}

		res,err:=client.GetStting()

		if err!=nil{
			log.Println(err.Error())
			c.Abort()
			return
		}

		ss:=struct {
			IsSucc bool `json:"isSucc"`
			Setting ezbuy.Setting `json:"item"`
		}{}

		if err:=res.BindJSON(&ss);err!=nil{
			c.JSON(200,gin.H{"isSucc":false,"error_msg":err.Error()})
			c.Abort()
			return
		}

		if !ss.IsSucc{
			c.JSON(200,gin.H{"isSucc":false,"error_msg":"请求服务器验证失败"})
			c.Abort()
			return
		}
		setting= &ss.Setting

		ez=&ezapi.EzbuyeApi{Cookie:setting.Cookie,ReqId:setting.Reqid,ShopName:setting.StoreName,Client:client}

	}
}

func main(){

	r:=gin.Default()

	r.OPTIONS("/local/:act",func(c *gin.Context){
		c.Header("Access-Control-Allow-Origin","*")
		c.Header("Access-Control-Allow-Headers","Origin, X-Requested-With, Content-Type, Accept, Connection, User-Agent, Cookie,Authorization")

		c.Status(200)
	})

	r.OPTIONS("/local/:act/:act1",func(c *gin.Context){
		c.Header("Access-Control-Allow-Origin","*")
		c.Header("Access-Control-Allow-Headers","Origin, X-Requested-With, Content-Type, Accept, Connection, User-Agent, Cookie ,Authorization")

		c.Status(200)
	})

	c:=r.Group("/local",Start())

	c.GET("/upitems",func(c *gin.Context){

		if setting.ItemsNum <1{
			c.JSON(200,gin.H{"isSucc":false,"error_msg":"商品数量不能小于1"})
			return
		}

		if setting.ItemsNum<=48{
			ez.GenList(setting.ItemsNum)
		}else{
			var i=0
			for {
				ez.GenList(i*48)
				if setting.ItemsNum<i*48{
					c.JSON(200,gin.H{"isSucc":true})
					return
				}
				i++
			}
		}
		c.JSON(200,gin.H{"isSucc":true})
	})

	c.GET("/runRefresh",func(c *gin.Context){
		if isRefresh {
			c.JSON(200,gin.H{"isSucc":false,"error_msg":"已经启动"})
			return
		}
		isRefresh=true
		go func(){
			ez.RefreshAll(setting)
		}()
		c.JSON(200,gin.H{"isSucc":true})
	})

	c.GET("/uporders",func(c *gin.Context){
		ez.GetOrders(0)

		c.JSON(200,gin.H{"isSucc":true})
	})

	c.GET("/downitem/:id",func(c *gin.Context){
		id:=configs.Int64(c.Param("id"))
		if res,err:=client.AlibabaItemGet(id);err==nil{

			data:=struct{
				Item alibaba.Item `json:"item"`
			}{}

			save_path:="D:/items/"+time.Now().Format("01.02")+"/"+time.Now().Format("15-0405")+"/"

			if err:=res.BindJSON(&data);err==nil{
				item:=data.Item

				if len(item.Images)>0{

					if false == util.IsFolder(save_path){
						os.MkdirAll(save_path,777)
					}


					for i,src:=range item.Images{
						if err:=images.Down(src,fmt.Sprint(save_path,i,filepath.Ext(src)));err!=nil{
							fmt.Println(err)
						}

					}
				}

				if len(item.SukImages)>0{
					sku_path:=save_path+"sku图片/"
					if false == util.IsFolder(sku_path){
						os.MkdirAll(sku_path,0)
					}
					for _,sku:=range item.SukImages{
						if err:=images.Down(sku.Original,fmt.Sprint(sku_path,filepath.Ext(sku.Name)));err!=nil{
							fmt.Println(err)
						}
					}
				}

				client.AlibabaItemSet(id,"savepath",save_path)
				c.JSON(200,gin.H{"isSucc":true,"data":save_path})
				return
			}
		}
		c.JSON(200,gin.H{"isSucc":false})

	})
	r.Run(":8000")
}