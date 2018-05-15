package doc

import (
	"helper/configs"
	"helper/dbs/mongodb"
	"fmt"
	"strings"
)
type Node struct{
  Index   string
  Content string
  CategoryId  int
  Sort    	int
}

type Category struct{
	Id int
	Index string
	Name string
	Pid int
	Sort int
	Nodes []*Node

}

type Manual struct{
	Id int
	Tag string
	Content []Node
	Categorys []Category
}

type Index struct {
	Index string `json:"index"`
	Path  string `json:"path"`
	Sort  int `json:"sort"`
	By string `json:"by"`
}

func CreateManual(tag string, categorys ...Category){
	mdb:=mongodb.Conn()
	manual:=mdb.C("manual")
	for _,category:=range categorys{
		_,err:=manual.Upsert(configs.M{"tag":tag},configs.M{"$set":configs.M{"tag":tag,category.Index:category}})
		fmt.Println(err)
	}
}

func GetManualIndex(tag string)([]Index){
	mdb:=mongodb.Conn()
	c:=struct{Indexs []Index}{}
	mdb.C("manual").Find(configs.M{"tag":tag}).One(&c)
	return c.Indexs
}

func AddNode(tag string, node *Node){
	mdb:=mongodb.Conn()
	mdb.C("manual").Update(configs.M{"tag":"redis"},configs.M{"$addToSet":configs.M{"nodes":node}})
}

func SetNodeContent(tag ,index, content string){
	mdb:=mongodb.Conn()
	mdb.C("manual").Update(configs.M{"tag":"redis","nodes.index":index},configs.M{"$set":configs.M{"nodes.$.content":content} })
}

func GetManualNode(tag ,index string) *Node{

	index=strings.ToUpper(strings.Replace(index,"_"," ",1))
	mdb:=mongodb.Conn()
	ss:=[]configs.M{
		configs.M{
			"$match": configs.M{"tag": "redis"},
		},configs.M{
			"$project":
			 configs.M{
				"nodes":
				configs.M{"$filter":
					configs.M{
						"input":"$nodes",
						"as":"node",
						"cond":configs.M{"$eq":[]string{"$$node.index",index}},
					}},
				 "_id": 0,
			 },
		},
	}
	v:= struct{
		Nodes []Node `bson:"nodes"`
	}{}

	mdb.C("manual").Pipe(ss).One(&v)
	if len(v.Nodes)>0{
		return &v.Nodes[0]
	}
	return &Node{}

}

func SetManualIndex(tag string ,indexs []*Index){
	mdb:=mongodb.Conn()
	manual:=mdb.C("manual")
	manual.Upsert(configs.M{"tag":tag},configs.M{"$set":configs.M{"tag":tag,"indexs":indexs}})
}