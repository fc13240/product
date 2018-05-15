package mail

import (
	"testing"
)

func TestSend(t *testing.T) {
	if err := Send("378590661@qq.com", "你好，刚", "<html><body><h1>验证码是:4476</h1></body></html>"); err != nil {
		t.Error(err.Error())
	}
}
