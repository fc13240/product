package product

import (
	"errors"
	"fmt"
	"helper/account"
	"helper/configs"
	"helper/dbs"
	timeutil "helper/time"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"product/gallery"
)

type Type struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	CnName string `json:"cnname"`
}

func GetChannelList() []Type {
	return []Type{
		{1, "Dresses", "连衣裙"},
		{2, "Coats", "大衣"},
		{3, "Wedding Dresses", "婚纱"},
		{4, "Shoes", "鞋子"},
		{5, "Bags", "包"},
	}
}


type AttrVal struct{
	Id  int `json:"id"`
	AttId  int  `json:"attid"`
	Value string  `json:"value"`
}

type Item struct {
	Id          int `json:"id"`
	Name        string  `redis:"name",json:"name"`
	EnName      string  `json:"en_name"`
	Sku         string  `redis:"sku",json:"sku"`
	Quant       int `json:"quant"`
	Channelid   int `json:"channel_id"`
	CategoryId  int `json:"category_id"`
	Desc        string	`json:"desc"`
	Headimg     string
	Images      []string
	Authorid    int
	Price       float32	`json:"price"`
	OldPrice    float32 `json:"old_price"`
	BuyingPrice float32 `json:"buying_price"`
	Author      *account.Account
	Addtime     time.Time
	IsUnder     bool
	Tempid      int
	EbayItemid  int
	WarehouseId int
	Length      float32 `json:"length"`
	Width       float32 `json:"width"`
	Height      float32 `json:"height"`
	Weight      float32 `json:"weight"`
	Attrs	   []AttrVal
	Saveprice  int
	Remarks     string `json:"remarks"`
	LabelLog   *LabelLog
}

func (item *Item) SavePrice() int{
	if item.Price<item.OldPrice{
		return int(100-item.Price/item.OldPrice*100)
	}
	return 0
}

func (item *Item) Save() (err error) {
	if item.Name == ""  {
		return errors.New("product name cen not be empty")
	}

	if item.Authorid==0{
		return errors.New("author id is empty")
	}

	data := configs.M{
		"`name`":       item.Name,
		"en_name":      item.EnName,
		"quant":        item.Quant,
		"`desc`":       item.Desc,
		"price":        item.Price,
		"old_price":    item.OldPrice,
		"channel_id":   item.Channelid,
		"up_time":      timeutil.String(),
		"warehouse_id": item.WarehouseId,
		"sku":          item.Sku,
		"length":       item.Length,
		"width":        item.Width,
		"height":       item.Height,
		"weight":       item.Weight,
		"buying_price": item.BuyingPrice,
		"category_id":item.CategoryId,
	}


	if item.Id == 0 || false == Exist(item.Id) {
		data.Add("add_time", timeutil.String())
		data.Add("author_id", item.Authorid)
		item.Id, err = dbs.Insert("product", data)
	} else {
		err = dbs.Update("product", data, fmt.Sprintf("id=%d and author_id=%d", item.Id,item.Authorid))
	}
	return err
}

func (item *Item) SetHeadImg(img string) {
	dbs.Update("product", configs.M{"headimg": filepath.Base(img), "up_time": timeutil.String()}, "id=?", item.Id)
}

func (item *Item) SetUnder() {
	dbs.Update("product", configs.M{"is_under": 1, "up_time": timeutil.String()}, "id=?", item.Id)
}

func (item *Item) SetUpper() {
	dbs.Update("product", configs.M{"is_under": 0, "up_time": timeutil.String()}, "id=?", item.Id)
}

//保存模板
func (item *Item) SetTemp(temp_id int, desc string) {
	dbs.Update("product", configs.M{"`temp_id`": temp_id, "`desc`": desc}, "id=?", item.Id)
}

func (item *Item) BindEbay(ebay_itemid int) error {
	return dbs.Update("product", configs.M{"ebay_itemid": ebay_itemid}, "id=?", item.Id)
}

//sku
func (item *Item) SetSku(sku string) error {
	return dbs.Update("product", configs.M{"sku": sku}, "id=?", item.Id)
}


//商品是否存在
func Exist(item_id int)bool{
	db:=dbs.Def()
	var id int
	db.One("SELECT id FROM product WHERE id=? LIMIT 1",item_id).Scan(&id)
	if id>0{
		return true
	}else{
		return false
	}
}



