package admin
import (
	"github.com/gin-gonic/gin"
	"helper/auth/other/qq"
	"fmt"
	"helper/redisCli"
	"helper/account"
	"helper/auth"
	"net/http"
	"log"
)


func QQcallBack(c * gin.Context){
	var err error
	r:=redisCli.Conn()

	if code,ok:=c.GetQuery("code");ok{

		var token auth.Token
		if token_value,ok:=c.GetQuery("state");ok{
			token=auth.Token{Value:token_value}
		}else{
			c.String(201,"认证失败")
			return
		}

		if token.IsExist() == false {
			c.String(202,"认证失败")
			return
		}

		if qq_token,err:=qq.GetAccessToken(code);err==nil{
			fmt.Println("Access Token:",qq_token.Value)

			if openId,err:=qq.GetOpenId(qq_token.Value);err==nil{
				fmt.Println("OPEN ID:",openId)

				r.Send("HMSET",fmt.Sprint("qq_openid:",openId),"access_token",qq_token.Value,"expires_in",qq_token.ExpiresIn)

				if qq_user_info,err:=qq.UserInfo(qq_token.Value,openId);err==nil{
					token.Set("qq_access_token",qq_token.Value,"qq_open_id",openId,"qq_nickname",qq_user_info.NickName,"qq_figureurl",qq_user_info.FigureURLQQ2,"qq_gender",qq_user_info.Gender)


					user:=fmt.Sprintf("%s@qqauth.com",openId)

					if account.IsExist(user) == false {
						account.New(user,openId,qq_user_info.NickName)
					}

					if user_info,err:=account.Verify(user,openId);err==nil{
						token.BindUser(user_info.Uid)

						c.Redirect(http.StatusMovedPermanently,"http://51helper.com/welcome")
						c.Abort()
						return
					}
				}else{
					log.Print(err)
				}

			}
		}

	}

	fmt.Println(err)

}

