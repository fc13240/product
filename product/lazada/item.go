package lazada
import(
	"helper/dbs/mongodb"
	"encoding/xml"
	"helper/configs"
	"product"
	"fmt"
	"strings"
	"time"
	"errors"
)


type Item struct{
	AuthorId int `json:"author_id" bson:"author_id"`
	CategoryId int `json:"category_id"`
	Name string  `json:"name" bson:"name" xml:"name"`
	Cnname string `json:"cnname" bson:"cnname"`
	Description string `json:"description" bson:"description"`
	ParentSku string `json:"parent_sku" bson:"parent_sku"`
	ShortDescription string `json:"short_description" bson:"short_description"`
	Skus []Sku `json:"skus" bson:"skus"`
	Attrs []ItemAttr `json:"other_attr"  bson:"other_attr"`
	Price float32 `json:"price"`
	PackageContent string `json:"package_content"  bson:"package_content"`
	PackageLength string `json:"package_length"  bson:"package_length"`
	PackageWidth string `json:"package_width"  bson:"package_width"`
	PackageHeight string `json:"package_height"  bson:"package_height"`
	PackageWeight string `json:"package_weight"  bson:"package_weight"`
	Sizes[] []string `json:"sizes"`
	Currency string `json:"currency"`
	SpecialFromDate		string `json:"special_from_date" xml:"special_from_date"`
	SpecialToDate		string `json:"special_to_date" xml:"special_to_date"`
	SpecialPrice	float32 `json:"special_price" xml:"special_price"`
}

type ItemAttr struct{
	Label string `json:"label"`
	Value product.AttrOptions `json:"value"`
}

func ItemCol() *mongodb.Collection{
	mdb:=mongodb.Conn()
	return mdb.C("lazada.item")
}

func SaveProduct(item *Item){
	col:=ItemCol()
	if total,_:=col.Find(configs.M{"author_id":item.AuthorId,"parent_sku":item.ParentSku}).Count();total==0{
		col.Insert(item)
	}else{
		col.Update(configs.M{"author_id":item.AuthorId,"parent_sku":item.ParentSku},item)
	}
}

func GetProduct(autor_id int,sku string)(item *Item){
	col:=ItemCol()
	item=&Item{}
	col.Find(configs.M{"author_id":autor_id,"parent_sku":sku}).One(&item)
	if item.Description =="" || item.Name =="" || item.Price == 0{
		pitem,_:=product.Get(sku)
		if item.Description==""{
			item.Description=pitem.Desc
		}

		if item.Name==""{
			item.Cnname=pitem.Name
		}

		if item.Price ==0{
			item.Price=pitem.Price
		}
	}
	item.Currency="RM"
	return item
}
type Product struct{
	PrimaryCategory int `xml:"PrimaryCategory"`
	Attributes struct {
		Content string `xml:",innerxml"`
	} `xml:"Attributes"`
	Skus []Sku `xml:"Skus>Sku"`
}

type Base struct{
	Name string `xml:"name"`
	ShortDescription string `xml:"short_description"`
	NameMs string `xml:"name_ms"`
	Description string `xml:"description"`
	DescriptionMs string `xml:"description_ms"`
	Brand string `xml:"brand"`
	WarrantyType string `xml:"warranty_type"`
}

type Dress struct{
	Base
	ClothingMaterial string `xml:"clothing_material"`
	DressLength string	`xml:"dress_length"`
	DressShape  string `xml:"dress_shape"`
}
func (item *Item)GetSizes()[]product.AttrOptions{
	return product.GetOneAttrSelectedOption(item.ParentSku,1)
}

