package ezbuy
import(
	"testing"
	"helper/webtest"
	"fmt"
)
func TestSaveSetting(t *testing.T){
	body:=`{
		"authorid":0,
		"checkneworders":true,
		"cookie":"",
		"minute":0,
		"num":0,
		"onrefresh":true,
		"reqid":"",
		"secrectkey":"K86LIY56WA33YC44",
		"skufirst":"",
		"store_cateid":0,
		"storeid":0,
		"storename":""}`
	h:=webtest.NewHeader()
	h.Add("Authorization","K86LIY56WA33YC44")
	res,_:=webtest.PostJson("http://192.168.8.119:8081/api/ezbuy.savesetting",h,body)
	fmt.Println(res.String())
}