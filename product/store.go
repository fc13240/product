package product
import(
	"helper/dbs"
	"helper/configs"
)
type Store struct{
	Id int `json:"id"`
	StoreId int `json:"store_id"`
	StoreName string `json:"stroe_name"`
	StoreSite string `json:"site"`
	GroupId int `json:"group_id"`
}

func NewStore()*Store{
	return &Store{}
}

func (store *Store)Save()(err error){
	db:=dbs.Def()
	if store.Id == 0{
		 if exist_store_id:=store.Exist();exist_store_id>0{  //if this site store_name exist 
			store.Id=exist_store_id
		 }
	}
	if store.Id>0 {
		err=db.Update("product_store",configs.M{
			"store_id":store.StoreId,
			"store_site":store.StoreSite,
			"store_name":store.StoreName,
		},"id=?",store.Id)
	}else{
		id,err:=db.Insert("product_store",configs.M{
			"store_id":store.StoreId,
			"store_site":store.StoreSite,
			"store_name":store.StoreName,
		})
		if err == nil{
			store.Id=id
		}
	}
	return err
}

func (store *Store)Exist()(id int){
	db:=dbs.Def()
	db.One("SELECT id FROM product_store WHERE store_name=? AND store_site=?",store.StoreName,store.StoreSite).Scan(&id)
	return id
}

func GetStores()(stores []Store){
	db:=dbs.Def()
	rows:=db.Rows("SELECT id,store_id,stroe_name FROM product_store")
	stores=[]Store{}
	defer rows.Close()
	for rows.Next(){
		store:=Store{}
		rows.Scan(&store.Id,&store.StoreId,&store.StoreName)
		stores=append(stores,store)
	}
	return stores
}