package main

import (
	"fmt"
	"time"
)

const base_format = "2006-01-02 15:04:05"
const myd_format = "2006-01-02"

func str2stamp(s string) int64 {
	t, _ := time.Parse(base_format, s)
	return t.Unix()
}

func stamp2str(n int64) string {
	ww := time.Unix(n, 0)
	return ww.Format(base_format)
}

func stamp2ymdstr(n int64) string {
	ww := time.Unix(n, 0)
	return ww.Format(myd_format)
}

// acquire now time.Time
func now() time.Time {
	return time.Now()
}

// acquire now
// 2019-01-01 10:00:00
func nowstr() string {
	return now().Format(base_format)
}

// acquire now timestamp
// 1547464640
func nowstamp() int64 {
	return now().Unix()
}

func main() {
	// timestr := "2018-01-01 10:00:00"
	timestr := nowstr()
	fmt.Printf("convert %q to timestamp: %v\n", timestr, str2stamp(timestr))

	timestamp := nowstamp()
	fmt.Printf("convert %v to timestr: %v\n", timestamp, stamp2str(timestamp))
	fmt.Printf("convert %v to y-m-d timestr : %v\n", timestamp, stamp2ymdstr(timestamp))
}
