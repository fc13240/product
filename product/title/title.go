package title

import (
	"helper/dbs"
	"helper/dbs/mongodb"
	"helper/configs"
	"strings"
)

//标签分类
type LabelCate struct{
	Id int `json:"id"`
	Name string `json:"name"`
}

func NewLabelCate(product_cid int,name string)(id int,err error){
	db:=dbs.Def()
	return db.Insert("product_title_labelcate",configs.M{"`name`":name,"product_cid":product_cid})
}

func GetLabelCateListing(filter configs.M)(data []LabelCate){
	data=[]LabelCate{}
	db:=dbs.Def()
	rows:=db.Rows("SELECT id,name FROM product_title_labelcate ORDER BY sort ASC")
	for rows.Next(){
		label:=LabelCate{}
		rows.Scan(&label.Id,&label.Name)
		data=append(data,label)
	}
	return data
}

type Lable struct {
	Id  int `json:"id"`
	Title string `json:"title"`
	CnTitle string `json:"cntitle"`
	Cid  int `json:"cid"`
}

func NewLabel(title,cntitle string,cid,id int)bool{
	l:=Lable{id,strings.Trim(title," "),strings.Trim(cntitle," "),cid}
	if l.Id>0{
		return l.Edit()
	}
	return l.Add()
}

func (l *Lable)Exist()bool{
	db:=dbs.Def()
	count:=0
	db.One("SELECT COUNT(*) FROM product_title_label WHERE cid=? AND label=?",l.Cid,l.Title).Scan(&count)
	if count>0{
		return true
	}else{
		return false
	}
}

func (l *Lable)Add()bool {
	if l.Title =="" || l.CnTitle ==""{
		return false
	}
	if l.Exist() == false {
	db:=dbs.Def()
	db.Insert("product_title_label",configs.M{"label":l.Title,"cnlabel":l.CnTitle,"cid":l.Cid})
		return true
	}
	return false
}

func (l *Lable)Edit()bool{
	if l.Title =="" || l.CnTitle ==""{
		return false
	}

	db:=dbs.Def()
	db.Update("product_title_label",configs.M{"label":l.Title,"cnlabel":l.CnTitle},"id=?",l.Id)
	return true
}

func GetLabels(cid int)(labels []Lable){
	db:=dbs.Def()
	labels=[]Lable{}
	rows:=db.Rows("SELECT id,label,cnlabel FROM product_title_label WHERE cid=?",cid)
	for rows.Next(){
		var l =Lable{Cid:cid}
		rows.Scan(&l.Id,&l.Title,&l.CnTitle)
		labels=append(labels,l)
	}
	return labels
}

type Example struct {
	Title string  `json:"title"`
	CnTitle string `json:"enTitle"`
	Images []string `json:"images"`
}

func (t *Example)Add(){
	mdb:=mongodb.Conn()
	mdb.C("title_keyword_example").Insert(&t)
}

func NewExamples(title,cntitle string,images []string){
	t:=Example{title,cntitle,images}
	t.Add()
}

func SearchTitle(q string) (labels []Lable){
	mdb:=mongodb.Conn()
	mdb.C("title_label").Find(configs.M{"title":configs.M{"$regex":q,"$options":"i"} }).All(&labels)
	titleSort(labels)
	return
}

func SearchCnTitle(q string) (labels []Lable){
	mdb:=mongodb.Conn()
	mdb.C("title_label").Find(configs.M{"cntitle":configs.M{"$regex":q,"$options":"i"} }).All(&labels)
	titleSort(labels)
	return
}

func titleSort(labels []Lable){
	 l:=Lable{}
	for i,_:=range labels{
		for j,_:=range labels{
			if len(labels[i].Title) < len(labels[j].Title){
				l=labels[i]
				labels[i]=labels[j]
				labels[j]=l
			}
		}
	}
}