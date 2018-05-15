package configs

import (
	"encoding/json"
	"fmt"
	"github.com/goconfig"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strconv"
)

type M map[string]interface{}


var Opt M
var config *goconfig.ConfigFile

func Ini(filename string) {
	var err error
	config, err = goconfig.LoadConfigFile(filename)
	if err != nil {
		panic(filename+" 配置文件不存在")
	}
}

func Get(key string) string {
	return GetSectionValue("global", key)
}

func GetSectionValue(section, key string) string {
	str, err := config.GetValue(section, key)
	if err != nil {
		log.Println(fmt.Sprintf("读取配置文件参数(%s)失败", key))
	}
	return str
}

func GetSection(section string) map[string]string {
	a, err := config.GetSection(section)
	if err != nil {
		log.Println(fmt.Sprintf("读取配置文件Section(%s)失败", section))
	}
	return a
}

func InitConfig(name string) {
	Opt = M{}
	if file, err := os.Open(name); err == nil {
		body, _ := ioutil.ReadAll(file)
		if err := json.Unmarshal(body, &Opt); err == nil {
		} else {
			fmt.Println("loading config error")
		}
	}
}

func (data M) String(k string) string {
	if str, ok := data[k].(string); ok {
		return str
	}
	return ""
}

func (data M) Get(k string) string {
	return data.String(k)
}

func (data M)Set(key string,v interface{})  {
	 data[key]=v
}

func (data M) Int(k string) int {
	return Int(data[k])
}

func (data M) Int64(k string) int64 {
	return Int64(data[k])
}

func (data M) GetInt(k string) int {
	return data.Int(k)
}

func (data M) Add(k string, v interface{}) {
	data[k] = v
}

func (data M) True(k string) bool {
	if b, ok := data[k].(bool); ok && b == true {
		return true
	}
	return false
}

func (data M) Val(k string) interface{} {
	return data[k]
}

func (data M) Ints(k string) []int {
	list := []int{}
	if v,ok:=data[k];ok{
		switch reflect.TypeOf(v).Kind() {
		case reflect.Slice:
			vs := reflect.ValueOf(v)
			for i := 0; i < vs.Len(); i++ {
				list = append(list, Int(vs.Index(i).Interface()))
			}
		}
	}
	return list
}

func (data M) Strings(k string) []string {
	return Strings(data[k])
}

func (data M) Float(k string) float32 {
	return Float(data[k])
}

func (m M) Encode() []byte {
	b, _ := json.Marshal(m)
	return b
}

func Int(vv interface{}) int {
	switch vv.(type) {
	case string:
		n, _ := strconv.Atoi(vv.(string))
		return n
	}

	if v, ok := vv.(float64); ok {
		return int(v)
	}

	if v, ok := vv.(int); ok {
		return v
	}

	return 0
}

func Int64(vv interface{}) int64 {
	switch vv.(type) {
	case string:
		n, _ := strconv.Atoi(vv.(string))
		return int64(n)
	}

	if v, ok := vv.(float64); ok {
		return int64(v)
	}

	if v, ok := vv.(int64); ok {
		return v
	}
	return 0
}

func Float(vv interface{}) float32 {

	if vv==""{
		return 0.0
	}
	if v, ok := vv.(float64); ok {
		return float32(v)
	}


	if v, ok := vv.(string); ok {

		if f, err := strconv.ParseFloat(v, 32); err == nil {
			return float32(f)
		} else {
			log.Println("err:", err)
		}
	}
	return 0.0
}

func Price(vv float32) string {
	return fmt.Sprintf("%0.2f", vv)
}


func Strings(v interface{}) []string {
	list := []string{}
	switch reflect.TypeOf(v).Kind() {
	case reflect.Slice:
		vs := reflect.ValueOf(v)
		for i := 0; i < vs.Len(); i++ {
			list = append(list, fmt.Sprint(vs.Index(i)))
		}
	}
	return list
}

