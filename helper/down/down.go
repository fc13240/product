package down

import (
	"github.com/garyburd/redigo/redis"
	"helper/configs"
	"helper/dbs"
	"helper/dbs/mongodb"
	"helper/redisCli"
	"helper/util"
	"io/ioutil"
	"labix.org/v2/mgo"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
	"fmt"
)

var (
	mu        sync.Mutex
	r         redis.Conn
	gridfs    *mgo.GridFS
	dbName    string
	IndexName string
)

func CreateIndex(){
	db:=dbs.Def()
	sql:=fmt.Sprintf(
		`create table if not exists %s (
			id int(11) NOT NULL AUTO_INCREMENT,
			url varchar(200) DEFAULT '',
			addtime datetime DEFAULT NULL,
			status tinyint(3) DEFAULT '0',
			md5 varchar(32) DEFAULT '',
			PRIMARY KEY (id)
			)ENGINE=MyISAM AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
	`,IndexName)
	db.Exec(sql)
}

func Push(url string) (id int, name string){

	name = util.Md5(url)

	if id = Exist(name); id > 0 {
		return id, name
	}

	var err error
	id, err = dbs.Insert(IndexName, configs.M{
		"`md5`":  name,
		"`url`":   url,
		"addtime": util.Datetime(),
	})

	if err != nil {
		log.Println(err)
		return 0, ""
	}
	return id, name
}

func Exist(name string) (id int) {
	sql:=fmt.Sprintf("SELECT id FROM %s WHERE name=%s", IndexName,name)
	dbs.One(sql).Scan(&id)
	return id
}

func EncoryUrl(url string) string {
	return util.Md5(url)
}

func GetBody(url string) (body []byte, err error) {
	md5:=EncoryUrl(url)

	file, e := mongodb.GridFS().Open(md5)
	db:=dbs.Def()
	if e != nil {
		db.Update(IndexName, configs.M{"status": 0}, "md5=?", md5)
		return body, e
	}
	defer file.Close()
	return ioutil.ReadAll(file)

}

func GetFile(name string) (file *mgo.GridFile, err error) {
	return mongodb.GridFS().Open(name)
}

func GetInfo(name string) (fileinfo FileInfo, err error) {
	dbs.One("SELECT `id`,`name`,`url`,`addtime`,`status`,`content_type`,`file_size` FROM bee_downs WHERE `name`=?", name).
		Scan(&fileinfo.Id, &fileinfo.Name, &fileinfo.Url, &fileinfo.Addtime, &fileinfo.Status, &fileinfo.ContentType, &fileinfo.Size)
	return fileinfo, nil
}

type FileInfo struct {
	Id          int
	Name        string
	Url         string
	Addtime     string
	Status      int
	ContentType string
	Size        int
}

func (file *FileInfo) CopyTo(dir string) {
	b, err := GetBody(file.Name)
	if err != nil {
		return
	}
	to := filepath.Join(dir, file.Name+file.Ext())

	new_file, err := os.Create(to)

	defer new_file.Close()

	if err == nil {
		new_file.Write(b)
	}
}

func (file *FileInfo) Ext() string {
	switch file.ContentType {
	case "image/jpeg":
		return ".jpg"
	case "text/html":
		return ".html"
	}
	return ""
}

func InitData() {
	index := "downs:"
	r := redisCli.Conn()

	redis.Int(r.Do("del", index))

	rows := dbs.Rows("SELECT url FROM bee_downs WHERE status=?", 0)
	for rows.Next() {
		var url string
		rows.Scan(&url)
		r.Send("rpush", index, url)
	}
	r.Flush()
}

func RunDown(maxline int) {
	InitData()
	var url string
	var err error
	time.Sleep(time.Second * 1)
	for i := 0; i <= maxline; i++ {
		go func(a int) {
			r := redisCli.Conn()
			for {
				log.Print("我是进程:", a)
				mu.Lock()
					url, err = redis.String(r.Do("lpop", "downs:"))
				mu.Unlock()
				if err != nil {
					log.Println("接受")
					return
				}
				Down(url)
			}
		}(i)
	}

	for {
		time.Sleep(time.Second * 1)
		r := redisCli.Conn()
		len, _ := redis.Int(r.Do("llen", "downs:"))
		log.Println("剩余", len)
	}
}

func Down(url string) {

	log.Println("开始下载")

	res, err := http.Get(url)

	if err != nil {
		log.Println("open src :", err.Error())
		return
	}

	defer res.Body.Close()

	name := util.Md5(url)

	file, err := mongodb.GridFS().Create(name)
	defer file.Close()
	if err != nil {
		log.Println("Mongo errro ", err)
		return
	}

	contentType := res.Header.Get("Content-Type")

	file.SetContentType(contentType)

	log.Println("下载完成.", res.ContentLength)

	if b, err := ioutil.ReadAll(res.Body); err == nil {

		if size, err := file.Write(b); err == nil {
			dbs.Update("bee_downs", configs.M{"content_type": contentType, "file_size": size, "status": 1}, "name=?", name)
		} else {
			dbs.Update("bee_downs", configs.M{"content_type": contentType, "file_size": size, "status": 2}, "name=?", name)
		}
		log.Println("保存成功.")
	} else {
		log.Println(err)
		return
	}
}
