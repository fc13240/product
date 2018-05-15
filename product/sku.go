package product

import (
	"helper/configs"
	"fmt"
	"helper/dbs"
	"helper/util"
	"helper/account"
	"time"
	"errors"
)

type Sku struct{
	Name string `json:"name"`
	ProductName string `json:"product_name"`
	Color string `json:"color"`
	Quant int `json:"quant"`
	Images []string `json:"images"`
}

//创建一个新的SKU
func NewSku()(sku *Sku,err error ){
	db:=dbs.Def()
	sku=&Sku{}
	last_id:=0

	db.One("SELECT MAX(id) FROM product_sku").Scan(&last_id)
	
	new_sku:=configs.Int(time.Now().Format("060102"))*100
	
	if new_sku>last_id{
		last_id=new_sku
	}else{
		last_id=last_id+1
	}

	sku.Name=fmt.Sprintf("BL%d",last_id)
	set:=configs.M{
		"create_date":util.Datetime(),
		"name":sku.Name,
		"id":last_id,
	}
	
	sku_isExist:=0

	db.One("select COUNT(*) FROM product_sku WHERE name=?",sku.Name).Scan(&sku_isExist)
	if sku_isExist>0{
		return sku,errors.New(" sku is exist:"+sku.Name)
	}
	if id,err:=db.Insert("product_sku",set);err!=nil || id==0{
		errors.New("create sku error "+err.Error())
	}
	return sku,nil
}

func (sku *Sku)AddBuySite(site string)error{
	 return AddSkuBuySite(sku.Name,site)
}

// 增加SKU采购地址
func AddSkuBuySite(sku,site string)error {
	db:=dbs.Def()
	buyurlMd5:=util.Md5(site)
	var id int 
	db.One("SELECT id FROM product_sku_buysite WHERE sku=? AND buyurl_md5=?",sku,buyurlMd5).Scan(&id)
	
	if id > 0{
		return errors.New("这个购买地址已经存在")
	}

	_,err:=db.Insert(
		"product_sku_buysite",
		configs.M{
			"sku":sku,
			"buyurl":site,
			"adddate":util.Datetime(),
			"sort":100,
			"buyurl_md5":buyurlMd5,
		})

	return err
		
}

//添加SKU使用日志
func AddSkuUploadLog(author *account.Account, sku string,store_id int){
	db:=dbs.Def()
	db.Insert("product_sku_uploadlog",configs.M{"sku":sku,"author_id":author.Uid, "store_id":store_id,"create_date":util.Datetime()})
}

type SkuUseStatus struct{
	StoreId int `json:"store_id"`
	StoreName string `json:"store_name"`
	UseDate string `json:"use_date"`
	UseNum int `json:"usenum"`
}

//获取使用的情况
func GetSkuUseStatus(author *account.Account,sku string) (list []SkuUseStatus){
	list=[]SkuUseStatus{}
	db:=dbs.Def()
	sql:="SELECT product_sku_uploadlog.store_id,product_store.store_name,product_sku_uploadlog.create_date,COUNT(product_sku_uploadlog.store_id) AS use_num  FROM product_sku_uploadlog LEFT JOIN product_store ON(product_store.store_id=product_sku_uploadlog.store_id) WHERE product_sku_uploadlog.sku=? GROUP BY product_sku_uploadlog.store_id"
	rows:=db.Rows(sql,sku)
	defer rows.Close()
	for rows.Next(){
		useStatus:=SkuUseStatus{}
		rows.Scan(&useStatus.StoreId,&useStatus.StoreName,&useStatus.UseDate,&useStatus.UseNum)
		list=append(list,useStatus)
	}
	return list
}