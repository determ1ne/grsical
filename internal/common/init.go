package common

import (
	"time"
)

func init() {
	var err error
	CSTLocation, err = time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
}
