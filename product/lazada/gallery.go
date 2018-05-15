package lazada

import (
	"io"
	"helper/webtest"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"bytes"
	"helper/dbs"
	"helper/configs"

	"helper/util"
)

func Images(sku string)(images []Image){
	db:=dbs.Def()
	rows:=db.Rows("SELECT code,url FROM lazada_images WHERE sku=?",sku)

	images=[]Image{}
	for rows.Next(){
		var code,url string
		rows.Scan(&code,&url)
		images=append(images,Image{code,url})
	}
	return images
}

func SaveUploadImageInfo(sku string,rel string)(img *Image ,err error) {
	client:=webtest.NewClient()
	req,_:=http.NewRequest("GET",rel,nil)

	resp,err:=client.Do(req)

	if err!=nil{
		return
	}
	body:=&bytes.Buffer{}
	io.Copy(body,resp.Body)
	defer resp.Body.Close()

	img,err=UploadImage(body)
	if err!=nil{
		return
	}
	db:=dbs.Def()
	 _,err=db.Insert("lazada_images",configs.M{
		"code":img.Code,
		"url":img.Url,
		"sku":sku,
		"rel":rel,
		"datetime":util.Datetime(),
	 })

	 if err!=nil{
		 return
	 }
	return
}

func UploadImage(body io.Reader)(*Image ,error){
	v:=NewReq("UploadImage")

	client:=webtest.NewClient()

	req,err:=http.NewRequest("POST","https://api.sellercenter.lazada.com.my?"+v.Encode(),body)
	if err!=nil{
		return nil,err
	}
	resp,err:=client.Do(req)
	if err!=nil{
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
		ErrorResponse ErrorResponse
	}{}

	if err:=json.Unmarshal(resp_body,&data);err!=nil{
		return nil,err
	}

	if err:=data.ErrorResponse.Error();err!=nil{
		return nil,err
	}


	return &data.SuccessResponse.Body.Image,nil
}
