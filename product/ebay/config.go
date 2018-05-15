package ebay

import (
	"encoding/xml"
	"fmt"
)

type ECall interface {
	Call(str string)
	Get()
}

type ERead interface {
	Inst() *Read
}

type E interface {
	Init()
}

type Read struct {
	Ack       string
	Message   string `xml:"Errors>ShortMessage"`
	Log       string `xml:"Errors>LongMessage"`
	Body      string
	FailTotal int
}

func (self *Read) Inst() *Read {
	return self
}

type Ebay struct {
	Lns   string `xml:"xmlns,attr"`
	Token string `xml:"RequesterCredentials>eBayAuthToken"`
}

func (self *Ebay) Init() {
	var account Account
	account.Init()
	self.Lns = "urn:ebay:apis:eBLBaseComponents"
	self.Token = account.Token
}

func R(e E, name string) {
	output, _ := xml.MarshalIndent(e, " ", "    ")
	fmt.Println(string(output))

}

func (self *Read) IsSuccess() bool {
	if self.Ack != "Failure" {
		return true
	} else {
		return false
	}
}

func Query(e E, name string, body ERead) {
	Q(e, name, body)
}

func Q(e E, name string, body ERead) {
	e.Init()
	output, _ := xml.MarshalIndent(e, " ", "    ")
	fmt.Println(string(output))
	Call(string(output), name, body)
}

func Up(e E, name string) *Read {
	e.Init()
	output, _ := xml.MarshalIndent(e, " ", "    ")
	fmt.Println(string(output))
	read := &Read{}
	Call(string(output), name, read)

	return read
}

//GetStoreRequest
type GetStoreRequest struct {
	Ebay
	Level int `xml:"LevelLimit"`
}

func (r *GetStoreRequest) C() {
	r.Level = 1
	R(r, "GetStore")
}

func (e Ebay) P() {
	output, _ := xml.MarshalIndent(e, " ", "    ")
	fmt.Println(string(output))
}

type GeteBayDetailsRequest struct {
	Ebay
}

func (r *GeteBayDetailsRequest) C() {
	R(r, "GeteBayDetails")
}

//GetMyeBaySellingRequest
type GetMyeBaySellingRequest struct {
	Ebay
	Sort   string `xml:"ActiveList>Sort"`
	PStart string `xml:"ActiveList>Pagination>EntriesPerPage"`
	Number string `xml:"ActiveList>Pagination>PageNumber"`
	Finds  string `xml:"OutputSelector"`
}

func (r *GetMyeBaySellingRequest) C() {
	r.Sort = "TimeLeft"
	r.PStart = "15"
	r.Number = "1"
	r.Finds = "ActiveList.ItemArray.Item.ItemID"
	R(r, "GetMyeBaySelling")
}
