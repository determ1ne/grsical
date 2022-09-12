package timetable

import (
	"time"
)

var CSTLocation *time.Location

func init() {
	var err error
	CSTLocation, err = time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
}
