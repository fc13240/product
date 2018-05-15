package template

import (
	"bytes"
	"fmt"
	"helper/configs"
	"helper/dbs"
	"html"
	"html/template"
	"io"
	"product"
)

type Template struct {
	Id      int
	Tag     string
	Content string
}

func (t *Template) Save() error {
	if t.ExistTag() {
		err := dbs.Update("product_template", configs.M{
			"content": t.Content,
		}, "id=?", t.Id)
		return err
	} else {
		id, err := dbs.Insert("product_template", configs.M{
			"tag":     t.Tag,
			"content": t.Content,
		})
		if err == nil {
			t.Id = id
		}
		return err
	}
}

func (t *Template) ExistTag() bool {
	var id int
	dbs.One("SELECT id FROM product_template WHERE tag=?", t.Tag).Scan(&id)
	if id > 0 {
		t.Id = id
		return true
	} else {
		return false
	}
}

func GetTemplete(tmp_id int) (*Template, error) {
	tmp := &Template{}
	err := dbs.One("SELECT id,tag,content FROM product_template WHERE id=?", tmp_id).Scan(&tmp.Id, &tmp.Tag, &tmp.Content)
	return tmp, err
}

func GetTempleteByTag(tag string) (*Template, error) {
	tmp := &Template{}
	err := dbs.One("SELECT id,tag,content FROM product_template WHERE tag=?", tag).Scan(&tmp.Id, &tmp.Tag, &tmp.Content)
	return tmp, err
}

func GetTemplateTags() []configs.M {
	list := []configs.M{}
	rows := dbs.Rows("SELECT id,tag FROM product_template")

	for rows.Next() {
		var id int
		var tag string
		rows.Scan(&id, &tag)
		list = append(list, configs.M{"id": id, "tag": tag})
	}
	return list
}

func GetTemplateContent(tag string) string {
	var content string
	dbs.One("SELECT content FROM product_template WHERE `tag`=?", tag).Scan(&content)
	return content
}

func SaveTemplateContent(tag string) {

}

func MergeContent(tmp_id int, body string) (content string, err error) {

	var tmp *Template
	if err == nil {
		tmp, err = GetTemplete(tmp_id)
		if err == nil {
			return tmp.Merge(body), nil
		}
	}
	return content, err
}

func MergeItem(item *product.Item) (content string, err error) {
	if err == nil {
		return MergeContent(item.Tempid, item.Desc)
	}
	return content, err
}

func MergeExec(wr io.Writer, item_id, tmp_id int) error {
	var item *product.Item
	var tmp *Template
	var err error

	item, err = product.IdGet(item_id)

	if err != nil {
		return err
	}

	tmp, err = GetTemplete(tmp_id)
	if err != nil {
		return err
	}

	temp := template.New("")

	_, e := temp.Parse(tmp.Content)

	if e != nil {
		return err
	} else {

		temp.Execute(wr, map[string]string{"body": html.UnescapeString(item.Desc)})
	}
	return nil

}

func (t *Template) Merge(body string) string {

	var buff bytes.Buffer

	temp := template.New("")

	//t.Funcs(template.FuncMap{"link": self.Link, "li": self.Li, "url": self.Url, "js": self.Js, "ads": self.Ads})

	_, e := temp.Parse(t.Content)

	if e != nil {
		fmt.Println(e)
	} else {
		temp.Execute(&buff, map[string]string{"body": body})
	}

	return html.UnescapeString(buff.String())
}
