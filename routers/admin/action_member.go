package admin

import (
	"fmt"
	"helper/account"
	"helper/configs"
	"helper/images"

	//"helper/net/view"
	//"strconv"
	"github.com/gin-gonic/gin"
	"helper/net/igin"

	//"github.com/frp/src/models/msg"
	//"os/user"
	"helper/auth"
	"member/follow"
	"helper/label"
)

func MyServeList(c *gin.Context){
	gh:=igin.H(c)
	items:=account.MyServeList(gh.Account().Uid)
	gh.Succ(gin.H{"items":items})
}

func AddFollowLabel(c *gin.Context){
	gh:=igin.H(c)
	data:=struct{
		Name string `json:"name"`
	}{}

	if err:=c.BindJSON(&data);err!=nil{
		gh.Fail(err.Error())
		return
	}
	author:=gh.Account()
	if lab:=label.Exist(data.Name);lab!=nil{
		if follow.LabelExist(author,lab.Id) {
			gh.Fail("你已经收藏了这个标签")
			return
		}
		follow.New(author,"label",lab)
		gh.Succ(nil)
	}else{
		gh.Fail("label not exist")
	}
}

func FollowLabels(c *gin.Context){
	gh:=igin.H(c)
	labels:=follow.Labels(gh.Account())
	gh.Succ(gin.H{"items":labels})
}

func FollowCancel(c *gin.Context){
	gh:=igin.H(c)
	data:=struct{
		LabelId int `json:"label_id"`
	}{}

	if err:=c.BindJSON(&data);err!=nil{
		gh.Fail(err.Error())
		return
	}

	if err:=follow.CancelLabel(gh.Account(),data.LabelId);err==nil{
		gh.Succ(nil)
	}else{
		gh.Fail(err.Error())
	}
}


func User(c *gin.Context) *auth.Customer{
	if user,ok:=c.Get("user");ok{
		userr:=user.(auth.Customer)
		return &userr
	}
	return &auth.Customer{}
}

func CreateAccount(c *gin.Context) {

	param:=struct{
		User string		`json:"user"`
		Password  string	`json:"password"`
		Nick string		`json:"nick"`
		VerifyCode string 	`json:"verify_code"`
	}{}

	c.BindJSON(&param)

	if param.User == "" || param.Password == "" {
		igin.Fail(c,"账号和密码不能为空")
		return
	}

	ch:=igin.H(c)

	if err := account.VerifyRegisterCode(ch.User().Token.Value, param.VerifyCode); err != nil {


		ch.Fail(err.Error())
		return
	}

	if user, err := account.New(param.User,param.Password,param.Nick); err == nil {
		User(c).SetAccount(user)
		GetAccountInfo(c)

	} else {
		ch.Fail(err.Error())
	}
}

//登录
func Login(c *gin.Context) {
	param:=struct{
		User string		`json:"user"`
		Password  string	`json:"password"`
	}{}
	gh:=igin.H(c)

	if err:=c.BindJSON(&param);err!=nil{
		gh.Fail(err.Error())
		return
	}

	if param.User == "" || param.Password == "" {

		gh.Fail("账号和密码不能为空")
		return
	}

	if account,err:=account.Verify(param.User,param.Password);err==nil{
		auth:=User(c)
		auth.SetAccount(account)
	}else{
		gh.Fail(err.Error())
		return
	}

	GetAccountInfo(c)
}

//修改密码
func RePassword(c *gin.Context) {

	param:=struct{
		OldPassword string	`json:"old_password"`
		Password  string	`json:"password"`
	}{}

	if err:=c.BindJSON(&param);err!=nil{
		igin.Fail(c,err.Error())
		return
	}

	if err := User(c).Account().RePassword(param.OldPassword, param.Password); err == nil {
		igin.Succ(c,nil)
	} else {
		igin.Fail(c,err.Error())
	}
}

//退出
func Logout(c *gin.Context) {
	if err:=User(c).Signout();err==nil {
		igin.Succ(c,nil)
	} else {
		igin.Fail(c,"退出失败")
	}
}

func GetAccountInfo(c *gin.Context) {
	user := User(c)
	if user.IsLogin() {
		info := user.Account()
		igin.Succ(c, gin.H{
			"item": gin.H{
				"uid":                info.Uid,
				"nick":               info.Nick,
				"sign":               info.Sign,
				"headimg":            info.Headimg,
				"unreadmessagecount": account.UnReadMessageCount(info), //未读消息
			},
		})
	}else{
		c.JSON(200,gin.H{"error_code":10000})
	}
}

