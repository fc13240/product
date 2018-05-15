package account

import (
	"helper/configs"
	"helper/dbs"
	"testing"
)

func TestEdit(t *testing.T) {
	info := Find(1)
	if !info.EditField("nick", "刚") {
		t.Error("test")
	}
}

func TestSendMail(t *testing.T) {
	//user := Find(36)
	//SendActiveMail(user)
}


func TestActive(t *testing.T) {
	var code1, add_date string

	dbs.One("SELECT code,add_date FROM account_active where uid=? ORDER BY add_date DESC", 1).Scan(&code1, &add_date)

	if len(code1) == 0 {
		t.Error("code1 is empty")
	}

	if err := VerifyActive(1, code1); err != nil {
		t.Errorf("VerifyActive(%d,%s) fail "+err.Error(), 1, code1)
	}

	if err := VerifyActive(1, "a"); err == nil {
		t.Errorf("VerifyActive(%d,%s) fail ", 1, "a")
	}
}

func TestSendMessage(t *testing.T) {
	user := Find(1)
	to_member := 2
	testData := []struct {
		ToId    int
		Content *MContent
	}{
		{1, MessageContent("hi", configs.M{"title": "同意", "api": "group.applconfirm", "id": 1})},
		{2, MessageContent("hi1", configs.M{"title": "同意", "api": "group.applconfirm", "id": 2})},
	}

	for _, d := range testData {
		if err := user.SendMessage(d.Content, d.ToId); err != nil {
			t.Errorf(" user.SendMessage(%f, %s,%s)", to_member, d.Content)
		}
	}
}

func TestToMessageList(t *testing.T) {
	user := Find(1)
	messages := ToMessageList(user)
	if len(messages) == 0 {
		t.Error("我的收件箱是空的")
	}
	for _, message := range messages {
		if message.ToId != user.Uid {
			t.Errorf("message.ToId(%d)!=user.Uid(%d)", message.ToId, user.Uid)
		}
	}
}

func TestFromMessageList(t *testing.T) {
	user := Find(1)
	messages := FromMessageList(user)
	if len(messages) == 0 {
		t.Error("我的是空的")
	}
	for _, message := range messages {
		if message.FromId != user.Uid {
			t.Errorf("message.FromId(%d)!=user.Uid(%d)", message.FromId, user.Uid)
		}
	}
}
