package wss

import (
	"product/ezbuy"
	"fmt"
	"log"
	"helper/configs"
	ezclient "product/ezbuy/client"
	ezapi "product/ezbuy/api"
)

var ezShopes []ezbuy.Setting
func Auth(mod string) (err error) {
	option:=configs.GetSection("api")
	
	client=&ezclient.HClient{Token:option["secrect_key"],Url:option["api"]}
	ezclient.ApiAddr=option["api"]
	res,err:=client.GetEzShopes()
	if err!=nil{
		log.Println(err)
		return err
	}
	data:=struct {
		IsSucc bool `json:"isSucc"`
		Shopes []ezbuy.Setting `json:"items"`
	}{}

	if err:=res.BindJSON(&data);err!=nil{
		return err
	}

	if !data.IsSucc{
		fmt.Println("请求服务器验证失败")
	}
	ezShopes=data.Shopes
	defer  res.Close()
	setting= &data.Shopes[0]
	ez=&ezapi.EzbuyeApi{Cookie:setting.Cookie,ReqId:setting.Reqid,ShopName:setting.StoreName,Client:client,SkuFirst:setting.SkuFirst}
	go CheckNewOrders()
	go Refresh()
	PrintRunStat()
	return nil
}

//切换店铺
func ChangeShop(storeId int) *ezbuy.Setting{
	for i,shop:=range ezShopes{
		if shop.StoreId == storeId{
			setting=&ezShopes[i]
			ez=&ezapi.EzbuyeApi{Cookie:setting.Cookie,ReqId:setting.Reqid,ShopName:setting.StoreName,Client:client,SkuFirst:shop.SkuFirst}
		}
	}
	return setting
}