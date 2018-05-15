package account

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"helper/crypto"
	"helper/dbs"
	"helper/mail"
	"helper/redisCli"
	"math/rand"
	"time"
)

type Account struct {
	Uid           int    `redis:"uid" json:"uid"`
	Nick          string `redis:"nick" json:"nick"`
	User          string `redis:"user" json:"user"`
	Sign          string `redis:"sign" json:"sign"`
	Headimg       string `redis:"headimg" json:"headimg" `
	password      string
	LastLoginTime time.Time
	RegisterTime  time.Time
	IsActive      bool
	Role	      string `json:"role"`
}

func Find(uid int) *Account {
	account := &Account{}
	dbs.One("SELECT uid,user,role,nick,is_active,sign,headimg FROM account WHERE uid=?", uid).
		Scan(&account.Uid, &account.User,&account.Role,&account.Nick,&account.IsActive,&account.Sign, &account.Headimg)
	return account
}

func Author(author_id int) *Account {
	account := &Account{}
	index := fmt.Sprint("author:", author_id)
	r := redisCli.Conn()
	defer r.Close()

	if exist, err := redis.Bool(r.Do("exists", index)); exist && err == nil { //如果有缓存，就读缓存
		if res,err:=redis.Values(r.Do("hgetall", index));err==nil{
			redis.ScanStruct(res,account)
			return account
		}
	}

	account = Find(author_id)
	r.Do("hmset", index, "uid", author_id, "nick", account.Nick, "user", account.User, "sign", account.Sign, "headimg", account.Headimg)
	return account

}

func New(user, password, nick string) (account *Account, err error) {

	account = &Account{}

	if IsExist(user) {
		return account, errors.New("账号已经存在")
	}

	stmt := dbs.Prepare("INSERT account(user,password,nick,register_time) VALUES(?,?,?,?)")
	defer stmt.Close()

	if re, err := stmt.Exec(user, crypto.Md5(password), nick, time.Now().Unix()); err == nil {
		if uid, err := re.LastInsertId(); err == nil {
			//SendRegisterMail(Find(int(uid)))
			account.User = user
			account.Nick = nick
			account.Uid = int(uid)
			return account, nil
		}
	}
	return account, errors.New("创建帐号失败")
}

func (self *Account) EditField(field string, value string) bool {
	if field == "nick" {
		dbs.Exec("UPDATE account SET nick=? WHERE uid=?", value, self.Uid)
	} else if field == "sign" {
		dbs.Exec("UPDATE account SET sign=? WHERE uid=?", value, self.Uid)
	} else if field == "headimg" {
		dbs.Exec("UPDATE account SET headimg=? WHERE uid=?", value, self.Uid)
	}
	return true
}

func (self *Account) UpHeadimg(name string) {
	self.EditField("headimg", name)
}

func (acc *Account) ISAdmin() bool {
	if acc.Role=="admin"{
		return true
	}
	return false

}

func IsExist(user string) bool {
	var uid int
	dbs.One("SELECT uid FROM account WHERE user=?", user).Scan(&uid)
	if uid > 0 {
		return true
	} else {
		return false
	}
}


func (self *Account) SetActive() {
	dbs.Exec("UPDATE account SET is_active=1 WHERE uid=?", self.Uid)
}

func (self *Account) SendMessage(content *MContent, to_member ...int) error {
	message := &Message{Content: content, FromId: self.Uid}
	return message.Send(to_member)
}

//更改密码
func (self *Account) RePassword(old_password, new_password string) error {

	if crypto.Md5(old_password) != self.password {
		return errors.New("原始密码不正确")
	}

	if len(new_password) < 5 {
		return errors.New("新密码不能小于6位数")
	}

	if old_password == new_password {
		return errors.New("新密码不能和老密码一样")
	}

	return dbs.Exec("UPDATE account SET password=? WHERE uid=?", crypto.Md5(new_password), self.Uid)
}

//验证帐号密码
func Verify(user, password string) (*Account, error) {
	account := &Account{}
	var old_password string
	var last_login_time int64
	dbs.One("SELECT uid,user,nick,last_login_time,is_active,password FROM account WHERE user=?", user).
		Scan(&account.Uid, &account.User, &account.Nick, &last_login_time, &account.IsActive, &old_password)

	if account.Uid < 1 {
		return account, errors.New("帐号不存在")
	}

	if crypto.Md5(password) != old_password {
		return account, errors.New("密码错误")
	}

	if account.IsActive == false {
		return account, errors.New("帐号没有激活")
	}

	if last_login_time > 0 {
		account.LastLoginTime = time.Unix(last_login_time, 0)
	}

	return account, nil
}

//发送激活邮件
func SendRegisterMail(u *Account) error {

	body := `<html>
				<body>
					<h5>Hi,%s,感谢您的注册</h5>
					<p>
					</p>
				</body>
			</html>
			`
	return mail.Send(u.User, "感谢您的注册", fmt.Sprintf(body, u.Nick))
}

//验证注册验证码
func VerifyRegisterCode(token, code string) error {
	r := redisCli.Conn()
	key := fmt.Sprint("token:", token, ":register_verify_code")

	if isExist, _ := redis.Bool(r.Do("exists", key)); isExist == false {
		return errors.New("注册验证码不存在")
	}

	if c, err := redis.String(r.Do("get", key)); err == nil && c == code {
		return nil
	} else {
		fmt.Println(c, ":", code)
	}

	return errors.New("注册验证码验证输入错误")
}

//发送验证码
func SendRegisterCodeMail(token, email string) error {
	r := redisCli.Conn()
	key := fmt.Sprint("token:", token, ":register_verify_code")

	if s, err := redis.Int(r.Do("ttl", key)); err == nil && s > 240 {
		return errors.New("您请求的太频繁，请稍等一会再试")
	}

	code := fmt.Sprint(rand.Int31n(12342678))

	r.Send("set", key, code)
	r.Send("expire", key, 300)

	r.Flush()

	body := `<html>
				<body>
					<h5>Hi,感谢您的访问</h5>
					<p>您的注册码是:<br/>
					<h3>%s</h3>
					</p>
				</body>
			</html>
			`
	return mail.Send(email, "51助手注册验证码", fmt.Sprintf(body, code))
}

//验证激活码
func VerifyActive(uid int, code string) error {

	user := Find(uid)

	if user.Uid == 0 {
		return errors.New("帐号不存在")
	}

	var code1, add_date string

	dbs.One("SELECT code,add_date FROM account_active where uid=? ORDER BY add_date DESC", uid).Scan(&code1, &add_date)

	if code != code1 {
		return errors.New("激活码不正确")
	}

	user.SetActive()

	dbs.Exec("DELETE FROM account_active WHERE uid=?", uid)

	return nil
}
