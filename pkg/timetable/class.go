package timetable

import (
	"fmt"
	"strings"
)

type Repeat uint8

const (
	EveryWeek Repeat = iota
	// TODO: 更多种课程类型
)

var RepeatDesc = map[Repeat]string{
	EveryWeek: "每周",
}

type Term uint8

const (
	Autumn Term = iota
	Winter
	AutumnWinter
	Spring
	Summer
	SpringSummer
)

var TermDesc = map[Term]string{
	Autumn:       "秋",
	Winter:       "冬",
	AutumnWinter: "秋冬",
	Spring:       "春",
	Summer:       "夏",
	SpringSummer: "春夏",
}

var ClassStart = map[int]int64{
	// 以分钟计
	1:  480,  // 8:00
	2:  530,  // 8:50
	3:  590,  // 9:50
	4:  640,  // 10:40
	5:  690,  // 11:30
	6:  795,  // 13:15
	7:  845,  // 14:05
	8:  895,  // 14:55
	9:  955,  // 15:55
	10: 1005, // 16:45
	11: 1110, // 18:30
	12: 1160, // 19:20
	13: 1210, // 20:10
}

type ClassDuration struct {
	Starts int
	Ends   int
}

// Class 代表一次课，并非代表一类课
type Class struct {
	Name        string
	Term        Term
	Repeat      Repeat
	Duration    ClassDuration
	Teacher     string
	Location    string
	DayOfWeek   int    // 星期一为1
	RawDuration string // 教务网的时间文本，作为冗余
}

func (c *Class) ToDesc() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("教师：%s\\n", c.Teacher))
	b.WriteString(fmt.Sprintf("时间安排：%s %s %s", TermDesc[c.Term], RepeatDesc[c.Repeat], c.RawDuration))
	return b.String()
}
