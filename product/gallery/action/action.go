package action
import(
	"github.com/gin-gonic/gin"
	"helper/net/igin"
	"product/gallery"
	"fmt"
	"helper/configs"
)

//上传图片
func UploadImage(c *gin.Context) {
	ig:=igin.H(c)
	sku:=c.PostForm("sku")
	sort:=configs.Int(c.PostForm("sort"))
	flag:=configs.Int(c.PostForm("flag"))
	if sku == "" {
		ig.Fail("请先保存在产品")
		return
	}

	file,_,err:=c.Request.FormFile("file");
	if err==nil{
		if name,err :=gallery.UploadImage(file); err == nil {
			info,err:=gallery.Save(sku,name,flag,sort)
			if err !=nil{
				ig.Fail(fmt.Sprint("保存图片信息失败:",err.Error()))
				return
			}
			ig.Succ(gin.H{"src":info.Src,"sort":info.Sort,"flag":info.Flag})
		} else {
			ig.Fail(fmt.Sprint("上传图片文件到服务器失败:",err.Error()))
		}
	}else{
		ig.Fail(fmt.Sprint("上传文件失败:",err.Error()))
	}
}

//图片列表
func Images(c *gin.Context) {
	ig:=igin.H(c)
	param:=struct{
		igin.ParamFilter
		igin.ParamPage
	}{}
	if err:=c.BindJSON(&param);err!=nil{
		ig.Fail(err.Error())
		return 
	}
	
	rows, total := gallery.Listing(param.Filter,param.Offset,param.Limit)
	flags := gallery.Flags()
	ig.Succ(gin.H{"items": rows, "total": total, "flags": flags})
}

//添加图片
func AddImage(c *gin.Context){
	ig:=igin.H(c)
	param:=struct{
		Src string `json:"src"`
		Sku string `json:"sku"`
	}{}
	if err:=c.BindJSON(&param);err!=nil{
		ig.Fail(err.Error())
		return 
	}
	 gallery.AddImage(param.Sku,param.Src)
	ig.Succ(nil)
}

//删除图片
func DelImage(c *gin.Context){
	ig:=igin.H(c)
	param:=struct{
		Ids []int `json:"ids"`
		Sku string `json:"sku"`
	}{}
	if err:=c.BindJSON(&param);err!=nil{
		ig.Fail(err.Error())
		return 
	}
	if err:=gallery.DelSkuImage(param.Sku,param.Ids);err==nil{
		ig.Succ(nil)
	}else{
		ig.Fail(err.Error())
	}
}

//设置图片标记
func SetImageFlag(c *gin.Context) {

	param:=struct{
		Ids []int `json:"ids"`
		Flags []int  `json:"flags"`
	}{}

	if err:=c.BindJSON(&param);err!=nil{
		igin.Fail(c,err.Error())
		return 
	}

	err := gallery.SetFlag(param.Ids, param.Flags...)
	if err == nil {
		igin.Succ(c,nil)
	} else {
		igin.Fail(c,err.Error())
	}
}