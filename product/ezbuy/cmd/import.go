package main
//ezbuy 导入命令，如果SKU不存在的情况下
import(
	"helper/text"
	"encoding/json"
	"fmt"
	"product/ezbuy"
	"helper/configs"
	"flag" 
func main(){
	var path,confg_path string
	var isUpdate bool
	var authorid int
	flag.StringVar(&confg_path, "c","","confing path")
	flag.BoolVar(&isUpdate,"update",false,"isUodate")
	flag.IntVar(&authorid,"authorid",0,"authorid")
	flag.Parse()
	if flag.Arg(0) == "h"{
		
		fmt.Printf(`
		导入命令，如果SKU已经存在的情况下不导入
		-c  confing path
		-update 只更新，根据sku.
		-authorid 指定作者
		`)
		return 
	}
	path=flag.Arg(0)

	configs.Ini(confg_path)

	text.ReadLine(path,func(line string){
		item:=ezbuy.Item{}
		json.Unmarshal([]byte(line),&item)
		old_authorid:=item.AuthorId
		if authorid >0 {
			item.AuthorId=authorid
		}
		if item.SKU!="" && isUpdate {
			ezbuy.ItemCol().Update(configs.M{"sku":item.SKU,"authorid":old_authorid},configs.M{"$set":item})
		}else if item.SKU!="" && ezbuy.SkuExist(item.SKU) == false {
			if err:=	ezbuy.ItemCol().Insert(item);err==nil{
				fmt.Println(item.SKU ," import succ")
			}
		}else{
			fmt.Println(item.SKU ," is exist")
		}
	})
}