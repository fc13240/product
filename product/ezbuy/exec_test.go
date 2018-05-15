package ezbuy

import (
	"github.com/tealeg/xlsx"
	"testing"
	"fmt"
	//"helper/configs"
	//"helper/configs"
	"helper/configs"
	"helper/dbs"
	_ "ainit"
)
//读取ez exec里的文档
func TestRead(t *testing.T){
	xlFile,_:=xlsx.OpenFile("D:/aaaa.xlsx")
	data:=map[int][]string{}
	for _, sheet := range xlFile.Sheets {

		for i, row := range sheet.Rows {
			for j, cell := range row.Cells {
				text,_ := cell.String()
				if text=="" && i==0{
					text,_=row.Cells[j+1].String()
				}
				c:=j/2
				data[c]=append(data[c],text)
			}
			fmt.Println("")
		}
	}
	//fmt.Println(data)
	var newData=map[int]map[int]string {}

	for _,d:=range data{
		j:=0
		row:=map[int]string{}
		id:=0
		for len(d[j:])>=2 {

			line:=d[j:j+2]

			if line[0] !="" {
				if line[0] == line[1]{
					id=configs.Int(line[0])

				}else{
					row_id:=configs.Int(line[0])
					row[row_id]=line[1]
				}
			}
			j=j+2
		}
		newData[id]=row
	}
	db:=dbs.Def()

	for cid,rows:=range newData{
		for att_id,att_label:=range rows{
			fmt.Println(configs.M{"cid":cid,"att_id":att_id,"att_name":att_label})
			db.Exec("INSERT INTO ezbuy.category_attr(cid,att_id,att_label)VALUES(?,?,?)",cid,att_id,att_label)
		}
	}
}

func TestExportAll(t *testing.T){
	
	ExportAll("D:/aaa.txt")
}
