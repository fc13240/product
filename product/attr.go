package product

import (
	"strings"
	"fmt"
	"helper/configs"
	"helper/dbs"
	"helper/util"
	"errors"
	"database/sql"
	
)

const (
	INPUT_NUMBER="number"
	INPUT_TEXT="text"
	INPUT_DATETIME="datetime"
	INPUT_DATE="date"
	INPUT_CHECKBOX="checkbox"
	INPUT_RADIO="radio"
	INPUT_OPTIONS="options"
)

const(
	Color=3//"colors"
	Size =1//"clothing_size"
)

func InputTypes()map[string]string{
	return map[string]string{
		INPUT_NUMBER:"数字类型",
		INPUT_TEXT:"输入框",
		INPUT_DATETIME:"时间类型",
		INPUT_DATE:"日期",
		INPUT_CHECKBOX:"多选项",
		INPUT_RADIO:"单选项",
		INPUT_OPTIONS:"下拉选项",
	}
}

type Attr struct{
	Id int `json:"id"`
	Label string  `json:"label"`
	InputType string `json:"input_type"`
	IsMandatory bool
	Name string `json:"name"`
	Options []AttrOptions `json:"options"`
	Type string `json:"type"`
	Config configs.M
}

func (attr *Attr)SetConfig(config string){
	attr.Config=configs.M{}
	if config!=""{
		for _,row:=range strings.Split(config,","){
			ss:=strings.Split(row,"=")
			attr.Config.Add(ss[0],ss[1])
		}
	}
}

//标记已经选项的
func (attr *Attr)FlagCheckedOption(sku string){
	db:=dbs.Def()
	rows:=db.Rows("SELECT valueid FROM product_attrval WHERE sku=? AND attid=?",sku,attr.Id)
	defer rows.Close()
	for rows.Next(){
		var valueid int 
		rows.Scan(&valueid)
		for i,_:=range attr.Options{
			if attr.Options[i].Id == valueid{
				attr.Options[i].IsChecked=true
			}
		}
	}
}

//新增属性
func NewAtt(label,name,InputType string)(id int,err error){
	db:=dbs.Def()
	if AttExist(label,name,0) {
		err=errors.New("属性label 或 name 已经存在")
		return 
	}
	return  db.Insert("product_attr",configs.M{"`label`":label,"`name`":name,"`input_type`":InputType})
}

//属性是否存在 编辑的时候 id>0
func AttExist(label,name string,id int)bool{
	db:=dbs.Def()
	total:=0
	db.One("SELECT COUNT(*) FROM product_attr WHERE label=? OR `name`=? ",label,name).Scan(&total)
	if total>0{
		return true
	}
	return false
}

//新增属性值
func NewAttOption(attid int ,value ,cnvalue string)(id int, err error){
	db:=dbs.Def()
	
	if AttrOptionExist(attid,value){
		err=errors.New("属性值已经存在")
		return 
	}
	return  db.Insert("product_attoptions",configs.M{"`attid`":attid,"`value`":value,"`cnvalue`":cnvalue})	
}

//属性值是否存在
func AttrOptionExist(attid int ,value string) bool{
	db:=dbs.Def()
	total:=0
	db.One("SELECT COUNT(*) FROM product_attoptions WHERE `attid`=? AND `value`=?",attid,value).Scan(&total)
	if total>0{
		return true
	}
	return false
}

type AttrOptions struct{
	Value string `json:"value"`
	Id int `json:"id"`
	Attid int `json:"attid"`
	CnValue string `json:"cnvalue"`
	IsChecked bool `json:"isChecked"`
}


//获取分类下面的属性列表
func GetCategoryAttrs(platform string,cate_id int,is_get_options bool)(attrs []Attr,err error){
	
	attrids:=[]int{}
	ids:=GetCategoryAttrIds(platform,cate_id)
	
	for attid,_:=range ids{
		attrids=append(attrids,attid)
	}
	if len(attrids) == 0{
		return attrs,errors.New("没有配置属性")
	} 
	attrs=GetAttrs(attrids...)
	
	new_attrs:=[]Attr{}
	if is_get_options {
		for i,attr:=range attrs{
			if attr.Type == "normal"{
				attrs[i].Options=attr.GetOptions()
				new_attrs=append(new_attrs,attrs[i])
			}
		}
	}else{
		for i,attr:=range attrs{
			if attr.Type == "normal"{
				new_attrs=append(new_attrs,attrs[i])
			}
		}
	}
	return new_attrs,nil
}

//获取分类下面的SKU属性列表
func GetCategorySkuAttrs(platform string,cate_id int,is_get_options bool)(attrs []Attr,err error){
	
	attrids:=[]int{}
	ids:=GetCategoryAttrIds(platform,cate_id)
	
	for attid,_:=range ids{
		attrids=append(attrids,attid)
	}
	if len(attrids) == 0{
		return attrs,errors.New("没有配置属性")
	} 
	attrs=GetAttrs(attrids...)
	new_attrs:=[]Attr{}
	if is_get_options {
		for i,attr:=range attrs{
			if attr.Type == "sku"{
				attrs[i].Options=attr.GetOptions()
				new_attrs=append(new_attrs,attrs[i])
			}
		}
	}
	return new_attrs,nil
}

