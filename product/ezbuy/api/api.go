package api
import (
	"net/http"
	"fmt"
	"helper/webtest"
	"time"
	"log"
	"helper/util"
	"product/ezbuy"
	"product/ezbuy/client"
	"encoding/json"
	"strings"
	"errors"
	"helper/configs"
	"qiniupkg.com/api.v7/kodocli"
	"io"
	"bytes"
)

type Sku struct{
	Attributes Attributes `json:"attributes"`
	Height int `json:"height" ,omitempty`
	IsOnSale bool `json:"isOnSale"`
	Length int `json:"length" ,omitempty`
	Name string `json:"name"`
	OriginalPrice float32 `json:"originalPrice"`
	Price float32 `json:"price"`
	Quantity int `json:"quantity"`
	SellerSkuId string `json:"sellerSkuId"`
	ShippingFee float32 `json:"shippingFee,omitempty"`
	SkuId int `json:"skuId,omitempty"`
	Volume Volume `json:"volume"`
	Weight float32 `json:"weight" ,omitempty`
	Width int `json:"width" ,omitempty`
}

type Attributes struct{
	ColorName string `json:"Color Name"`
	Size string `json:"Size(Clothes),omitempty"`
	Material string `json:"材质（服饰）,omitempty"`
}

type Volume struct{
	Length int `json:"length"`
	Width int `json:"width"`
	Height int `json:"height"`
}

type SkuPropsValues struct {
	ValueId int `json:"valueId"`
	ValueName string `json:"valueName"`
	Image string `json:"image"`

}

type SkuProps struct{
	PropId int `json:"propId"`
	PropName string `json:"propName"`
	Values []SkuPropsValues `json:"values"`
}

type ProductDetail struct{
	Attributes Attributes `json:"attributes"`
	CategoryId int `json:"categoryId"`
	Committed bool `json:"committed"`
	Description string `json:"description"`
	EnName string `json:"enName"`
	ForceOffSale bool `json:"forceOffSale"`
	Images []string `json:"images"`
	IsOnSale bool `json:"isOnSale"`
	IsStocked  bool `json:"isStocked"`
	Name string `json:"name"`
	OriginCode string `json:"originCode"`
	Pid int `json:"pid"`
	PrimaryImage string `json:"primaryImage"`
	SellType int `json:"sellType"`
	ShipmentInfo int `json:"shipmentInfo"`
	SkuProps []SkuProps `json:"skuProps"`
	Source int `json:"source"`
	Url string `json:"url"`
}

type EzbuyeApi struct {
	Cookie string
	ReqId string
	ShopName  string
	Client *client.HClient
	SkuFirst string
}

func NewEzbuyeApi(setting *ezbuy.Setting) ( *EzbuyeApi){
	return &EzbuyeApi{
		Cookie:setting.Cookie,
		ReqId:setting.Reqid,
		ShopName:setting.StoreName,
		SkuFirst:setting.SkuFirst,
	}
}

func (ez *EzbuyeApi)RefreshAll(setting *ezbuy.Setting){

	data:=struct{
		IsSucc bool `json:"isSucc"`
		Items []ezbuy.Item `json:"items"`
		Total  int `json:"total"`
	}{}

	if setting.Minute< 1{
		log.Println("循环时间不能小于1")
		return
	}

	if setting.Num>20 || setting.Num<1{
		log.Println("每次循环的商品数量不能大于20，小于1")
		return
	}

	for {
		ezbuy.RunSta.RefreshRunSta="on"
		res,err:=ez.Client.GetItems(0,setting.Num)

		if err!=nil{
			log.Println(err.Error())
			return
		}

		if err:=res.BindJSON(&data);err!=nil{
			log.Println("读取商品失败",err.Error())
			return
		}

		if len(data.Items) == 0{
			log.Println("没有更新的商品....")
			time.Sleep(time.Second*30)
			continue
		}

		fmt.Println("开始更新,更新数量:",len(data.Items),"总数:",data.Total)
		for _,item:=range data.Items{
			time.Sleep(time.Second*1)
			ez.Refresh(item,false)
			time.Sleep(time.Second*1)
			ez.Refresh(item,true)
		}
		log.Println("等待，",setting.Minute,"分钟....")
		time.Sleep(time.Second*time.Duration(setting.Minute*60))
	}
}


