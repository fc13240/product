package wss
import (
	"helper/webtest"
	"product/alibaba"
	"product/ezbuy"
	"helper/util"
	"os"
	"fmt"
	"path/filepath"
	"time"
	"helper/images"

	"errors"
	"helper/configs"
)

type Detail struct{
	Colors []int
	Desc string
	Title string
	Sizes []int
	Id  string
}

func FillEZ(m *Message){
	if ezMessage == nil{
		m.Fail("ezbuy没有打开")
		return
	}
	data,err:=Parse(m.Cmd.Args.Int64("id"))

	if err!=nil{
		m.Fail(err.Error())
		return
	}
	ezMessage.Client.WriteJSON(configs.M{"isSucc":true, "data":data,"cmd":"fill"})
	m.Succ(data)
}

func Parse(id int64 )(detal *Detail ,err error) {
	if res, err := client.AlibabaItemGet(id); err == nil {

		data := struct {
			Item alibaba.Item `json:"item"`
		}{}

		if err := res.BindJSON(&data); err == nil {
			item := data.Item

			detal = &Detail{Title: item.Title, Desc: item.Desc}

			for _, color := range item.Colors {
				if color := ezbuy.GetColor(color.Name); color.ID > -1 {
					detal.Colors = append(detal.Colors, color.ID)
				}
			}

			for _, size := range item.Sizes {
				if size := ezbuy.GetSize(size.Name); size.ID>0 {
					detal.Sizes = append(detal.Sizes, size.ID)
				}
			}
			return detal,nil
		}else{
			return nil,err
		}
	}else{
		return nil,err
	}
	return nil,errors.New("undefined")
}

func Down(m *Message){
	
var res *webtest.Result
var err error
	id:=m.Cmd.Args.Int64("id")

	sku:=m.Cmd.Args.Get("sku")
	if id > 0{
		res,err=client.AlibabaItemGet(id)
	}else if(sku!=""){
		res,err=client.AlibabaItemGetBySku(sku)
	}else{
		m.Fail("失败")
		return 
	}

	if err==nil{
		data:=struct{
			Item alibaba.Item `json:"item"`
		}{}
	
		if err:=res.BindJSON(&data);err==nil{
			item:=data.Item
			save_path:="E:/items/"+time.Now().Format("01.02")+"/"+fmt.Sprint(item.Title)+"/"
			if len(item.Images)>0{
				
				if false == util.IsFolder(save_path){
					os.MkdirAll(save_path,777)
				}
				
				for i,src:=range item.Images{
					if err:=images.Down(src,fmt.Sprint(save_path,i,filepath.Ext(src)));err!=nil{
						m.Fail(err.Error())
					}
				}
			}

			if len(item.SukImages)>0{
				sku_path:=save_path+"sku图片/"
				if false == util.IsFolder(sku_path){
					os.MkdirAll(sku_path,0)
				}
				for _,sku:=range item.SukImages{
					if err:=images.Down(sku.Original,fmt.Sprint(sku_path,filepath.Ext(sku.Name)));err!=nil{
						m.Fail(err.Error())
					}
				}
			}
			if _,err:=client.AlibabaItemSet(id,"savepath",save_path);err!=nil{
				m.Fail(err.Error())
			}else{
				m.Succ(configs.M{"savepath":save_path,"succ":true})
			}
		}else{
			m.Fail(err.Error())
		}
	}else{
		m.Fail(err.Error())
	}
}

//更新订单列表
func UpdateOrder(m *Message){
	if res,err:=ez.GetOrders(0);err!=nil{

		data:=struct {
			Total int `json:"total"`
			Orders []ezbuy.Order `json:"data"`
		}{}

		if err:=res.BindJSON(&data);err!=nil{
			m.Fail(err.Error())
			return
		}
		ez.Client.SaveOrders(ezbuy.Encode(configs.M{"total":data.Total,"data":data.Orders}))
		m.Fail(err.Error())
	}else{
		m.Succ(res.String())
	}
}

func UpItems(m *Message){
	var i=0
	for {
		if items,err:=ez.GenList(i*40);err==nil{
			if len(items) <40 {
				m.SendTxt("更新商品完成")
				m.SendEnd()
				return
			}
			m.SendTxt(fmt.Sprint("更新第",(i+1)*40))
		}else{
			m.SendEnd("error:"+err.Error())
			return
		}
		i++
	}
}

//获取未编辑产品
func UserPruoductsFromSource(m *Message){
	if res,err:=ez.UserProductsFromSource(20);err==nil{
		data:=configs.M{}
		res.BindJSON(&data)
		m.Succ(data)
	}else{
		m.Fail(err.Error())
	}
}

//提交修改未编辑的产品
func OnSale(m *Message){
	pid:=m.Cmd.Args.GetInt("pid")
	err:=ez.UserUnCommitedProductDetail(pid)
	if err==nil{
		m.Succ(configs.M{})
	}else{
		m.Fail(err.Error())
	}
}