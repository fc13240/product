package crontab
import(
	"testing"
	"fmt"
	_ "ainit"
	"time"
	"helper/configs"

)

func TestSynEzbuyListing(t *testing.T){
	for page:=1;page<=100;page++{
	count:=0
		time.Sleep(time.Second*1)
		fmt.Println("start ",page)
	rows,err:=SynEzbuyListing(page)
	if err!=nil{
		t.Error(err.Error())
		return 
	}
	
	for _,row:=range rows{
		count++
		if row.Exist() == false{
			fmt.Println("ADD SUCC")
			row.Save()
		}else{
			fmt.Println("EXIST ",row.Name)
			row.Save()
		}
	}
	if count<40 {
		return 
	}
	}
	fmt.Println("OK")
}

func TestGetDetail(t *testing.T){
	GetDetail(7793710)
}

func TestNewProduct(t *testing.T){
	err:=SynToRemote(52265337)
	fmt.Println(err)
}

func TestDown(t *testing.T){
	Down(52138998)
}

func TestSynRemote(t *testing.T){
	SynToRemote(52138998)
}

func TestDowns(t *testing.T){
	rows,total:=Listing(configs.M{"categoryid":96},0,100)
	fmt.Println("Total:",total)
	for _,row:=range rows{
		time.Sleep(time.Second*1)
		if err:=Down(row.Pid);err==nil{
			fmt.Println(row.Name)
		}else{
			fmt.Println(err)
		}
	}
}