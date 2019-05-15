package doc

import (
	"testing"

	"git.cloud.top/pei_hainan/ooxml/doc"
)

func TestMain(t *testing.T) {

	tabledata := make(map[string]interface{})
	tabledata["label"] = []string{"a", "b", "c", "d"}
	tabledata["data"] = [][]string{
		[]string{"A1", "B2", "C3", "D4"},
		[]string{"A1", "B2", "C3", "D4"},
		[]string{"A1", "B2", "C3", "D4"},
		[]string{"A1", "B2", "C3", "D4"},
	}

	qr := doc.Unit{
		Title:    "Scene One Title",
		Subtitle: "Scene One subtitle",
		Pretb:    "pre desc xxxxxxxx",
		Tbdata:   tabledata,
		Subtb:    "sub desc xxxxxxxx",
		Preimg:   "pre img desc xxxx",
		Imgpath:  "../image/bar.png",
		Subimg:   "sub img desc xxxx",
	}

	ww := make([]doc.Unit, 2)
	ww[0] = qr
	ww[1] = qr

	doc.ChartData("/tmp/nosaynever", ww)
}
