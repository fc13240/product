package label

import (
	"helper/dbs"
	"helper/util"
	"helper/configs"
	"strings"
	"helper/redisCli"
	"github.com/garyburd/redigo/redis"
	"errors"
	"fmt"
	"time"
)

type Label struct{
	Name string `json:"name"`
	Id int	`json:"id"`
}

func Exist(name string) *Label{
	var label *Label
	db:=dbs.Def()
	var id int
	db.One("SELECT id FROM label where md5=?",util.Md5(name)).Scan(&id)
	if id>0{
		label=&Label{Name:name,Id:id}
	}
	return label
}

func Search(name string,limit int)(labels []Label){
	db:=dbs.Def()
	sql:=fmt.Sprintf("SELECT id,name FROM label where name LIKE '%s%s' LIMIT %d",name,"%",limit)
	fmt.Println(sql)
	rows:=db.Rows(sql)
	labels=[]Label{}
	for rows.Next(){
		var id int
		var name string
		rows.Scan(&id,&name)
		labels=append(labels,Label{Id:id,Name:name})
	}
	return labels
}

func New(name string) (lab *Label, err error){
	name=strings.Trim(name," ")

	if name ==""{
		err=errors.New("label不能为空")
		return lab,err
	}

	if old_lab:=Exist(name);old_lab!=nil{
		return old_lab,err
	}

	id,err:=dbs.Insert("label",configs.M{"name":name,"md5":util.Md5(name)})

	if err!=nil{
		err=errors.New(fmt.Sprint("创建标签,",name,",失败",err.Error()))
	}else{
		lab=&Label{Id:id,Name:name}
	}
	return lab,err
}

type Result struct {
	Id int `json:"id"`
	Name string `json:"name"`
	ResultNum int `json:"result_num"`
	UpdateTime time.Time `json:"update_time"`
}

func GetResultSort(limit int)[]Result{
	r:=redisCli.Conn()
	ids,_:=redis.Ints(r.Do("ZRANGEBYSCORE labels_resultnum_sort:",0,limit))
	data:=[]Result{}
	for _,label_id:=range ids{
		res:=Result{}
		vv,_:=redis.Values(r.Do("hgetall",fmt.Sprint("label:",label_id)))
		redis.ScanStruct(vv,&res)
		data=append(data,res)
	}
	return data
}

func DelRel(){

}
