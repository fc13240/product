package doc

import (
	"helper/dbs/mongodb"
	"time"

	"helper/label"
	"helper/redisCli"
	"github.com/garyburd/redigo/redis"
	"qiniupkg.com/x/errors.v7"
	"helper/configs"
	"fmt"
	"helper/history"
	"member/favorite"
)

type Doc struct{
	Id int		`bson:"id" json:"id"`
	AuthorId   int  	 `bson:"author_id" json:"author_id"`
	Code string 		 `bson:"code",json:"code"`
	Title string		 `bson:"title" json:"title"`
	Content string	     `bson:"content" json:"content"`
	Labels []label.Label `bson:"labels" json:"labels"`
	Addtime time.Time 	 `bson:"addtime" json:"addtime"`
	ShortDesc string 	`bson:"short_desc" json:"short_desc"`
	Uptime time.Time `bson:"uptime",json:"uptime"`
	Readnum int `bson:"readnum" json:"readnum"`
	SourceUrl string `bson:"source_url" json:"source_url"`
	SourceName string `bson:"source_name" json:"source_name"`
	ContentType string `bson:"content_type" json:"content_type"`
	Ismarkdown bool `bson:"ismarkdown" json:"ismarkdown"`
	MarkdownContent string `bson:"markdown_content" json:"markdown_content"`
	IsMe  bool `json:"isme"`
}

var History history.Flag = "article"

const (
	FavoriteFlag ="article"
)

func Get(id int)*Doc{
	cc:=mongodb.Conn()
	info:=Doc{}
	col:=cc.C("doc")

	col.Find(configs.M{"id":id}).One(&info)

	if info.SourceName == ""{
		info.SourceName="未知"
	}

	if info.ContentType==""{
		info.ContentType="article"
	}
	return &info
}

func(doc *Doc)ReadnumIncr(){
	cc:=mongodb.Conn().C("doc")
	cc.Upsert(configs.M{"id":doc.Id},configs.M{"$inc":configs.M{"readnum":1}})
}

func NewId()(doc_id int){
	r:=redisCli.Conn()
	doc_id,_=redis.Int(r.Do("incr","doc_id"))
	return doc_id
}

func (doc *Doc)Save()error{
	db:=mongodb.Conn()
	doc.Addtime=time.Now()

	return db.C("doc").Insert(&doc)
}

func (doc *Doc)Update() (err error){
	db:=mongodb.Conn()
	err=db.C("doc").Update(configs.M{"id":doc.Id},
		configs.M{"$set":configs.M{
			"title":doc.Title,
			"short_desc":doc.ShortDesc,
			"source_url":doc.SourceUrl,
			"content":doc.Content,
			"labels":doc.Labels,
			"uptime":time.Now(),
			"content_type":doc.ContentType,
			"ismarkdown":doc.Ismarkdown,
			"markdown_content":doc.MarkdownContent,
		}})
	return err
}

func New(author_id int,title,short_desc,content string, labels []label.Label)(dd *Doc,err error){
	redisCli.Conn()
	dd=&Doc{
		AuthorId:author_id,
		Id:NewId(),
		Title:title,
		ShortDesc:short_desc,
		Content:content,
		Labels:labels,
		Addtime:time.Now(),
	}
	if err:=dd.Save();err!=nil{
		err=errors.New(fmt.Sprint("创建文档失败",err.Error()))
	}
	return dd,err
}

func Listing(m configs.M,skip,limit int,sort ...string)(items []Doc,total int){
	items =[]Doc{}
	c:=mongodb.Conn().C("doc")
	c.Find(m).Select(configs.M{"title":true,"id":true,"labels":true,"addtime":true,"short_desc":true}).Skip(skip).Limit(limit).Sort(sort...).All(&items)
	total,_=c.Find(m).Count()
	return items,total
}

func CancelFav(uid,id int)error{
	return favorite.New(uid,FavoriteFlag).Del(id)
}

//文章收藏记录
func FavListing(m configs.M,skip,limit int,sort ...string)(items []Doc,total int){
	items =[]Doc{}
	autorh_id:=m.Int("author_id")

	fav:=favorite.New(autorh_id,FavoriteFlag)
	ids:=fav.Limit(skip,limit)

	if len(ids) == 0{ //如果收藏数量=0
		return
	}

	total=fav.Count()

	c:=mongodb.Conn().C("doc")
	filter:=configs.M{
		"id":configs.M{"$in":ids},
	}

	c.Find(filter).Select(configs.M{"title":true,"id":true,"labels":true,"addtime":true,"short_desc":true}).Skip(skip).Limit(limit).Sort(sort...).All(&items)
	newSortItems:=[]Doc{}
	if len(items)>0{
		for _,id:=range ids{
			for _,item:= range items{
				if item.Id == id{
					newSortItems=append(newSortItems,item)
				}
			}
		}
	}
	return newSortItems,total
}

func Rows(m configs.M,limit int,sort ...string)(items []Doc){
	items =[]Doc{}
	c:=mongodb.Conn().C("doc")
	c.Find(m).Select(configs.M{"title":true,"id":true,"labels":true,"addtime":true,"short_desc":true}).Limit(limit).Sort(sort...).All(&items)
	return items
}

func Rowss(m configs.M,limit int,sort ...string)(items []Doc){
	items =[]Doc{}
	c:=mongodb.Conn().C("doc")
	c.Find(m).Select(configs.M{"title":true,"id":true,"labels":true,"addtime":true,"short_desc":true,"content":true}).Limit(limit).Sort(sort...).All(&items)

	return items
}

func Update(id int,title,short_desc,content string, labels []label.Label)(dd *Doc,err error){
	dd=&Doc{Id:id,Title:title,ShortDesc:short_desc, Content:content,Labels:labels}
	if err:=dd.Update();err!=nil{
		err=errors.New(fmt.Sprint("保存文档失败",err.Error()))
	}
	return dd,err
}