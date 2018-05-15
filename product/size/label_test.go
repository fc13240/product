package size 
import(
	_"ainit"
	"testing"
	"fmt"
)

func TestSave(t *testing.T){
	size:=Size{}
	size.ItemId=100
	ss:=GetLabels()
	
	size.Head=fmt.Sprintf("%s,%s,%s,%s",ss[0],ss[1],ss[2],ss[3])
	size.Rows=append(size.Rows,"22,33,44,55","11,22,33,11")
	fmt.Println(Save(size))
}

func TestGet(t *testing.T){
	size,_:=Get(100)
	fmt.Println(size.Head)
	for _,row:=range size.Rows{
		fmt.Println(row)
	}

}