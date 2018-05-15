package lazada

import (
	"fmt"
	"errors"
	"helper/configs"
	"html"
)

var (
	ApiUrl string
	UID string
	ApiKey string
)

func init(){
	config:=configs.GetSection("lazada")
	ApiUrl=config["apiAddr"]
	UID=config["UID"]
	ApiKey=config["apiKey"]
}

type Category struct{
	Name string `json:"name"`
	Var bool `json:"var"`
	CategoryId uint64 `josn:"categoryId"`
	Leaf bool `json:"leaf"`
	Children []Category  `json:"children"`
}

func (c *Category)PrintChild(flag string ){
	if len(c.Children)> 0{
		for _,child:=range c.Children{
			fmt.Println(flag,child.Name,child.CategoryId)
			child.PrintChild(flag+flag)
		}
	}
}

type Image struct {
	Code string
	Url string
}

type ErrorResponse struct{
	Body struct{
		Errors []map[string]string
	}
}

func (ers *ErrorResponse)Error()error{
	if len(ers.Body.Errors)==0{
		return nil
	}
	str:=""
	for k,m:=range ers.Body.Errors{
		str+=fmt.Sprintf("%d:%s\n",k,m)
	}
	return errors.New(str)
}

type Sku struct{
	Images		[]string `json:"Images" xml:"Images>Image"`
	SellerSku	string `json:"SellerSku" xml:"SellerSku"`
	ColorFamily	string `json:"color_family" xml:"color_family"`
	PackageContent string `json:"package_content" xml:"package_content"`
	PackageHeight	string `json:"package_height" xml:"package_height" `
	PackageLength	string `json:"package_length" xml:"package_length"`
	PackageWeight	string `json:"package_weight" xml:"package_weight"`
	PackageWidth	string `json:"package_width" xml:"package_width"`
	Price		float32 `json:"price" xml:"price"`
	Quantity		int `json:"quantity" xml:"quantity"`
	Size	string `json:"size" xml:"size"`
	TaxClass string `json:"tax_class" xml:"tax_class"`
	SpecialFromDate		string `json:"special_from_date" xml:"special_from_date"`
	SpecialToDate		string `json:"special_to_date" xml:"special_to_date"`
	SpecialPrice	float32 `json:"special_price" xml:"special_price"`
}
type Specia struct {
	SpecialFromDate		string `json:"special_from_date" xml:"special_from_date"`
	SpecialFromTime		string `json:"special_from_time" xml:"special_from_time"`
	SpecialPrice	float32 `json:"special_price" xml:"special_price"`
	SpecialTimeFormat string `json:"special_time_format" xml:"special_time_format"`
	SpecialToDate 	string `json:"special_to_date" xml:"special_to_date"`
	SpecialToTime	string `josn:"special_to_time" xml:"special_to_time"`
}

type SkuOLd struct{
	Available int `json:"Available"`
	Images		[]string `json:"Images"`
	SellerSku	string `json:"SellerSku"`
	ShopSku		string `json:"ShopSku"`
	Status		string `json:"Status"`
	Url	string `json:"Url"`
	ColorFamily	string `json:"color_family" xml:"color_family"`
	PackageContent string `json:"package_content"`
	PackageHeight	float32 `json:"package_height" `
	PackageLength	float32 `json:"package_length"`
	PackageWeight	float32 `json:"package_weight"`
	PackageWidth	float32 `json:"package_width"`
	Price		float32 `json:"price"`
	Quantity		int `json:"quantity"`
	Size	string `json:"size"`
	SpecialFromDate		string `json:"special_from_date"`
	SpecialFromTime		string `json:"special_from_time"`
	SpecialPrice	float32 `json:"special_price"`
	SpecialTimeFormat string `json:"special_time_format"`
	SpecialToDate 	string `json:"special_to_date"`
	SpecialToTime	string `josn:"special_to_time"`
	TaxClass string `json:"tax_class"`
}

type XmlAttrs []string
func (attrs *XmlAttrs)Put(k,v string){
	
	*attrs=append(*attrs,fmt.Sprintf("<%s>%s</%s>\n",k,html.EscapeString(v),k))
}
func (attrs XmlAttrs)Encode()string{
	
	ss:=""

	for _,s:=range attrs{
		ss+=string(s)
	}
	return ss
}

type Attribute struct{
	Name string `json:"name"`
	Label string `json:"label"`
	IsMandatory int `json:"isMandatory"`
	Options []configs.M `json:"options"`
	InputType string `json:"inputType"`
	AttributeType string `json:"attributeType"`
}