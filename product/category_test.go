package product 
import(
	_ "ainit"
	"testing"
)

func TestSetCateAtt(t *testing.T){
	if cat,err:=GetCategory(101);err==nil{
		attrs1:=map[int]int{1:1,3:3}
		if err=cat.SetAttr(attrs1);err!=nil{
			t.Error(err)
		}
		attrs2:=map[int]int{2:2,3:4}
		if err=cat.SetAttr(attrs2);err!=nil{
			t.Error(err)
		}
	}else{
		t.Error(err)
	}
}

func TestGetPlatformCategory(t *testing.T){
	opt:=GetPlatformCategorys(100,"lazada")
	if len(opt) == 1 && opt[7275].Name !="" {
		t.Log("OK",opt[0])
	}else{
		t.Error(opt)
	}

}

func TestGetPlatformCategorySelected(t *testing.T){
	GetPlatformCategorySelected(100,"lazada")
}