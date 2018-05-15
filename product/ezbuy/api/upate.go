package api 
import(
	"helper/webtest"
	"fmt"
	"encoding/json"
	"sort"
	"strings"
)
const (
	EzApi="https://webapi.ezbuy.com/api/EzSeller"
)

type EzItem struct{
	Base ProductDetail `json:"base"`
	Skus []Sku `json:"skus"`
	Sku string `json:"sku,omitempty"`
}

//获取商品详情
func (api *EzbuyeApi)UserProductDetail(productId int)(ezItem *EzItem,err error) {
	body:=`{"productId":%d}`
	body=fmt.Sprintf(body,productId)
	res,err:=webtest.PostJson(EzApi+"/UserProductDetail",api.sellerHeader(),body)
	
	defer res.Close()
	if err==nil{
		err=res.BindJSON(&ezItem)
	}else{
		return ezItem,err
	}
	return
}

//提交修改
func (api *EzbuyeApi)UserProductUpdate(ezItem *EzItem)(new_item *EzItem, err error){

	if b,err:=json.Marshal(ezItem);err==nil{
		body:=string(b)
		body=strings.Replace(body,"base","data",1)

		res,err:=webtest.PostJson(EzApi+"/UserProductUpdate",api.sellerHeader(),body)
		defer res.Close()
			
		if err==nil{
			new_item=&EzItem{}
			err=res.BindJSON(&new_item)
		}else{
			return nil,err
		}
	}
	return 
}

func (api *EzbuyeApi)Add(pid int){
	ezItem,_:=api.UserProductDetail(pid)
	ezItem.Base.Source=1
	ezItem.Base.ForceOffSale=true
	ezItem.Base.Pid=0
	ezItem.Base.EnName=""
	for i:=range ezItem.Skus {  //把size 和 color 转换过来，以前程序问题
		sku:=&ezItem.Skus[i]
		sku.SkuId=0
	}
	api.UserProductUpdate(ezItem)
}

func (api *EzbuyeApi)UserProductList(offset,limit int)(*webtest.Result,error){
	body:=`{"filter":{"isOnSale":true, "sellType": 1, "soldCountSortType": 0, "minPrice": null, "maxPrice": null,"minCreateDate":null},
	"offset":%d,"limit":%d}`
	body=fmt.Sprintf(body,offset,limit)
	return webtest.PostJson(EzApi+"/UserProductList",api.sellerHeader(),body);
}


func (api *EzbuyeApi)UserProductUpdateId(product_id int){
	if ezitem,err:=api.UserProductDetail(product_id);err==nil{
		api.UserProductUpdate(ezitem)
	}
}

func(api *EzbuyeApi)UserProductQuickUpdate(product_id int){
	body:=`{"productId":[37050414], "change": {"sellType":1}}`
	res,_:=webtest.PostJson(EzApi+"/UserProductQuickUpdate",api.sellerHeader(),body)
	fmt.Println(res.String())
}

func (api *EzbuyeApi)UserProductUpdateSize(product_id int){
	if ezitem,err:=api.UserProductDetail(product_id);err==nil{
		s:=[]string{"SM","6XL","2XL","3XL","170","175","180","185","190","25","wx4","4XL","XXXXL","ss","XS","XXS","M","L","XL","XXL","均码","XXXL","XXXS"}
		isUpdate:=false
		for i:=range ezitem.Skus {  //把size 和 color 转换过来，以前程序问题
			sku:=&ezitem.Skus[i]
			if sort.SearchStrings(s,sku.Attributes.ColorName) != 23{
				old:=sku.Attributes.ColorName 
				sku.Attributes.ColorName = sku.Attributes.Size
				sku.Attributes.Size=old
				isUpdate=true
				fmt.Println(sku.Attributes)
			}	
		}
		if isUpdate{
			fmt.Println("更新。。。。")
			api.UserProductUpdate(ezitem)
		}
	}

}

func(api *EzbuyeApi)UserProductBatchDelete(pid ...int){
	body:=struct{
		Pids []int `json:"pids"`
	}{pid}
	b,_:=json.Marshal(body)
	res,_:=webtest.PostJson(EzApi+"/UserProductBatchDelete",api.sellerHeader(),string(b))
	fmt.Println(res.String())
}