package shopee
import(
	"fmt"
)
var header =[]string{
	"ps_category_list_id",
	"ps_product_name",
	"ps_product_description",
	"ps_price",
	"ps_stock",
	"ps_product_weight",
	"ps_days_to_ship",
	"ps_sku_ref_no_parent",
	"ps_mass_upload_variation_help",
	"ps_variation 1 ps_variation_sku",
	"ps_variation 1 ps_variation_name",
	"ps_variation 1 ps_variation_price",
	"ps_variation 1 ps_variation_stock",
	"ps_variation 2 ps_variation_sku",
	"ps_variation 2 ps_variation_name",
	"ps_variation 2 ps_variation_price",
	"ps_variation 2 ps_variation_stock",
	"ps_variation 3 ps_variation_sku",
	"ps_variation 3 ps_variation_name",
	"ps_variation 3 ps_variation_price",
	"ps_variation 3 ps_variation_stock",
	"ps_variation 4 ps_variation_sku",
	"ps_variation 4 ps_variation_name",
	"ps_variation 4 ps_variation_price",
	"ps_variation 4 ps_variation_stock",
	"ps_variation 5 ps_variation_sku",
	"ps_variation 5 ps_variation_name",
	"ps_variation 5 ps_variation_price",
	"ps_variation 5 ps_variation_stock",
	"ps_variation 6 ps_variation_sku",
	"ps_variation 6 ps_variation_name",
	"ps_variation 6 ps_variation_price",
	"ps_variation 6 ps_variation_stock",
	"ps_variation 7 ps_variation_sku",
	"ps_variation 7 ps_variation_name",
	"ps_variation 7 ps_variation_price",
	"ps_variation 7 ps_variation_stock",
	"ps_variation 8 ps_variation_sku",
	"ps_variation 8 ps_variation_name",
	"ps_variation 8 ps_variation_price",
	"ps_variation 8 ps_variation_stock",
	"ps_variation 9 ps_variation_sku",
	"ps_variation 9 ps_variation_name",
	"ps_variation 9 ps_variation_price",
	"ps_variation 9 ps_variation_stock",
	"ps_variation 10 ps_variation_sku",
	"ps_variation 10 ps_variation_name",
	"ps_variation 10 ps_variation_price",
	"ps_variation 10 ps_variation_stock",
	"ps_variation 11 ps_variation_sku",
	"ps_variation 11 ps_variation_name",
	"ps_variation 11 ps_variation_price",
	"ps_variation 11 ps_variation_stock",
	"ps_variation 12 ps_variation_sku",
	"ps_variation 12 ps_variation_name",
	"ps_variation 12 ps_variation_price",
	"ps_variation 12 ps_variation_stock",
	"ps_variation 13 ps_variation_sku",
	"ps_variation 13 ps_variation_name",
	"ps_variation 13 ps_variation_price",
	"ps_variation 13 ps_variation_stock",
	"ps_variation 14 ps_variation_sku",
	"ps_variation 14 ps_variation_name",
	"ps_variation 14 ps_variation_price",
	"ps_variation 14 ps_variation_stock",
	"ps_variation 15 ps_variation_sku",
	"ps_variation 15 ps_variation_name",
	"ps_variation 15 ps_variation_price",
	"ps_variation 15 ps_variation_stock",
	"ps_variation 16 ps_variation_sku",
	"ps_variation 16 ps_variation_name",
	"ps_variation 16 ps_variation_price",
	"ps_variation 16 ps_variation_stock",
	"ps_variation 17 ps_variation_sku",
	"ps_variation 17 ps_variation_name",
	"ps_variation 17 ps_variation_price",
	"ps_variation 17 ps_variation_stock",
	"ps_variation 18 ps_variation_sku",
	"ps_variation 18 ps_variation_name",
	"ps_variation 18 ps_variation_price",
	"ps_variation 18 ps_variation_stock",
	"ps_variation 19 ps_variation_sku",
	"ps_variation 19 ps_variation_name",
	"ps_variation 19 ps_variation_price",
	"ps_variation 19 ps_variation_stock",
	"ps_variation 20 ps_variation_sku",
	"ps_variation 20 ps_variation_name",
	"ps_variation 20 ps_variation_price",
	"ps_variation 20 ps_variation_stock",
	"ps_img_1",
	"ps_img_2",
	"ps_img_3",
	"ps_img_4",
	"ps_img_5",
	"ps_img_6",
	"ps_img_7",
	"ps_img_8",
	"ps_img_9",
}


type Row struct{
	CategoryID int
	ProductName string
	ProductDescription string
	Price float32
	Stock int
	ProductWeight int
	ShipOutIn string
	ParentSKUReferenceNo string
	Specs []Variation
    Images []string
}

func NewRow(item *Item)(data []interface{}){
	data=[]interface{}{
		item.Cid,
		item.Name,
		item.Desc,
		item.Price,
		200,
		item.Weight,
		5,
		item.Sku,
		"-",
	}
	fmt.Println(item)
	for i:=0;i<=19;i++{
		if len(item.Variations)>i{
			spec:=item.Variations[i]
			data=append(data,spec.Sku,spec.Name,spec.Price,spec.Stock)
		}else{
			data=append(data,"","","","")
		}
	}

	for i:=0;i<=9;i++{
		if len(item.Images)>i{
			image:=item.Images[i]
			data=append(data,image)
		}else{
			data=append(data,"")
		}
	}
	return data
}