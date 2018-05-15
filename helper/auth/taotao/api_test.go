package taobao

import (
	"testing"
	"helper/configs"
	"fmt"
)

func TestRequest(t *testing.T) {
	res,_:=Request("taobao.tbk.uatm.favorites.item.get",configs.M{
		"fields":  "num_iid,title,tk_rate,click_url",
		"adzone_id":"129028491",
		"favorites_id":"5725068",
		})
		fmt.Println(res)
}


func TestRequest2(t *testing.T) {
	res,_:=Request("taobao.tbk.uatm.favorites.get",configs.M{
		"fields":"favorites_title,favorites_id,type",
		"page_size":"20",
	})
	fmt.Println(res)
}



func TestSearch(t *testing.T){
	res,_:=Request("taobao.tbk.item.get",
		configs.M{
			"fields":"click_url,num_iid,title,pict_url,small_images,reserve_price,zk_final_price,user_type,provcity,item_url,seller_id,volume,nick",
			"q":"SUPOR/苏泊尔 CYSB50YCW10D-100电压力锅家用双胆5L",
			},
		)
	fmt.Println(res)
}

func TestS( t *testing.T){
	res,_:=Request("taobao.tbk.item.convert",configs.M{
		"fields":"num_iid,click_url",
		"num_iids":"535950639471",
		"adzone_id":"129028491",
	})

	fmt.Println(res)
}