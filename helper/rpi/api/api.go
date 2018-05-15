package api

import (
	"helper/net/igin"
)

func init() {
	r:=igin.R
	r.POST("/api/rpi.tank.register",Register)
}
