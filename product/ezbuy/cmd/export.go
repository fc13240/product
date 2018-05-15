package main
//ezbuy 导入命令，如果SKU不存在的情况下
import(
	
	
	"fmt"
	"product/ezbuy"
	"helper/configs"
	"flag" 
)

func main(){
	var confg_path string
	flag.StringVar(&confg_path, "c","","confing path")

	flag.Parse()
	
	
	if flag.Arg(0) == "h"{
		fmt.Printf(`
				导出所有ezbuy 产品命令，
				-c  	confing path
			`)
		return 
	}
	configs.Ini(confg_path)
	ezbuy.ExportAll(flag.Arg(0))
