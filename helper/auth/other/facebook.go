package other

import (
	"helper/dbs"
	"helper/account"
	"fmt"
	"helper/configs"
	"helper/util"

	"qiniupkg.com/x/errors.v7"
)

type Plus struct{
  Uid int
  Flag int
  Puid int64
}

type Loginer interface{
	Login()(*account.Account, error)
}

func (plus Plus)Insert()(id int ,err error){
	return  dbs.Insert("account_plus",configs.M{"uid":plus.Uid,"flag":plus.Flag,"puid":plus.Puid,"addtime":util.Datetime()})
}
type  FaceBook struct{
	 Token string	`json:"accessToken"`
	 ExpiresIn int	`json:"expiresIn"`
	 SignedRequest	string `json:"signedRequest"`
	 UserID int64	`json:"userID"`
	 Nick string 	`json:"nick"`
}

func (facebook *FaceBook)Login()(*account.Account, error){

	if facebook.UserID == 0{
		return nil,errors.New("user id is empty")
	}

	if facebook.Nick == ""{
		return nil,errors.New("nick is empty")
	}

	if uid,ok:=facebook.Exist();ok{
		account:=account.Find(uid)
		return account,nil
	}else{
		return facebook.signup()
	}
}

func (facebook *FaceBook) signup()(*account.Account, error){
	username:=fmt.Sprint(facebook.UserID,"@facebook.com")
	user,err:=account.New(username,"null",facebook.Nick)
	if err==nil{
		plus:=&Plus{Uid:user.Uid,Flag:1,Puid:facebook.UserID}
		plus.Insert()
	}
	return user,err
}


func (facebook *FaceBook)Exist()(uid int ,exist bool){
	dbs.One("SELECT uid FROM account_plus WHERE puid=?",facebook.UserID).Scan(&uid)

	if uid>0{
		return uid,true
	}
	return uid,false
}