func (ez *EzbuyeApi)HotRefresh(){
	data:=struct{
		IsSucc bool `json:"isSucc"`
		Items []ezbuy.Item `json:"items"`
		Total  int `json:"total"`
	}{}

	res,err:=ez.Client.GetHotItems(0,50)

	if err!=nil{
		log.Println(err.Error())
		return
	}

	if err:=res.BindJSON(&data);err!=nil{
		log.Println("读取商品失败",err.Error())
		return
	}

	if len(data.Items) == 0{
		log.Println("没有更新的商品....")
		time.Sleep(time.Second*30)
	}

	for {
		for _,item:=range data.Items{
			ez.Refresh(item,false)
			ez.Refresh(item,true)
			time.Sleep(time.Second*1)
		}
		time.Sleep(time.Second*30)
	}
}

// set product on sale or off the this
func (ez *EzbuyeApi)Refresh(item ezbuy.Item,flag bool) bool {
	defer func() {
		if e:=recover();e!=nil{
			log.Println(e)
		}
	}()
	item_id:=item.Id
	if item.Id==0{
		return false
	}

//	body:=`{"change":{"isOnSale":%t},"productId":%d}`
	body:=struct{
		Change struct{
			IsOnSale bool `json:"isOnSale"`
		} `json:"change"`
		ProductId []int `json:"productId"`
	}{}
	body.Change.IsOnSale=flag
	body.ProductId=[]int{item_id}
	//body=fmt.Sprintf(body,flag,item_id)
	b,_:=json.Marshal(body)

	res,_:=webtest.PostJson("https://webapi.ezbuy.com/api/EzSeller/UserProductQuickUpdate?",ez.sellerHeader(),string(b))
	defer res.Close()
	if flag {
		fmt.Println(item.Name," 上架 SUCC" ,"上次更新时间:",item.Update)
		re,_:=ez.Client.UpItemField(item.Id,"update",util.Datetime())
		fmt.Println(re)
	}else{
		fmt.Println(item.Name,"下架 SUCC")
	}
	return true
}

//上传批量上传编辑好的商品
//首先要调用 UserUnCommitedProductDetail
func (api *EzbuyeApi)userUnCommitedProductUpdate(detail ProductDetail,skus []Sku){
	body := `{"data":%s,"skus":%s}`
	detail_encode,_:=json.Marshal(detail)
	skus_encode,_:=json.Marshal(skus)

	body=fmt.Sprintf(body,string(detail_encode),string(skus_encode))
	body=strings.Replace(body,"Material","材质（服饰）",-1)

	webtest.PostJson("https://webapi.ezbuy.com/api/EzSeller/UserUnCommitedProductUpdate?", api.sellerHeader(), body)
}

//获取一个批量上传未编辑的商品详细，并且编辑好，提交..
func (api *EzbuyeApi)UserUnCommitedProductDetail(pid int )error {

	body := `{"productId":%d}`
	body=fmt.Sprintf(body,pid)
	res, err := webtest.PostJson("https://webapi.ezbuy.com/api/EzSeller/UserUnCommitedProductDetail?", api.sellerHeader(), body)
	if err!=nil{
		return err
	}
	var data=struct{
		Data ProductDetail `json:"base"`
		Skus []Sku `json:"skus"`
	}{}

	err=res.BindJSON(&data)

	if err!=nil{
		fmt.Println(err)
	}else{
		//fmt.Println(data)
	}
	item,_:=api.Client.GetItem(data.Skus[0].SellerSkuId)
	def_sku:=data.Skus[0]
	def_sku.Height=def_sku.Volume.Height
	def_sku.Width=def_sku.Volume.Width

	var material *ezbuy.Material
	if len(data.Data.SkuProps) ==3{ //材质处理
		material_obj:=data.Data.SkuProps[2].Values[0]
		material=&ezbuy.Material{ID:material_obj.ValueId,Name:material_obj.ValueName}
	}
	fmt.Println(data.Data.SkuProps)
	for i,_:=range data.Data.SkuProps{
		curr:=data.Data.SkuProps[i]
		if curr.PropId == 122535 { //颜色属性
			
			for i,_:=range curr.Values{
				curr.Values[i].Image=item.GetColorImage(curr.Values[i].ValueName) 
			}
		}
	}
	
	
	new_skus:=[]Sku{}

	addSku:=func(size SkuPropsValues){
		skuProps:=data.Data.SkuProps
		for i,_:=range skuProps{
			curr:=skuProps[i]
			fmt.Println("PropId:",curr.PropId)
		if curr.PropId == 122535 { //颜色属性
		
		for _,color:=range skuProps[i].Values {
			fmt.Println(color)
			var attributes Attributes
			var name string
			if size.ValueName==""{
				name=fmt.Sprintf("%s",color.ValueName)
				attributes=Attributes{color.ValueName,"",""}
			}else{
				name=fmt.Sprintf("%s; %s",size.ValueName,color.ValueName)
				attributes=Attributes{color.ValueName,size.ValueName,""}
			}
			if material!=nil{
				name=fmt.Sprintf("%s; %s; %s",size.ValueName,color.ValueName,material.Name)
				attributes=Attributes{color.ValueName,size.ValueName,material.Name}
			}

			sku_name:=strings.Replace(strings.ToUpper(def_sku.SellerSkuId),"BL",api.SkuFirst,-1)

			new_sku:=Sku{
				Name:name,
				Attributes:attributes,
				Height:def_sku.Volume.Height,
				IsOnSale:def_sku.IsOnSale,
				Length:def_sku.Length,
				OriginalPrice:def_sku.OriginalPrice,
				Price:def_sku.Price,
				Quantity:def_sku.Quantity,
				SellerSkuId:sku_name,
				//ShippingFee:def_sku.ShippingFee,
				//SkuId:def_sku.SkuId,
				Volume:def_sku.Volume,
				Weight:def_sku.Weight,
				Width:def_sku.Volume.Width,
			}
			new_skus=append(new_skus,new_sku)
		}
		}}
	}
	if len(data.Data.SkuProps)>1 && len(data.Data.SkuProps[0].Values)>0 { 
		for _,size:=range data.Data.SkuProps[0].Values{
			addSku(size)
		}
	}else{
		addSku(SkuPropsValues{ValueName:""})
	}
	fmt.Println("开始上传....")

	api.userUnCommitedProductUpdate(data.Data,new_skus)
	go func(pid int){
		time.Sleep(time.Second*3)
		api.Refresh(ezbuy.Item{Id:pid},true)
	}(pid)
	return err
}

