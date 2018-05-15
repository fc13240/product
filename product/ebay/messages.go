package ebay

import (
	"encoding/xml"
	"log"
	"time"
)

type RMessages struct {
	Read
	Items []RMessageItem `xml:"Messages>Message"`
}

type RMessageItem struct {
	Sender          string
	RecipientUserID string
	SendToName      string
	Subject         string
	MessageId       string `xml:"MessageID"`
	Text            string
}

func (self *RMessages) Dow() {

	type Req struct {
		XMLName     xml.Name `xml:"GetMyMessages"`
		DetailLevel string
		Ebay
		StartTime string
	}

	type ReqDetail struct {
		Req
		Ebay
		MessageID []string `xml:"MessageIDs>MessageID"`
	}
	startTime := time.Now().Add(-3600 * 48 * time.Second).String()
	Q(&Req{DetailLevel: "ReturnHeaders", StartTime: startTime}, "GetMyMessages", self)

	if self.IsSuccess() {

		var ids []string

		for _, item := range self.Items {
			if false == item.Exist() {
				ids = append(ids, item.MessageId)
			}
		}

		Q(&ReqDetail{Req: Req{DetailLevel: "ReturnMessages", StartTime: startTime}, MessageID: ids}, "GetMyMessages", self)

		if self.IsSuccess() {
			for _, item := range self.Items {
				item.Save()
			}
			log.Print("Message 更新成功")
		} else {
			log.Print("下载Messages详细出错")
		}
	} else {
		log.Print("下载Messages列表出错" + self.Message)
	}
}

func (self *RMessageItem) Exist() bool {
	return false
}

func (self *RMessageItem) Save() {}
