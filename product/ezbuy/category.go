package ezbuy

import (
	"helper/dbs/mongodb"
	"helper/account"
	"helper/configs"
)

var myCateCol *mongodb.Collection

func MyCateCol() *mongodb.Collection{
	if myCateCol==nil{
		myCateCol=mongodb.Conn().C("ezbuy.mycategorys")
	}
	return myCateCol
}

var cateCol *mongodb.Collection

func CateCol() *mongodb.Collection{
	if cateCol==nil{
		cateCol=mongodb.Conn().C("ezbuy.categorys")
	}
	return cateCol
}

func GetCategorys(parent_id int)(categorys []Category){
	CateCol().Find(configs.M{"pid":parent_id}).All(&categorys)
	return categorys
}

func AddToMyCategory(authorid ,id int ,name string)error{
	_,err:= MyCateCol().Upsert(configs.M{"authorid":authorid,"id":id},configs.M{"id":id,"name":name,"authorid":authorid})
	return err
}

type Category struct{
   Id   int	`json:"id" bson:"id"`
   Name string `json:"name"`
   Pid int `json:"pid"`
}

func GetMyCategorys(author *account.Account)(categorys []Category){
	MyCateCol().Find(configs.M{"authorid":author.Uid}).All(&categorys)
	return categorys
}