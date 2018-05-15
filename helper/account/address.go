package account

import (
	"errors"
	"helper/configs"
	"helper/dbs"
	"log"
)

var countrys map[string]string

type ShippingAddress struct {
	Id           int    `json:"id"`
	Uid          int    `json:"uid"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Address      string `json:"address"`
	OtherAddress string `json:"other_address"`
	CountryCode  string `json:"country_code"`
	CountryName  string `json:"country_name"`
	City         string `json:"city"`
	State        string `json:"state"`
	Zip          string `json:"zip"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
}

func (shipping *ShippingAddress) Del() bool {
	if shipping.Id > 0 {
		dbs.Exec("DELETE FROM account_address WHERE id=?", shipping.Id)
		return true
	}
	return false
}

//国家
func Countrys(code string)string{
	if countrys == nil{
		rows:=dbs.Rows("SELECT code,name FROM country")
		countrys=map[string]string{}
		for rows.Next(){
			var id,name string
			rows.Scan(&id,&name)
			countrys[id]=name
		}
	}
	return countrys[code]
}


func GetShippingAddress(uid int) []ShippingAddress {
	var items []ShippingAddress

	rows := dbs.Rows("SELECT id,first_name,last_name,address,other_address,country_code,city,state,zip,phone FROM account_address WHERE uid=?", uid)
	for rows.Next() {
		var shipping ShippingAddress
		rows.Scan(&shipping.Id, &shipping.FirstName, &shipping.LastName, &shipping.Address, &shipping.OtherAddress, &shipping.CountryCode, &shipping.City, &shipping.State, &shipping.Zip, &shipping.Phone)


		shipping.CountryName=Countrys(shipping.CountryCode)

		items = append(items, shipping)
	}
	return items
}

func GetShippingInfo(shipping_id int) (*ShippingAddress, error) {
	shipping := &ShippingAddress{}
	dbs.One("SELECT id,uid,first_name,last_name,address,other_address,country_code,city,state,zip,phone FROM account_address WHERE id=?", shipping_id).
		Scan(&shipping.Id, &shipping.Uid, &shipping.FirstName, &shipping.LastName, &shipping.Address, &shipping.OtherAddress, &shipping.CountryCode, &shipping.City, &shipping.State, &shipping.Zip, &shipping.Phone)
	if shipping.Id < 1 {
		log.Printf("shipping address not exist shipping_id(%d)", shipping_id)
		return nil, errors.New("shipping address not exist")
	}
	return shipping, nil
}

func SaveShippingAddress(shipping *ShippingAddress) error {
	data := configs.M{
		"uid":           shipping.Uid,
		"first_name":    shipping.FirstName,
		"last_name":     shipping.LastName,
		"address":       shipping.Address,
		"other_address": shipping.OtherAddress,
		"country_code":  shipping.CountryCode,
		"city":          shipping.City,
		"state":         shipping.State,
		"zip":           shipping.Zip,
		"phone":         shipping.Phone,
	}
	var err error
	if shipping.Id > 0 {
		err = dbs.Update("account_address", data, "id=?", shipping.Id)
		return err
	} else {
		shipping.Id, err = dbs.Insert("account_address", data)
	}
	return err
}
