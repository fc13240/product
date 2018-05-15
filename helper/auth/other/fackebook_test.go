package other

import (
	_"ainit"
	"testing"
	"fmt"
)

func TestFaceBook_Login(t *testing.T) {
	faccbook:=FaceBook{Token:"123456dfasdf",ExpiresIn:3600,UserID:15678564,Nick:"XUEGANG LI"}
	fmt.Println(faccbook.Login())
	t.Fail()
}