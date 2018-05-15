package redisCli

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"helper/configs"
	"log"
	"time"
)

var link *redis.Pool

func newPool(host,password,port string) *redis.Pool {

	return &redis.Pool{
		MaxIdle:     10,
		MaxActive:   2000, // max number of connections
		IdleTimeout: 5 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", host, port))

			if err != nil {
				log.Println("redis 连接错误:", err.Error())
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					log.Printf("auth error:", err.Error())
				}
			}
			c.Do("SELECT", 0)
			return c, nil
		},
	}
}

func Conn() redis.Conn {
	if link == nil{
		options := configs.GetSection("redis")
		host := options["host"]
		password := options["password"]
		port := options["port"]
		link = newPool(host,password,port)
	}
	return link.Get()
}

func Close(){
	if link != nil{
		link.Close()
	}
}

func ActiveCount() int {
	if link !=nil{
		return link.ActiveCount()
	}else{
		return 0
	}
}


func New(host,password,port string)*redis.Pool{
	return newPool(host,password,port)
}

func Get(name string) *Value {
	c := link.Get()
	defer c.Close()
	value := &Value{}
	if v, err := c.Do("GET", name); err != nil {
		log.Println("redis 操作失败:", err)
		return nil
	} else {
		value.Val = v
		return value
	}
}

func Del(key string) bool {
	c := link.Get()
	defer c.Close()
	if _, err := c.Do("DEL", key); err != nil {
		log.Println(err)
		return false
	}
	return true
}

type Value struct {
	Val interface{}
}

func (self *Value) Int() int {
	v, _ := redis.Uint64(self.Val, nil)
	return int(v)
}

func (self *Value) String() string {
	s, _ := redis.String(self.Val, nil)
	return s
}
