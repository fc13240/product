package dbs

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"helper/configs"
	"helper/util"
	"log"
	"regexp"
	"strings"
	"errors"
)

var link *sql.DB = nil

func Conn() {
	var err error
	options := configs.GetSection("mysql")
	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", options["user"], options["password"], options["host"], options["port"], options["dbname"])

	link, err = sql.Open("mysql", url)
	if err != nil {
		link.Close()
		log.Println("数据库连接失败", err.Error())
		link = nil
	} else {
		log.Println("连接mysql... Succ")
	}
}

func Def() *Connect{

	if link == nil {
		 Conn()
	}
	if err:=link.Ping();err!=nil{
		Conn()
	}
	return &Connect{link: link}
}


func New(s string) (*Connect, error) {
	conn, err := sql.Open("mysql", s+"?charset=utf8")

	if err != nil {
		return nil, err
	}

	return &Connect{link: conn}, nil
}

type Connect struct {
	link *sql.DB
}

func (db *Connect) One(sql string, args ...interface{}) *sql.Row {
	return db.link.QueryRow(sql, args...)
}

func (db *Connect) Rows(query string, args ...interface{}) *sql.Rows {
	rows, errs := db.link.Query(query, args...)
	if errs != nil {
		log.Println("connect failing mysql.")
		log.Println(errs.Error())
	}
	return rows
}

func (db *Connect) Insert(tab string, data configs.M) (id int,err error) {
	fields := []string{}
	values := []interface{}{}
	q := []string{}

	for field, v := range data {
		fields = append(fields, field)
		q = append(q, "?")
		values = append(values, v)
	}

	sql := "INSERT INTO `" + tab + "`(" + strings.Join(fields, ",") + ")VALUE(" + strings.Join(q, ",") + ")"
	
	stmt := db.Prepare(sql)

	defer stmt.Close()
	if r, err := stmt.Exec(values...); err == nil {
		var v int64
		v, err= r.LastInsertId()
		if err!=nil || v==0{
			return 0,errors.New(fmt.Sprint("insert err:",err))
		}
		id=int(v)
	} else {
		err=errors.New(fmt.Sprintf("insert error :%s sql:%s",err.Error(),sql))
	}
	return id,err
}

func (db *Connect) Update(tab string, data configs.M, where string, args ...interface{}) error {
	values := []interface{}{}
	fields := []string{}
	for field, v := range data {
		fields = append(fields, field+"=?")
		values = append(values, v)
	}
	values = append(values, args...)

	return db.Exec("UPDATE `"+tab+"` SET "+strings.Join(fields, ",")+" WHERE "+where, values...)
}

func (db *Connect) Prepare(query string) *sql.Stmt {
	stmt, err := db.link.Prepare(query)
	if err != nil {
		log.Println(err)
	}
	return stmt
}

func (db *Connect)Count(sql string, args ...interface{}) int {
	if m, _ := regexp.MatchString("COUNT", sql); m {
		re, _ := regexp.Compile("SELECT.*FROM")
		sql = re.ReplaceAllString(sql, "SELECT COUNT(*) AS co FROM")
	}
	var n int
	db.One(sql, args...).Scan(&n)
	return n
}

func (db *Connect) Exec(sql string, args ...interface{}) error {
	stmt, err := db.link.Prepare(sql)

	defer stmt.Close()

	if err != nil {
		return err
	}
	if result, err := stmt.Exec(args...); err == nil {
		if _, err := result.RowsAffected(); err != nil {
			return err
		}
	}else{
		return err
	}
	return nil
}

func One(sql string, args ...interface{}) *sql.Row {
	return link.QueryRow(sql, args...)
}

func Rows(query string, args ...interface{}) *sql.Rows {
	rows, errs := link.Query(query, args...)
	if errs != nil {
		log.Println("connect failing mysql.")
		log.Println(errs.Error())
	}
	return rows
}

func Prepare(query string) *sql.Stmt {
	stmt, err := link.Prepare(query)
	if err != nil {
		log.Println(err)
	}
	return stmt
}

func Close() {
	link.Close()
	link = nil
}

func Exec(sql string, args ...interface{}) error {
	stmt, err := link.Prepare(sql)
	defer stmt.Close()
	if err != nil {
		log.Println(err.Error())
		return err
	}
	if result, err := stmt.Exec(args...); err == nil {
		if _, err := result.RowsAffected(); err != nil {
			log.Println("Failing", err.Error())
			return err
		}
	}
	return nil
}

func Insert(tab string, data configs.M) (int, error) {
	fields := []string{}
	values := []interface{}{}
	q := []string{}

	for field, v := range data {
		fields = append(fields, field)
		q = append(q, "?")
		values = append(values, v)
	}

	sql := "INSERT INTO `" + tab + "`(" + strings.Join(fields, ",") + ")VALUE(" + strings.Join(q, ",") + ")"
	stmt := Prepare(sql)

	defer stmt.Close()
	if r, err := stmt.Exec(values...); err == nil {
		v, _ := r.LastInsertId()
		return int(v), nil
	} else {
		return 0, err
	}
}

