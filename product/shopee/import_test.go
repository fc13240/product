package shopee
import (
	_ "ainit"
	"testing"
	"helper/account"
	"fmt"
	"strings"
)

func TestExec(t *testing.T){
	ss:="ps_category_list_id ps_product_name	ps_product_description	ps_price	ps_stock	ps_product_weight	ps_days_to_ship	ps_sku_ref_no_parent	ps_mass_upload_variation_help	ps_variation 1 ps_variation_sku	ps_variation 1 ps_variation_name	ps_variation 1 ps_variation_price	ps_variation 1 ps_variation_stock	ps_variation 2 ps_variation_sku	ps_variation 2 ps_variation_name	ps_variation 2 ps_variation_price	ps_variation 2 ps_variation_stock	ps_variation 3 ps_variation_sku	ps_variation 3 ps_variation_name	ps_variation 3 ps_variation_price	ps_variation 3 ps_variation_stock	ps_variation 4 ps_variation_sku	ps_variation 4 ps_variation_name	ps_variation 4 ps_variation_price	ps_variation 4 ps_variation_stock	ps_variation 5 ps_variation_sku	ps_variation 5 ps_variation_name	ps_variation 5 ps_variation_price	ps_variation 5 ps_variation_stock	ps_variation 6 ps_variation_sku	ps_variation 6 ps_variation_name	ps_variation 6 ps_variation_price	ps_variation 6 ps_variation_stock	ps_variation 7 ps_variation_sku	ps_variation 7 ps_variation_name	ps_variation 7 ps_variation_price	ps_variation 7 ps_variation_stock	ps_variation 8 ps_variation_sku	ps_variation 8 ps_variation_name	ps_variation 8 ps_variation_price	ps_variation 8 ps_variation_stock	ps_variation 9 ps_variation_sku	ps_variation 9 ps_variation_name	ps_variation 9 ps_variation_price	ps_variation 9 ps_variation_stock	ps_variation 10 ps_variation_sku	ps_variation 10 ps_variation_name	ps_variation 10 ps_variation_price	ps_variation 10 ps_variation_stock	ps_variation 11 ps_variation_sku	ps_variation 11 ps_variation_name	ps_variation 11 ps_variation_price	ps_variation 11 ps_variation_stock	ps_variation 12 ps_variation_sku	ps_variation 12 ps_variation_name	ps_variation 12 ps_variation_price	ps_variation 12 ps_variation_stock	ps_variation 13 ps_variation_sku	ps_variation 13 ps_variation_name	ps_variation 13 ps_variation_price	ps_variation 13 ps_variation_stock	ps_variation 14 ps_variation_sku	ps_variation 14 ps_variation_name	ps_variation 14 ps_variation_price	ps_variation 14 ps_variation_stock	ps_variation 15 ps_variation_sku	ps_variation 15 ps_variation_name	ps_variation 15 ps_variation_price	ps_variation 15 ps_variation_stock	ps_variation 16 ps_variation_sku	ps_variation 16 ps_variation_name	ps_variation 16 ps_variation_price	ps_variation 16 ps_variation_stock	ps_variation 17 ps_variation_sku	ps_variation 17 ps_variation_name	ps_variation 17 ps_variation_price	ps_variation 17 ps_variation_stock	ps_variation 18 ps_variation_sku	ps_variation 18 ps_variation_name	ps_variation 18 ps_variation_price	ps_variation 18 ps_variation_stock	ps_variation 19 ps_variation_sku	ps_variation 19 ps_variation_name	ps_variation 19 ps_variation_price	ps_variation 19 ps_variation_stock	ps_variation 20 ps_variation_sku	ps_variation 20 ps_variation_name	ps_variation 20 ps_variation_price	ps_variation 20 ps_variation_stock	ps_img_1	ps_img_2	ps_img_3	ps_img_4	ps_img_5	ps_img_6	ps_img_7	ps_img_8	ps_img_9"
	sss:=strings.Split(ss,"	")
	for _,s:=range sss{
		fmt.Printf(`"%s",`,s)
		fmt.Println()
	}
	item,err:=Get("Bl20002225")
	if err==nil{
		auth:=account.Find(21)
		if _,err:=Export(auth,item);err!=nil{
			t.Error(err.Error())
		}
	}
}