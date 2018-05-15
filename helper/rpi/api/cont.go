package api

import (

	"helper/rpi"
	"github.com/gin-gonic/gin"
	"helper/configs"
	"helper/net/igin"

)

func Register(c *gin.Context) {
	data:=configs.M{}
	gh:=igin.H(c)

	if err:=c.BindJSON(&data);err!=nil{
		gh.Fail(err.Error())
		return
	}
	tank := &rpi.Tank{
		Name: data.Get("name"),
		In1:  data.Get("in1"),
		In2:  data.Get("in2"),
		In3:  data.Get("in3"),
		In4:  data.Get("in4"),
		En1:  data.Get("en1"),
		En2:  data.Get("en2"),
	}

	if err := tank.Register(); err == nil {
		gh.Succ(nil)
	} else {
		gh.Fail(err.Error())
	}
}

func Edit(c *gin.Context) {
	data:=configs.M{}
	gh:=igin.H(c)

	if err:=c.BindJSON(&data);err!=nil{
		gh.Fail(err.Error())
		return
	}

	tank := &rpi.Tank{
		In1: data.Get("in1"),
		In2: data.Get("in2"),
		In3: data.Get("in3"),
		In4: data.Get("in4"),
		En1: data.Get("en1"),
		En2: data.Get("en2"),
	}

	if err := tank.Save(); err != nil {
		gh.Succ(nil)
	} else {
		gh.Fail(err.Error())
	}
}

func Get(c *gin.Context) {
	gh:=igin.H(c)
	tank.G

}
