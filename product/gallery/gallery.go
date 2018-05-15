package gallery

import (
	"errors"
	"fmt"
	"helper/account"
	"helper/configs"
	"helper/dbs"
	"helper/util"
	"io"
	"helper/dbs/mongodb"
	"labix.org/v2/mgo"
	"helper/filestore"
)
const (
	HeadImageFlag   = 1 //橱窗图片
	DetailImageFlag = 2 //详情图片
	SkuImageFlag    = 4 //Sku图片
	CoverImageFlage =7 //封面
	SOURCE_QN =1  //七牛
	SOURCE_ALI=2  //阿里
)

var store *mgo.GridFS

func Store() *mgo.GridFS{
	if store == nil{
		store=mongodb.NewGridFs("images")
	}
	return store
}

type ImageInfo struct {
	Id       int
	Sku      string    `json:"sku"`
	Name     string `json:"name"`
	Addtime  string `json:"addtime"`
	Src      string `json:"src"`
	AuthorId int
	Author   *account.Account `json:"author"`
	Flag     int              `json:"flag"`
	Sort     int `json:"sort"`
	Source   int `json:"source"`
	Label string `json:"label"`
}

type Flag struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

func Flags()[]Flag {
	return []Flag{
		{Id:HeadImageFlag, Title: "橱窗图"},
		{Id:DetailImageFlag, Title:"详情图"},
		{Id:SkuImageFlag,  Title:"Sku图片"},
		{Id:CoverImageFlage,Title:"封面"},
	}
}

func UploadImage(file io.Reader)(name string,err error){
	return filestore.Save(file)
}


func Save(sku,src string,flag,sort int) (*ImageInfo, error) {
	data := configs.M{
		"sku":     sku,
		"add_time": util.Datetime(),
		"sort":     sort,
		"src":  src,
		"source":SOURCE_ALI,
	}
	db:=dbs.Def()
	id,err:=db.Insert("product_image",data)
	
	return &ImageInfo{
		Id:id,
		Src:src,
		Sku:sku,
		Sort:sort,
		Source:SOURCE_ALI,
	},err
}

//添加图片
func AddImage(sku string ,srcs ...string ){
	for i,src:=range srcs{
		Save(sku,src,HeadImageFlag,i)
	}
}

//删除一个Sku的相关图片
func DelSkuImage(sku string,ids []int )error{
	db:=dbs.Def()
	sql:=fmt.Sprintf("DELETE FROM product_image WHERE id IN(%s) AND sku='%s'",util.IntJoin(ids,","),sku)
	return db.Exec(sql)
}

func (res *ImageInfo) SetFlag(f ...int) {
	for _, v := range f {
		res.Flag = res.Flag | v
	}
	dbs.Update("product_image", configs.M{"flag": res.Flag}, "id=?", res.Id)
}

//设置图片标签
func SetFlag(ids []int, f ...int)error{

	if len(ids) == 0 {
		return errors.New("请至少选择一项")
	}
	wehre:=dbs.NewWhere()
	
	wehre.AndIntIn("id",ids)

	sql:="SELECT id,flag,src,sku FROM product_image" + wehre.ToString()

	rows := dbs.Rows( sql)
	defer rows.Close()
	for rows.Next() {
		res := ImageInfo{}
		rows.Scan(&res.Id,&res.Flag,&res.Src,&res.Sku)
		if f[0] == CoverImageFlage {
			sql:=fmt.Sprintf("UPDATE product SET headimg='%s' WHERE sku='%s'",res.Src,res.Sku)
	
			return dbs.Exec(sql)
		}
		res.SetFlag(f...)
	}
	return nil
}

//主图
func (res *ImageInfo) IsHeadImage() bool {
	if res.Flag & HeadImageFlag == HeadImageFlag {
		return true
	}
	return false
}

//详细页图
func (res *ImageInfo) IsDetailImage() bool {
	if res.Flag&DetailImageFlag == DetailImageFlag {
		return true
	}
	return false
}

//SKU图
func (res *ImageInfo) IsSkuImage() bool {
	if res.Flag&SkuImageFlag == SkuImageFlag {
		return true
	}
	return false
}

func (res *ImageInfo) Flags() []int {
	f := []int{}
	flags := Flags()
	for _, flag := range flags {

		if flag.Id&res.Flag == flag.Id {

			f = append(f, flag.Id)
		}
	}
	return f
}


func GetHeadImages(sku string)(images []string){

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

func Listing(filter configs.M, offset, rowCount int) ([]configs.M, int) {
	var total = 0
	where := dbs.NewWhere()

	if sku := filter.Get("sku");sku!="" {
		where.And("sku='%s'",sku) 
	}
	dbs.One(fmt.Sprint("SELECT count(*) FROM product_image ", where.ToString())).Scan(&total)

	sql := fmt.Sprint("SELECT id,src,sku,flag,source,sort,add_time FROM product_image ", where.ToString(), " ORDER BY add_time DESC", dbs.Limit(offset, rowCount))

	data := []configs.M{}

	rows := dbs.Rows(sql)

	for rows.Next() {
		res := ImageInfo{}
		rows.Scan(&res.Id, &res.Src,&res.Sku,&res.Flag,&res.Source,&res.Sort,&res.Addtime)
		item := configs.M{
			"id":            res.Id,
			"name":          res.Name,
			"sku":        	 res.Sku,
			"addtime":       res.Addtime,
			"author_id":     res.AuthorId,
			"author":        account.Author(res.AuthorId),
			"addr":          res.Src,
			"flag":          res.Flag,
			"isheadimage":   res.IsHeadImage(),
			"isdetailimage": res.IsDetailImage(),
			"isskuimage":    res.IsSkuImage(),
		}
		data = append(data, item)
	}
	return data, total
}