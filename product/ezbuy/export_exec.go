package ezbuy

import (
	"helper/configs"

	"github.com/tealeg/xlsx"
	"strings"
	"helper/account"
	"errors"
	"os"
	"path/filepath"
	"time"
	"fmt"
)

var header =[]string{"店铺ID","店铺名称","CID","类目","商品英文名","商品中文名","商品描述","商品主图","是否敏感品","是否单品","货号","原价","售卖价","库存",
	"重量","长","宽","高","颜色图片","属性名ID","属性名1","属性值ID","属性值1","属性名ID","属性名2","属性值ID","属性值2","属性名ID","属性名3",
	"属性值ID","属性值3","属性名ID","属性名4","属性值ID","属性值4","属性名ID","属性名5","属性值ID","属性值5","属性名ID","属性名6","属性值ID","属性值6","属性名ID",
	"属性名7","属性值ID","属性值7","属性名ID","属性名8","属性值ID","属性值8","属性名ID","属性名9","属性值ID","属性值9","属性名ID","属性名10","属性值ID","属性值10"}

func ExportAll(path string ){
	c:=ItemCol()
	res:=[]Item{}
	c.Find(nil).All(&res)
	f,err:=os.Create(path)
	if err!=nil{
		fmt.Println(err)
	}
	defer f.Close()
	for _,row:=range res{
		f.WriteString(string(row.Encode())+"\n")
	}
}

func Export(author *account.Account,store_id int,items ...*Item)(filename string,err error) {
	f:=xlsx.NewFile()
	f.AddSheet("类目信息")
	sheet,_:=f.AddSheet("商品")
	sheet.AddRow()
	sheet.AddRow().WriteSlice(&header,-1)

	//获取店铺信息
	store:=GetSetting(author,store_id)
	if store.StoreId == 0{
		return filename,errors.New("store id is empty")
	}

	if store.StoreName == ""{
		return filename,errors.New("store name is empty")
	}

	for _,item:=range items{
		
		row:=sheet.AddRow()
		data:=[]interface{}{
			store.StoreId,
			store.StoreName,
			item.CID,
			item.CName,
			item.Name,
			item.CNName,
			item.Desc,
			strings.Join(item.Images,","),
			"N",
			"N",
			item.SKU,
			item.OldPrice,
			item.Price,
			item.Quant,
			item.Weight,
			item.Length,
			item.Width,
			item.Height,
			strings.Join(item.SkuImages,","),
		}
		row.WriteSlice(&data,-1)
		addRowSizes(row,item.Sizes,item.CID)
		addRowColors(row,item.Colors)
		addRowMaterials(row,item.Materials)
		if item.Styles.Value>0 {
			addRowStyles(row,item.Styles.Value)
		}
	}
	upload_dir:=configs.GetSection("ezbuy")["upload_dir"]

	curr_dir:=fmt.Sprint("ezitem","/",time.Now().Format("200612"))

	full_dir:=filepath.Join(upload_dir,curr_dir)

	if _,err:=os.Stat(full_dir);err !=nil{
		os.Mkdir(full_dir,777)
	}

	filename=items[0].SKU+".xlsx"
	err=f.Save(filepath.Join(full_dir,filename))

	filename=fmt.Sprint(curr_dir,"/",filename)

	if err!=nil{
		return filename,err
	}
	return filename,nil
}

func addRowColors(row *xlsx.Row,colors []Color){
	for _,color:=range colors{
		row.AddCell().SetInt(122535)
		row.AddCell().SetString("Color Name")
		row.AddCell().SetInt(color.ID)
		row.AddCell().SetString(color.Name)
	}
}

func addRowSizes(row *xlsx.Row,sizes []Size,cid int ){
	for _,size:=range sizes {

		if cid == 113{  //热裤分类
			row.AddCell().SetInt(58308)
			row.AddCell().SetString("尺码")
		}else{
			row.AddCell().SetInt(137236)
			row.AddCell().SetString("Size(Clothes)")
		}
		row.AddCell().SetInt(size.ID)
		row.AddCell().SetString(size.Name)
	}
}

func addRowMaterials(row *xlsx.Row,materials []int){
	for _,id:=range materials {
		row.AddCell().SetInt(137504)
		row.AddCell().SetString("材质（服饰）")
		row.AddCell().SetInt(id)
		row.AddCell().SetString(Materials[id])
	}
}

func addRowStyles(row *xlsx.Row,styles ...int){
	for _,id:=range styles {
		row.AddCell().SetInt(137526)
		row.AddCell().SetString("风格（女装）")
		row.AddCell().SetInt(id)
		row.AddCell().SetString(Styles[id])
	}
}


