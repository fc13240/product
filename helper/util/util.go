package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
	"time"
	"encoding/base64"
	"os"

	"bytes"

	"strconv"
)

func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func RemoveHtml(src string) string {
	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)

	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")

	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")

	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "\n")

	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")

	re, _ = regexp.Compile("&nbsp;")
	src = re.ReplaceAllString(src, "")

	return strings.TrimSpace(src)
}

func StringToNumber(s string) {

}

func NumberToString(n interface{}) string {
	return fmt.Sprintf("%d", n)
}

func IntJoin(a []int, sep string)string{
	ss:=[]string{}
	for _,v:=range a{
		ss=append(ss,fmt.Sprint(v))
	}
	return strings.Join(ss,sep)
}


func IntKJoin(a map[int]int, sep string)string{
	ss:=[]string{}
	for i,_:=range a{
		ss=append(ss,fmt.Sprint(i))
	}
	return strings.Join(ss,sep)
}

func Datetime() string {

	return time.Now().Format("2006-01-02 15:04:05")
}

func Time(value string) time.Time {
	if t, e := time.Parse("2006-01-02 15:04:05", value); e == nil {
		return t
	} else {
		return time.Now()
	}
}

func UrlEncode(url string )string{
	return base64.StdEncoding.EncodeToString([]byte(url))
}

func UrlUncode(url string)(string,error){
	if url,err:= base64.StdEncoding.DecodeString(url);err==nil{
		return string(url),nil
	}else{
		return "",err
	}
}

func U2S(str string) (to string, err error) {

	buf := bytes.NewBuffer(nil)
	i, j := 0, len(str)
	for i < j {
		x := i + 6
		if x > j {
			buf.WriteString(str[i:])
			break
		}
		if str[i] == '\\' && str[i+1] == 'u' {
			hex := str[i+2 : x]
			r, err := strconv.ParseUint(hex, 16, 64)
			if err == nil {
				buf.WriteRune(rune(r))
			} else {
				buf.WriteString(str[i:x])
			}
			i = x
		} else {
			buf.WriteByte(str[i])
			i++
		}
	}
	return buf.String(),nil
}

//是否是文件夹
func IsFolder(path string) bool{
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func ElapseTime(t time.Time) string {

	since := time.Since(t)
	hours := since.Hours()

	if month := int(hours / (24 * 30)); month > 0 {
		return fmt.Sprintf("%d%s", month, "个月前")
	}

	if day := int(hours / (24)); day > 0 {
		return fmt.Sprintf("%d%s", day, "天前")
	}

	if h := int(hours); h > 0 {
		return fmt.Sprintf("%d%s", h, "小时前")
	}

	if min := int(since.Minutes()); min > 1 {
		return fmt.Sprintf("%d%s", min, "分钟前")
	}

	return fmt.Sprintf("%d%s", int(since.Seconds()), "秒前")
}