//获取批量上传的待编辑的商品列表
func (api *EzbuyeApi)UserProductsFromSource(limit int )(*webtest.Result,error) {
	body:=`{"src": 4, "committed": false, "offset": 0, "limit": %d}`
	body=fmt.Sprintf(body,limit)
	return webtest.PostJson("https://webapi.ezbuy.com/api/EzSeller/UserProductsFromSource?", api.sellerHeader(), body)
}

//公共头部
func (api *EzbuyeApi)sellerHeader() http.Header{
	head := http.Header{}
	head.Set("origin","https://ezseller.ezbuy.com")
	head.Set("referer","https://ezseller.ezbuy.com")
	head.Set("cookie", api.Cookie)
	head.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.91 Safari/537.36")
	head.Set("x-req-id", api.ReqId)
	head.Set("access-control-allow-headers", "X-Requested-With,Content-Type,Ajax")
	head.Set("access-control-allow-origin", "https://ezseller.ezbuy.com")
	head.Set("access-control-allow-credentials", "true")
	head.Set("vary", "Accept-Encoding")
	head.Set("Host","webapi.ezbuy.com")
	return head
}

//get ez product list
func (api *EzbuyeApi)GenList(offset int) ([]ezbuy.Item,error) {

	body := `{
	"filter":{
		"isOnSale":true,
		"isStocked":null,
		"maxCreateDate":null,
		"maxPrice":null,
		"minCreateDate":null,
		"minPrice":null,
		"productName":"",
		"sellType":1,
		"soldCountSortType":0
	},
	"limit":40,
	"offset":%d
	}
	`
	body = fmt.Sprintf(body, offset)
	res, err := webtest.PostJson("https://webapi.ezbuy.com/api/EzSeller/UserProductList?", api.sellerHeader(), body)
	defer res.Close()
	if err != nil {
		return nil,err
	}
	param:=struct {
		Total int `json:"total"`
		Items []ezbuy.Item `json:"products"`
	}{}

	if err:=res.BindJSON(&param);err!=nil{
		return  nil,err
	}

	if len(param.Items) == 0{
		return nil,errors.New("请求ezbuy失败"+res.String())
	}

	if res1, err := api.Client.SaveItems(ezbuy.Encode(param)); err != nil {
		return  nil,err
	} else {
		defer res1.Close()
		re := ezbuy.Result{}
		if err := res1.BindJSON(&re); err == nil {
			if !re.IsSucc {
				return  nil,errors.New(re.ErrorMsg)
			}
			return param.Items,nil
		} else {
			return nil,err
		}
	}
}

