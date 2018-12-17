package main

import (
	"fmt"
	"time"
)

func main() {

	t1 := "2017-03-15 13:30:39"
	t2 := "2017-03-17 13:31:39"

	fmt.Printf("The variable is %#v\n", its(t1, t2))
}

func its(t1, t2 string) (s []string) {

	f1, err := time.Parse("2006-01-02 15:04:05", t1)
	checkError(err)

	f2, err := time.Parse("2006-01-02 15:04:05", t2)
	checkError(err)

	bd, _ := time.ParseDuration("-24h")
	ad, _ := time.ParseDuration("24h")

	lastday := f2.Add(bd)

	for {
		afd := f1.Add(ad)
		if afd.After(lastday) {
			break
		}
		s = append(s, afd.Format("2006-01-02"))
		f1 = afd
	}

	return
}

func checkError(err error) {
	if err != nil {
		fmt.Printf("The err is %#v\n", err)
	}
}