func Replace(tab string, data configs.M) (int, error) {
	fields := []string{}
	values := []interface{}{}
	q := []string{}

	for field, v := range data {
		fields = append(fields, field)
		q = append(q, "?")
		values = append(values, v)
	}

	sql := "REPLACE INTO `" + tab + "`(" + strings.Join(fields, ",") + ")VALUE(" + strings.Join(q, ",") + ")"
	stmt := Prepare(sql)

	defer stmt.Close()
	if r, err := stmt.Exec(values...); err == nil {
		v, _ := r.LastInsertId()
		return int(v), nil
	} else {
		return 0, err
	}
}

func Update(tab string, data configs.M, where string, args ...interface{}) error {
	values := []interface{}{}
	fields := []string{}
	for field, v := range data {
		fields = append(fields, field+"=?")
		values = append(values, v)
	}
	values = append(values, args...)

	return Exec("UPDATE `"+tab+"` SET "+strings.Join(fields, ",")+" WHERE "+where, values...)
}

func Limit(offset, rowCount int) string {
	if rowCount == -1 {
		return " "
	}
	return " LIMIT " + util.NumberToString(offset) + "," + util.NumberToString(rowCount)
}

func Count(sql string, args ...interface{}) int {
	if m, _ := regexp.MatchString("COUNT", sql); m {
		re, _ := regexp.Compile("SELECT.*FROM")
		sql = re.ReplaceAllString(sql, "SELECT COUNT(*) AS co FROM")
	}
	var n int
	One(sql, args...).Scan(&n)
	return n
}

type Tx struct {
	link *sql.Tx
}

func (db *Tx) One(sql string, args ...interface{}) *sql.Row {
	return db.link.QueryRow(sql, args...)
}

func (db *Tx) Rows(query string, args ...interface{}) *sql.Rows {
	rows, errs := db.link.Query(query, args...)
	if errs != nil {
		log.Println("connect failing mysql.")
		log.Println(errs.Error())
	}
	return rows
}

func (db *Tx) Prepare(query string) *sql.Stmt {
	stmt, err := db.link.Prepare(query)
	if err != nil {
		log.Println(err)
	}
	return stmt
}

func (db *Tx) Exec(sql string, args ...interface{}) error {
	stmt, err := db.link.Prepare(sql)
	defer stmt.Close()

	if err != nil {
		log.Println(err.Error())
		return err
	}
	if result, err := stmt.Exec(args...); err == nil {
		if _, err := result.RowsAffected(); err != nil {
			log.Println("Failing", err.Error())
			return err
		}
	}
	return nil
}

func (db *Tx) Insert(tab string, data configs.M) (int, error) {
	fields := []string{}
	values := []interface{}{}
	q := []string{}

	for field, v := range data {
		fields = append(fields, field)
		q = append(q, "?")
		values = append(values, v)
	}

	sql := "INSERT INTO `" + tab + "`(" + strings.Join(fields, ",") + ")VALUE(" + strings.Join(q, ",") + ")"
	stmt := db.Prepare(sql)

	defer stmt.Close()
	if r, err := stmt.Exec(values...); err == nil {
		v, _ := r.LastInsertId()
		return int(v), nil
	} else {
		return 0, err
	}
}

func (db *Tx) Update(tab string, data configs.M, where string, args ...interface{}) error {
	values := []interface{}{}
	fields := []string{}
	for field, v := range data {
		fields = append(fields, field+"=?")
		values = append(values, v)
	}
	values = append(values, args...)

	return db.Exec("UPDATE `"+tab+"` SET "+strings.Join(fields, ",")+" WHERE "+where, values...)
}

func (db *Tx) Limit(offset, rowCount int) string {
	return " LIMIT " + util.NumberToString(offset) + "," + util.NumberToString(rowCount)
}

func (db *Tx) Count(sql string, args ...interface{}) int {
	if m, _ := regexp.MatchString("COUNT", sql); m {
		re, _ := regexp.Compile("SELECT.*FROM")
		sql = re.ReplaceAllString(sql, "SELECT COUNT(*) AS co FROM")
	}
	var n int
	db.One(sql, args...).Scan(&n)
	return n
}

func (db *Tx) Commit() {
	db.link.Commit()
}

func (db *Tx) Rollback() {
	db.link.Rollback()
}

func (c *Connect)Begin() (*Tx, error) {
	tx, err := c.link.Begin()
	return &Tx{link: tx}, err
}

type Where struct{
	sql []string
}

func (w *Where)And(s string ,f ... interface{}) *Where{
	s=fmt.Sprintf(s,f...)
	w.sql=append(w.sql, " AND ( "+s+" ) ")
	return w
}

func (w *Where)AndIn(s string ,in []string) *Where{
	ss:=[]string{}
	for _,v:=range in{
		ss=append(ss,fmt.Sprintf("'%s'",v))
	}
	w.sql=append(w.sql, fmt.Sprintf(" AND %s IN ( %s ) ",s,strings.Join(ss,",") ))
	return w
}

func (w *Where)AndIntIn(s string ,in[]int)*Where{
	w.sql=append(w.sql, fmt.Sprintf(" AND %s IN ( %s ) ",s,util.IntJoin(in,",") ))
	return w
}

func (w *Where)ToString()string {
	return strings.Join(w.sql," ")
}


func NewWhere()*Where{
	w:= &Where{}
	w.sql=append(w.sql," WHERE 1=1 ")
	return w
}