func (item *Item) GetAttrIds(type_name string) []int {
	rows := dbs.Rows("SELECT value_id FROM product_attr WHERE product_id=? AND `type`=?", item.Id, type_name)
	var ids []int
	defer rows.Close()
	for rows.Next() {
		var id int
		rows.Scan(&id)
		ids = append(ids, id)
	}
	return ids
}



func (item *Item) GetSmallImages() []configs.M {
	list := []configs.M{}

	return list
	rows := dbs.Rows("SELECT product_image_samll.`id`,product_image_samll.`sort`,product_image.`name`,`product_image_samll`.`name` AS small_name,product_image_samll.`image_id` "+
		"FROM  product_image  left join  product_image_samll ON (product_image.id=product_image_samll.image_id)  WHERE product_image.item_id=?  "+
		"AND typeid=90120 ORDER BY product_image.`sort`", item.Id)
	defer rows.Close()
	for rows.Next() {
		var id, image_id, sort int
		var name, small_name string

		rows.Scan(&id, &sort, &name, &small_name, &image_id)

		list = append(list, configs.M{
			"id":         id,
			"sort":       sort,
			"name":       JoinImgPath(item.Id, name),
			"small_name": JoinImgPath(item.Id, small_name),
			"image_id":   image_id,
		})
	}
	return list
}

//返回产品规格
func (item *Item)GetOptVal(att_id int)(opts []AttrOptions){
	selected:=GetAttrVal(item.Sku)
	opts=[]AttrOptions{}
	if _,ok:=selected[att_id];ok {
		att:=GetAttr(att_id)
		opts=att.GetOptions(selected[att_id]...)
	}
	return opts
}


//返回图片路径
func JoinImgPath(product_id int, name string) string {
	if name == "" {
		return ""
	}
	return strings.Replace(filepath.Join("product", fmt.Sprintf("%d/%d", configs.Int(product_id/1000), product_id), name), "\\", "/", -1)
}

func (item *Item) Tags() map[string]int {
	return GetItemSaleTags(item.Id)
}

func Get(sku string) (*Item, error) {
	return get(fmt.Sprintf("sku='%s'",sku))
}

func IdGet(id int) (*Item, error) {
	return get(fmt.Sprintf("id=%d",id))
}

func get(where string)(item *Item,err error){
	item = &Item{}
	var add_time string
	sql:=fmt.Sprint("SELECT `id`,`sku`,`author_id`,`channel_id`,`category_id`,`name`,`en_name`,`quant`,`desc`,`headimg`,`add_time`,`ebay_itemid`,`price`,`old_price`,`temp_id`,`length`,`width`,`height`,`weight`,`buying_price` FROM product WHERE ",where," LIMIT 1")
	dbs.One(sql).Scan(&item.Id, &item.Sku, &item.Authorid, &item.Channelid,&item.CategoryId, &item.Name, &item.EnName, &item.Quant, &item.Desc, &item.Headimg, &add_time, &item.EbayItemid, &item.Price, &item.OldPrice, &item.Tempid, &item.Length,&item.Width,&item.Height,&item.Weight,&item.BuyingPrice)
	item.Addtime, _ = time.Parse("2006-01-02 15:04:05", add_time)

	if item.Id > 0 {
		item.Saveprice=item.SavePrice()
		item.Images=gallery.GetHeadImages(item.Sku)

		if len(item.Images)>0{
			item.Headimg=item.Images[0]
		}
		return item, nil
	}
	return item, errors.New("没有这个商品")
}

func GetBySku(sku string,author *account.Account) (*Item, error){
	return get(fmt.Sprintf("sku='%s' AND author_id=%d",sku,author.Uid))
}

func SkuExist(sku string,author *account.Account)(int ,bool){
	where:=fmt.Sprintf("sku='%s' AND author_id=%d",sku,author.Uid)
	sql:=fmt.Sprint("SELECT `id`  FROM product WHERE ",where," LIMIT 1")
	var id int
	dbs.One(sql).Scan(&id)
	if id>0{
		return id,true
	}else{
		return 0,false
	}
}

