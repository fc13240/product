package shopee

import (
	"helper/configs"

	"github.com/tealeg/xlsx"
	
	"helper/account"

	"os"
	"path/filepath"
	"time"
	"fmt"
	
)

func Export(author *account.Account,items ...*Item)(filename string,err error) {
//	f:=xlsx.NewFile()
	f,err:=xlsx.OpenFile("D:/Shopee.xlsm")
	if err!=nil{
		fmt.Println(err)
		return "",err
	}
	sheet:=f.Sheets[0]
	sheet.AddRow()
	
	sheet.AddRow().WriteSlice(&header,-1)

	for _,item:=range items{
		row:=sheet.AddRow()
		data:=NewRow(item)
		row.WriteSlice(&data,-1)
	}
	upload_dir:=configs.Get("download_dir")

	curr_dir:=fmt.Sprint("shopee","/",time.Now().Format("200612"))

	full_dir:=filepath.Join(upload_dir,curr_dir)

	if _,err:=os.Stat(full_dir);err !=nil{
		os.Mkdir(full_dir,777)
	}

	filename=time.Now().Format("150405")+".xlsx"
	err=f.Save(filepath.Join(full_dir,filename))

	filename=fmt.Sprint(curr_dir,"/",filename)

	if err!=nil{
		return filename,err
	}
	return filename,nil
}


