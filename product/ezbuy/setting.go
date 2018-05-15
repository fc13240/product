package ezbuy

import (
	"fmt"
	"helper/configs"
	"helper/dbs/mongodb"
	"helper/account"
	"helper/redisCli"
	"github.com/garyburd/redigo/redis"
)

var settingCol *mongodb.Collection

func SettingCol() *mongodb.Collection{
	if settingCol==nil{
		settingCol=mongodb.Conn().C("ezbuy.setting")
	}
	return settingCol
}

type Setting struct{
	Cookie string `json:"cookie"`
	Reqid string `json:"reqid"`
	Num int `json:"num"`
	Minute int `json:"minute"`
	AuthorId int `json:"authorid" bson:"authorid"`
	StoreId int `json:"storeid" bson:"storeid"`
	StoreName string `json:"storename" bson:"storename"`
	//ItemsNum int `json:"itemnum"`
	CheckNewOrders bool `json:"checkneworders"`
	OnRefresh  bool `json:"onrefresh"`
	SecrectKey string `json:"secrectkey"`
	SkuFirst   string `json:"skufirst"`
	StoreCateId int `json:"store_cateid" bson:"store_cateid"`
	SellerSkuId string `json:"sellerskuid"`
}

//生成SKU
func (setting *Setting)GenSku()string{
	r:=redisCli.Conn()
	skuincr,_:=redis.Int(r.Do("incr","skuincr:"))
	return fmt.Sprint(setting.SkuFirst,skuincr)
}

//保存设置
func SaveSetting(v Setting,author *account.Account){
	v.AuthorId=author.Uid
	col:=SettingCol()
	if n,_:=col.Find(configs.M{"authorid":author.Uid,"storeid":v.StoreId}).Count();n>0{
		col.Upsert(configs.M{"authorid":author.Uid,"storeid":v.StoreId},v)
	}else{
		v.AuthorId=author.Uid
		col.Insert(v)
	}
}

//获取ez配置
func GetSetting(author *account.Account,StoreId int ) *Setting{
	setting:=Setting{}
	SettingCol().Find(configs.M{"authorid":author.Uid,"storeid":StoreId}).One(&setting)
	return &setting
}

//获取EZ列表
func MyStores(author *account.Account)(rows []Setting){
	SettingCol().Find(configs.M{"authorid":author.Uid}).All(&rows)
	return rows
}