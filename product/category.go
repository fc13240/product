package product

import (
	"helper/configs"
	"helper/dbs"
)
type Category struct{
	Id int  `json:"id"`
	Name string `json:"name"`
	PId int   `json:"pid"`
	Level int  `json:"level"`
	EzbuyId int `json:"ezbuyid"`
	ShopeeId int   `json:"shopeeid"`
	LazadaId int `json:"lazadaid"`
	Platform string `json:"platform"`
}

func AddCategory(cate *Category)(id int, err error){
	db:=dbs.Def()
	if cate.Id>0{
		err=db.Update("product_category",
			configs.M{"name":cate.Name,"pid":cate.PId,"ezbuyid":cate.EzbuyId,"shopeeid":cate.ShopeeId,"lazadaid":cate.LazadaId},
			"id=?",cate.Id)
			
		id=cate.Id
		return 
	}
	id,err= db.Insert("product_category",
		configs.M{"name":cate.Name,"pid":cate.PId,"ezbuyid":cate.EzbuyId,"shopeeid":cate.ShopeeId,"lazadaid":cate.LazadaId})
		return id,err
	}

func CategoryExist(name string)bool{
	db:=dbs.Def()
	count:=0
	db.One("SELECT COUNT(*) FROM product_category WHERE name=?",name).Scan(&count)
	if count>0{
		return true
	}else{
		return false
	}
}

func Categorys(pid int)(items []Category){
	db:=dbs.Def()
	rows:=db.Rows("SELECT id,pid,name,level,ezbuyid,shopeeid,lazadaid FROM product_category WHERE pid=?",pid)
	items=[]Category{}
	for rows.Next(){
		var c Category
		rows.Scan(&c.Id,&c.PId,&c.Name,&c.Level,&c.EzbuyId,&c.ShopeeId,&c.LazadaId)
		items=append(items,c)
	}
	return items
}

//获取单个分类
func GetCategory(cid int)(c *Category,err error){
	db:=dbs.Def()
	c=&Category{}
	err=db.One("SELECT id,pid,name,level,ezbuyid,shopeeid,lazadaid FROM product_category WHERE id=?",cid).
	Scan(&c.Id,&c.PId,&c.Name,&c.Level,&c.EzbuyId,&c.ShopeeId,&c.LazadaId)
	return 
}

//设置分类属性
func (c *Category)SetAttr(attrs map[int]int)error{
	db:=dbs.Def()
	if tx,err:=db.Begin();err==nil{
		err=db.Exec("DELETE FROM product_category_attr WHERE cid=? AND platform=?",c.Id,c.Platform)
		for attid,sort:=range attrs{
			_,err=db.Insert("product_category_attr",configs.M{"cid":c.Id,"attid":attid,"sort":sort,"platform":c.Platform})
		}
		if err==nil{
			tx.Commit()
		}else{
			tx.Rollback()
		}
		return err
	}else{
		return err
	}
}

//分类的属性ID
func GetCategoryAttrIds(platform string ,cid int)(attrs map[int]int){
	db:=dbs.Def()
	where:=dbs.NewWhere()
	where.And(" cid=%d AND platform='%s' ",cid,platform)
	
	rows:=db.Rows("SELECT attid,sort FROM product_category_attr "+where.ToString())
	defer rows.Close()
	attrs=map[int]int{}
	for rows.Next(){
		var attid,sort int
		rows.Scan(&attid,&sort)
		attrs[attid]=sort
	}
	return attrs
}


func CategoryListing(filter configs.M,offset,limit int)(items []Category,count int){
	db:=dbs.Def()
	sql:="SELECT id,name,level,ezbuyid,shopeeid,lazadaid FROM product_category"
	rows:=db.Rows(sql)
	defer rows.Close()

	count=db.Count(sql)
	items=[]Category{}
	for rows.Next(){
		var c Category
		rows.Scan(&c.Id,&c.Name,&c.Level,&c.EzbuyId,&c.ShopeeId,&c.LazadaId)
		items=append(items,c)
	}
	return items,count
}

type PlatformProductCategory struct{
	Id int `json:"id"` 
	Name  string `json:"name"`
}

func GetPlatformCategorys(cid int ,platform string) (data []PlatformProductCategory){
	data=[]PlatformProductCategory{}
	db:=dbs.Def()
	where:=dbs.NewWhere()
	where.And("cid=%d AND platform='%s'",cid,platform)
	sql:="SELECT platform_cid,platform_cname FROM platform_product_category"+where.ToString()
	
	rows:=db.Rows(sql)

	defer rows.Close()

	for rows.Next(){
		var (
			id int
			name string
		)
		rows.Scan(&id,&name)
		data=append(data,PlatformProductCategory{id,name})
	}

	return data
}

func GetPlatformCategorySelected(platform_cid int ,platform string)(data []configs.M){
	data=[]configs.M{}
	db:=dbs.Def()

	if platform == "base"{ //如果不属于平台
		cate,_:=GetCategory(platform_cid)
		data=append(data,configs.M{"label":cate.Name,"value":cate.Id})
		return 
	}


	//如果属于平台分类。也就是说有两级分类

	where:=dbs.NewWhere()
	where.And("platform_cid=%d AND platform='%s'",platform_cid,platform)
	sql:="SELECT cid,platform_cid,platform_cname FROM platform_product_category"+where.ToString()
	
	var (
		cid int
		platform_cname string
	)

	db.One(sql).Scan(&cid,&platform_cid,&platform_cname)

	if cid>0{
		cate,_:=GetCategory(cid)
		data=append(data,configs.M{"label":cate.Name,"value":cate.Id})
		data=append(data,configs.M{"label":platform_cname,"value":platform_cid})
	}

	return data
}