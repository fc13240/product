package manage

import (
	"helper/dbs"
	"helper/account"
	"fmt"
	"time"
)
type Member int

func (*Member) Listing(offset,limit int)(accounts []account.Account ,total int){
	db:=dbs.Def()
	accounts=[]account.Account{}
	sql:=fmt.Sprintf("SELECT uid,user,nick,register_time,last_login_time from member ORDER BY register_time DESC")
	rows:=db.Rows(sql+dbs.Limit(offset,limit))

	for rows.Next(){
		var uid int
		var register_time,last_login_time int64
		var user,nick string
		rows.Scan(&uid,&user,&nick,&register_time,&last_login_time)

		member:=account.Account{
			Uid:uid,
			User:user,
			Nick:nick,
			RegisterTime:time.Unix(register_time,0),
			LastLoginTime:time.Unix(last_login_time,0 ),
		}
		accounts=append(accounts,member)
	}
	db.One("SELECT count(*) from member").Scan(&total)
	return accounts,total
}

func (*Member)Servers(){

}
