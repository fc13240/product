package auth

import (
	"github.com/garyburd/redigo/redis"
	"helper/redisCli"
	"fmt"
	"helper/account"
	"log"
	//"helper/util"
//	"errors"
	"helper/crypto"
	"time"
	"math/rand"
)
type Token string

func (token Token) Set(vv ...interface{}){
	r:=redisCli.Conn()
	index:=tokenIndex(token)
	r.Send("HMSET",append([]interface{}{index},vv...)...)
	r.Send("EXPIRE",index,3600)
	r.Flush()
	defer r.Close()
}

func (token Token)Get(f string)string{
	r:=redisCli.Conn()
	defer r.Close()
	s,err:=redis.String(r.Do("HGET",tokenIndex(token),f))
	if err!=nil{
		log.Printf("set token  fail:%s,error:%s",f,err.Error())
	}
	return s
}


func (token Token)Clean(){
	r:=redisCli.Conn()
	defer r.Close()
	r.Send("DEL",tokenIndex(token))
	r.Flush()
}

func (token Token)GetInt(f string)int{
	r:=redisCli.Conn()
	defer r.Close()
	n,_:=redis.Int(r.Do("HGET",tokenIndex(token),f))
	return n
}

func (t *Token)BindUser(uid int){
	t.Set("uid",uid)
}

func (t *Token)UnbindUser(){
	t.Set("uid",0)
}

func (token *Token)IsLogin()bool {
	uid:=token.GetInt("uid")
	if uid>0{
		return true
	}
	return false
}
func (token Token)Uid()int {
	return token.GetInt("uid")
}

func(t Token)IsExist()(isexist bool){
	r:=redisCli.Conn()
	isexist,_=redis.Bool(r.Do("EXISTS",tokenIndex(t) ))
	defer r.Close()
	return
}
func tokenIndex(token Token)string{
	return "token:"+string(token)
}

func (t *Token)User()(author *account.Account,err error){
	uid:=t.Uid()
	author=account.Author(uid)
	return
}

func NewToken()(token Token){
	token=Token(crypto.Hmac(fmt.Sprint(time.Now().Unix()),fmt.Sprint(rand.Intn(123456))).Md5())
	return token
}

type Customer struct{
	Token 	*Token
	Uid int
	Nick string

}
/*

func (customer *Customer)Set(k string ,value interface{}) (err error){
	r:=redisCli.Conn()
	token:=fmt.Sprint("token:",customer.Token.Value)
	_,err=r.Do("HSET",token,k,value)
	if err!=nil{
		log.Println("设置token 数据失败",err.Error(),fmt.Sprint(token," ",k," ",value))
	}
	defer r.Close()
	return err
}

func VerifyToken(token string)bool{
	customer.Token=&Token{Value:token}
	if Exist(token){
		conn:=redisCli.Conn()
		defer conn.Close()

		vv,err:=redis.Values(conn.Do("hgetall", "token:"+customer.Token.Value))
		if err!=nil{
			log.Println(err.Error())
		}

		redis.ScanStruct(vv,customer.Token)
		if customer.Token.Uid>0 {
			conn.Send("set",fmt.Sprint("author:",customer.Token.Uid,":",customer.Token.Value),util.Datetime())
			conn.Send("hset",fmt.Sprint("author:",customer.Token.Uid),"visit_time",util.Datetime())
			conn.Flush()
		}
		return true
	}
	return false
}


func (customer *Customer)IsLogin()bool{
	if customer.Token == nil{
		log.Println("token is empty.")
		return false
	}
	if customer.Uid()>0{
		if uid,err:=redis.Int(customer.get("uid"));uid>0{
			return true
		}else{
			log.Println("读取缓存失败",err)
			return false
		}
	}
	return false
}

func(customer *Customer)Signout()error{
	customer.Token.Uid=0
	return customer.Set("uid",0)
}

func (customer *Customer)get(key string)(interface{},error){
	r:=redisCli.Conn()
	defer r.Close()
	token:=fmt.Sprint("token:",customer.Token.Value)
	v,err:= r.Do("hget",token,key)
	if err!=nil{
		log.Println("获取token 值失败",err)
	}
	return v,err
}


func (customer *Customer)Account() *account.Account{
	if customer._account == nil && customer.Token.Uid>0{
		customer._account=account.Find(customer.Token.Uid)
	}
	return customer._account
}

func (token Token) Exist() (exist bool){
	conn:=redisCli.Conn()
	defer conn.Close()
	var err error
	exist,err=redis.Bool(conn.Do("exists",fmt.Sprint("token:",token)))
	if err!=nil{
		log.Println(" check tocken error  auth.go ",err.Error())
	}
	return exist
}*/