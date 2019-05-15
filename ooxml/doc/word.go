package doc

import (
	"fmt"
	"io"
	"log"
	"path/filepath"
	"strings"
	"time"

	// "baliance.com/gooxml/color"
	"github.com/unidoc/unioffice/color"
	"github.com/unidoc/unioffice/common"
	"github.com/unidoc/unioffice/document"
	"github.com/unidoc/unioffice/measurement"
	"github.com/unidoc/unioffice/schema/soo/wml"

	tp "git.cloud.top/report/topreport"
)

// 各种图的大小
const (
	AXISX = 450
	AXISY = 250

	BARX = 450
	BARY = 300

	PIEX = 250
	PIEY = 250

	DONUTX = 250
	DONUTY = 250

	STACKEDBARX = 300
	STACKEDBARY = 300
)

var Hdr document.Header
var Ftr document.Footer

type srpDoc struct {
	d *document.Document
}

// 初始化 header, footer 数据
func (s *srpDoc) initHeaderFooter(h, f string) {

	Hdr = s.d.AddHeader()
	para := Hdr.AddParagraph()
	run := para.AddRun()

	// 灰色下划线 以及 header 内容
	run.Properties().SetUnderline(wml.ST_UnderlineSingle, color.Gray)
	run.AddText(h)

	Ftr = s.d.AddFooter()
	para = Ftr.AddParagraph()
	para.Properties().AddTabStop(6*measurement.Inch, wml.ST_TabJcRight, wml.ST_TabTlcNone)
	run = para.AddRun()
	run.AddText(f)
	run.AddTab()
	run.AddText("第 ")
	run.AddField(document.FieldCurrentPage)
	run.AddText("页 共")
	run.AddField(document.FieldNumberOfPages)
	run.AddText(" 页")
}

// 首页数据
// imgpath :logo path
// title subtitle company reporter rtime

func (s *srpDoc) setFirstPage(h Home) {

	if h.Logopath == "" {
		h = defaultHome
	}

	img, err := common.ImageFromFile(h.Logopath)
	checkError(err)

	iref, err := s.d.AddImage(img)
	checkError(err)

	para := s.d.AddParagraph()
	anchored, err := para.AddRun().AddDrawingAnchored(iref)
	anchored.SetOffset(1*measurement.Inch, 1*measurement.Inch)

	para = s.d.AddParagraph()
	para.Properties().SetAlignment(wml.ST_JcCenter)
	para.SetStyle("Title")

	run := para.AddRun()
	nBlank(run, 3)
	run.Properties().SetBold(true)
	run.AddText(h.Title)

	para = s.d.AddParagraph()
	para.Properties().SetAlignment(wml.ST_JcCenter)

	run = para.AddRun()
	// 设置 14px 字体
	run.Properties().SetSize(measurement.Pixel72 * 14)
	nBlank(run, 1)
	run.AddText(h.Subtitle)

	nBlank(run, 16)

	para = s.d.AddParagraph()
	para.Properties().SetAlignment(wml.ST_JcCenter)
	run = para.AddRun()
	run.Properties().SetSize(measurement.Pixel72 * 14)
	run.AddText(h.Company)
	nBlank(run, 1)

	run.AddText(h.Reporter)
	nBlank(run, 1)
	run.AddText(h.Rtime)

	s.setHeaderFooter()
}

// 设置分页符并设置页眉页脚
func (s *srpDoc) setHeaderFooter() {
	para := s.d.AddParagraph()
	sect := para.Properties().AddSection(wml.ST_SectionMarkNextPage)
	sect.SetHeader(Hdr, wml.ST_HdrFtrDefault)
	sect.SetFooter(Ftr, wml.ST_HdrFtrDefault)
}

func nBlank(r document.Run, n int) {
	m := 0
	for m <= n {
		r.AddBreak()
		m++
	}
}

type Home struct {
	Logopath, Title, Subtitle string
	Company, Reporter, Rtime  string
}

var defaultHome = Home{
	Logopath: "/var/report_template/logo.jpg",
	Title:    "TopSec",
	Subtitle: "TopSec",
	Company:  fmt.Sprintf("公司: %s", "TopSec"),
	Reporter: fmt.Sprintf("报告人: %s", "Topsec"),
	Rtime:    fmt.Sprintf("时间: %s", time.Now().Format("2006-01-02")),
}

// 初始化加载格式模板文件
func InitConfig(h Home) srpDoc {

	tdoc, err := document.OpenTemplate("/var/report_template/template.docx")
	checkError(err)

	doc := srpDoc{d: tdoc}

	doc.initHeaderFooter("Header", "Footer")

	doc.setFirstPage(h)

	// Headers and footers are not immediately associated with a document as a
	// document can have multiple headers and footers for different sections.
	doc.d.BodySection().SetHeader(Hdr, wml.ST_HdrFtrDefault)
	doc.d.BodySection().SetFooter(Ftr, wml.ST_HdrFtrDefault)

	return doc
}

func checkError(err error) {
	if err != nil {
		log.Printf("ooxml doc error : %s\n", err)
	}
}

// 实现该接口
type Reporter interface {
	Render(r ReportData, w io.Writer) error
	RenderToFile(r ReportData, fp string) error
	Name() string
}

