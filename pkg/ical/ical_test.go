package ical

import (
	"fmt"
	"testing"
	"time"
)

func TestVEvent(t *testing.T) {
	//print(time.Now().UTC().Format("20060102T150405Z"))
	fmt.Printf("%s\n", time.Now().String())
	a, _ := time.ParseInLocation("20060102", "20220908", time.Local)
	fmt.Printf("%s\n", a.String())
}
