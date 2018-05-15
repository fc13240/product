package extends

import (
	"helper/net/igin"
	"github.com/gin-gonic/gin"
)

var FrontRoute *gin.RouterGroup
var ToolGroup *gin.RouterGroup
func init(){
	FrontRoute=igin.R.Group("/api")
	ToolGroup=FrontRoute.Group("/tool")
}