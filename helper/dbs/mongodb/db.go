package mongodb

import (
	"fmt"
	"helper/configs"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

var (
	conn         *Mogno
	mongo_url    string
	mongo_dbname string
)

func Conn() *Mogno {
	if conn == nil {
		opt := configs.GetSection("mongo")
		mongo_url = opt["url"]
		mongo_dbname = opt["dbname"]
		conn = NewConn(mongo_url, mongo_dbname)
	}else if err:=conn.Db.Ping();err!=nil{
		opt := configs.GetSection("mongo")
		mongo_url = opt["url"]
		mongo_dbname = opt["dbname"]
		conn = NewConn(mongo_url, mongo_dbname)
	}
	return conn
}

func Rexp(reg,opt string)bson.RegEx{
	return bson.RegEx{Pattern:reg,Options:opt}
}

type Mogno struct {
	url        string
	Db         *mgo.Session
	Database   *mgo.Database
	Collection *mgo.Collection
}

type Collection struct {
	*mgo.Collection
}

type MyCollectioner interface {
	C() *Collection
}

type Result struct {
	Items []configs.M
	Item  configs.M
	index int
	col   *Collection
}

func (r *Result) Next() bool {
	if len(r.Items) > r.index {
		r.Item = r.Items[r.index]
		r.index++
		return true
	} else {
		r.index = 0
		return false
	}
}

func (r *Result) Scan(vv *configs.M) {
	*vv = r.Item
}

func (r *Result) UpSet(vv configs.M) {
	if v, ok := r.Item["_id"].(bson.ObjectId); ok {
		r.col.UpsertId(v, configs.M{"$set": vv})
	}
}

func (self *Collection) Query(query configs.M, page, rowCount int) *mgo.Query {
	skip := (page - 1) * rowCount
	return self.Find(query).Skip(skip).Limit(rowCount)
}

func (self *Collection) ItemsCall(query configs.M, call func(items *Result), page, rowCount int) {

	count, _ := self.Find(query).Count()
	for {
		skip := (page - 1) * rowCount
		if skip >= count {
			return
		}
		items := &Result{col: self}
		self.Find(query).Skip(skip).Limit(rowCount).All(&items.Items)

		call(items)
		page++
	}
}

func C(colldction string) *Collection {
	col := conn.Database.C(colldction)
	return &Collection{col}
}

func GridFS(name ...string) *mgo.GridFS {
	return conn.Database.GridFS("fs")
}

func NewGridFs(name string) *mgo.GridFS {
	return conn.Database.GridFS(name)
}

func (mdb *Mogno) New(dbName string) *Mogno {
	mdb.Connect(mongo_url)
	mdb.DB(dbName)
	return mdb
}

func (mdb *Mogno) GridFS(name ...string) *mgo.GridFS {
	return mdb.Database.GridFS("fs")
}

func NewConn(url, dbName string) *Mogno {
	mdb := Mogno{}
	if mdb.Connect(url) == false {
		fmt.Println("mongodb connect failing")
	}
	
	mdb.DB(dbName)
	return &mdb
}

func (mdb *Mogno) C(colldction string) *Collection {
	col := mdb.Database.C(colldction)
	return &Collection{col}
}

func (mdb *Mogno) Connect(url string) bool {
	db, err := mgo.Dial(url)
	if err != nil {
		fmt.Println("connect failing,", err.Error())
		return false
	}
	mdb.Db = db
	return true
}

func (mdb *Mogno) DB(db string) {
	mdb.Database = mdb.Db.DB(db)
}
