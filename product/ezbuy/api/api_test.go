package api

import (
	"testing"
	"fmt"
	"encoding/json"
	"time"
)

func TestEzbuyeApi_UserProductsFromSource(t *testing.T) {
	cli:=&EzbuyeApi{
		ReqId:"FNTQ9zlUa2sAAC0j",
		Cookie:"_ga=GA1.2.1360315065.1497319248; 65_ezseller=9FB34CA91A7BD303324CC86CC8AE5D11A3C54391643E6C94A76B0F8422F52307,seller:1daGgr:h4YKNMKJxTOE22W-dA4WHBPp6f4",
	}
	cli.UserProductsFromSource(5)
	t.Error("test")
}



func TestUserUnCommitedProductUpdate(t *testing.T) {
	cli:=&EzbuyeApi{
		ReqId:"FNTQ9zlUa2sAAC0j",
		Cookie:"_ga=GA1.2.1360315065.1497319248; 65_ezseller=9FB34CA91A7BD303324CC86CC8AE5D11A3C54391643E6C94A76B0F8422F52307,seller:1daGgr:h4YKNMKJxTOE22W-dA4WHBPp6f4",
	}

	cli.UserUnCommitedProductDetail(14636418)

	t.Error("test")
}


func TestEzbuyeApi_UserUploadProducts(t *testing.T) {
	cli:=&EzbuyeApi{
		ReqId:"FNTQ9zlUa2sAAC0j",
		Cookie:"_ga=GA1.2.1360315065.1497319248; 65_ezseller=9FB34CA91A7BD303324CC86CC8AE5D11A3C54391643E6C94A76B0F8422F52307,seller:1daGgr:h4YKNMKJxTOE22W-dA4WHBPp6f4",
	}
	fmt.Println(cli.UserUploadProducts("FsAuhXzD7JGcLusu6jHSyjNB6vuB"))
}


func TestEzbuyeApi_UploadExec(t *testing.T) {
	cli:=&EzbuyeApi{
		ReqId:"FPUAFY3Eg8MAADLf",
		Cookie:"GA1.2.205130300.1509717626;65_ezseller=9FB34CA91A7BD303324CC86CC8AE5D11A3C54391643E6C94A76B0F8422F52307,seller:1eC0J6:PyWQV9YYQd37pyhJ0Ry-Ye8z07I",
	}
	cli.UploadExec("http://g.com:8081/download/ezitem/2017118/Bl22225.xlsx")
}

var api=&EzbuyeApi{
	ReqId:"FPUAFY3Eg8MAADLf",
	Cookie:"GA1.2.205130300.1509717626;65_ezseller=9FB34CA91A7BD303324CC86CC8AE5D11A3C54391643E6C94A76B0F8422F52307,seller:1eC0J6:PyWQV9YYQd37pyhJ0Ry-Ye8z07I",
}

func TestUserProductsFromSource(t *testing.T){
	data:=struct {
		Total int `json:"total"`
		Result []struct{
				Pid int `json:"pid"`
				Name string `json:"name"`
		} `json:"results"`
	}{}
	if res,err:=api.UserProductsFromSource(5);err==nil{

		res.BindJSON(&data)
		if data.Total >0 {
				fmt.Println(data.Result)
		}

	}

}

func TestUserProductDetail(t *testing.T){
	ezItem,err:=api.UserProductDetail(36686953)
	fmt.Println(ezItem,err)
	b,_:=json.Marshal(&ezItem)
	fmt.Println(string(b))
}

func TestUserProductUpdateId(t *testing.T){
	api.UserProductQuickUpdate(34553249)
}

func TestAdd(t *testing.T){
	api.Add(34553249)
}

func TestUserProductList(t *testing.T){
	var api=&EzbuyeApi{
		ReqId:"FPUAFY3Eg8M2D2Lf",
		Cookie:`_ga=GA1.2.857665726.1516867509; mp_69914a78e8e19916b5776e4841ed85e4_mixpanel=%7B%22distinct_id%22%3A%20%221614655f24e534-01f3a8364cb014-32637401-fa000-1614655f24f74d%22%2C%22%24initial_referrer%22%3A%20%22%24direct%22%2C%22%24initial_referring_domain%22%3A%20%22%24direct%22%7D; mp_672068bd75cd6425e155c3eb64d96a6e_mixpanel=%7B%22distinct_id%22%3A%20%221614655f25636b-0dc0c96d4df592-32637401-fa000-1614655f2575ba%22%2C%22%24initial_referrer%22%3A%20%22%24direct%22%2C%22%24initial_referring_domain%22%3A%20%22%24direct%22%7D; _gid=GA1.2.1352825671.1523172512; 65_ezseller="461432A01FE4B13F3402BFE4F00C14BAB02219F54B62ED61A76B0F8422F52307,seller:1f5Lis:NAM0epKsEWUybyjb69CJhx8yHB4"`}
	list:=struct{
		Products []struct{
			Pid int  `json:"pid"`
		} `json:"products"`
	}{}
	for _,i:=range []int{0,1,2,3}{
		if res,err:=api.UserProductList(i*40,40);err==nil{
			res.BindJSON(&list)
			fmt.Println(list)
			for _,p:=range list.Products{
				api.UserProductUpdateId(p.Pid)
				time.Sleep(time.Second*1)
			}
			res.Close()
		}
	}
}

func TestUserProtductListAll(t *testing.T){
	list:=struct{
		Products []struct{
			Pid int  `json:"pid"`
		} `json:"products"`
	}{}
	res,_:=api.UserProductList(1100,20)
	res.BindJSON(&list)
	fmt.Println(list.Products)
	for _,p:=range list.Products{
		api.Add(p.Pid)
		time.Sleep(time.Second*10)
	}
}