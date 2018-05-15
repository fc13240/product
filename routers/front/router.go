package front
import (
	"main/front/extends"
	_ "tools/ip"
	_ "tools/conv"
	_ "tools/time"
	_"tools/encoding"
	_ "tools/crypto"
	_ "member/api"
)

func init(){
	frontRoute:=extends.FrontRoute
	var docGroup=frontRoute.Group("/doc")
	docGroup.POST("/label.search",SearchLabel)
	docGroup.POST("/listing", Listing)
	docGroup.GET("/get/:id", Get)
	docGroup.POST("/historys", Historys)
	docGroup.POST("/favorites.add", AddFav)
	docGroup.POST("/favorites.cancel", CancelFav)

	var toolGroup=extends.ToolGroup
	toolGroup.POST("get.alibaba.product.images",toolAction.GetAlbabaProductImages)
	toolGroup.POST("md5",toolAction.Md5)
	toolGroup.POST("json.format",toolAction.JsonFormat)

	frontRoute.GET( "/get/:code",toolAction.Get)
	frontRoute.POST("/req/:id",toolAction.Req)
	frontRoute.POST("/apps",toolAction.Apps)

}