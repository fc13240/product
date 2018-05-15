package filestore

import (
	"fmt"
	"helper/filestore"
	"os"
	"testing"
)

func TestSave(t *testing.T) {
	if file, err := os.Open("D:/z.txt"); err == nil {

		defer file.Close()
		f, _ := file.Stat()
		fmt.Println("SIZE:", f.Size())
		if name, err := filestore.Save(file); err == nil {
			fmt.Println("okFFF", name)
			t.Log("OK", name)
		} else {
			t.Error(err)
		}
	} else {
		t.Error(err)
	}

}
