package ezbuy
import(
	"helper/text"
	"encoding/json"
	"fmt"
)
func Import(){
	text.ReadLine("D:/20180104本地mongodb.txt",func(line string){
		item:=Item{}
		json.Unmarshal([]byte(line),&item)
		fmt.Println(item)
	})
}