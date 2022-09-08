package timetable

type Repeat uint8

const (
	EveryWeek Repeat = iota
	// TODO: 更多种课程类型
)

type Term uint8

const (
	Autumn Term = iota
	Winter
	AutumnWinter
	Spring
	Summer
	SpringSummer
)

var ClassStart = map[int]int{
	1:  800,
	2:  850,
	3:  950,
	4:  1040,
	5:  1130,
	6:  1315,
	7:  1405,
	8:  1455,
	9:  1555,
	10: 1645,
	11: 1830,
	12: 1920,
	13: 2010,
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
