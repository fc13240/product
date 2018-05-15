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

type SPU struct{
	Name string `json:"name"`
}

//创建一个新的SPU
func NewSpu()(spu *SPU,err error ){
	db:=dbs.Def()
	spu=&SPU{}
	last_id:=0
	
	db.One("SELECT MAX(id) FROM product_spu").Scan(&last_id)
	
	//创建一个当天的年月日为前缀的SPU
	new_spu:=configs.Int(time.Now().Format("060102"))*100
	
	if new_spu>last_id{
		last_id=new_spu
	}else{
		last_id=last_id+1
	}

	spu.Name=fmt.Sprintf("BL%d",last_id)
	set:=configs.M{
		"create_date":util.Datetime(),
		"name":spu.Name,
		"id":last_id,
	}

	spu_isExist:=0

	db.One("select COUNT(*) FROM product_spu WHERE name=?",spu.Name).Scan(&spu_isExist)
	if spu_isExist>0{
		return spu,errors.New("SPU 已经存在"+spu.Name)
	}
	if id,err:=db.Insert("product_spu",set);err!=nil{
		err=errors.New("create spu error "+err.Error())
		return spu,err
	}else{
		fmt.Println(id)
	}
	return spu,nil
}

func (spu *SPU)AddBuySite(site string)error{
	 return AddSPUBuySite(spu.Name,site)
}

// 增加SPU采购地址
func AddSPUBuySite(spu,site string)error{
	db:=dbs.Def()
	buyurlMd5:=util.Md5(site)
	var id int 
	db.One("SELECT id FROM product_spu_buysite WHERE spu=? AND buyurl_md5=?",spu,buyurlMd5).Scan(&id)
	
	if id > 0{
		return errors.New("这个购买地址已经存在")
	}

	_,err:=db.Insert(
		"product_spu_buysite",
		configs.M{
			"spu":spu,
			"buyurl":site,
			"adddate":util.Datetime(),
			"sort":100,
			"buyurl_md5":buyurlMd5,
		})

	return err
		
}

//添加SPU使用日志
func AddSPUUploadLog(author *account.Account, spu string,store_id int){
	db:=dbs.Def()
	db.Insert("product_spu_uploadlog",configs.M{"spu":spu,"author_id":author.Uid, "store_id":store_id,"create_date":util.Datetime()})
}

type SPUUseStatus struct{
	StoreId int `json:"store_id"`
	StoreName string `json:"store_name"`
	UseDate string `json:"use_date"`
	UseNum int `json:"usenum"`
}

//获取SPU使用的情况,上传的各个店铺的情况
func GetSpuUseStatus(author *account.Account,spu string) (list []SPUUseStatus){
	list=[]SPUUseStatus{}
	db:=dbs.Def()
	sql:="SELECT product_spu_uploadlog.store_id,product_store.store_name,product_spu_uploadlog.create_date,COUNT(product_spu_uploadlog.store_id) AS use_num  FROM product_spu_uploadlog LEFT JOIN product_store ON(product_store.store_id=product_spu_uploadlog.store_id) WHERE product_spu_uploadlog.spu=? GROUP BY product_spu_uploadlog.store_id"
	rows:=db.Rows(sql,spu)
	defer rows.Close()
	for rows.Next(){
		useStatus:=SPUUseStatus{}
		rows.Scan(&useStatus.StoreId,&useStatus.StoreName,&useStatus.UseDate,&useStatus.UseNum)
		list=append(list,useStatus)
	}
	return list
}