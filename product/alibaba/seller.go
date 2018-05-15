package alibaba

import (
	"helper/dbs"
	"helper/util"
	"helper/configs"
	"log"
	"fmt"
)

type Seller struct{
	Id int `json:"id"`
	CompanyName string `json:"company_name"`
	CompanyUrl string `json:"company_url"`
	Flag string `json:"flag"`
	AuthorId int `json:"author_id"`
	Sku   string `json:"sku"`
	ItemNum int `json:"itemnum"`
}

func (s *Seller)Save()(seller_id int ){
	db:=dbs.Def()
	md5:=util.Md5(s.CompanyUrl)

	db.One("SELECT seller_id FROM alibaba_seller WHERE md5=?",md5).Scan(&seller_id)

	if seller_id>0{
		return seller_id
	}else{
		var err error
		if s.Id,err=db.Insert("alibaba_seller",configs.M{"company_name":s.CompanyName,"company_url":s.CompanyUrl,"md5":md5,"author_id":s.AuthorId});err==nil{
			return s.Id
		}else{
			log.Println("店铺信息保存失败",s)
		}
	}
	return seller_id
}

func (s *Seller)ItemCount() (total int) {
	total,_=Col().Find(configs.M{"seller_id":s.Id}).Count()
	return total
}

func SellerListing(author_id int,offset,limit int)(sellers []Seller,total int){
	db:=dbs.Def()
	sellers=[]Seller{}
	
	sql:=fmt.Sprintf("SELECT seller_id,company_name,company_url,author_id,sku FROM alibaba_seller WHERE author_id=%d",author_id)

	if author_id == 0{
		sql=fmt.Sprintf("SELECT seller_id,company_name,company_url,author_id,sku FROM alibaba_seller")
		total,_=Col().Find(nil).Count()
	}else{
		total,_=Col().Find(configs.M{"seller_id":author_id}).Count()
	}
	

	rows:=db.Rows(sql+dbs.Limit(offset,limit))

	for rows.Next(){
		var id, author_id int
		var company_name, company_url ,sku string
		rows.Scan(&id,&company_name,&company_url,&author_id,&sku)
		seller:=Seller{Id:id,CompanyName:company_name,CompanyUrl:company_url,Sku:sku}
		seller.ItemNum=seller.ItemCount()
		sellers=append(sellers,seller)
	}
	//dbs.One("SELECT count(*) FROM alibaba_seller WHERE author_id=?",author_id).Scan(&total)
	return sellers,total
}