func (item *Item)Create()error{
	req:=NewReq("CreateProduct")
	sizes:=item.GetSizes()
	Skus:=[]Sku{}

	//促销价格
	if item.SpecialPrice >0 { //如果设置了促销价格
		now:=time.Now()
		item.SpecialFromDate=now.Format("2006-01-02 15:04")
		item.SpecialToDate=now.Add(time.Hour*24*30).Format("2006-01-02 15:04")
	}

	for i,_:=range item.Skus{
		curr_sku:=item.Skus[i]
		curr_sku.Price=item.Price
		curr_sku.PackageContent=item.PackageContent
		curr_sku.PackageHeight=item.PackageHeight
		curr_sku.PackageLength=item.PackageLength
		curr_sku.PackageWeight=item.PackageWeight
		curr_sku.PackageWidth=item.PackageWidth

		curr_sku.SpecialPrice=item.SpecialPrice
		curr_sku.SpecialFromDate=item.SpecialFromDate
		curr_sku.SpecialToDate=item.SpecialToDate
		if len(sizes)>0{
			for i,size:=range sizes{
				new_sku:=Sku{
					Images:curr_sku.Images,
					Size:size.Value,
					ColorFamily:curr_sku.ColorFamily,
					PackageContent:curr_sku.PackageContent,
					PackageHeight:curr_sku.PackageHeight,
					PackageLength:curr_sku.PackageLength,
					PackageWeight:curr_sku.PackageWeight,
					PackageWidth:curr_sku.PackageWidth,
					Price:curr_sku.Price,
					Quantity:curr_sku.Quantity,	
					TaxClass :curr_sku.TaxClass,
					SellerSku:curr_sku.SellerSku+fmt.Sprint(i+1),
					SpecialPrice:curr_sku.SpecialPrice,
					SpecialFromDate:curr_sku.SpecialFromDate,
					SpecialToDate:curr_sku.SpecialToDate,
					
				}
				Skus=append(Skus,new_sku)
			}
		}
	}
	itemDetail:=Product{Skus:Skus,PrimaryCategory:item.CategoryId}

	type Request struct{
		Product Product 
	}

	//atts:=map[string]string{}
	selectedAttrs:=product.GetAttrSelectedOption(item.ParentSku)
	
	var attrs XmlAttrs

	
	attrs.Put("name",item.Name)
	attrs.Put("short_description",ToHighlights(item.ShortDescription))
	attrs.Put("description",item.Description)
	for _,att:=range selectedAttrs{
		att_name:=att.Config.Get("lazada")
		
		if att_name==""{
			continue
		}

		if len(att.Options)>0{
			//atts[att.Name]=att.Options[0].Value
			fmt.Println(att_name,att.Options[0].Value)
			attrs.Put(att_name,att.Options[0].Value)
		}
	}
	/*
	dress:=Dress{
		Base:Base{
			Brand:"No Brand",
			Name:item.Name,
			NameMs:item.Name,
			ShortDescription:item.ShortDescription,
			Description:item.Description,
		},
		ClothingMaterial:atts["clothing_material"],
		DressLength:atts["dress_length"],
		DressShape:atts["dress_shape"],
	}
	*/
	itemDetail.Attributes.Content=attrs.Encode()
	bb,_:=xml.Marshal(Request{itemDetail})

	output, _ := xml.MarshalIndent(Request{itemDetail}, "  ", "    ")
	fmt.Println(string(output))
	  
	//data.Add(ss+"[Attributes][kid_years]","Kids (6-10yrs)")
	//req.PostTest(data)
	res,_:=req.Post(bb)

	 resData:=struct{
		ErrorResponse struct{
			Body struct{
				Errors []struct{
					Field string 
					Message string
					SellerSku string
				}
			}
		}
		SuccessResponse struct{
			Body struct{
				Warnings []struct{
					Field string 
					Message string
					SellerSku string
				}
			}
		}
	}{}
	res.BindJSON(&resData)
	fmt.Println(resData)
	
	if len(resData.ErrorResponse.Body.Errors) == 0{
		return nil
	}

	return errors.New(resData.ErrorResponse.Body.Errors[0].Message)
	
}



func ToHighlights(desc string)string{
	rows:=strings.Split(desc,"\n")
	var new_desc string
	for _,line:=range rows{
		if line!=""{
			new_desc+="<li>"+strings.Trim(line," ")+"</li>"
		}
	}
	return fmt.Sprint("<ul>\n",new_desc,"</ul>")
}