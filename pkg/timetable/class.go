package timetable

import (
	"fmt"
	"strings"
	"time"
)

type Repeat uint8

const (
	EveryWeek Repeat = iota
	SingleWeek
	DoubleWeek
	// TODO: 更多种课程类型
)

var RepeatDesc = map[Repeat]string{
	EveryWeek:  "每周",
	SingleWeek: "单周",
	DoubleWeek: "双周",
}

type Semester uint8

const (
	Autumn Semester = iota
	Winter
	AutumnWinter
	Spring
	Summer
	SpringSummer
)

var SemesterDesc = map[Semester]string{
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
	3:  600,  // 10:00
	4:  650,  // 10:50
	5:  700,  // 11:40
	6:  805,  // 13:25
	7:  855,  // 14:15
	8:  905,  // 15:05
	9:  975,  // 16:15
	10: 1025, // 17:05
	11: 1130, // 18:50
	12: 1180, // 19:40
	13: 1230, // 20:30
	14: 1280, // 21:20
	15: 1330, // 22:10
}

type ClassDuration struct {
	Starts int
	Ends   int
}

// Class 代表一次课，并非代表一类课
type Class struct {
	Name        string
	Semester    Semester
	Repeat      Repeat
	Duration    ClassDuration
	Teacher     string
	Location    string
	DayOfWeek   int    // 星期一为1
	RawDuration string // 教务网的时间文本，作为冗余
	date        time.Time
	tweakDesc   string
}

type Exam struct {
	Name      string
	Semester  string
	ID        string
	Region    string
	StartTime time.Time
	EndTime   time.Time
	Location  string
	SeatNo    string
	Remark    string
}

func (c *Class) ToDesc(weekNo int) string {
	var b strings.Builder
	switch c.Repeat {
	case SingleWeek:
		b.WriteString(fmt.Sprintf("%s\\n", RepeatDesc[SingleWeek]))
	case DoubleWeek:
		b.WriteString(fmt.Sprintf("%s\\n", RepeatDesc[DoubleWeek]))
	}
	b.WriteString(fmt.Sprintf("第 %d 周\\n", weekNo+1))
	b.WriteString(fmt.Sprintf("教师：%s\\n", c.Teacher))
	b.WriteString(fmt.Sprintf("时间安排：%s %s %s", SemesterDesc[c.Semester], RepeatDesc[c.Repeat], c.RawDuration))
	if c.tweakDesc != "" {
		b.WriteString(fmt.Sprintf("\\n%s", c.tweakDesc))
	}
	return b.String()
}

func (exam *Exam) ToDesc() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%s %s学期", exam.Name, exam.Semester))
	b.WriteString(fmt.Sprintf("课程号：%s", exam.ID))
	if exam.SeatNo != "" {
		b.WriteString(fmt.Sprintf("座位号: %s\\n", exam.SeatNo))
	}
	if exam.Remark != "" {
		b.WriteString(fmt.Sprintf("备注: %s\\n", exam.Remark))
	}
	return b.String()
}