//搜索
func Search(param configs.M, offset, rowCount int,sort string ) ([]*Item, int) {
	db:=dbs.Def()
	var items []*Item

	var orderBy = "id DESC"
	where := dbs.NewWhere()
	if param.String("search_type")  == "name" {
		where.And("`name` LIKE '%s'","%"+param.Get("value")+"%" )
	}else if param.String("search_type") == "sku" {
		where.And("`sku` LIKE '%s'",param.Get("value")+"%")
	}
	if param.Int("is_under") == 1 {
		where.And(" is_under =1 ")
	} else {
		where.And(" is_under =0 ")
	}



	if author_id:=param.Int("author_id");author_id>0{
		where.And("`author_id`=%d ",author_id)
	}

	if item_ids:=param.Ints("in");len(item_ids)>0 {
		where.And(" `id` IN(%s)",strings.Trim(strings.Replace(fmt.Sprint(item_ids), " ",",", -1), "[]"))
	}

	if cate_id:=param.Int("cate_id");cate_id>0{
		where.And(" category_id = %d",cate_id)
	}

	switch sort {
	case "update":
		orderBy = "up_time DESC "
	case "min_price":
		orderBy = "price ASC"
	case "max_price":
		orderBy = "price DESC"
	}

	var sql = "SELECT `id`,`author_id`,`channel_id`,`name`,`sku`,`quant`,`headimg`,`add_time`,`price`,`old_price`,`warehouse_id` FROM product " + where.ToString()
	total := db.Count("SELECT COUNT(id) FROM product " + where.ToString())

	rows := db.Rows(fmt.Sprint(sql, " ORDER BY ", orderBy, " ", dbs.Limit(offset, rowCount)))
	defer rows.Close()
	for rows.Next() {
		item := &Item{}
		var add_time string
		rows.Scan(&item.Id, &item.Authorid, &item.Channelid, &item.Name, &item.Sku, &item.Quant,  &item.Headimg, &add_time, &item.Price, &item.OldPrice, &item.WarehouseId)
		item.Addtime, _ = time.Parse("2006-01-02 15:04:05", add_time)
		items = append(items, item)

		item.Saveprice=item.SavePrice()

		item.Images=gallery.GetHeadImages(item.Sku)

		if len(item.Images)>0{
			item.Headimg=item.Images[0]
		}

		if label,ok:=FirstLabelLog(item.Id);ok{
			item.LabelLog=label
		}
	}
	defer rows.Close()
	return items, total
}

//保存图片路径
func addImgPath(path, small_path string, sort, product_id int) (int, error) {
	id, err := dbs.Insert("product_image",
		configs.M{
			"path":       path,
			"small_path": small_path,
			"add_time":   timeutil.String(),
			"sort":       sort,
		})
	return id, err
}

//保存图片路径
func SaveImgPath(name, small_name string, sort, product_id int) (int, error) {
	data := configs.M{
		"name":       name,
		"small_name": small_name,
		"add_time":   timeutil.String(),
		"sort":       sort,
		"product_id": product_id,
	}
	db:=dbs.Def()
	if product_id > 0 {
		db.Exec("DELETE FROM product_image WHERE product_id=? AND sort =?", product_id, sort)
	}
	id, err := dbs.Insert("product_image", data)
	return id, err
}

//图片和产品绑定
func (item *Item) BindImg(images []int) {
	ids := []string{}
	for _, id := range images {
		ids = append(ids, strconv.Itoa(id))
	}
	db:=dbs.Def()
	db.Exec("DELETE FROM product_image WHERE product_id=? AND id NOT IN("+strings.Join(ids, ",")+")", item.Id)
	db.Exec("UPDATE product_image SET product_id=? WHERE id IN("+strings.Join(ids, ",")+")", item.Id)
}

//del
func DelItem(author *account.Account,ids ...int)error{
	ss := []string{}
	for _, id := range ids {
		ss = append(ss, strconv.Itoa(id))
	}
	db:=dbs.Def()
	return db.Exec("DELETE FROM product WHERE id IN(?) and author_id=?",strings.Join(ss, ","),author.Uid)
}


func (item *Item)SaveField(field string,value interface{} )(error){
	set:=configs.M{}
	set[field]=value
	return dbs.Update("product",set,"sku=?",item.Sku)

}

func GetImages(sku string)(images []string){

	images=[]string{}
	db:=dbs.Def()

	rows:=db.Rows("SELECT src FROM product_image WHERE sku=? ORDER BY sort,id ",sku)
	for rows.Next(){
		var name string
		rows.Scan(&name)
		images=append(images,name)
	}
	return images
}
