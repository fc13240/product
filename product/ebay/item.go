package ebay

import (
	"encoding/xml"
	"log"
	//"strings"
	"time"
)

type RItemBase struct {
	Id         string `bson:"_id"`
	ItemID     string
	Quantity   string
	Title      string
	SKU        string
	Url        string `xml:"ListingDetails>ViewItemURL"`
	StartTime  Xtime  `xml:"ListingDetails>StartTime"`
	Price      string `xml:"SellingStatus>CurrentPrice"`
	Update     string
	WatchCount string
	Gallery    string `xml:"PictureDetails>GalleryURL"`
}

func (self *RItemBase) DowDetail() {
	/*
		d := ItemDetail{ItemID: self.ItemID}
		data := d.Dow()
		dbs.C("ebay.items").Update(dbs.Col{"sku": self.SKU}, dbs.Col{"$set": dbs.Col{"endtime": data.EndTime, "quantitysold": data.QuantitySold, "starttime": data.StartTime}})
	*/
}

type Xtime struct {
	time.Time
}

func (x *Xtime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	d.DecodeElement(&v, &start)
	if pase, err := time.Parse(time.RFC3339, v); err == nil {
		*x = Xtime{pase}
		return nil
	} else {
		return err
	}
}

type RItem struct {
	Read
	ItemID         string
	Quantity       string
	Title          string
	Url            string `xml:"Item>ListingDetails>ViewItemURL"`
	StartTime      Xtime  `xml:"Item>ListingDetails>StartTime"`
	EndTime        Xtime  `xml:"Item>ListingDetails>EndTime"`
	QuantitySold   string `xml:"Item>SellingStatus>QuantitySold"`
	Price          string `xml:"StartPrice"`
	Currency       string `xml:"Currency"`
	Description    string `xml:"Item>Description"`
	Duration       string `xml:"ListingDuration"`
	Type           string `xml:"ListingType"`
	Location       string `xml:"Location"`
	PaymentMethods string `xml:"PayMentMethods"`
	PayPal         string `xml:"PayPalEmailAddress"`
	Site           string `xml:"Site"`
}

type RItems struct {
	Read
	Items []RItemBase `xml:"ActiveList>ItemArray>Item"`
}

type Items struct {
	XMLName xml.Name `xml:"GetMyeBaySelling"`
	Ebay
	Include string `xml:"ActiveList>Include"`
	Notes   string `xml:"ActiveList>IncludeNotes"`
	//Type    string `xml:"ActiveList>ListingType"`
	PStart int `xml:"ActiveList>Pagination>EntriesPerPage"`
	Number int `xml:"ActiveList>Pagination>PageNumber"`
}

func (self *Items) Dow(f func(item RItemBase)) {

	self.Include = "true"
	self.Notes = "true"
	//r.Type = "Auction,FixedPriceItem,StoresFixedPrice"
	self.PStart = 30
	self.Number = 1

	rs := RItems{}

	Query(self, "GetMyeBaySelling", &rs)

	log.Print(rs.Body)

	if rs.IsSuccess() {
		for _, item := range rs.Items {
			f(item)
			log.Print("下载 Ebay " + item.Title + " Success")
		}
	} else {
		log.Print("下载Ebay 数据错误")
	}
}

func (self *RItemBase) Save() {
	/*
		self.Update = time.Now().Format("2006-01-02 15:04:05")
		if false == self.Exist() {
			self.Id = dbs.NewId()
			dbs.C("ebay.items").Insert(self)
		} else {
			dbs.C("ebay.items").Update(dbs.Col{"sku": self.SKU}, dbs.Col{"$set": self})
		}
	*/
}

func (self *RItemBase) Exist() bool {
	/*
		if v, _ := dbs.C("ebay.items").Find(dbs.Col{"sku": self.SKU}).Count(); v > 0 {
			return true
		}
		return false
	*/
	return false
}

type ItemDetail struct {
	XMLName xml.Name `xml:"GetItem"`
	Ebay
	ItemID      string `xml:"ItemID"`
	DetailLevel string `xml:"DetailLevel"`
}

func (self *ItemDetail) Dow() RItem {
	self.DetailLevel = "ItemReturnDescription"
	body := RItem{}
	Query(self, "GetItem", &body)

	//endtime := strings.Replace(body.EndTime, "T", " ", 1)
	//body.EndTime = strings.Replace(endtime, ".000Z", "", 1)

	//stime := strings.Replace(body.StartTime, "T", " ", 1)
	//body.StartTime = strings.Replace(stime, ".000Z", "", 1)

	return body
}

type Categorie struct {
	Id               string `xml:"CategoryID"`
	Level            string `xml:"CategoryLevel"`
	Name             string `xml:"CategoryName"`
	ParentID         string `xml:"CategoryParentID"`
	BestOfferEnabled string `xml:"BestOfferEnabled"`
	AutoPayEnabled   string `xml:"AutoPayEnabled"`
}
type A struct {
	Ebay
	XMLName xml.Name `xml:"GetCategoriesRequest"`

	CategorySiteID int
	DetailLevel    string
	//LevelLimit     int
}

func (self *A) Q() {
	/*
		list := struct {
			Read
			CategoryArray []Categorie `xml:"CategoryArray>Category"`
		}{}
		self.CategorySiteID = 0
		self.DetailLevel = "ReturnAll"
		//self.LevelLimit = 1
		Query(self, "GetCategories", &list)
		i := 0
		if list.IsSuccess() {
			for _, item := range list.CategoryArray {
				dbs.C("ebay.cate").Insert(item)
				i++
			}
			log.Print("Total", i)
		} else {
			log.Print("Faling")
		}
	*/
}

func GetCategories() {
	req := A{}
	req.Q()
	//return list.CategoryArray

}
