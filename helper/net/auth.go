package net

import (
	//"fmt"
	"github.com/garyburd/redigo/redis"
	"helper/account"
	"helper/crypto"
	"helper/redisCli"
	"log"
	"net/http"
	"strconv"
	"fmt"
)

type Auth struct {
	isLogin    bool
	Account    *account.Account
	tokenValue string
}

func (self *Auth) Verify(re *http.Request) bool {
	r := redisCli.Conn()
	if token, err := re.Cookie("token"); err == nil {
		self.tokenValue = token.Value
		if uid, err := redis.Int(r.Do("HGET", "token:"+token.Value,"uid")); err == nil {

			if uid > 0 {
				self.Account = account.Find(uid)
				self.isLogin = true
				return true
			}
		} else {
			log.Println(err)
		}
		return false
	} else {
		self.Account = &account.Account{Uid: 0, Nick: "游客"}
	}
	return false
}

func (self *Auth) Token() string {
	return self.tokenValue
}

func (self *Auth) Login(user, password string) error {

	account, err := account.Verify(user, password)

	if err != nil {
		return err
	}

	r := redisCli.Conn()

	token_index := account.User + strconv.Itoa(account.Uid)
	var token string

	if self.tokenValue!=""{
		token = self.tokenValue
	} else {
		token = crypto.Md5(token_index)
	}
	if err := r.Send("HSET", "token:"+token, "uid",account.Uid); err != nil {
		return err
	}
	r.Flush()

	self.Account = account
	//act.SetCookie("token", token)

	return nil
}

func (self *Auth)IsLogin()bool{
	return self.isLogin
}

func (self *Auth) SetLogin(account *account.Account) {
	r := redisCli.Conn()
	self.Account = account
	r.Do("HSET", "token:"+self.Token(),"uid", account.Uid)
}

func (self *Auth) Logout() bool {
	r := redisCli.Conn()
	r.Do("del", fmt.Sprint("token:", self.tokenValue))
	return false
}
