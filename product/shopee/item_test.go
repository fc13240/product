package shopee
import(
	_ "ainit"
	"testing"
	"fmt"
)
func TestCopy(t *testing.T){
	item,_:=Paste("MBL1")
	fmt.Println(item)
}