//获取单个属性
func GetAttr(id int)(attr *Attr){
	attr=&Attr{Config:configs.M{}}
	db:=dbs.Def()
	config:=""
	db.One("SELECT id,label,name,input_type,config FROM product_attr WHERE id=?",id).
	Scan(&attr.Id,&attr.Label,&attr.Name, &attr.InputType,&config)
	attr.SetConfig(config)
	return attr
}



//属性列表
func GetAttrs(ids ...int)(attrs []Attr){
	attrs=[]Attr{}
	db:=dbs.Def()
	var sql string
	
	if len(ids)>0{
		sql=fmt.Sprintf("SELECT id,label,name,input_type,att_type,config FROM product_attr WHERE id IN(%s)",util.IntJoin(ids,","))
	}else{
		sql=fmt.Sprintf("SELECT id,label,name,input_type,att_type,config FROM product_attr")
	}

	rows:=db.Rows(sql)
	defer rows.Close()
	for rows.Next(){
		attr:=Attr{}
		config:=""
		rows.Scan(&attr.Id,&attr.Label,&attr.Name, &attr.InputType,&attr.Type,&config)
		attr.SetConfig(config)

		attrs=append(attrs,attr)
	}
	return attrs
}

//属性值
//selected 只获得选择（selected）的部分
func (att *Attr)GetOptions(selected ...int) (options []AttrOptions){
	db:=dbs.Def()
	att.Options=[]AttrOptions{}
	ss:="SELECT id,attid,value,cnvalue FROM product_attoptions WHERE attid=?"
	var rows *sql.Rows
	
	if len(selected)>0{
		ss=ss+ " AND id IN("+util.IntJoin(selected,",")+") "
		rows=db.Rows(ss,att.Id)
	}else{
		rows=db.Rows(ss,att.Id)
	}
	defer rows.Close()
	for rows.Next(){
		opt:=AttrOptions{}
		rows.Scan(&opt.Id,&opt.Attid,&opt.Value,&opt.CnValue)
		att.Options=append(att.Options,opt)
	}
	return att.Options
}

//获取一个SKU，一下属性的选项
func GetOneAttrSelectedOption(sku string ,attid int)(options []AttrOptions){
	 selectd:=[]int{}
	rows:=dbs.Rows("SELECT valueid FROM product_attrval WHERE sku=? AND  attid=?",sku,attid)
	for rows.Next(){
		var id int 
		rows.Scan(&id)
		selectd=append(selectd,id)
	}
	options=[]AttrOptions{}
	//
	if len(selectd) == 0{
		return options
	}
	att:=GetAttr(attid)
	return att.GetOptions(selectd...)
}

func GetAttrSelectedOption(sku string)(atts []Attr){
	atts=[]Attr{}
	
   rows:=dbs.Rows("SELECT attid,valueid FROM product_attrval WHERE sku=?",sku)
   for rows.Next(){
	   var attid,valueid int 
	   rows.Scan(&attid,&valueid)
	   att:=GetAttr(attid)
	   att.Options= att.GetOptions(valueid)
	   atts=append(atts,*att)
   }
   return atts
}

//获取一个SKU所有option,除sku选项值以为
func GetSkuSelectedOptions(sku string,filter configs.M)(opts []AttrOptions){
	db:=dbs.Def()
	if to:=filter.Get("filter");to == "color"{

	}
	sql:=`SELECT opt.id,opt.attid,opt.value,opt.cnvalue FROM product_attrval
			LEFT JOIN product_attoptions AS opt ON(opt.id = product_attrval.valueid)
		WHERE product_attrval.sku='%s' AND product_attrval.attid NOT IN(%d,%d)`

	sql=fmt.Sprintf(sql,sku,Color,Size)
	
	rows:=db.Rows(sql)

	defer rows.Close()
	
	opts=[]AttrOptions{}

	for rows.Next(){
		opt:=AttrOptions{}
		rows.Scan(&opt.Id,&opt.Attid,&opt.Value,&opt.CnValue)
		opts=append(opts,opt)
	}
	return opts
}

//保存产品属性
func SaveAttrVal(sku string,attid int,opts []AttrOptions){
	db:=dbs.Def()
	db.Exec("DELETE FROM product_attrval WHERE sku=? AND attid=?",sku,attid)
	for _,opt:=range opts{
		db.Exec("INSERT INTO product_attrval(attid,valueid,value,sku)VALUE(?,?,?,?)",attid,opt.Id,opt.Value,sku)
	}
}


//获取产品属性
func GetAttrVal(sku string)(selected map[int][]int){
	selected=map[int][]int{}
	db:=dbs.Def()
	
	rows:=db.Rows("SELECT attid,valueid FROM product_attrval WHERE sku=?",sku)
	var attid,valueid int
	for rows.Next(){
		rows.Scan(&attid,&valueid)
		if selected[attid] == nil {
			selected[attid]=[]int{}
		}
		selected[attid]=append(selected[attid],valueid)
	}
	return 
}

type AttributeInfo struct{

}
func SaveAttributeInfo(sku string ,){}

func GetColorList(selected ...int) (opts []AttrOptions) {
	return GetAttr(Color).GetOptions(selected...)
}
                             
func GetSizeList(selected ...int) (opts []AttrOptions){
	
	return GetAttr(Size).GetOptions(selected...)
}