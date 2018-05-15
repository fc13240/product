package product

import (
	"helper/dbs"
	"helper/dbs/mongodb"
	"helper/configs"
	"time"
	"encoding/json"
)

var (
	logCol *mongodb.Collection
)

//标签
type Label struct{
	Id int `json:"id"`
	Name string `json:"name"`
	Color string `json:"color"`
}

//日志
type LabelLog struct{
	LabelId int `json:"labelid"`
	LabelName string `json:"labelname"`
	Remarks string `json:"remarks"`
	Addtime time.Time `json:"addtime"`
	Itemid int `json:"ordernum"`
}

//关联产品
type LabelRel struct{
	LabelId int
	Rid	int
}

func SaveLabelRel(rid int,label_id... int){

	CancelLabel(rid) //取消之前的关联，从新关联
	db:=dbs.Def()
	for _,lid:=range label_id{
		db.Insert("product_label_rel",configs.M{"label_id":lid,"rid":rid})
	}
}

//产品是否有关联标签
func IsRelLabel(rid int)bool{
	db:=dbs.Def()
	var total=0
	db.One("SELECT COUNT(*) FROM product_label_rel WHERE rid=?",rid).Scan(&total)
	if total>0{
		return true
	}else{
		return false
	}
}

//格式化时间
func (d *LabelLog) MarshalJSON() ([]byte, error) {
	type Alias LabelLog
	return json.Marshal(&struct {
		*Alias
		Addtime string `json:"addtime"`
	}{
		Alias: (*Alias)(d),
		Addtime: d.Addtime.Format("01/02 15:04"),
	})
}

func Labels()(labels []*Label){
	db:=dbs.Def()
	labels=[]*Label{}
	rows:=db.Rows("SELECT id,name,color FROM product_label")
	for rows.Next() {
		var id int
		var name,color string
		rows.Scan(&id,&name,&color)
		labels=append(labels,&Label{Id:id,Name:name,Color:color})
	}
	return labels
}

func GetLabel(id int) (label *Label,err error){
	db:=dbs.Def()
	label=&Label{}
	err=db.One("SELECT id,name FROM product_label WHERE id=? LIMIT 1",id).Scan(&label.Id,&label.Name)
	return label,err
}

func  LogCol() *mongodb.Collection{
	if logCol==nil{
		logCol=mongodb.Conn().C("product.labellog")
	}
	return logCol
}

func (l *Label)AddLog(content string,item_ids ...int){
	for _,itemid:=range item_ids{
		SaveLabelRel(itemid,l.Id)

		log:=configs.M{
			"labelid":l.Id,
			"labelname":l.Name,
			"remarks":content,
			"addtime":time.Now(),
			"itemid":itemid,
		}
		LogCol().Insert(log)
	}
}

//取消关联
func CancelLabel(rid int ){
	db:=dbs.Def()
	db.Exec("DELETE FROM product_label_rel WHERE rid=?",rid)
}

//获取第一个标签日志，如果有绑定标签
func FirstLabelLog(itemid int)(log *LabelLog,ok bool){
	if IsRelLabel(itemid) == false{
		ok=false
		return
	}

	log=&LabelLog{}
	LogCol().Find(configs.M{"itemid":itemid}).Sort("-addtime").One(log)

	if log.LabelId>0{
		return log,true
	}else{
		return log,false
	}
}

func ItemLabelLogs(itemid int)(logs []*LabelLog){
	logs=[]*LabelLog{}
	LogCol().Find(configs.M{"itemid":itemid}).Sort("-addtime").All(logs)
	return logs
}



