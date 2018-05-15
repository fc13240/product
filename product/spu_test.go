package product 
import(
	_ "ainit"
	"testing"
	"sync"
	"fmt"
	
)

func TestNewSpu(t *testing.T){
	var wg sync.WaitGroup
	var n int = 100
	//var l sync.Locker
	var error_count=0
	wg.Add(n)
	
	for i:=range make([]int,n){
		go func(i int){
			if _,err:=NewSpu();err!=nil{
				fmt.Println(i,err.Error())
				error_count++
			}
			wg.Done()
		
		}(i)
	}
	wg.Wait()
	fmt.Println(error_count)
}