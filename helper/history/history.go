package history

import (
	"helper/dbs/mongodb"
	"helper/configs"
	"helper/util"
)

type Flag string

func Col() *mongodb.Collection{
	return mongodb.Conn().C("history")
}

func (flag * Flag)Add(authorid int ,source_url,source_title string,labels ...string) error{
	md5:=util.Md5(source_url)
	col:=Col()
	count,_:=col.Find(configs.M{"md5":md5,"authorid":authorid}).Count()

	datetime:=util.Datetime()
	var err error
	if count > 0{
		err=col.Update(configs.M{"md5":md5,"authorid":authorid},
			configs.M{
				"$inc":configs.M{"visit_number":1},
				"$set":configs.M{
				"last_visit":datetime,
				"source_title":source_title,
				},
			})
	}else{
		data:=configs.M{
			"authorid":authorid,
			"source_url":source_url,
			"flag":flag,
			"addtime":datetime,
			"last_visit":datetime,
			"md5":md5,
			"visit_number":1,
			"source_title":source_title,
		}
		err=col.Insert(data)
	}
	return err
}

func (flag * Flag)Listing(filter configs.M,offset,limit int ,sort ...string )(rows []configs.M,count int ){
	col:=Col()
	rows=[]configs.M{}
	col.Find(filter).Skip(offset).Limit(limit).Sort(sort...).All(&rows)
	count,_=col.Find(filter).Count()

	return
}