//get order
func (ez *EzbuyeApi)GetOrders(offset int)(*webtest.Result,error){
	head:=http.Header{}
	body:=`{
	"dataType":"new",
	"language":"en",
	"limit":10,
	"offset":%d,
	"filter":{}}`
	body=fmt.Sprintf(body,offset)
	head.Set("origin","https://ezseller.ezbuy.com")
	head.Set("referer","https://ezseller.ezbuy.com")
	head.Set("cookie",ez.Cookie)
	head.Set("user-agent","Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/46.0.2490.80 Safari/537.36")
	head.Set("x-req-id",ez.ReqId)
	head.Set("access-control-allow-headers","X-Requested-With,Content-Type,Ajax")
	head.Set("access-control-allow-origin","https://ezseller.ezbuy.com")
	head.Set("access-control-allow-credentials","true")
	head.Set("vary","Accept-Encoding")
	res,err:=webtest.PostJson("https://webapi.ezbuy.com/api/EzSeller/UserOrderList",head,body)

	if err!=nil{
		return  nil,err
	}
	return res,nil
}

//提交文件key
func (ez *EzbuyeApi) UserUploadProducts(fileKey string)(*webtest.Result,error){
	body:=fmt.Sprintf(`{"fileKey":"%s"}`,fileKey)
	return webtest.PostJson("https://webapi.ezbuy.com/api/EzSeller/UserUploadProducts",ez.sellerHeader(),body)
}

//get upload file tokne
func (ez *EzbuyeApi) GetUploadToken()string {
	body:=`{}`
	res,_:=webtest.PostJson("https://webapi.ezbuy.com/api/homepage/AdminHomepage/GetUploadInfo",ez.sellerHeader(),body)
	info:=struct {
		Token string `json:"token"`
		BaseUrl string `json:"baseUrl"`
	}{}
	res.BindJSON(&info)
	return info.Token
}

func (ez *EzbuyeApi)Upload(file io.Reader)(string,error){
	buff := &bytes.Buffer{}
	filesize,_:=io.Copy(buff,file)
	var ret kodocli.PutRet
	uploader := kodocli.NewUploader(0, nil)
	err:= uploader.PutWithoutKey(nil, &ret,ez.GetUploadToken(), buff, filesize, nil)
	return ret.Key,err
}

//upload exec to ezbuy
func  (ez *EzbuyeApi)UploadExec(url string) (string, error) {
	//构建一个uploader
	log.Println("上传："+url)
	zone := 0
	response, e := http.Get(url)
	if e != nil {
		log.Fatal(e)
	}
	fmt.Println(ez)
	token:=ez.GetUploadToken()

	log.Println("Ez Token:",token)

	uploader := kodocli.NewUploader(zone, nil)

	var ret kodocli.PutRet

	buff := &bytes.Buffer{}

	filesize, _ := buff.ReadFrom(response.Body)

	err:= uploader.PutWithoutKey(nil, &ret,token, buff, filesize, nil)

	defer response.Body.Close()

	if err != nil {
		log.Println("上传失败:",err)
		return "", err
	}
	log.Println("上传文件成功,Key:",ret.Key)
	res,err:=ez.UserUploadProducts(ret.Key)

	str:=res.String()
	log.Println("提交到EZ结果:",str)

	if str == "[]"{
		go ez.UploadLastProductStatus()
		return "succ",nil
	}else{
		return str,errors.New(str)
	}

}

//check new order and send message to email
func (ez *EzbuyeApi) CheckNewOrders() error {
	if res,err:=ez.GetOrders(0);err==nil{
		data:=struct {
			Total int `json:"total"`
			Orders []ezbuy.Order `json:"data"`
		}{}

		if err:=res.BindJSON(&data);err!=nil{
			return err
		}

		log.Println("订单总数:",data.Total)

		if data.Total== 0{
			return nil
		}
		res,err:=ez.Client.CheckNewOrders(ezbuy.Encode(configs.M{"total":data.Total,"data":data.Orders}))

		resdata:=struct {
			IsSucc bool `json:"isSucc"`
			NewOrderNum int `json:"newOrderNum"`
		}{}
		if err!=nil{
			fmt.Println(res.String())
		}
		res.BindJSON(&resdata)
		log.Println("新订单数:",resdata.NewOrderNum)
		return nil
	}else{
		return err
	}
}

//update new last 5 products status to sale list
func (ez *EzbuyeApi)UploadLastProductStatus(){

	data:=struct {
		Total int `json:"total"`
		Result []struct{
			Pid int `json:"pid"`
			Name string `json:"name"`
		} `json:"results"`
	}{}
	
	if res,err:=ez.UserProductsFromSource(5);err==nil{
		res.BindJSON(&data)
		fmt.Printf("有%d 款商品有上传到销售列表",data.Total)
		if data.Total >0 {
			for _,p:=range data.Result{
				fmt.Println(p.Name," 上传到销售列表 ",ez.UserUnCommitedProductDetail(p.Pid)) //上架

			}
		}
	}
}



