package main

import (
	"fmt"
	"time"

	"github.com/mergestat/timediff"
)

func main() {
	// layout := "2025-04-03 17:44:09.699692 +0700 +07 m=+0.000194084"
	layout := "2006-01-02 15:04:05 MST"
	curTime := time.Now()
	strTime := curTime.Format(layout)
	str1 := timediff.TimeDiff(curTime.Add(-10 * time.Second))
	fmt.Println(curTime.String())
	fmt.Println(str1) // a few seconds ago
	// fmt.Println(curTime.Format("2006-01-02 15:04:05 MST"))
	fmt.Printf("%T\t%v\n", strTime, strTime)
	// fmt.Println(time.Parse(layout, strTime))

	cvtTime, _ := time.Parse(layout, strTime)
	fmt.Println(cvtTime)
	fmt.Println(timediff.TimeDiff(cvtTime.Add(-10 * time.Second)))

}
