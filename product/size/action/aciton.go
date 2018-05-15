package action
import(
	"github.com/gin-gonic/gin"
	"helper/net/igin"
	"product/size"
	"helper/configs"
)

func SaveSizeTemplate(c *gin.Context){
	ig:=igin.H(c)
	param:=size.Size{}
	var err error
	
	if err:=c.BindJSON(&param);err==nil{
		if id,err:=size.Save(param);err==nil{
			ig.Succ(gin.H{"data":id})
		}
	}

	if err!=nil{
		ig.Fail(err.Error)
	}
}

func GetSizeTemplate(c *gin.Context){
	item_id:=configs.Int(c.Param("item_id"))
	ig:=igin.H(c)
	if data,err:=size.Get(item_id);err==nil{
		ig.Succ(gin.H{"data":data})
	}else{
		ig.Fail(err.Error())
	}


}