package warehouse

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/goquery"
	"helper/account"
	"helper/configs"
	"helper/crypto"
	"helper/dbs"
	"helper/images"
	"helper/util"
	"mahonia"
	"os"
	"path/filepath"
	"strings"

)

var (
	warehouse_base = "./"
	enc            = mahonia.NewDecoder("gbk")
)



type ProductWarehouse struct {
	Id        int `json:"id"`
	Author    *account.Account
	BuyPrice  float32
	Title     string            `json:"title"`
	BuyAddr   string            `json:"buyAddr"`
	SaveAddr  string            `json:"saveAddr"`
	BaseAddr  string            `json:"baseAddr"`
	Addtime   string            `json:"addtime"`
	SukImages []ProductSkuImage `json:"sukImages"`
	Images    []string          `json:"images"`
	Md5       string
	Desc      string 	    `json:"desc"`
	Descurl   string            `json:"desurl"`
}

type ProductSkuImage struct {
	Name          string `json:"name"`
	Original      string `json:"original"`
	Preview       string `json:"preview"`
	LocalOriginal string `json:"localOriginal"`
	LocalPreview  string `json:"localPreview"`
}

func GetWarehouseProductInfo(item_id int) (*ProductWarehouse, error) {
	var err error
	pro := &ProductWarehouse{SaveAddr: warehouse_base}
	item, e := Get(item_id)

	if e != nil {
		err = e
	}

	if item.WarehouseId == 0 {
		err = errors.New("仓库ID不存在")
	}

	if err != nil {
		return nil, err
	}

	var sku_images_str, images_str string

	dbs.One("select id,buy_addr,sku_imags,images,addtime FROM product_warehouse WHERE id=?", item.WarehouseId).Scan(&pro.Id, &pro.BuyAddr, &sku_images_str, &images_str, &pro.Addtime)

	if e := json.Unmarshal([]byte(sku_images_str), &pro.SukImages); e != nil {
		return pro, e
	}

	pro.Images = strings.Split(images_str, ",")
	return pro, nil
}

func SearchWarehouse() {

}

func (p *ProductWarehouse) New() (err error) {
	md5 := crypto.Md5(p.BuyAddr)

	dbs.One("select id FROM product_warehouse WHERE md5=?", md5).Scan(&p.Id)

	createFolder := func() {
		p.BaseAddr = fmt.Sprint(p.Id)
		os.Mkdir(filepath.Join(warehouse_base, p.BaseAddr), 777)
	}

	if p.Id > 0 {
		createFolder()
		return nil
	}

	p.Id, err = dbs.Insert("product_warehouse",
		configs.M{
			"title":     p.Title,
			"buy_addr":  p.BuyAddr,
			"save_addr": p.SaveAddr,
			"buy_price": p.BuyPrice,
			"md5":       md5,
			"author_id": p.Author.Uid,
			"addtime":   util.Datetime(),
		})

	if err == nil {
		createFolder()

		item := Item{
			Id:          0,
			Name:      p.Title,
			WarehouseId: p.Id,
			Author:      p.Author,
		}
		if err := item.Save(); err != nil {
			return err
		}
	}
	return err
}

func (p *ProductWarehouse) Set(k string, v interface{}) {
	dbs.Update("product_warehouse", configs.M{k: v}, "id=?", p.Id)
}

func (p *ProductWarehouse) DownAllImage() {
	url := strings.Replace(p.BuyAddr, "/offer/", "/pic/", 1)

	doc, err := goquery.NewDocument(url)

	if err == nil {
		doc.Find("#dt-bp-tab-nav ul li").Each(func(i int, s *goquery.Selection) {
			src, _ := s.Attr("data-img")

			_, filename := filepath.Split(src)
			base_save := filepath.Join(p.BaseAddr, filename)
			save_path := filepath.Join(warehouse_base, base_save)

			images.Down(src, save_path)

			p.Images = append(p.Images, base_save)
		})
	}

	if len(p.Images) > 0 {
		p.Set("images", strings.Join(p.Images, ","))
	}
}

//1688下载图片
func DownToWarehouse(url string, author *account.Account) error {

	pro := &ProductWarehouse{BuyAddr: url}

	doc, err := goquery.NewDocument(url)

	if err != nil {
		return err
	}
	//price := doc.Find(".price-now").Text()
	pro.BuyPrice = 0.0
	pro.Title = enc.ConvertString(doc.Find(".d-title").Text())
	pro.Author = author
	if err := pro.New(); err != nil {
		return err
	}

	pro.DownAllImage()
	base_sku_path := filepath.Join(pro.BaseAddr, "sku图")
	sku_pic_path := filepath.Join(warehouse_base, base_sku_path)

	os.Mkdir(sku_pic_path, 777)

	doc.Find(".desc-lazyload-container img").Each(func(i int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		fmt.Println(src)
	})


	desc_url,ok:=doc.Find("#desc-lazyload-container").Attr("data-tfs-url")
	if ok {
		pro.Descurl=desc_url
	}

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
			return
		}

		skuImage := ProductSkuImage{
			Name:          attName.Name,
			Preview:       attImg.Preview,
			Original:      attImg.Original,
			LocalOriginal: fmt.Sprint(base_sku_path, "/", attName.Name+"_大图", ".jpg"),
			LocalPreview:  fmt.Sprint(base_sku_path, "/", attName.Name+"_小图", ".jpg"),
		}

		images.Down(attImg.Preview, filepath.Join(warehouse_base, skuImage.LocalPreview))
		images.Down(attImg.Original, filepath.Join(warehouse_base, skuImage.LocalOriginal))

		pro.SukImages = append(pro.SukImages, skuImage)

	})

	doc.Find(".table-sku tr").Each(func(i int, s *goquery.Selection) {

		title, _ := s.Find(".name span").Attr("title")
		attimgData, _ := s.Find(".name span").Attr("data-imgs")

		attImg := struct {
			Preview  string `json:"preview"`
			Original string `json:"original"`
		}{}
		title = filertName(enc.ConvertString(title))
		json.Unmarshal([]byte(attimgData), &attImg)

		if attImg.Preview == "" || attImg.Original == "" {
			return
		}
		skuImage := ProductSkuImage{
			Name:          title,
			Preview:       attImg.Preview,
			Original:      attImg.Original,
			LocalOriginal: fmt.Sprint(base_sku_path, "/", title+"_大图", ".jpg"),
			LocalPreview:  fmt.Sprint(base_sku_path, "/", title+"_小图", ".jpg"),
		}
		images.Down(attImg.Preview, filepath.Join(warehouse_base, skuImage.LocalPreview))
		images.Down(attImg.Original, filepath.Join(warehouse_base, skuImage.LocalOriginal))

		pro.SukImages = append(pro.SukImages, skuImage)
	})

	if b, err := json.Marshal(pro.SukImages); err == nil {
		pro.Set("sku_imags", string(b))
	} else {
		fmt.Println(err)
	}
	return nil
}

func filertName(title string) string {
	title = strings.Replace(title, "#", "", -1)
	title = strings.Replace(title, " ", "", -1)
	return title
}
