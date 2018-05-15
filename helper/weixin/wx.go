package weixin

import (
	"bytes"
	"comm"
	"crypto/sha1"
	"dbs"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"g"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

const (
	token                = "18126503782"
	appID                = "wxb23d99910cd582f3"
	appSecret            = "90937b4a4b6e9558c82960d64395ae0c"
	accessTokenFetchUrl  = "https://api.weixin.qq.com/cgi-bin/token"
	customServicePostUrl = "https://api.weixin.qq.com/cgi-bin/message/custom/send"
)

var accessToken string
var openID string = "ocsYks0tEp5W7IDpz2JZQI7Iy7C4"

type WX string

func (self WX) Init() {
	comm.Route("/weixin/token.hmtl", func(act comm.Action) {
		r := act.R()
		r.ParseForm()
		timestamp := strings.Join(r.Form["timestamp"], "")
		nonce := strings.Join(r.Form["nonce"], "")
		signatureGen := makeSignature(timestamp, nonce)

		signatureIn := strings.Join(r.Form["signature"], "")
		if signatureGen != signatureIn {
			return
		}
		echostr := strings.Join(r.Form["echostr"], "")
		fmt.Fprintf(act.W(), echostr)

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
			return
		}
		fmt.Println(string(body))
		requestBody := &TextRequestBody{}
		xml.Unmarshal(body, requestBody)
		fmt.Println(requestBody)

		msg := "你好"
		err = pushCustomMsg(accessToken, openID, msg)
		if err != nil {
			log.Println("Push custom service message err:", err)
			return
		}
	})

	comm.Route("/wx/gentoken", func(act comm.Action) {
		self.GetToken()
	})

	comm.Route("/wx/users", func(act comm.Action) {
		self.GetUsers()
	})

	comm.Route("/wx/send", func(act comm.Action) {
		err := pushCustomMsg(accessToken, openID, act.Get("msg"))
		if err != nil {
			log.Println("Push custom service message err:", err)
			return
		}
	})
}

func (self *WX) C() *dbs.Collection {
	return dbs.C("weixin")
}

func makeSignature(timestamp, nonce string) string {
	sl := []string{token, timestamp, nonce}
	sort.Strings(sl)
	s := sha1.New()
	io.WriteString(s, strings.Join(sl, ""))
	return fmt.Sprintf("%x", s.Sum(nil))
}

type TextRequestBody struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   time.Duration
	MsgType      string
	Content      string
	MsgId        int
}

func push(call_name, path string, param interface{}) error {

	body, err := json.MarshalIndent(param, " ", "  ")

	if err != nil {
		return err
	}

	postReq, err := http.NewRequest("POST",
		strings.Join([]string{"https://api.weixin.qq.com/cgi-bin/", call_name, "?access_token=", accessToken}, ""),
		bytes.NewReader(body))
	if err != nil {
		return err
	}

	postReq.Header.Set("Content-Type", "application/json; encoding=utf-8")

	client := &http.Client{}
	resp, err := client.Do(postReq)
	if err != nil {
		return err
	}
	resp.Body.Close()

	return nil
}

func get(url string, data interface{}) error {
	requestLine := Join("https://api.weixin.qq.com/cgi-bin/", url)

	resp, err := http.Get(requestLine)

	if err != nil || resp.StatusCode != http.StatusOK {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}
	fmt.Println(string(body))
	err = json.Unmarshal(body, data)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	var wx WX
	token := Token{}

	wx.G("token", &token)
	fmt.Println("init token", token.Access_token)
	accessToken = token.Access_token
}

type Token struct {
	Access_token string
	Expires_in   int
	Errcode      string
	Errmsg       string
	UpDate       string
}

func (self *WX) GetToken() {
	url := Join("token", "?grant_type=client_credential&appid=", appID, "&secret=", appSecret)

	data := Token{}

	if err := get(url, &data); err != nil {
		g.E("获取Token失败")
		fmt.Println(err)
	}

	if data.Errmsg == "" {
		accessToken = data.Access_token
		data.UpDate = time.Now().Format("2006-1-2 3:4:5")
		self.Set(data, "token")
	}

	fmt.Println("tooken:", data)
}

type Users struct {
	total       int      //关注该公众账号的总用户数
	count       int      //拉取的OPENID个数，最大值为10000
	data        []string //列表数据，OPENID的列表
	next_openid string   //拉取列表的后一个用户的OPENID

}

func (self *WX) GetUsers() {
	url := Join("user/get", "?access_token=", accessToken, "next_openid")

	users := Users{}

	if err := get(url, &users); err != nil {
		g.E("获取用户列表失败")
		fmt.Println(err)
		return
	}
	self.Set(users, "users")
}

func (self *WX) Set(data interface{}, t string) error {
	self.C().Upsert(dbs.Col{"type": t}, dbs.Col{"data": data, "type": t})
	return nil
}

func (self *WX) G(t string, data interface{}) error {

	item := struct{ Data interface{} }{}

	self.C().Find(dbs.Col{"type": t}).Select(dbs.Col{"data": 1, "_id": 0}).One(&item)

	v, _ := json.Marshal(item.Data)

	err := json.Unmarshal(v, &data)

	if err != nil {
		return err
	}

	return nil
}

func Join(s ...string) string {
	return strings.Join(s, "")
}

type CustomServiceMsg struct {
	ToUser  string         `json:"touser"`
	MsgType string         `json:"msgtype"`
	Text    TextMsgContent `json:"text"`
}

type TextMsgContent struct {
	Content string `json:"content"`
}

func pushCustomMsg(accessToken, toUser, msg string) error {
	csMsg := &CustomServiceMsg{
		ToUser:  toUser,
		MsgType: "text",
		Text:    TextMsgContent{Content: msg},
	}

	body, err := json.MarshalIndent(csMsg, " ", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(body))

	postReq, err := http.NewRequest("POST",
		strings.Join([]string{customServicePostUrl, "?access_token=", accessToken}, ""),
		bytes.NewReader(body))
	if err != nil {
		return err
	}

	postReq.Header.Set("Content-Type", "application/json; encoding=utf-8")

	client := &http.Client{}
	resp, err := client.Do(postReq)
	if err != nil {
		return err
	}
	resp.Body.Close()

	return nil
}
