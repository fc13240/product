
package action
import(
	"github.com/gin-gonic/gin"
	"helper/net/igin"
	"product"
)

func BindEbay(c *gin.Context) {
	param:=struct{
		Itemid int `json:"item_id"`
		EbayItemid int `json:"ebay_itemid"`
	}{}

	c.BindJSON(&param)
	if item, err := product.IdGet(param.Itemid); err == nil {
		item.BindEbay(param.EbayItemid)
	} else {
		igin.Fail(c,err.Error())
		return
	}
	igin.Succ(c,nil)

}