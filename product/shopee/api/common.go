package api
import(
	"encoding/json"
	"helper/crypto"
	"helper/webtest"
	"helper/configs"
	"time"
	"fmt"
)

var (
	ApiAddr string
	Secret string
	ParterId int
	Shopid int
)

func init(){
	 config:=configs.GetSection("shopee")
	 ApiAddr=config["apiAddr"]
	 ParterId=configs.Int(config["parterId"])
	 Shopid=configs.Int(config["shopid"])
	 Secret=config["secret"]
}


func POST(addr string,param configs.M) (*webtest.Result,error){
	timestamp:=time.Now().Unix()
	
	addr=ApiAddr+addr
	
	common_param:=configs.M{
		"partner_id":ParterId,
		"shopid":Shopid,
		"timestamp":timestamp,
	}
	//bb,_:=json.Marshal(common_param)
	param.Add("partner_id",common_param["partner_id"])
	param.Add("shopid",common_param["shopid"])
	param.Add("timestamp",common_param["timestamp"])
	
	body,err:=json.Marshal(param)
	fmt.Println(err)
	content:=string(body)
	h:=webtest.NewHeader()
	secret_content:=addr+"|"+content
	h.Add("Authorization",crypto.Hmac(secret_content,Secret).Sha256())
	
	return webtest.PostJson(addr,h,content)

}