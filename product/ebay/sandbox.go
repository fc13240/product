package ebay

import (
	"fmt"
)

type Account struct {
	DevId  string
	AppId  string
	CertId string
	Token  string
	AppUrl string
}

var Mod string = "dev"

func ProductionMod() {
	Mod = "production"
}

func (account *Account) Init() *Account {
	if Mod == "dev" {
		fmt.Println(Mod)
		account.d()
	} else if Mod == "production" {
		fmt.Println(Mod)
		account.p()
	}

	return account
}

func (account *Account) d() {
}

func (account *Account) p() {
}
