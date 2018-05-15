package html

import (
	"html/template"
	"strings"
	"fmt"
	"time"
	"helper/label"
)

var base string ="/"


func Html (x string) template.HTML{ return template.HTML(x) }

func SetBase(dir string){
	base=dir
}

func Js(files ...string) template.HTML{
	var ss []string
	for _,file:=range files{
		ss=append(ss,"<script src=\""+base+"scripts/"+file+"\"></script>")
	}
	return template.HTML(strings.Join(ss,"\n"))
}

func Css(files ...string) template.HTML{
	var ss []string
	for _,file:=range files{
		ss=append(ss,"<link rel=\"stylesheet\" href=\""+base+"styles/"+file+"\">")
	}
	return template.HTML(strings.Join(ss,"\n"))
}

func Url(url string)string{
	return url
}

func R(url string)string{
	if strings.Contains(url,"?"){
		return fmt.Sprint(url,"&r=",time.Now().Unix())
	}
	return fmt.Sprint(url,"?r=",time.Now().Unix())
}


func LabelsToStr(labels []label.Label)string{
	d:=[]string{}
	for _,l:=range labels{
		d=append(d,l.Name)
	}
	return strings.Join(d,",")
}

func Time(layout  string  ,tt time.Time) string {
	t:=time.Unix( tt.Unix(),0)
	return t.Format(layout)
}

var FuncMap=template.FuncMap{
	"html":Html,
	"css":Css,
	"js":Js,
	"url":Url,
	"r":R,
	"labeltoStr":LabelsToStr,
	"t":Time,
}