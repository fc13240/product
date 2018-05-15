package passport
import(
	"helper/auth"
	"regexp"
	"github.com/gin-gonic/gin"
	"log"
)
func Auth() gin.HandlerFunc{
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Println("异常错误", r)
			}
		}()
		ok,_:=regexp.MatchString("/api/",c.Request.URL.Path)
		if ok{
			token:=auth.Token(c.Request.Header.Get("Authorization"))
			
		//	customer:=auth.Customer{}
		//	customer.VerifyToken(token)

			if token=="" || token.IsLogin()== false {
				switch c.Request.URL.Path {
				case "/api/account/create":
				case "/api/account/login":
				case "/api/account/sendregistercode":
				case "/api/item.margecontent":
				default:
					c.String(401,"not login")
					c.Abort()
				}
			}
			c.Set("token",token)
		}else{
			c.String(400,"error")
			c.Abort()
		}
	}
}