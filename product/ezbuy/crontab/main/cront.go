package main
//定时上传任务
import(
	"product/ezbuy/crontab"
	_ "ainit"
	"helper/configs"
	"helper/redisCli"
	"github.com/garyburd/redigo/redis"
	"log"
	"time"
)

var k="ezbuy:syn:crontab"
func main(){
	r:=redisCli.Conn()
	for {
		if n,_:=redis.Int(r.Do("LLEN",k));n==0{
			log.Println("没有要上传的产品")
			time.Sleep(time.Second*30)
			continue
		}else{
			log.Println("还剩余:",n)
		}
		syn()
		time.Sleep(time.Second*300)
	}
}

func syn(){
	r:=redisCli.Conn()
	defer r.Close()
	for i:=0;i<4;i++{
		pid,_:=redis.Int(r.Do("LPOP",k))
		if pid<=0 {
			continue 
		}
		time.Sleep(time.Second*5)
		if item,err:=crontab.Down(pid);err==nil{
			log.Println("开始上传一个:",item.Base.Name)
			crontab.SynToRemote(pid)
		}
	}
}

func init(){
	r:=redisCli.Conn()
	defer r.Close()
	if n,_:=redis.Int(r.Do("LLEN",k));n>0{
		log.Println("还有:",n)
		return 
	}//"infor":true
	rows,_:=crontab.Listing(configs.M{},0,10000)
	for _,row:=range rows{
		r.Send("RPUSH",k,row.Pid)
	}
	r.Flush()
}