package mongodb

import (
	mdb "helper/dbs/mongodb"
	"labix.org/v2/mgo/bson"
	"helper/configs"
)
var (
	link *mdb.Mogno
)

func conn()*mdb.Mogno{
	if link == nil{
		link=mdb.Conn()
	}
	return link
}

func Cols()[]string{
	db:=conn()
	result:=struct{
		Retval []string `bson:"retval"`
	}{}
	db.Database.Run(bson.D{{"eval","db.getCollectionNames()"}},&result)
	return result.Retval
}

func CommendRun(cmd string)configs.M{
	db:=conn()
	result:=configs.M{}
	db.Database.Run(bson.D{{"eval",cmd}},&result)
	return result
}

func Listing(name string ,find configs.M,offset ,limit int ,sort ... string  )( result []configs.M ,total int){
	db:=conn()
	db.C(name).Find(find).Skip(offset).Limit(limit).Sort(sort...).All(&result)
	total,_=db.C(name).Count()
	return
}