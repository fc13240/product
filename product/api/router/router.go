package router

import (
	"github.com/gin-gonic/gin"
	"product/api"
)


type Comp struct{
	Method int
	Url string
	Func gin.HandlerFunc
}
const (
	GET =1
	POST=2
)
var RR=[]Comp{
	{GET,"sdfadf",api.Save},
	{GET,"sdfadf",api.Save},
	{GET,"sdfadf",api.Save},
	{GET,"sdfadf",api.Save},
	{GET,"sdfadf",api.Save},
	{GET,"sdfadf",api.Save},
	{GET,"sdfadf",api.Save},
}
