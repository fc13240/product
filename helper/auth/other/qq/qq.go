package qq

import (
	"fmt"
	"helper/webtest"
	"net/url"
	"regexp"
	"strings"
	"net/http"
	"io/ioutil"
	"helper/configs"

)

var (
	// AppID appID
	AppID ="100274858"
	AppKey ="06177675a23100f818baf661e132c9b4"
	UrlQQ      = "https://graph.qq.com"
	UrlQQOAuth = UrlQQ + "/oauth2.0"
	redirectUrl="http://api.51helper.com/login/callback/qq"
)

type QQ struct {
		Ret             int `json:"ret"`                // 返回码
		Msg             string `json:"msg"`                // 如果ret<0，会有相应的错误信息提示，返回数据全部用UTF-8编码。
		NickName        string `json:"nickname"`           // 用户在QQ空间的昵称。
		FigureURL       string `json:"figureurl"`          // 大小为30×30像素的QQ空间头像URL。
		FigureURL1      string `json:"figureurl_1"`        // 大小为50×50像素的QQ空间头像URL。
		FigureURL2      string `json:"figureurl_2"`        // 大小为100×100像素的QQ空间头像URL。
		FigureURLQQ1    string `json:"figureurl_qq_1"`     // 大小为40×40像素的QQ头像URL。
		FigureURLQQ2    string `json:"figureurl_qq_2"`     // 大小为100×100像素的QQ头像URL。需要注意，不是所有的用户都拥有QQ的100x100的头像，但40x40像素则是一定会有。
		Gender          string `json:"gender"`             // 性别。 如果获取不到则默认返回"男"
		ISYellowVip     string `json:"is_yellow_vip"`      // 标识用户是否为黄钻用户（0：不是；1：是）。
		Vip             string `json:"vip"`                // 标识用户是否为黄钻用户（0：不是；1：是）
		YelloVipLevel   string `json:"yellow_vip_level"`   // 黄钻等级
		Level           string `json:"level"`              // 黄钻等级
		IsYellowYearVip string `json:"is_yellow_year_vip"` // 标识是否为年费黄钻用户（0：不是； 1：是）
	}


// User user
func UserInfo(accessToken, openID string) (QQ, error) {
	result := QQ{}
	url := fmt.Sprintf("https://graph.qq.com/user/get_user_info?access_token=%s&oauth_consumer_key=%s&openid=%s",
		accessToken,
		AppID,
		openID,
	)
	res, err:=webtest.Get(url)


	if err:=res.BindJSON(&result);err!=nil{
		fmt.Println("解析 josn error ：",err)
	}

	result.NickName=strings.Trim(result.NickName," ")
	fmt.Println(result)
	if err != nil {
		return result, err
	}

	return result, err
}


func GetAuthorizationCodeUrl(state, scope string) string {
	v := url.Values{}
	v.Add("response_type", "code")
	v.Add("client_id", AppID)
	v.Add("redirect_uri", redirectUrl)
	v.Add("state", state)
	v.Add("scope", scope)
	return UrlQQOAuth + "/authorize?" + v.Encode()
}

type AccessToken struct{
	Value string
	ExpiresIn int
	RefreshToken string
}

func GetAccessToken(authCode string) (token AccessToken, err error) {
	v := url.Values{}
	v.Add("grant_type", "authorization_code")
	v.Add("client_id", AppID)
	v.Add("client_secret", AppKey)
	v.Add("code", authCode)
	v.Add("redirect_uri", redirectUrl)

	reqUrl := UrlQQOAuth + "/token?" + v.Encode()
	token=AccessToken{}

	if respContent, err := qqGet(reqUrl); err == nil {
		if values, err := url.ParseQuery(string(respContent)); err == nil {
			token.Value = values.Get("access_token")
			token.ExpiresIn = configs.Int(values.Get("expires_in"))
			token.RefreshToken = values.Get("refresh_token")
		}
	}
	return
}

func RefreshToken(appId, appKey, refreshToken string) (access_token, expires_in, refresh_token string, err error) {
	v := url.Values{}
	v.Add("grant_type", "refresh_token")
	v.Add("client_id", appId)
	v.Add("client_secret", appKey)
	v.Add("refresh_token", refreshToken)

	reqUrl := UrlQQOAuth + "/token?" + v.Encode()

	if respContent, err := qqGet(reqUrl); err == nil {
		if values, err := url.ParseQuery(string(respContent)); err == nil {
			access_token = values.Get("access_token")
			expires_in = values.Get("expires_in")
			refresh_token = values.Get("refresh_token")
		}
	}
	return
}

func GetOpenId(accessToken string) (string, error) {
	reqUrl := UrlQQOAuth + "/me?access_token=" + accessToken

	var err error
	var respContent []byte

	if respContent, err = qqGet(reqUrl); err == nil {
		if openId, err := extractDataByRegex(string(respContent), `"openid":"(.*?)"`); err == nil {
			return openId, nil
		}
	}

	return "", err
}

func extractDataByRegex(content, query string) (string, error) {
	rx := regexp.MustCompile(query)
	value := rx.FindStringSubmatch(content)

	if len(value) == 0 {
		return "", fmt.Errorf("正则表达式没有匹配到内容:(%s)", query)
	}

	return strings.TrimSpace(value[1]), nil
}

func qqGet(reqUrl string) ([]byte, error) {
	var err error
	if resp, err := http.Get(reqUrl); err == nil {
		defer resp.Body.Close()

		if content, err := ioutil.ReadAll(resp.Body); err == nil {
			//先测试返回的是否是ReturnError
			if values, err := url.ParseQuery(string(content)); err == nil {
				code := values.Get("code")
				msg := values.Get("msg")
				if len(code) > 0 && len(msg) > 0 {
					return nil, fmt.Errorf("Request %s failed with code %s. Error message is '%s'.",
						reqUrl, code, msg)
				}
			}

			return content, nil
		}
	}

	return nil, err
}