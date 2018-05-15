package ebay

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

func Call(str, name string, read ERead) {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, 5*time.Second)
				if err != nil {
					return nil, err
				}

				c.SetDeadline(time.Now().Add(7 * time.Second))
				return c, nil
			},
		},
	}
	xmlStr := `<?xml version="1.0" encoding="UTF-8" ?>` + str
	var account Account
	account.Init()

	req, rerr := http.NewRequest("POST", account.AppUrl, strings.NewReader(xmlStr))

	if rerr != nil {
		fmt.Println("Fatal error", rerr.Error())
	}
	req.Header.Add("Content-Type", "text/xml")
	req.Header.Add("X-EBAY-API-COMPATIBILITY-LEVEL", "737")

	req.Header.Add("X-EBAY-API-DEV-NAME", account.DevId)
	req.Header.Add("X-EBAY-API-APP-NAME", account.AppId)
	req.Header.Add("X-EBAY-API-CERT-NAME", account.CertId)
	req.Header.Add("X-EBAY-API-CALL-NAME", name)
	req.Header.Add("X-EBAY-API-SITEID", "0")

	response, err := client.Do(req)
	base := read.Inst()

	if err != nil {
		base.Ack = "Failure"
		base.Message = "请求超时失败"
		base.Log = err.Error()
		base.FailTotal += 1
		if base.FailTotal < 4 {
			log.Print("失败第:" + fmt.Sprintf("%d", base.FailTotal) + "继续请求中....")
			Call(str, name, read)
		}
	} else if response.StatusCode == 200 {
		bodyByte, _ := ioutil.ReadAll(response.Body)
		xml.Unmarshal(bodyByte, read)
		base.Body = string(bodyByte)
	} else {
		base.Ack = "Failure"
		base.Message = "请求失败,未知的错误"
		base.Log = "请求失败"
	}
	read = base
}
