package net

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"helper/configs"
	"log"
	"net/http"
	"strconv"

)

var (
	def *Router
)

func init() {
 	r:= mux.NewRouter()
	http.Handle("/", r)
	def=&Router{r}

}

func Host(host string) *Router {
	r:=def.r.Host(host).Subrouter()
	return &Router{r}
}

type Router struct {
	r *mux.Router
}

func Do(pattern string, handler func(act *Act), params ...string) {
	def.Do(pattern, handler, params...)
}

func (router *Router) Do(pattern string, handler func(act *Act), params ...string) {
	router.r.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Headers", "Content-type")
		w.Header().Add("Access-Control-Allow-Origin", "*")
		act := &Act{}
		act.W = w
		act.R = r
		isAuth := act.IsAuth()
		if len(params) > 0 {
			for _, val := range params {
				switch val {
				case "auth":
					if false == isAuth {
						act.Fail("没有登录", 10000)
						return
					}
				}
			}
		}
		handler(act)
	})
}

type Act struct {
	W    http.ResponseWriter
	R    *http.Request
	Auth Auth
}

func (act *Act) IsAuth() bool {
	return act.Auth.Verify(act.R)
}

//IsSucc输出成功JSON data_other 附加的输出
func (act *Act) Succ(data_other ...configs.M) {
	data := configs.M{"isSucc": 1}
	if len(data_other) > 0 {
		for key, value := range data_other[0] {
			data[key] = value
		}
	}
	act.WJson(data)
	return
}

//IsFail输出失败JSON data_other 附加的输出
func (act *Act) Fail(error_msg string, error_code ...int) {
	data := configs.M{"isSucc": 0, "error_msg": error_msg, "error_code": 0}
	if len(error_code) > 0 {
		data["error_code"] = error_code[0]
	}
	act.WJson(data)
	return
}

//WJson 输出JSON数据
func (act *Act) WJson(data interface{}) {
	json, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		act.Fail("解析错误")
	} else {
		act.W.Header().Set("content-type", "application/json")
		act.W.Write(json)
	}
}


func ResDir(url, path string) {
	def.r.PathPrefix(url).Handler(http.StripPrefix(url, http.FileServer(http.Dir(path))))
}

func (act *Act) GetInt(param string) int {
	def_num := 0
	if s := act.Get(param); len(s) > 0 {
		if n, e := strconv.Atoi(s); e == nil {
			return n
		}
		return def_num
	}
	return def_num
}

func (act *Act) ParseJson() (*configs.M, error) {
	data := &configs.M{}
	decoder := json.NewDecoder(act.R.Body)
	if err := decoder.Decode(data); err != nil {
		return data, err
	}
	return data, nil
}

func (act *Act) BindJson(vv interface{}) error{
	decoder := json.NewDecoder(act.R.Body)
	if err := decoder.Decode(vv); err != nil {
		return  err
	}
	return  nil
}

func (act *Act) JSON(code int,data interface{}){
	json, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		act.Fail("解析错误")
	} else {
		act.W.Header().Set("content-type", "application/json")
		act.W.Write(json)
	}
}

func (act *Act) Get(param string) string {
	return act.R.URL.Query().Get(param)
}

func (act *Act) SetCookie(key, value string) {
	http.SetCookie(act.W, &http.Cookie{Name: key, Value: value, Path: "/"})
}

func (act *Act) Vars() map[string]string {
	return mux.Vars(act.R)
}
