package account

import (
	"encoding/json"
	"helper/configs"
	"helper/dbs"
	"log"
	"time"
)

type MContent struct {
	Value string    `json:"value"`
	Event configs.M `json:"event"`
}

func (mc *MContent) String() string {
	return mc.Value
}

func (mc *MContent) Set(k string, v interface{}) *MContent {
	mc.Event[k] = v
	return mc
}

func (mc *MContent) EventToString() string {
	if data, err := json.Marshal(&mc.Event); err != nil {
		log.Println("解析event error", mc.Event)
		return ""
	} else {
		return string(data)
	}
}

func MessageContent(value string, event configs.M) *MContent {
	return &MContent{Value: value, Event: event}
}

type Message struct {
	Id      int `json:"id"`
	Content *MContent
	AddTime time.Time
	FromId  int
	ToId    int
}

func (m *Message) Del() error {
	return dbs.Exec("DELETE FROM message_to WHERE message_id=? AND to_mid=?", m.Id, m.ToId)
}

func (m *Message) Send(toids []int) error {
	stmt := dbs.Prepare("INSERT INTO message(content,event,From_mid,add_time)VALUES(?,?,?,?)")

	defer stmt.Close()

	add_time := time.Now().Format("2006-01-02 15:04:05")

	if re, err := stmt.Exec(m.Content.String(), m.Content.EventToString(), m.FromId, add_time); err == nil {
		message_id, err := re.LastInsertId()

		if err != nil {
			return err
		}

		for _, toid := range toids {
			dbs.Exec("INSERT INTO message_to(message_id,to_mid)VALUES(?,?)", message_id, toid)
		}

	} else {
		log.Println("add message failing .", err.Error())
		return err
	}
	return nil
}

//未读消息数量
func UnReadMessageCount(user *Account) int {
	count := 0
	dbs.One("SELECT COUNT(to_id) FROM message_to WHERE to_mid=? AND read_date IS NULL", user.Uid).Scan(&count)
	return count
}

func FromMessageList(user *Account) []Message {
	stmt := dbs.Rows("SELECT id,content,event,from_mid,add_time FROM message WHERE from_mid=?", user.Uid)
	defer stmt.Close()
	var messages []Message
	for stmt.Next() {
		msg := Message{}
		addtime, content, event := "", "", ""

		stmt.Scan(&msg.Id, &content, &event, &msg.FromId, &addtime)
		msg.AddTime, _ = time.Parse("2006-01-02 15:04:05", addtime)

		e := configs.M{}
		json.Unmarshal([]byte(event), &e)

		msg.Content = MessageContent(content, e)
		messages = append(messages, msg)
	}
	return messages
}

func ToMessageList(user *Account) []Message {
	stmt := dbs.Rows("SELECT id,content,event,from_mid,add_time,to_mid FROM message RIGHT JOIN message_to ON(message_to.message_id=message.id) WHERE message_to.to_mid=?", user.Uid)
	var messages []Message
	for stmt.Next() {
		msg := Message{}
		addtime, content, event := "", "", ""
		stmt.Scan(&msg.Id, &content, &event, &msg.FromId, &addtime, &msg.ToId)
		msg.AddTime, _ = time.Parse("2006-01-02 15:04:05", addtime)

		e := configs.M{}
		json.Unmarshal([]byte(event), &e)
		msg.Content = MessageContent(content, e)

		messages = append(messages, msg)
	}
	return messages
}
