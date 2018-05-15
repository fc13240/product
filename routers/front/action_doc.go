package front

import (
	"helper/label"
	"helper/doc"
	"helper/configs"
	"helper/net/igin"
	"github.com/gin-gonic/gin"

	"strings"
	"member/favorite"

	"fmt"
)

func Listing(c *gin.Context){
	ig:=igin.H(c)
	param:=configs.M{}
	data:=struct{
		igin.ParamPage
		Label  string `json:"label"`
		igin.ParamFilter
		AuthorId int
	}{}

	if err:=c.BindJSON(&data);err!=nil{
		igin.Fail(c,err.Error())
		return
	}

	if data.Label!=""{
		data.Label=strings.Replace(data.Label,"  "," ",-1)
		data.Label=strings.Replace(data.Label,";"," ",-1)
		data.Label=strings.Replace(data.Label,","," ",-1)
		data.Label=strings.Trim(data.Label,"")
		labels:=strings.Split(data.Label," ")
		if len(labels)>0{
			param["labels.name"]=configs.M{"$in":labels}
		}
	}

	if data.Limit == 0{
		data.Limit=20
	}
	var list []doc.Doc
	var total int
	if search_type,ok:=c.GetQuery("search_type");ok {
		token,err:=ig.Token()

		if err!=nil{
			ig.Fail(err)
			return
		}

		if token.Uid == 0{
			ig.Fail("没有登录")
			return
		}

		switch search_type {
		case "favorites":
			param.Set("author_id",token.Uid)
			list ,total= doc.FavListing(param, data.Offset,data.Limit,"-id")
		}
	}else{
		list ,total= doc.Listing(param, data.Offset,data.Limit,"-id")
	}
	ig.Succ(gin.H{"items":list,"total":total})
}

func Get(c *gin.Context){
	id:=configs.Int(c.Param("id"))
	info:= doc.Get(id)

	ig:=igin.H(c)

	info.ReadnumIncr()

	data:=configs.M{
		"title":info.Title,
		"short_desc":info.ShortDesc,
		"addtime":info.Addtime.Format("2006/01/02 15:04"),
		"content":info.Content,
		"code":info.Code,
		"source_name":info.SourceName,
		"labels":info.Labels,
		"readnum":info.Readnum,
		"id":info.Id,
		"ismarkdown":info.Ismarkdown,
		"markdown_content":info.MarkdownContent,
	}

	if token,err:=ig.Token();err==nil && token.IsLogin(){
		doc.History.Add(token.Uid,fmt.Sprint("/article?id=%d",info.Id),info.Title)

		if token.Uid == info.AuthorId{
			info.IsMe=true
			data["isme"]=true
		}
	}

	ig.Succ(gin.H{"data":gin.H{"item":data}})
}

func Historys(c *gin.Context){

	gh:=igin.H(c)

	data:=struct{
		Offset int `json:"offset"`
		Limit int `json:"limit"`
	}{}

	if err:=c.BindJSON(&data);err!=nil{
		igin.Fail(c,err.Error())
		return
	}

	filter:=configs.M{"authorid":gh.Account().Uid}
	rows,count:=doc.History.Listing(filter,data.Offset,data.Limit,"-adddate")

	gh.Succ(gin.H{"items":rows,"count":count})
}

func AddFav(c *gin.Context){
	ig:=igin.H(c)
	token,_:=ig.Token()
	if token.Uid == 0{
		ig.Fail("请先登录")
		return
	}
	parame:=struct{
		Id int `json:"id"`
	}{}

	if err:=c.BindJSON(&parame);err==nil{
		fav:=favorite.New(token.Uid,"article")
		fav.Add(parame.Id)
		ig.Succ(nil)
	}else{
		ig.Fail(err.Error())
	}

}

func CancelFav(c *gin.Context){
	ig:=igin.H(c)
	token,_:=ig.Token()
	if token.Uid == 0{
		ig.Fail("请先登录")
		return
	}
	parame:=struct{
		Id int `json:"id"`
	}{}

	if err:=c.BindJSON(&parame);err==nil{
		if err=doc.CancelFav(token.Uid,parame.Id);err==nil{
			ig.Succ(nil)
		}else{
			ig.Fail(err.Error())
		}

	}else{
		ig.Fail(err.Error())
	}
}

func SearchLabel(c *gin.Context){
	gh:=igin.H(c)

	param:= struct {
		Name string `json:"name"`
	}{}

	if err:=c.BindJSON(&param);err!=nil{
		igin.Fail(c,err.Error())
		return
	}

	labels:=label.Search(param.Name,15)
	gh.Succ(gin.H{"items":labels})
}