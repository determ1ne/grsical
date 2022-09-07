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
