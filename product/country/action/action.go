package action
import(
	"github.com/gin-gonic/gin"
	"helper/net/igin"
	"product/country"
)

//货币列表
func GetCurrencys(c *gin.Context){
	ig:=igin.H(c)
	ig.Succ(gin.H{"data":country.GetCurrencys()})
}
