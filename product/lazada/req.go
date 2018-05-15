package lazada

import (
	"fmt"
	"time"
	"helper/webtest"
	"crypto/sha256"
	"net/url"
	"crypto/hmac"
	"encoding/hex"
)

func NewReq(action_name string)*Req{
	v:=&url.Values{}
	v.Add("Action",action_name)
	return &Req{Params: v,Url:"https://api.sellercenter.lazada.com.my"}
}


type Req struct{
	Url string
	Params *url.Values
}


func (a *Req)SetParam(k,v string){
	a.Params.Set(k,v)
}

func (api *Req)Encode() string {

	v:=api.Params

	h := hmac.New(sha256.New, []byte(ApiKey))
	v.Add("Timestamp",time.Now().Format(time.RFC3339))
	v.Add("Version","1.0")
	v.Add("UserID",UID)
	v.Add("Format","json")

	h.Write([]byte(v.Encode()))

	signature:=hex.EncodeToString(h.Sum(nil))
	v.Add("Signature",signature)
	return api.Params.Encode()
}

func (api *Req)BindJson(body interface{} )error{
	if res,err:=webtest.Get("https://api.sellercenter.lazada.com.my?"+api.Encode());err==nil{
		return res.BindJSON(&body)
		
	}else{
		return err
	}
}

func (api *Req)Get(dd  interface{}) (error){
	res,err:=webtest.Get("https://api.sellercenter.lazada.com.my?"+api.Encode())
	if err==nil{
		res.BindJSON(&dd)
	}
	return err
}

func (api *Req)Post(data []byte)(*webtest.Result,error){
	return webtest.Post(api.Url+"?"+api.Encode(),data)

}

func (api *Req)PostTest(data url.Values){
	res,_:=webtest.PostForm("http://gg?"+api.Encode(),data)
	fmt.Println(res.String())
}