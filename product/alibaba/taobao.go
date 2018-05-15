package alibaba


import (
	"helper/account"
	"helper/configs"
	"github.com/goquery"
	"net/url"
	"helper/util"
	"fmt"
	"regexp"
	"strings"
	"errors"
)

type Taobao struct{

}

func (*Taobao)Down(addr string,author *account.Account)(item *Item,err error){

	addr=strings.Trim(addr," ")
	addr=strings.Trim(addr,"\t")
	addr=strings.Trim(addr,"\n")
	defer func(){
		if r:=recover();r!=nil{
			err=errors.New(fmt.Sprint("下载淘宝商品异常",addr ,"error:",r))
		}
	}()

	item=&Item{Source:"taobao",BuyAddr:addr,Authorid:author.Uid}

	if item.Exist(){
		return nil,errors.New("商品已经存在")
	}

	val,_:=url.ParseQuery(addr)

	item.Id = configs.Int64(val.Get("id"))

	doc,err:=goquery.NewDocument(addr)

	if err!=nil{
		return nil,err
	}

	text:=doc.Find("script").First().Text()
	text=enc.ConvertString(text)


	shopname_reg:=regexp.MustCompile("shopName\\s*:\\s'(.*)'")

	if s:=shopname_reg.FindStringSubmatch(text);len(s)>1{
		item.CompanyName,_=util.U2S(s[1])
	}

	reg:=regexp.MustCompile("url\\s:\\s'(.*)\\.taobao\\.com/'")

	if s:=reg.FindStringSubmatch(text);len(s)>1{
		item.CompanyUrl=fmt.Sprint("https:",s[1],".taobao.com")
	}


	title_reg:=regexp.MustCompile("title\\s*:\\s'(.*)'")

	if s:=title_reg.FindStringSubmatch(text);len(s)>1{

		item.Title,_=util.U2S(s[1])
		fmt.Println(item.Title)
	}

	images_reg:=regexp.MustCompile("auctionImages\\s*:\\s\\[(.*)\\]")

	if s:=images_reg.FindStringSubmatch(text);len(s)>1{

		s:=strings.Replace(s[1],`"`,"",-1)
		item.Images=strings.Split(s,",")
	}

	desc,_:=doc.Find(".tb-item-info").Html()


	seller:=&Seller{CompanyName:item.CompanyName,CompanyUrl:item.CompanyUrl}
	item.SellerId=seller.Save()

	item.Desc=enc.ConvertString(desc)

	item.Save()
	return item,err
}