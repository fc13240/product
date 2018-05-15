package alibaba

import (
	"github.com/goquery"
	"encoding/json"
	"fmt"
	"helper/configs"
	"strings"
	"net/http"
	"helper/webtest"
	"path/filepath"
	//"os"
	"helper/account"
	"errors"
	"time"
	"mahonia"
	"log"
	"regexp"
	"net/url"
)
var (
	enc   = mahonia.NewDecoder("gbk")
)
func (p *Item) DownAllImage() {
	url := strings.Replace(p.BuyAddr, "/offer/", "/pic/", 1)

	doc, err := goquery.NewDocument(url)

	if err == nil {
		doc.Find("#dt-bp-tab-nav ul li").Each(func(i int, s *goquery.Selection) {
			src, _ := s.Attr("data-img")
			
			/*
			_, filename := filepath.Split(src)
			base_save := filepath.Join(p.BaseAddr, filename)
			save_path := filepath.Join(warehouse_base, base_save)
			images.Down(src, save_path)
			 */
			p.Images = append(p.Images, src)
		})
	}

	if len(p.Images) > 0 {
		Col().Update(configs.M{"id":p.Id},configs.M{"$set":configs.M{"images":p.Images}})
	}
}

//获取地址中的商品id
func GetUrlId(addr string )int64{
	uinfo,err:=url.ParseRequestURI(addr)
	if err!=nil{
		log.Println(err)
		return 0
	}
	switch uinfo.Host {
	case "item.taobao.com":
		val,_:=url.ParseQuery(addr)
		return configs.Int64(val.Get("id"))

	case "detail.1688.com":
		reg:=regexp.MustCompile("/offer/(\\d*)\\.html");
		info:=reg.FindStringSubmatch(addr)
		if len(info)>1{
			return configs.Int64(info[1])
		}
	}
	return 0
}


