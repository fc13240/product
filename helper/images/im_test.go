package images

import (
	"log"
	"testing"
)

func TestOpen(t *testing.T) {
	im, _ := Open("D:/web/demo/helper/uploads/product/20160721102505.png")
	log.Println(im.Scale("aaaa", 100, 100))
	//Create("D:/ss.jpg").Post("http://d.com/api/account.upheadimg")

}
