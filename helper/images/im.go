package images

import (
	"encoding/base32"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"resize"
	"strconv"
	"strings"
	"time"
	"io"

)

const images_upload_dir = "./uploads"

type Img struct {
	Name string
	Path string
	Ext  string
	Href string
	head *multipart.FileHeader
	file multipart.File
}

func Create(filename string) *Img {
	img := &Img{Name: filename}
	img.file, _ = os.Open(img.Name)
	return img
}

func Open(filename string) (*Img, error) {
	ext := filepath.Ext(filename)
	path := filepath.Dir(filename)
	name := filepath.Base(filename)
	name = strings.Replace(name, ext, "", 1)
	img := &Img{Name: name, Ext: ext, Path: path}
	var err error
	img.file, err = os.Open(filename)
	return img, err
}

func (im *Img) Save() error {
	defer im.file.Close()

	body, err := ioutil.ReadAll(im.file)
	if err != nil {
		return err
	}

	file, err := os.Create(im.FullName())
	if err != nil {
		return err
	}

	if _,err:=file.Write(body); err!=nil{
		return err
	}
	defer file.Close()

	return nil
}

func (im *Img) Body() []byte {
	body, _ := ioutil.ReadAll(im.file)
	return body
}



func (im *Img) FullName() string {
	fullPath := filepath.Join(im.Path, im.Name+im.Ext)
	return fullPath
}

func RmExt(name string) string {
	return strings.Replace(name, filepath.Ext(name), "", 1)
}

func (im *Img) String() string {
	return im.Name + im.Ext
}

func (im *Img) SaveAs(path, new_name string) string {

	if new_name == "" {
		new_name = RandName() + im.Ext
	} else {
		new_name = new_name + im.Ext
	}
	file, err := os.Create(filepath.Join(path, new_name))
	body, _ := ioutil.ReadAll(im.file)
	E(err)
	file.Write(body)
	defer im.file.Close()
	defer file.Close()
	return new_name
}

func (im *Img) WriteFrom(r *http.Request, name string) error {
	var err error
	im.file, im.head, err = r.FormFile(name)
	if err != nil {
		fmt.Println(err)
		return err
	}
	im.Ext = filepath.Ext(im.head.Filename)
	return nil
}

func New(path, name string) *Img {
	if name == "" {
		name = RandName()
	} else {
		name = RmExt(name)
	}
	return &Img{Path: path, Name: name}
}

func (self *Img) Decode() (image.Image, error) {
	im,err:=Decode(self.file,self.Ext)
	E(err)
	return im, err
}

func Decode(file io.Reader,ext string) (image.Image, error) {
	var im image.Image
	var err error

	switch ext  {
	case ".jpg":
		im, err = jpeg.Decode(file)
	case ".png":
		im, err = png.Decode(file)
	case ".gif":
		im, err = gif.Decode(file)
	default:
		errors.New("无效的图片")
	}
	return im, err
}

//等比例缩放
func Resize(filename string, width, hight int) string {
	return ResizeRename(filename, RandName(), width, hight)
}

func ResizeRename(filename, new_name string, width, hight int) string {

	im, err := Open(filename)

	new_name = RmExt(new_name)

	file_path, _, file_ext := Info(filename)

	if err != nil {
		fmt.Println("打开原始图片失败")
		return ""
	}

	decode, err := im.Decode()

	if err != nil {
		fmt.Println("图片解码失败")
		return ""
	}

	dst := resize.Thumbnail(uint(width), uint(hight), decode, resize.NearestNeighbor)

	new_file_name := new_name + file_ext
	scale_name := filepath.Join(file_path, new_file_name)

	file, _ := os.Create(scale_name)

	switch file_ext {
	case ".jpg":
		jpeg.Encode(file, dst, &jpeg.Options{100})
	case ".png":
		png.Encode(file, dst)
	case ".gif":
		gif.Encode(file, dst, &gif.Options{256, nil, nil})
	}
	defer file.Close()
	defer im.file.Close()

	return new_file_name
}

//POST一张图片
func (self *Img) Post(uri string) error {

	body, _ := ioutil.ReadAll(self.file)

	data := base32.StdEncoding.EncodeToString(body)

	resp, err := http.PostForm(uri, url.Values{"file": {data}})

	defer resp.Body.Close()

	if err != nil {
		return err
	}

	if _, err = ioutil.ReadAll(resp.Body); err != nil {
		return err
	}

	return nil
}

//文件名是否
func (self *Img) Exist() bool {
	return false
}

//接受一个base32编码了的图片
func PostDecode(filename, data string) error {
	img := Img{Name: filename}

	if img.Exist() {
		errors.New("file name exist.")
	}

	new_file, err := os.Create(filename)

	if err != nil {
		return err
	}

	body, _ := base32.StdEncoding.DecodeString(data)
	new_file.Write(body)
	defer new_file.Close()
	return nil
}

func E(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

// 随机生成文件名
func RandName() string {
	rand.Seed(time.Now().UnixNano())
	return strconv.Itoa(rand.Int())
}

func Info(filename string) (string, string, string) {
	ext := filepath.Ext(filename)
	path := filepath.Dir(filename)
	name := filepath.Base(filename)
	name = strings.Replace(name, ext, "", 1)
	return path, name, ext
}

//下载图片
func Down(img_url, new_path string) error {
	res, err := http.Get(img_url)
	if err != nil {
		return err
	}
	file, _ := os.Create(new_path)
	defer file.Close()
	if b, err := ioutil.ReadAll(res.Body); err == nil {
		res.Body.Close()
		file.Write(b)
		return nil
	} else {
		fmt.Println(err)
		return err
	}
}

func DownTo(img_url, new_path string)error{
	res, err := http.Get(img_url)
	if err != nil {
		return err
	}
	file, _ := os.Create(new_path)
	defer file.Close()
	if b, err := ioutil.ReadAll(res.Body); err == nil {
		res.Body.Close()
		file.Write(b)
		return nil
	} else {
		fmt.Println(err)
		return err
	}
}

