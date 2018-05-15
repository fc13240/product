package account

import (
	"helper/dbs"
	"fmt"
)

type Serve struct{
	Router string `json:"router"`
	Code string `json:"code"`
	Name string `json:"name"`
	Desc string `json:"desc"`
}

func MyServeList(author_id int )(list []Serve){
	db:=dbs.Def()
	sql:=fmt.Sprintf("select router,serve.code,serve.name from account_serve as ser right join serve ON(ser.`code`=serve.`code` AND serve.is_enable=1) where ser.author_id=%d",author_id)
	rows:=db.Rows(sql)
	list=[]Serve{}
	for rows.Next(){
		var name ,code,router string
		rows.Scan(&router,&code,&name)
		list=append(list,Serve{Name:name,Code:code,Router:router})
	}
	return list
}

func Servers() (list []Serve){
	db:=dbs.Def()
	list=[]Serve{}
	sql:=fmt.Sprint("SELECT `code`,`name`,`router`,`desc` FROM serve")
	rows:=db.Rows(sql)
	for rows.Next(){
		ser:=Serve{}
		rows.Scan(&ser.Code,&ser.Name,&ser.Router,&ser.Desc)
		list=append(list,ser)
	}
	return
}