//1688下载图片
func Down(url string, author *account.Account) (item *Item,err error) {
	url=strings.Trim(url," ")
	url=strings.Trim(url,"\t")
	url=strings.Trim(url,"\n")
	id:=GetUrlId(url)
	if id == 0{
		return nil,errors.New("商品ID为0")
	}
	item=&Item{BuyAddr:url,Id:id,Source:"alibaba",Addtime:time.Now(),Authorid:author.Uid}

	if item.Exist(){
		return nil,errors.New("商品已经存在")
	}
	header:=http.Header{}
	header.Add("user-agent","Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36")
	header.Add("referer","https://sec.1688.com/query.htm?smApp=laputa&smPolicy=laputa-detail_1688_com_page-anti_Spider-html-checklogin&smCharset=GBK&smTag=MTEzLjg3LjE2My43NywsZDIzZTgxOTlmODM3NDY2MmJhZDlmMWE2Mzg4N2M0ZTU%3D&smReturn=https%3A%2F%2Fdetail.1688.com%2Foffer%2F552930363292.html&smSign=yJa%2FgiDv1ZUaLUn1295dwQ%3D%3D`)")
	header.Add("cookie",`JSESSIONID=N8yY12h-zD9YiWNUSvsD17Jsk7-6MClbRQ-6yf; __cn_logon__=false; cna=fyAPEim9vwICAXFXoS4OVYiY; ad_prefer="2017/08/07 18:08:36"; h_keys="%u8170%u9760"; ali_ab=113.87.161.46.1502100515845.7; alicnweb=touch_tb_at%3D1502100468058; UM_distinctid=15dbc2a4bdc170-0aa23e5f35ff7-474a0521-1fa400-15dbc2a4bddb68; CNZZDATA1253659577=1608202206-1502097831-https%253A%252F%252Fs.1688.com%252F%7C1502097831; _csrf_token=1502100523853; _tmp_ck_0="tlNt2b1nlHVeJJI2PLuRCRbSv2eNTRazPJl%2BG6tBRJVK5YfUp7GjBYBUQfviJWSirf%2BX%2FmrLsTG9l7hiNm50NiD8V%2FDptTtUlQecj7JqPgYVEpMreZVQJTKNp%2Bes0O2G9gKk5lCAdTQbyG5EGR0xo8WLbSh2fjmHX7hZNcxXW0stnv7vSvLVT41r9mULFSzPGnH4X6cCvVeFy%2FbA1qDkORAGdtR0ZHaQTIXAuhxs3M3l99MivPetrEcDNMRVJKfFO2jPbDb2oPDfIPcdZNxAGMPzXS%2F1xmyWG90KLweEXS9KWo8s12%2BaX8yqdvTAjeut3Tp8xQnidPEvagRBkE6CQhEX595uo6ZEd3qMhsiYH5wh2t4AeeyOD5piW1WB8fUm5kFw8DtfRxo%3D"; _uab_collina=150210047570662329073693; isg=Anl5FIXEBOD81tjDUhtMV6DliOWTLmxRtJ0Sv5uu9aAfIpm049Z9COdw0hEu; _umdata=0712F33290AB8A6D079C8720652EEA99AC181AFCD359DEBF2AA5E65CB7833CCB38920D71556939C2CD43AD3E795C914CD8143F55D2BC63EEA01723C7323A9C33`)
	res,err:=webtest.HGet(url,header)

	if err != nil {
		log.Println("res error",err)
		return item,err
	}
	//doc, err :=goquery.NewDocument(url)
	doc, err :=goquery.NewDocumentFromResponse(res.Resp)

	if err != nil {
		log.Println(url,err)
		return item,err
	}

	item.BuyPrice = 0.0
	item.Title = enc.ConvertString(doc.Find(".d-title").Text())

	//item.BuyPrice=configs.Float(doc.Find(".price-original-sku .value").Text())

	if desc_url,ok:=doc.Find("#desc-lazyload-container").Attr("data-tfs-url");ok{
		item.Descurl=desc_url
		if desc,err:=webtest.Get(item.Descurl);err==nil{

			content:=enc.ConvertString(desc.String())
			if ok,_:=regexp.MatchString("^var desc=",content);ok{
				content=strings.TrimLeft(content,"var desc=")
				content=strings.TrimRight(content,";")
				item.Desc=content
			}else{

				dd:=struct{
					Content  string `json:"content"`
				}{}

				ss:=strings.TrimLeft(content,"var offer_details=")
				ss=strings.TrimRight(ss,";")

				if err:=json.Unmarshal([]byte(ss),&dd);err!=nil{
					return item,errors.New("内容解析失败:"+err.Error())
				}
				item.Desc=dd.Content
			}
		}else{
			return item,errors.New("获取内容失败:"+err.Error())
		}
	}
	company:=doc.Find(".company-name")

	if company_name,ok:=company.Attr("title");ok{
		item.CompanyName=enc.ConvertString(company_name)
	}

	if company_url,ok:=company.Parent().Attr("href");ok{
		item.CompanyUrl=company_url
	}

	seller:=&Seller{CompanyName:item.CompanyName,CompanyUrl:item.CompanyUrl}
	item.SellerId=seller.Save()

	item.Save()

	item.DownAllImage()
	base_sku_path := filepath.Join(item.BaseAddr, "sku图")

	//sku_pic_path := filepath.Join(warehouse_base, base_sku_path)

	//os.Mkdir(sku_pic_path, 777)

	doc.Find(".desc-lazyload-container img").Each(func(i int, s *goquery.Selection) {
		//src, _ := s.Attr("src")
		//fmt.Println(src)
	})

	colors:=[]Color{}
	sizes:=[]Size{}

	doc.Find(".list-leading li").Each(func(i int, s *goquery.Selection) {
		attData, _ := s.Find("div").Attr("data-unit-config")
		attimgData, _ := s.Find("div").Attr("data-imgs")

		attName := struct {
			Name string `json:"name"`
		}{}

		attImg := struct {
			Preview  string `json:"preview"`
			Original string `json:"original"`
		}{}

		json.Unmarshal([]byte(enc.ConvertString(attData)), &attName)
		json.Unmarshal([]byte(attimgData), &attImg)

		if attImg.Preview == "" || attImg.Original == "" {
			if attName.Name!=""{
				colors=append(colors,Color{attName.Name})
			}
			return
		}

		skuImage := SkuImage{
			Name:          attName.Name,
			Preview:       attImg.Preview,
			Original:      attImg.Original,
			LocalOriginal: fmt.Sprint(base_sku_path, "/", attName.Name+"_大图", ".jpg"),
			LocalPreview:  fmt.Sprint(base_sku_path, "/", attName.Name+"_小图", ".jpg"),
		}

		//images.Down(attImg.Preview, filepath.Join(warehouse_base, skuImage.LocalPreview))
		//images.Down(attImg.Original, filepath.Join(warehouse_base, skuImage.LocalOriginal))

		item.SukImages = append(item.SukImages, skuImage)

	})

	doc.Find(".table-sku tr").Each(func(i int, s *goquery.Selection) {

		title, _ := s.Find(".name span").Attr("title")
		attimgData, _ := s.Find(".name span").Attr("data-imgs")


		if d,ok:=s.Attr("data-sku-config");ok{

			s:=struct{Name string `bson:"Name" json:"skuName"`}{}
			if err:=json.Unmarshal([]byte(enc.ConvertString(d)),&s);err==nil{
				sizes=append(sizes,Size{s.Name})
			}else{
				log.Println("获取尺寸失败",err.Error())
			}
		}
		attImg := struct {
			Preview  string `json:"preview"`
			Original string `json:"original"`
		}{}

		title = filertName(enc.ConvertString(title))
		json.Unmarshal([]byte(attimgData), &attImg)

		if attImg.Preview == "" || attImg.Original == "" {
			return
		}

		skuImage := SkuImage{
			Name:          title,
			Preview:       attImg.Preview,
			Original:      attImg.Original,
			LocalOriginal: fmt.Sprint(base_sku_path, "/", title+"_大图", ".jpg"),
			LocalPreview:  fmt.Sprint(base_sku_path, "/", title+"_小图", ".jpg"),
		}
		//images.Down(attImg.Preview, filepath.Join(warehouse_base, skuImage.LocalPreview))
		//images.Down(attImg.Original, filepath.Join(warehouse_base, skuImage.LocalOriginal))
		item.SukImages = append(item.SukImages, skuImage)
	})

	Col().Update(configs.M{"id":item.Id},configs.M{"$set":configs.M{"sukimages":item.SukImages,"colors":colors,"sizes":sizes}})

	return item,nil
}