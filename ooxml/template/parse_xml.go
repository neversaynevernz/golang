package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"text/template"
	"time"

	"github.com/axgle/mahonia"
	"unicode/utf8"

	"image"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

func toUTF8(str string) string {
	if !utf8.ValidString(str) {
		return mahonia.NewDecoder("GBK").ConvertString(str)
	}
	fmt.Printf("The variable is %#v\n", str)
	return str
}

func RandomStr() string {
	rand.Seed(time.Now().Unix())
	return fmt.Sprintf("%d", rand.Intn(1000))
}

func GetKey(m map[string]interface{}, key string) interface{} {
	return m[key]
}

func adjustImage(s string, w int) (hsize int) {

	file, err := os.Open(s)
	defer file.Close()

	img, _, err := image.DecodeConfig(file)
	checkErr(err)

	width := img.Width
	height := img.Height

	wpercent := (float64(w) / float64(width))
	hsize = int(float64(height) * float64(wpercent))
	return hsize
}

func main() {

	/*
		info := map[string]interface{}{
			"Han Meimei": true,
			"LiLei":      false,
			"slice": [][]string{
				[]string{"1", "2", "3"},
			},
			"maps": []map[string]interface{}{
				map[string]interface{}{
					"key1": "value1",
					"key2": "value2",
					"key3": "value3",
				},
			},
		}
		fmt.Printf("The info is %#v\n", info)
	*/

	img1 := "/pics/bar.png"
	img2 := "/pics/pie.png"
	img3 := "/pics/line.png"

	// fmt.Printf("datas is %#v\n", datas)

	dd := make(map[string]interface{})

	// 表头输出顺序
	var xx []string
	xx = []string{"index", "s_srcip", "records", "s_dstip"}
	dd["order"] = xx

	// 表头翻译
	label := make(map[string]string)
	label["index"] = "序号"
	label["s_srcip"] = "源地址"
	label["s_dstip"] = "目的地址"
	label["records"] = "事件数"
	dd["label"] = label

	mm := make(map[string]interface{})
	mm["s_srcip"] = "192.168.18.1"
	mm["index"] = "序号1"
	mm["s_dstip"] = "192.168.76.18"
	mm["records"] = 8

	datas := make([]map[string]interface{}, 20)
	for i := 0; i < 20; i++ {
		datas[i] = mm
	}
	dd["datas"] = datas
	dd["pics"] = []interface{}{picData(img1, "A"), picData(img2, "B"), picData(img3, "C")}

	scenedatas := []interface{}{dd}

	ss := make(map[string]interface{})
	ss["sceneName"] = "scene name"
	ss["scenePattern"] = []string{
		"self define test follow your heart",
		"self define test follow your heart",
		"self define test follow your heart",
	}
	ss["sceneData"] = scenedatas

	logo := "/pics/logo.jpg"
	bs, err := ioutil.ReadFile(logo)
	checkErr(err)

	// 模板中限定 logo 大小
	reportData := map[string]interface{}{
		"LOGOPATH":     logo,
		"LOGODATA":     base64.StdEncoding.EncodeToString(bs),
		"TITLE":        "审计报告",
		"SUBTITLE":     "报表",
		"COMPANY":      "天融信公司",
		"REPORTER":     "天融信TAAA",
		"CREATETIME":   time.Now().Format("2006-01-02"),
		"REPORTHEADER": "My 安全报告",
		"pagedatas": []interface{}{
			ss,
			ss,
			ss,
		},
	}

	createReport(reportData, "/", ".doc")
}

func nowYMDStr() string {
	ss := time.Now().Unix()
	ww := time.Unix(ss, 0)
	return ww.Format("2006-01-02")
}

func createReport(reportData map[string]interface{}, outpath, format string) {

	fn := "NeverSayNever-" + nowYMDStr() + format
	filename := path.Join(outpath, fn)

	f := open(filename)
	defer f.Close()

	// 注册模板函数
	funcMap := template.FuncMap{"RandomStr": RandomStr, "GetKey": GetKey}
	t := template.New("T1").Funcs(funcMap)
	t = template.Must(t.ParseFiles("/var/report_template/word.tpl"))
	err := t.ExecuteTemplate(f, "word.tpl", reportData)
	checkErr(err)
}

func picData(fname, desc string) map[string]interface{} {
	d := make(map[string]interface{})
	d["title"] = "TOPSEC"
	d["height"] = adjustImage(fname, 420)

	bs, err := ioutil.ReadFile(fname)
	checkErr(err)
	d["data"] = base64.StdEncoding.EncodeToString(bs)
	d["pathname"] = fname
	d["width"] = 420
	d["desc"] = desc
	return d
}

func checkErr(err error) {
	if err != nil {
		fmt.Printf("The err is %v\n", err)
	}
}

func open(filepath string) *os.File {
	e, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)
	checkErr(err)
	return e
}
