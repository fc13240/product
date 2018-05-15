package client

import (
	"helper/webtest"
	"fmt"
	"net/http"
	"encoding/json"
	"helper/configs"
	"product/ezbuy"
	"log"
)
var ApiAddr string

type HClient struct{
	Token string
	Url string
}

func (cli *HClient)UpItemField(id int,field,value string)(*webtest.Result,error){
	body:=`{"id":%d,"field":"%s","value":"%s"}`
	body=fmt.Sprintf(body,id,field,value)
	return webtest.PostJson(fmt.Sprint(cli.Url,"ezbuy.setitem"),cli.header(),body)
}

func (cli *HClient)GetItems(offset,limit int)(*webtest.Result,error){
	body:=`{"offset":%d,"limit":%d,"orderBy":"update","issale":true}`
	return webtest.PostJson(fmt.Sprint(cli.Url,"ezbuy.items"),cli.header(),fmt.Sprintf(body,offset,limit))
}

func (cli *HClient)GetItem(sku string)(item *ezbuy.Item,err error){
	item =&ezbuy.Item{}
	res,err:= webtest.HGet(fmt.Sprint(cli.Url,"ezbuy.get/"+sku),cli.header())
	if err==nil{
		ss:=struct{
			Item *ezbuy.Item
		}{}
		if err:=res.BindJSON(&ss);err==nil{
			item=ss.Item
		}
	}
	return 
}

func (cli *HClient)GetHotItems(offset,limit int)(*webtest.Result,error){
body:=`{"offset":%d,"limit":%d,"orderBy":"update","ishost":true}`
return webtest.PostJson(fmt.Sprint(cli.Url,"ezbuy.items"),cli.header(),fmt.Sprintf(body,offset,limit))
}

func (cli *HClient)SaveItems(body string)(*webtest.Result,error){
	return webtest.PostJson(fmt.Sprint(cli.Url,"ezbuy.saveitems"),cli.header(),body)
}

func (cli *HClient)SaveOrders(body string)(*webtest.Result,error){
	return webtest.PostJson(fmt.Sprint(cli.Url,"ezbuy.saveorders"),cli.header(),body)
}

func (cli *HClient)CheckNewOrders(body string)(*webtest.Result,error){
	return webtest.PostJson(fmt.Sprint(cli.Url,"ezbuy.checkneworders"),cli.header(),body)
}

func (cli *HClient)header()http.Header{
	head:= http.Header{}
	head.Set("Authorization",cli.Token)
	return head
}

func (cli *HClient)GetStting()(*webtest.Result,error){
	return webtest.HGet(fmt.Sprint(cli.Url,"ezbuy.getsetting"),cli.header())
}

func (cli *HClient)GetEzStore(store_id int)(store *ezbuy.Setting){
	res,err:=webtest.HGet(fmt.Sprint(cli.Url,"ezbuy/store/",store_id),cli.header())

	if err!=nil{
		log.Panicln(err,"请求店铺信息失败")
	}
	ss:=struct{
		isSucc bool `json:"isSucc"`
		Item ezbuy.Setting `json:"item"`
	}{}

	if err:=res.BindJSON(&ss);err!=nil{
		log.Panicln(err,"解析店铺JSON失败")
	}
	
	return &ss.Item
}

func (cli *HClient)GetEzShopes()(*webtest.Result,error){
	return webtest.HGet(fmt.Sprint(cli.Url,"ezbuy/mystores"),cli.header())
}

func (cli *HClient)AlibabaItemGet(id int64)(*webtest.Result,error){
	return webtest.HGet(fmt.Sprint(cli.Url,"alibaba.getsource?id=",id),cli.header())
}

func (cli *HClient)AlibabaItemGetBySku(sku string)(*webtest.Result,error){
	return webtest.HGet(fmt.Sprint(cli.Url,"alibaba.getsource?sku=",sku),cli.header())
}

func (cli *HClient)AlibabaItemSet(id int64,filed,value string)(*webtest.Result,error){
	body:=`{"id":%d,"field":"%s","value":"%s"}`
	body=fmt.Sprintf(body,id,filed,value)
	return webtest.PostJson(fmt.Sprint(cli.Url,"alibaba.set"),cli.header(),body)
}

func (cli *HClient)AlibabaAdd(param configs.M) (data *webtest.Result,err error ){
	b,err:=json.Marshal(param)
	if err != nil{
		return nil,err
	}
	return webtest.PostJson(fmt.Sprint(cli.Url,"alibaba.add"),cli.header(),string(b))
}


func (cli *HClient)CheckItemExist(param configs.M) (data *webtest.Result,err error ){
	b,err:=json.Marshal(param)
	if err != nil{
		return nil,err
	}
	fmt.Println(fmt.Sprint(cli.Url,"alibaba/checkitemexist"))
	return webtest.PostJson(fmt.Sprint(cli.Url,"alibaba/checkUrlExist"),cli.header(),string(b))
}

func (cli *HClient)SaveUserProductsFromSource(data string)(*webtest.Result,error){
	return webtest.PostJson(fmt.Sprint(cli.Url,"ezbuy.saveUserProductsFromSource"),cli.header(),data)
}