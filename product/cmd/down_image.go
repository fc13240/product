package main

import (
	_ "ainit"
	"helper/redisCli"
	"github.com/garyburd/redigo/redis"
	"encoding/json"
	"net/http"
	"fmt"
	"helper/util"
	"helper/dbs"
	"helper/dbs/mongodb"
	"helper/configs"
	"time"
	"image/jpeg"
	"resize"

)

//download the image from a url save to mongodb and update product_resource table.


func init(){
	mongodb.Conn()

}

type Resource struct{
	Id   int `json:"id"`
	Url string `json:"url"`
}


//检查
func main(){
	r:=redisCli.Conn()
	for {
		str,err:=redis.String(r.Do("lpop", "wait_down_images:"))
		if err!=nil{
			time.Sleep(time.Second*3)
		}
		res:=Resource{}
		json.Unmarshal([]byte(str),&res)
		fmt.Println(res)
		if err:=SaveByUrl(res);err == nil{
			fmt.Println("OK")
		}else{
			fmt.Println(
				err.Error())
		}
	}
}

func SaveByUrl(re Resource) error{

	var store =mongodb.NewGridFs("image_small")
	url:=re.Url

	res, err :=http.Get(url)

	if err != nil {
		return err
	}

	defer res.Body.Close()
	fmt.Println(res.Header)
	name := util.Md5(url)

	file, err :=store.Create(name)

	defer file.Close()

	if err != nil {
		return err
	}

	var fs=mongodb.NewGridFs("image_small")

	if n,_:=fs.Find(configs.M{"filename":name}).Count();n>0{
		dbs.Update("product_resource",configs.M{"name": name},"id=?",re.Id)
		return nil
	}


	if im,err:=jpeg.Decode(res.Body);err==nil{
		dst := resize.Thumbnail(240, 240, im, resize.NearestNeighbor)

		file,_:=fs.Create(name)

		if err:=jpeg.Encode(file,dst,&jpeg.Options{100});err!=nil{
			return err
		}


		defer file.Close()
	}else{
		return err
	}

	dbs.Update("product_resource",configs.M{"name": name},"id=?",re.Id)
	return nil
}
