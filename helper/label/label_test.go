package label

import (
	_"ainit"
	"testing"
	"fmt"
)

func TestSearch(t *testing.T)  {
	labels:=Search("b",10)
	if len(labels) == 0{
		t.Fail()
		fmt.Println("fial")
	}else{
		fmt.Println("result num:", len(labels))
		fmt.Println(labels)
	}
}
