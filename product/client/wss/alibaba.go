package wss

import (
	"helper/configs"
)

func AddItem(m *Message){
	args:=m.Cmd.Args

	data,err:=client.AlibabaAdd(args)
	mdata:=configs.M{}

	data.BindJSON(&mdata)

	if err!=nil{
		m.Fail(err.Error())
	}else{
		m.Succ(mdata)
	}
}

func CheckItemExist(m *Message){
	args:=m.Cmd.Args
	data,err:=client.CheckItemExist(args)
	mdata:=configs.M{}
	data.BindJSON(&mdata)
	if err!=nil{
		m.Fail(err.Error())
	}else{
		m.Succ(mdata)
	}
}