package lazada

import (
	"testing"
	_"ainit"
)

func TestSaveUploadImageInfo(t *testing.T) {
	_,err:=SaveUploadImageInfo("11","https://cbu01.alicdn.com/img/ibank/2017/460/478/4145874064_1921047407.jpg")
	if err!=nil{
		t.Error(err.Error())
	}
}

func TestImages(t *testing.T) {

	t.Logf("len %d",len(Images("11")))

}