package down

import (
	"fmt"
	"helper/configs"
	"helper/dbs"
	"helper/dbs/mongodb"
	"helper/down"
	"testing"
)

var (
	url = "http://d32hv4kumpgy27.cloudfront.net/image/floryday/500_685/da/4b/da4bcdc5eb63d58d091a6cda9ff74c64.jpg"
)

func init() {
	configs.Ini("D:/config.ini")
	mongodb.Conn()
	dbs.Conn()

}

func TestDown(t *testing.T) {
	down.Down(url)
}

func TestGetContent(t *testing.T) {
	name := down.EncoryUrl(url)
	if body, err := down.GetBody(name); err != nil {
		t.Error(err)
	} else {

	}
}