//保存帐号信息
func SaveAccountInfo(c *gin.Context) {
	account := User(c).Account()


	param:= struct {
		Field string `json:"field"`
		Value  string `json:"value"`
	}{}

	if err:=c.BindJSON(&param);err!=nil{
		igin.Fail(c,err.Error())
		return
	}

	if account.EditField(param.Field, param.Value) {
		igin.Succ(c,nil)
	} else {
		igin.Fail(c,"保存失败")
	}
}

//保存帐号头像
func SaveAccountHeadimg(c *gin.Context) {
	user:=User(c)
	im := images.New("./uploads/headimg", fmt.Sprintf("%d", user.Account().Uid))
	if err := im.WriteFrom(c.Request, "file"); err == nil {
		if err := im.Save(); err == nil {
			images.Resize(im.FullName(), 160, 160)
			user.Account().UpHeadimg(im.Name)
		}
	} else {
		igin.Fail(c,err.Error())
	}
}
/*
//保存地址
func SaveShippingAddress(act *gin.Context) {
	data, err := act.ParseJson()

	if err != nil {
		act.Fail(err.Error())
		return
	}

	shipping := &account.ShippingAddress{
		Id:           data.Int("id"),
		Uid:          act.Auth.Account.Uid,
		FirstName:    data.Get("first_name"),
		LastName:     data.Get("last_name"),
		Address:      data.Get("address"),
		OtherAddress: data.Get("other_address"),
		CountryCode:  data.Get("country_id"),
		City:         data.Get("city"),
		State:        data.Get("state"),
		Zip:          data.Get("zip"),
		Phone:        data.Get("phone"),
		Email:        data.Get("email"),
	}

	if err := account.SaveShippingAddress(shipping); err == nil {
		act.Succ(configs.M{"id": shipping.Id})
	} else {
		act.Fail(err.Error())
	}
}

func GetShippingAddress(act *gin.Context) {
	items := account.GetShippingAddress(act.Auth.Account.Uid)
	act.Succ(configs.M{"items": items})
}

func DelShippingAddress(act *gin.Context) {
	if shipping, err := account.GetShippingInfo(act.GetInt("id")); err == nil && shipping.Uid == act.Auth.Account.Uid {
		shipping.Del()
		act.Succ()
	} else {
		act.Fail("del failing")
	}
}


//激活
func ActiveAccount(act *gin.Context) {
	vars := act.Vars()
	code := vars["code"]
	uid, _ := strconv.Atoi(vars["uid"])
	act.W.Header().Set("Content-Type", "text/html")

	if err := account.VerifyActive(uid, code); err == nil {
		view.New(act).Show("verify_succ.html")
	} else {
		view.New(act).Show("verify_fail.html")
	}
}
*/

//发送消息
func SendMessage(c *gin.Context) {


	param:=struct{
		Content string `json:"content"`
		To int `json:"to"`
	}{}
	c.BindJSON(&param)

	if len(param.Content) == 0 {
		igin.Fail(c,"内容不能为空")
		return
	}
	user:=User(c)
	mContent := account.MessageContent(param.Content, configs.M{})
	if err := user.Account().SendMessage(mContent, param.To); err == nil {
		igin.Succ(c,nil)
	} else {
		igin.Fail(c,err.Error())
	}
}

//删除消息
func DelMessage(c *gin.Context) {
	param:= struct {
		Id int `json:"id"`
	}{}

	c.BindJSON(&param)

	user:=User(c)
	msg := account.Message{Id:param.Id, ToId: user.Account().Uid}
	if err := msg.Del(); err == nil {
		igin.Succ(c,nil)
	} else {
		igin.Fail(c,"删除失败")
	}
}

func MessageList(c *gin.Context) {
	user:=User(c)
	list := account.ToMessageList(user.Account())
	igin.Succ(c,gin.H{
		"items":              list,
		"unreadmessagecount": account.UnReadMessageCount(user.Account())}) //未读消息
}

func SendRegisterCodeMail(c *gin.Context) {

	param:= struct {
		Email string `json:"email"`
	}{}

	c.BindJSON(&param)

	if account.IsExist(param.Email) {
		igin.Fail(c,"帐号已经存在")
		return
	}

	user:=User(c)

	if err := account.SendRegisterCodeMail(user.Token.Value, param.Email); err == nil {
		igin.Succ(c,nil)
	} else {
		igin.Fail(c,err.Error())
	}

}
