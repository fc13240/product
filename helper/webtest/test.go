package webtest

import (
	"encoding/json"

	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"bytes"

	"regexp"
	"errors"
	"net"
	"time"
	"mime/multipart"
	"io"
	"os"
)

func Get(r string) (*Result, error) {
	resp, err := http.Get(r)
	res := &Result{}
	if err == nil {
		res.Resp = resp
	}
	return res, err

}

func HGet(r string, header http.Header) (*Result, error) {
	client := NewClient()
	res := &Result{}
	req, err := http.NewRequest("GET", r, nil)
	if err == nil {
		req.Header = header
		resp, err := client.Do(req)
		if err == nil {
			res.Resp = resp
		}
	}
	return res, err
}

func PostForm(r string, data url.Values) (*Result, error) {

	resp, err := http.PostForm(r, data)
	res := &Result{}
	if err == nil {
		res.Resp = resp

	}
	return res, err

}

func HPostForm(r string, header http.Header, data url.Values) (*Result, error) {
	client := NewClient()

	req, err := http.NewRequest("POST", r, strings.NewReader(data.Encode()))
	res := &Result{}
	if err == nil {
		req.Header = header
		resp, err := client.Do(req)

		if err == nil {
			res.Resp = resp

		}
	}
	return res, err
}

func Post(r string, data []byte) (*Result, error) {
	client := NewClient()

	req, err := http.NewRequest("POST", r, bytes.NewBuffer(data))
	res := &Result{}
	if err == nil {
		resp, err := client.Do(req)

		if err == nil {
			res.Resp = resp

		}
	}
	return res, err
}

type UpdateFile struct{
	FileName string 
	File io.Reader
	Path string 
}
func AddFile(filename  ,path string)*UpdateFile{
	file,err:=os.Open(path)
	if err!=nil{
		log.Panicln(err)
	}
	return &UpdateFile{filename,file,path}
}
func PostFile(r string, header http.Header,values url.Values,files ... *UpdateFile)(*Result, error) {
	
	b := &bytes.Buffer{}
	
	bw:= multipart.NewWriter(b)
	for _,f:=range files{
		fw, err := bw.CreateFormFile(f.FileName,f.Path)
		if err != nil {
			log.Println("error writing to buffer")
		}
		len,_:=io.Copy(fw,f.File)
		log.Println("len:",len/1024," KB")
	}
	
	for field,_:=range values {
		fw,_:=bw.CreateFormField(field)

		log.Println( field,":",values.Get(field) )
		fw.Write([]byte( values.Get(field) ))
	}
	
	header.Add("Content-Type",bw.FormDataContentType())

	bw.Close()

	client := NewClient()

	req, err := http.NewRequest("POST",r,b)
	res := &Result{}
	if err == nil {
		req.Header = header
		resp, err := client.Do(req)
		if err == nil {
			res.Resp = resp
		}
	}
	return res, err
}

func NewClient()*http.Client{
	return &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(30 * time.Second)
				c, err := net.DialTimeout(netw, addr, time.Second*20)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
		},

	}
}

func NewHeader()http.Header{
	return http.Header{}
}

func NewValue() url.Values{
	return url.Values{}
}

func PostJson(r string, header http.Header,body string) (*Result, error){
	client := NewClient()

	header.Add("Content-Type","application/json;charset=UTF-8")

	req, err := http.NewRequest("POST", r, bytes.NewBuffer([]byte(body)))
	res := &Result{}

	if err == nil {
		req.Header = header

		resp, err := client.Do(req)
		if err != nil {
			return nil,err
		}
		res.Resp = resp
		if resp.Status!="200 OK"{
			err=errors.New(resp.Status)
		}

	}

	return res, err

}

type Result struct {
	Resp *http.Response
}

func (r *Result) String() string {
	if body, err := ioutil.ReadAll(r.Resp.Body); err == nil {
		return string(body)
	} else {
		log.Println("解析失败", err)
	}
	return ""
}

func (r *Result) Body() []byte {
	if body, err := ioutil.ReadAll(r.Resp.Body); err == nil {
		return body
	} else {
		log.Println("解析失败", err)
	}
	return []byte("")
}



func (r *Result) ParseJson(data interface{}) error {

	decoder := json.NewDecoder(r.Resp.Body)
	if err := decoder.Decode(data); err != nil {
		return err
	}
	return nil
}

func (r *Result) BindJSON(data interface{}) error {

	decoder := json.NewDecoder(r.Resp.Body)
	if err := decoder.Decode(data); err != nil {
		return err
	}
	return nil
}

func (r *Result) Close() {
	defer func() {
		if p := recover(); p != nil {
			log.Println("关闭网络异常。。。。。",p)
		}
	}()
	r.Resp.Body.Close()
}

type Params struct {
	data  []string
	valus []string
	Value url.Values
}

func (p *Params) Set(k, v string) *Params {
	p.data = append(p.data, k+"="+v)
	p.valus = append(p.valus, v)
	if len(p.Value) == 0 {
		p.Value = url.Values{}
	}
	p.Value.Add(k, v)
	return p
}

func (p *Params) Join() string {
	return strings.Join(p.data, "&")
}

func (p *Params) Values() []string {
	return p.valus
}

func TrimHtml(body string )(out_body string) {

	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	body = re.ReplaceAllStringFunc(body, strings.ToLower)

	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	body = re.ReplaceAllString(body, "")

	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	body = re.ReplaceAllString(body, "")

	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	body = re.ReplaceAllString(body, "\n")

	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	body = re.ReplaceAllString(body, "\n")

	return strings.TrimSpace(body)
}


func Substr(s string, l int) string {
	if len(s) <= l {
		return s
	}
	ss, sl, rl, rs := "", 0, 0, []rune(s)
	for _, r := range rs {
		rint := int(r)
		if rint < 128 {
			rl = 1
		} else {
			rl = 2
		}

		if sl + rl > l {
			break
		}
		sl += rl
		ss += string(r)
	}
	return ss
}


func IsEmpty(s string, b bool) string {
	if b {
		return ""
	}
	return s
}