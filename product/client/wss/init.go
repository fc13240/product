package wss
import (
	"fmt"
	"net/http"
	"product/ezbuy"
	"product/ezbuy/api"
	ezclient "product/ezbuy/client"

)
var (
	token string
	ezcall *Call
	hcall *Call
	albabaCall *Call
	client *ezclient.HClient
 	setting *ezbuy.Setting
 	ez *api.EzbuyeApi
	isRefresh =false
	ezMessage *Message
	aliMessage *Message
	runSta=ezbuy.RunSta
)
const(
	On  ="on"
	Off ="off"
)
func init() {
	hcall=&Call{}
	if err:=Auth("api");err!=nil{
		fmt.Println(err)
		return
	}
	hcall.Add("refresh",OnRefresh)
	hcall.Add("updateorder",UpdateOrder)
	hcall.Add("down",Down)
	hcall.Add("fillez",FillEZ)
	hcall.Add("upitems",UpItems)
	hcall.Add("userProductsFromSource",UserPruoductsFromSource)
	hcall.Add("onSale",OnSale)
	//上传一个产品到yz
	hcall.Add("uploadExec", func(m *Message){
		url:=m.Cmd.Args.Get("url")
		store_id:=m.Cmd.Args.Int("store_id")
		
		store:=ChangeShop(store_id)
		
		_,err:=ez.UploadExec(url)
		if err==nil{
			m.Succ("上传到: "+store.StoreName+" 成功")
		}else{
			m.Fail(err.Error())
		}
	})

	http.HandleFunc("/cli", Start)

	ezcall=&Call{}
	
	http.HandleFunc("/ezcli", EzStart)

	albabaCall=&Call{}

	albabaCall.Add("additem",AddItem)
	albabaCall.Add("checkItemExist",CheckItemExist)

	http.HandleFunc("/1688",Alibaba)

	port:=":1826"
	fmt.Println(fmt.Println("Start OK",port))
	err := http.ListenAndServeTLS(port, "./cert.crt","./key.key",nil)
	if err != nil {
		fmt.Println(err)
	}
}

func OnRefresh(m *Message){
	Refresh()
}


