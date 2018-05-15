package view

import (
	"bytes"
	"helper/configs"
	"helper/net"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type View struct {
	BasePath string
	path     string
	name     string
	act      *net.Act
	tem      *template.Template
	funs     template.FuncMap
}

func unescaped(x string) interface{} { return template.HTML(x) }

func New(act *net.Act) *View {
	view := &View{}
	view.act = act
	return view
}

func (self *View) Show(file string, data ...configs.M) {
	var err error

	self.name = file
	self.tem = template.New(filepath.Base(filepath.Join(self.BasePath,file)))

	helper := &Helper{View: self}

	self.tem.Funcs(template.FuncMap{
		"html":    unescaped,
		"url":     helper.Url,
		"js":      helper.Js,
		"css":     helper.Css,
		"include": helper.Include,
	})

	tmp, err := self.tem.ParseFiles(filepath.Join(self.BasePath,file))
	if err != nil {
		log.Println("解析出錯", err.Error())
	}
	//Files(self.GetFile())

	if len(data) > 0 {
		err = tmp.Execute(self.act.W, data[0])
	} else {
		err = tmp.Execute(self.act.W, nil)
	}
	if err != nil {
		log.Println(err)
	}
}

func (self *View) AddFun(name string, fun interface{}) {
	if len(self.funs) == 0 {
		self.funs = template.FuncMap{}
	}
	self.funs[name] = fun

}

type Helper struct {
	View *View
}

func (self *Helper) Include(names string) interface{} {
	//var err error
	var buff bytes.Buffer
	file, _ := os.Open(filepath.Join(self.View.BasePath,names))
	defer file.Close()
	b, _ := ioutil.ReadAll(file)

	funs := template.FuncMap{
		"html":    unescaped,
		"url":     self.Url,
		"js":      self.Js,
		"css":     self.Css,
		"include": self.Include,
	}
	t, err := template.New("").Funcs(funs).Parse(string(b))
	if err != nil {
		log.Println(err)
	}
	t.Execute(&buff, nil)

	return unescaped(buff.String())

}

func (self *Helper) Url(str ...string) string {
	if len(str) > 0 {
		return "/" + str[0]
	} else {
		return "/"
	}
}

func (self *Helper) Js(names ...template.HTML) template.HTML {
	var str template.HTML
	for _, v := range names {
		str += "<script type='text/javascript' src='/scripts/" + v + "'></script>"
	}
	return str
}

func (self *Helper) Css(names ...template.HTML) template.HTML {
	var str template.HTML
	for _, v := range names {
		return "<link rel='stylesheet' type='text/css' href='/styles/" + v + "'>"
	}
	return str
}
