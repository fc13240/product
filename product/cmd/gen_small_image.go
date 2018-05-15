package main

import (
	_"ainit"
	"helper/configs"
	"helper/dbs/mongodb"
	"resize"
	"image/jpeg"
	"io"
	"fmt"
	"helper/dbs"
	"os"
	"bytes"

	"time"
)
var conn *mongodb.Mogno

func init(){
	conn=mongodb.Conn()
}

func main(){

	cc()

	old:=mongodb.NewGridFs("images")


	res:=mongodb.Result{}

	conn.C("images").Find(nil).All(&res.Items)

	data:=[]configs.M{}
	for res.Next(){
		c:=configs.M{}
		res.Scan(&c)
		data=append(data,c)

	}
	for _,v:=range data{

		head_image:=v.Get("head_image")

		if file,err:=old.Open(head_image);err==nil{
			Resize(file,head_image,240,240)
			file.Close()
		}

		listing:=v.Strings("listing")

		for _,filename:=range listing{
			if file,err:=old.Open(filename);err==nil{

				Resize(file,filename,340,340)
				Resize(file,"sm"+filename,50,70)
				file.Close()
				fmt.Println("OK")
			}
		}
	}
}
type Size struct{
	Flag string
	Width ,Height int
}

var(
	 sizes= []Size{
		 Size{"h",315,420},
		 Size{"sm",360,480},
		 Size{"sec",120,160},
		 Size{"ms",45,60},
	 }
)

func cc(){
	db:=dbs.Def()

	upload_dir:=configs.Get("upload_dir")
	for {
		rows:=db.Rows("SELECT rel_id,name FROM product_resource ORDER BY id DESC LIMIT 100")
		for rows.Next(){
			var id int
			var name string
			rows.Scan(&id,&name)
			filename:=fmt.Sprint(upload_dir,id,"/",name)
			if file,err:=os.Open(filename);err==nil{
				fmt.Println("OK ",filename)
				 b:=new(bytes.Buffer)
				 b.ReadFrom(file)
				for _,size:=range sizes{
					Resize(bytes.NewBuffer(b.Bytes()),fmt.Sprint(size.Flag,name),size.Width,size.Height)
				}
				file.Close()
				b.Reset()
			}else{
				fmt.Println(err)
			}
		}
		rows.Close()
		time.Sleep(time.Second *10)
	}
}

func Resize(oldfile io.Reader, filename string, width, hight int)  {

	var fs=mongodb.NewGridFs("image_small")
	if n,_:=fs.Find(configs.M{"filename":filename}).Count();n>0{
		return
	}
	if im,err:=jpeg.Decode(oldfile);err==nil{

		dst := resize.Thumbnail(uint(width), uint(hight), im, resize.NearestNeighbor)

		file,_:=fs.Create(filename)

		jpeg.Encode(file, dst, &jpeg.Options{100})

		defer file.Close()
	}else{
		fmt.Println(err)
	}
}