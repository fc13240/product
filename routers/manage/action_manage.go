package manage

import (
	"helper/manage"
	"helper/dbs"
	"github.com/gin-gonic/gin"
	"helper/net/igin"
	"helper/account"
	"helper/configs"
	"helper/redisCli"
	"github.com/garyburd/redigo/redis"
	"fmt"
)

func Console(c *gin.Context){
	gh:=igin.H(c)

	var console manage.Console

	redisactiveCount:=console.RedisActiveCount()

	gh.Succ(gin.H{"redisactivecount":redisactiveCount})


}

func Members(c *gin.Context){
	var db=dbs.Def()
	gh:=igin.H(c)

	data:= struct {
		Offset int `json:"offset"`
		Limit int  `json:"limit"`
	}{}

	if err:=c.BindJSON(&data);err!=nil{
		gh.Fail(err.Error())
		return
	}

	memberes := []configs.M{}

	rows:=db.Rows("SELECT uid,user,role,nick,is_active,sign,headimg FROM account "+dbs.Limit(data.Offset,data.Limit))

	var total int

	db.One("SELECT COUNT(*) FROM account").Scan(&total)
	r:=redisCli.Conn()
	defer r.Close()
	for rows.Next(){
		item:=account.Account{}
		rows.Scan(&item.Uid, &item.User,&item.Role,&item.Nick,&item.IsActive,&item.Sign, &item.Headimg)
		visit_time,_:=redis.String(r.Do("hget",fmt.Sprint("author:",item.Uid),"visit_time"))

		memberes=append(memberes,configs.M{
			"uid":item.Uid,
			"user":item.User,
			"role":item.Role,
			"nick":item.Nick,
			"serves":account.MyServeList(item.Uid),
			"visit_time":visit_time,
		})
	}

	gh.Succ(gin.H{"items":memberes,"total":total})
}

func Servers(c *gin.Context){

	gh:=igin.H(c)
	gh.Succ(gin.H{"items":account.Servers()})
}