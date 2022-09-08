package ical

import (
	"crypto/sha1"
	"fmt"
	"strings"
	"time"
)

type DateTime struct {
	Year   int
	Month  int
	Day    int
	Hour   int
	Minute int
	Second int
}

func (d DateTime) ToISOString() string {
	return fmt.Sprintf("%04d%02d%02dT%02d%02d%02d", d.Year, d.Month, d.Day, d.Hour, d.Minute, d.Second)
}

const dateLayoutUTC = "20060102T150405Z"

type VEvent struct {
	Description string
	Summary     string
	Location    string
	StartTime   DateTime
	EndTime     DateTime
}

func (v *VEvent) GetHash() string {
	m := sha1.New()
	m.Write([]byte(v.Description))
	m.Write([]byte(v.Summary))
	m.Write([]byte(v.Location))
	m.Write([]byte(v.StartTime.ToISOString()))

	return fmt.Sprintf("%x", m.Sum([]byte("")))
}

func (v *VEvent) String() string {
	var b strings.Builder
	utcStr := time.Now().UTC().Format(dateLayoutUTC)
	etStr := v.EndTime.ToISOString()
	stStr := v.StartTime.ToISOString()

	b.WriteString(fmt.Sprintf("BEGIN:VEVENT\nCLASS:PUBLIC\nCREATED:%s\n", utcStr))
	if v.Description != "" {
		b.WriteString(fmt.Sprintf("DESCRIPTION: %s\n", v.Description))
	}
	b.WriteString(fmt.Sprintf("DTEND;TZID=\"China Standard Time\":%s\n", etStr))
	b.WriteString(fmt.Sprintf("DTSTAMP:%s\n", utcStr))
	b.WriteString(fmt.Sprintf("DTSTART;TZID=\"China Standard Time\":%s\n", stStr))
	b.WriteString(fmt.Sprintf("LAST-MODIFIED:%s", utcStr))
	if v.Location != "" {
		b.WriteString(fmt.Sprintf("LOCATION: %s\n", v.Location))
	}
	b.WriteString(fmt.Sprintf("SEQUENCE:0\nSUMMARY;LANGUAGE=zh-cn:{Summary}\nTRANSP:OPAQUE\nUID:%s\n", v.GetHash()))
	b.WriteString("BEGIN:VALARM\nTRIGGER:-PT15M\nACTION:DISPLAY\nDESCRIPTION:提醒\nEND:VALARM\nEND:VEVENT\n")

	return b.String()
}

type VCalendar struct {
	VEvents []*VEvent
}

func (v *VCalendar) GetICS(icalName string) string {
	var b strings.Builder
	if icalName == "" {
		icalName = "GRSICAL 课程表"
	}
	b.WriteString(fmt.Sprintf("BEGIN:VCALENDAR\nX-WR-CALNAME:%s\nX-APPLE-CALENDAR-COLOR:#2BBFF0\nPRODID:-//Azuk Workshop//Ejector 0.2//EN\nVERSION:2.0\nMETHOD:PUBLISH\nBEGIN:VTIMEZONE\nTZID:China Standard Time\nBEGIN:STANDARD\nDTSTART:16010101T000000\nTZOFFSETFROM:+0800\nTZOFFSETTO:+0800\nEND:STANDARD\nEND:VTIMEZONE\n", icalName))
	for _, v := range v.VEvents {
		b.WriteString(v.String())
	}
	b.WriteString("END:VCALENDAR\n")
	return b.String()
}
