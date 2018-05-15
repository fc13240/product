package action
import(
	"github.com/gin-gonic/gin"
	"helper/net/igin"
	"product/title"
	"helper/configs"
	"strings"
)
func AddTitleLabel(c *gin.Context){
	ig:=igin.H(c)
	param:=struct{
		Id int `json:"id"`
		Label string `json:"title"`
		Cnlabel string `json:"cntitle"`
		Cid int `json:"cid"`
	}{}
	
	if err:=c.BindJSON(&param);err!=nil{
		ig.Fail(err)
		return
	}

	if title.NewLabel(param.Label,param.Cnlabel,param.Cid,param.Id){
		ig.Succ(nil)
	}else{
		ig.Fail("添加失败")
	}
}

func AddTitleLabelCate(c *gin.Context){
	ig:=igin.H(c)

	param:=struct{
		Name string `json:"name"`
		ProductCid int `json:"product_cid"`
	}{}

	if err:=c.BindJSON(&param);err!=nil{
		ig.Fail(err)
		return
	}

	if id,err:= title.NewLabelCate(param.ProductCid,param.Name);err==nil{
		ig.Succ(gin.H{"data":id})
	}else{
		ig.Fail(err.Error())
	}
}

func GetTitleLabels(c *gin.Context){
	cid:=configs.Int(c.Param("cid"))
	labels:=title.GetLabels(cid)
	igin.H(c).Succ(gin.H{"data":labels})
}

func GetTitleLabelCateListing(c *gin.Context){
	filter:=configs.M{}
	labels:=title.GetLabelCateListing(filter)
	igin.H(c).Succ(gin.H{"data":labels})
}



//标题关键字搜索
func TitleKeywordSearch(c *gin.Context){
	q,_:=c.GetQuery("q")
	ig:=igin.H(c)
	if search_type,ok:=c.GetQuery("search_type");ok && search_type =="cn"{
		ig.Succ(gin.H{"data":title.SearchCnTitle(q)})
		return
	}else{
		ig.Succ(gin.H{"data":title.SearchTitle(q)})
		return
	}
}

func TitleKeywordImport(c *gin.Context){
	ig:=igin.H(c)
	parmas:=struct{
		Content string `json:"content"`
	}{}
	
	cid:=configs.Int(c.Param("cid"))
	if err:=c.BindJSON(&parmas);err==nil{
		ss:=strings.Split(parmas.Content,"\n\n")
		for _,s:=range ss{
			row:=strings.Split(s,"\n")
			title.NewLabel(row[0],row[1],cid,0)
		}
		ig.Succ(nil)
	}else{
		ig.Fail(err.Error())
	}
}