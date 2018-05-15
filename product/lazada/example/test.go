package main

import (
	"time"
	"helper/webtest"
	"crypto/sha256"
	"net/url"
	"fmt"
	"crypto/hmac"
	"encoding/hex"

	"helper/configs"
	"helper/dbs/mongodb"
	_"ainit"
	"mime/multipart"
	"bytes"
	"os"
	"io"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"product/lazada"
)

var (
	ApiUrl=lazada.ApiKey
	UID=lazada.UID
	ApiKey=lazada.ApiKey
)

func main()  {
	GetCategoryAttributes(1740)
}

func GetProducts(){
	data:=struct{
		SuccessResponse	struct{
			Head configs.M `json:"Head"`
			Body struct{
				TotalProducts int
				Products []configs.M
			}
		}
	}{}
	Get("GetProducts",&data)

	mdb:=mongodb.Conn()
	c:=mdb.C("lazada.items")

	for _,pro:=range data.SuccessResponse.Body.Products{
		c.Insert(pro)
	}
	fmt.Println(data)
}

func GetCategoryTree(){
	data:=struct{
		SuccessResponse	struct{
			Head configs.M `json:"Head"`
			Body []Category `json:"Body"`
		}
	}{}
	Get("GetCategoryTree",&data)
	mdb:=mongodb.Conn()
	c:=mdb.C("lazada.category")
	for _,pc:=range data.SuccessResponse.Body{
		c.Insert(pc)
	}
}

type SuccessResponse struct{
	Head configs.M `json:"Head"`
	Body []configs.M `json:"Body"`
}

type Category struct{
	Name string `json:"name"`
	Var bool `json:"var"`
	CategoryId uint64 `josn:"categoryId"`
	Leaf bool `json:"leaf"`
	Children []Category  `json:"children"`
}

func CommonUrlParam(action string,other map[string]string)*url.Values{
	h := hmac.New(sha256.New, []byte(ApiKey))
	v:=&url.Values{}

	v.Add("Timestamp",time.Now().Format(time.RFC3339))
	v.Add("Version","1.0")
	v.Add("UserID",UID)
	v.Add("Format","json")
	v.Add("Action",action)

	if other!=nil{
		for key,val:=range other{
			v.Add(key,val)
		}
	}

	h.Write([]byte(v.Encode()))

	signature:=hex.EncodeToString(h.Sum(nil))
	v.Add("Signature",signature)
	return v
}

func NewReq(action_name string)*Req{
	v:=&url.Values{}
	v.Add("Action",action_name)
	return &Req{Params: v}
}

type Req struct{
	Params *url.Values
}

func (a *Req)SetParam(k,v string ){
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
		res.BindJSON(&body)
		return nil
	}else{
		return err
	}
}

func (api *Req)Get() (*webtest.Result,error){

	return webtest.Get("https://api.sellercenter.lazada.com.my?"+api.Encode())

}


func Get(action string,body interface{} )error{

	v:=NewReq(action)

	if res,err:=webtest.Get("https://api.sellercenter.lazada.com.my?"+v.Encode());err==nil{

	res.BindJSON(&body)
	return nil
	}else{
		return err
	}
}

type Image struct {
	Code string
	Url string
}

func UploadImage()(*Image ,error){
	bodyBuf:=&bytes.Buffer{}
	bodyWriter:=multipart.NewWriter(bodyBuf)
	f,err:=os.Open("D:/4337387070_1887077446.jpg")
	if err!=nil{
		return nil,err
	}
	defer f.Close()

	io.Copy(bodyBuf,f)

	defer bodyWriter.Close()

	v:=CommonUrlParam("UploadImage",nil)

	client:=webtest.NewClient()

	req,err:=http.NewRequest("POST","https://api.sellercenter.lazada.com.my?"+v.Encode(),bodyBuf)
	if err!=nil{
		return nil,err
	}
	resp,err:=client.Do(req)
	if err!=nil{
		fmt.Println(err)
		return nil,err
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil,err
	}
	data:=struct{
		SuccessResponse struct{
			Body struct{
				Image Image
			}
		}
	}{}

	if err:=json.Unmarshal(resp_body,&data);err!=nil{
		return nil,err
	}
	return &data.SuccessResponse.Body.Image,nil
}

func GetCategoryAttributes(cid int){
	data:=struct{
		SuccessResponse struct{
			Body []configs.M
		}
	}{}
	req:=NewReq("GetCategoryAttributes")
	req.SetParam("PrimaryCategory",fmt.Sprint(cid))

	mdb:=mongodb.Conn()
	c:=mdb.C("lazada.categoryattrs")

	length_name:="( Fashion/Women/Clothing/Dresses )服饰/女士/服装/连衣裙"
	name:="连衣裙"

	if err:=req.BindJson(&data);err==nil{
		c.Insert(configs.M{"cid":cid,"name":name,"length_name":length_name,"attrs":data.SuccessResponse.Body})
	}

}