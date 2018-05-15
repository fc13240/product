package lazada
import(
	"testing"
	"product"
	"fmt"
)
func TestGetProduct(t *testing.T){
	pro:=GetProduct(21,"Bl2000162")
	fmt.Println(ToHighlights(pro.ShortDescription))
}

func TestCreate(t *testing.T){
//	f,_:=os.Open("D:/aa.jpg")
//	fmt.Println(UploadImage(f))
	GetProduct(21,"Bl2000162").Create()
}



func TestGet1Attr(t *testing.T){
	
	data:=struct{
		SuccessResponse struct{
			Body []Attribute `json:"Body"`
		}
	}{}
	req:=NewReq("GetCategoryAttributes")
	req.SetParam("PrimaryCategory","1740")
		
	if err:=req.BindJson(&data);err==nil{
		for _,attr:=range data.SuccessResponse.Body{
			if attr.AttributeType!="sku"  {
			//fmt.Println(attr.Name)
			if attr.Name == "dress_shape"{
			for _,v:=range attr.Options{//sleeves,dress_shape 5,collar_type,fa_pattern,clothing_material
				att_id:=7
				opt_name:=v.Get("name")
				if product.AttrOptionExist(att_id,opt_name) == false {
					fmt.Println(opt_name," not exist")
					product.NewAttOption(att_id,v.Get("name"),"")
				}

			}
		
		}
		//fmt.Println(attr.Label,"---------------", attr.Name)
			
		}
		}
	}else{
		fmt.Println(err)
	}
}


func TestGetCategorys(t *testing.T){
	data:=struct{
		SuccessResponse struct{
			Body []Category
		}
	}{}

	req:=NewReq("GetCategoryTree")
	
	if err:=req.BindJson(&data);err==nil{
		for _,m:=range data.SuccessResponse.Body{
		
			if m.CategoryId == 1819{
				
				if len(m.Children)>0{
					m.PrintChild("----")
				}
			}
		}
	}
	
}