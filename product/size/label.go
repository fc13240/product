package size
import(
	"helper/dbs"
	"encoding/json"
	"helper/configs"
	"errors"
)
func GetLabels()[]string{
	return []string{
		"胸围 Bust Girth",
		"腰围 Waist",
		"臀围 Hips",
		"袖长 Sleeve Length",
		"裙长 Length",
		"衣长 Length",
		"肩宽 Shoulder",
		
	}
}
type Size struct{
	Id int `json:"id"`
	ItemId int `json:"item_id"`
	Head string `json:"head"`
	Rows []string `json:"rows"`
}

func Save(size Size)(  int ,  error){
	if b,err:=json.Marshal(size);err==nil{
		db:=dbs.Def()
		if size.Id>0{
			err=db.Update("product_size_template",configs.M{"body":b},"id=?",size.Id)
		}else{
			size.Id,err=db.Insert("product_size_template",configs.M{"body":b,"item_id":size.ItemId})
		}
		return size.Id,err
	}else{
		return size.Id,err
	}
}

func Get(item_id int)(size Size,err error){
	db:=dbs.Def()
	var body []byte
	var id int 
	db.One("SELECT id,body FROM product_size_template WHERE item_id =?",item_id).Scan(&id,&body)
	size=Size{}

	if id == 0{
		err=errors.New("not exist")
		return 
	}
	if err:=json.Unmarshal(body,&size);err==nil{
		size.Id=id
		return size,nil
	}else{
		return size,err
	}
}