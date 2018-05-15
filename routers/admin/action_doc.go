package admin

import (
	"helper/label"
	"helper/doc"
	"helper/configs"
	"helper/net/igin"
	"github.com/gin-gonic/gin"
	"strings"

)


func Listing(c *gin.Context){
	hc:=igin.H(c)
	author:=hc.Account()

	param:=configs.M{}

	if false == author.ISAdmin() {
		param["author_id"]=author.Uid
	}

	data:=struct{
		Offset int `json:"offset"`
		Limit int `json:"limit"`
		Label  string `json:"label"`
		AuthorId int
	}{}

	if err:=c.BindJSON(&data);err!=nil{
		igin.Fail(c,err.Error())
		return
	}

	if data.Label!=""{
		labels:=strings.Split(data.Label," ")
		if len(labels)>0{
			param["labels.name"]=configs.M{"$in":labels}
		}
	}

	list ,total:= doc.Listing(param,data.Offset,data.Limit,"-id")
	hc.Succ(gin.H{"items":list,"total":total})
}

func MyList(c *gin.Context){
	ig:=igin.H(c)
	token,err:=ig.Token()
	if err!=nil{
		ig.Fail(err.Error())
		return
	}
	if  token.IsLogin() == false{
		ig.Fail("没有登录")
		return
	}

	author:=ig.Account()

	param:=configs.M{}

	param["author_id"]=author.Uid

	data:=struct{
		Offset int `json:"offset"`
		Limit int `json:"limit"`
		Label  string `json:"label"`
	}{}

	if err:=c.BindJSON(&data);err!=nil{
		igin.Fail(c,err.Error())
		return
	}

	if data.Label!=""{
		labels:=strings.Split(data.Label," ")
		if len(labels)>0{
			param["labels.name"]=configs.M{"$in":labels}
		}
	}

	list ,total:= doc.Listing(param,data.Offset,data.Limit,"-id")
	for i,_:=range list{
		list[i].IsMe=true
	}
	ig.Succ(gin.H{"items":list,"total":total})
}

func Get(c *gin.Context){
	id:=configs.Int(c.Param("id"))
	doc:= doc.Get(id)
	hc:=igin.H(c)
	hc.Succ(gin.H{"item":doc})
}

func addLabel(c *gin.Context){

	param:= struct {
		Name string `json:"name"`
	}{}

	if err:=c.BindJSON(&param);err!=nil{
		igin.Fail(c,err.Error())
		return
	}

	labe,err:=label.New(param.Name)
	if err!=nil{
		igin.Fail(c,err.Error())
		return
	}

	igin.Succ(c,gin.H{"item":labe})
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



func Create(c *gin.Context){

	param:= struct {
		Id int `json:"id"`
		Title string `json:"title"`
		Code string 	`json:"code"`
		Content string `json:"content"`
		ContentType string `json:"content_type"`
		SocureUrl string `json:"socure_url"`
		ShortDesc string `json:"short_desc"`
		Labels []label.Label `json:"labels"`
		Ismarkdown bool `json:"ismarkdown"`
		MarkdownContent string `json:"markdown_content"`
	}{}

	if err:=c.BindJSON(&param);err!=nil{
		igin.Fail(c,err.Error())
		return ;
	}

	var err error
	hc:=igin.H(c)

	info:=&doc.Doc{
		Id:param.Id,
		Title:param.Title,
		ShortDesc:param.ShortDesc,
		Content:param.Content,
		Labels:param.Labels,
		SourceUrl:param.SocureUrl,
		ContentType:param.ContentType,
		Code:param.Code,
		Ismarkdown:param.Ismarkdown,
		MarkdownContent:param.MarkdownContent,
	}

	if param.Id>0{
		err=info.Update()
	}else{
		author:=hc.Account()
		info.AuthorId=author.Uid
		info.Id=doc.NewId()
		err=info.Save()
	}

	if err==nil{
		igin.Succ(c,gin.H{"id":info.Id})
	}else{
		igin.Fail(c,err.Error())
	}
}