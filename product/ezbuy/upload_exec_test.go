package ezbuy

import (
	_ "ainit"
	"testing"


	"helper/account"
	"fmt"
)

func TestUploadExec(t *testing.T) {
	fmt.Println(UploadExec(&account.Account{}))
}
