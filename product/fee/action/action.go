package action
import(
	"github.com/gin-gonic/gin"
	"helper/net/igin"
	"product/fee"
)

//费用列表
func Listing(c *gin.Context){
	ig:=igin.H(c)
	platform,ok:=c.GetQuery("platform")
	if ok == false{
		platform="common"
	}
	sku:=c.Param("sku")
	feesOpt,rate_count:=fees.Listing(sku,platform)
	ig.Succ(gin.H{"rate_count":rate_count,"fees":feesOpt})
}

//编辑费用
func Edit(c *gin.Context){
	
}


//添加费用
func Add(c *gin.Context){
	
}