type DocReport struct {
	*srpDoc
}

func (s *DocReport) Render(r tp.ReportData, w io.Writer) error     {}
func (s *DocReport) RenderToFile(r tp.ReportData, fp string) error {}
func (s *DocReport) Name() string                                  {}

func (s *DocReport) Render(r tp.ReportData, w io.Writer) error {
	// func ChartData(name string, scenes []Unit) {

	var h Home
	doc := InitConfig(h)

	// 分页符个数
	count := 0
	ass := len(scenes)

	for i, data := range scenes {

		log.Print("@@@@", i)
		log.Print("@@@@ data ", data)

		para := doc.d.AddParagraph()
		para.SetStyle("Title")
		para.Properties().SetAlignment(wml.ST_JcCenter)
		para.AddRun().AddText(data.Title)

		para = doc.d.AddParagraph()
		para.Properties().SetAlignment(wml.ST_JcCenter)
		para.SetStyle("Subtitle")
		para.AddRun().AddText(data.Subtitle)

		// 表格数据处理
		if data.Tbdata != nil {

			para = doc.d.AddParagraph()
			run := para.AddRun()
			run.AddText(data.Pretb)
			// nBlank(run, 1)

			// using a pre-defined table style
			table := doc.d.AddTable()
			table.Properties().SetWidthPercent(95)
			table.Properties().SetAlignment(wml.ST_JcTableCenter)
			table.Properties().SetStyle("GridTable4-Accent1")

			// 处理具体数据
			labelx := len(data.Tbdata["label"].([]string))
			rowy := len(data.Tbdata["data"].([][]string))

			row := table.AddRow()
			for r := 0; r < labelx; r++ {
				cell := row.AddCell()
				cell.AddParagraph().AddRun().AddText(data.Tbdata["label"].([]string)[r])
			}

			for r := 0; r < rowy; r++ {
				row := table.AddRow()
				for c := 0; c < labelx; c++ {
					cell := row.AddCell()
					cell.AddParagraph().AddRun().AddText(data.Tbdata["data"].([][]string)[r][c])
				}
			}

			para = doc.d.AddParagraph()
			run = para.AddRun()
			// nBlank(run, 1)
			run.AddText(data.Subtb)
		}

		if data.Imgpath != "" {

			para = doc.d.AddParagraph()
			run := para.AddRun()
			nBlank(run, 1)
			run.AddText(data.Preimg)

			img1, err := common.ImageFromFile(data.Imgpath)
			if err != nil {
				log.Fatalf("unable to create image: %s", err)
			}

			fp := filepath.Base(img1.Path)
			switch {
			case strings.Contains(fp, "axis"):
				img1.Size.X = AXISX
				img1.Size.Y = AXISY
			case strings.Contains(fp, "pie"):
				img1.Size.X = PIEX
				img1.Size.Y = PIEY
			case strings.Contains(fp, "stacked"):
				img1.Size.X = STACKEDBARX
				img1.Size.Y = STACKEDBARY
			case strings.Contains(fp, "bar"):
				img1.Size.X = BARX
				img1.Size.Y = BARY
			case strings.Contains(fp, "donut"):
				img1.Size.X = DONUTX
				img1.Size.Y = DONUTY
			}

			img1ref, err := doc.d.AddImage(img1)
			if err != nil {
				log.Fatalf("unable to add image to document: %s", err)
			}

			para := doc.d.AddParagraph()
			para.Properties().SetAlignment(wml.ST_JcCenter)
			_, err = para.AddRun().AddDrawingInline(img1ref)

			if err != nil {
				log.Fatalf("unable to add image: %s", err)
			}

			para = doc.d.AddParagraph()
			run = para.AddRun()
			// nBlank(run, 1)
			run.AddText(data.Subimg)
		}

		if count < ass-1 {
			doc.setHeaderFooter()
		}
		count = count + 1
	}

	doc.d.SaveToFile(name + ".docx")
	fmt.Printf("save xxxx.docx done")
}

// unit 报表数据的一个单元 包含一个统计场景
type Unit struct {
	Title, Subtitle         string                 // 标题 副标题
	Pretb, Subtb            string                 // 表格上下描述
	Preimg, Subimg, Imgpath string                 // 图片路径以及上下描述
	Tbdata                  map[string]interface{} //表格数据
}

func demo() {

	tabledata := make(map[string]interface{})
	tabledata["label"] = []string{"a", "b", "c", "d"}
	tabledata["data"] = [][]string{
		[]string{"A1", "B2", "C3", "D4"},
		[]string{"A1", "B2", "C3", "D4"},
		[]string{"A1", "B2", "C3", "D4"},
		[]string{"A1", "B2", "C3", "D4"},
	}

	qr := Unit{
		Title:    "Scene One Title",
		Subtitle: "Scene One subtitle",
		Pretb:    "pre desc xxxxxxxx",
		Tbdata:   tabledata,
		Subtb:    "sub desc xxxxxxxx",
		Preimg:   "pre img desc xxxx",
		Imgpath:  "../image/bar.png",
		Subimg:   "sub img desc xxxx",
	}

	ww := make([]Unit, 2)
	ww[0] = qr
	ww[1] = qr

	ChartData("nosaynever", ww)
}
