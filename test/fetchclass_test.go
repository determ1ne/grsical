package test

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"grs-ical/pkg/ical"
	"grs-ical/pkg/timetable"
	"grs-ical/pkg/zjuapi"
	"os"
	"testing"
	"time"
)

func TestFetchClass(t *testing.T) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	u, ue := os.LookupEnv("GRSICAL_USERNAME")
	p, pe := os.LookupEnv("GRSICAL_PASSWORD")
	if (ue && pe) == false {
		log.Info().Msg("username or password not set, skip")
		t.SkipNow()
	}

	c := zjuapi.NewClient()
	err := c.Login(zjuapi.GrsLoginUrl, u, p)
	if err != nil {
		log.Error().Msg(err.Error())
		t.FailNow()
	}

	r, err := c.FetchTimetable(2022, zjuapi.Autumn)
	if err != nil {
		log.Error().Msg(err.Error())
		t.FailNow()
	}

	table, err := timetable.GetTable(r)
	if err != nil {
		log.Error().Msg(err.Error())
		t.FailNow()
	}
	cl, err := timetable.ParseTable(table)
	if err != nil {
		log.Error().Msg(err.Error())
		t.FailNow()
	}

	fm, err := time.ParseInLocation("20060102", "20220912", time.Local)
	if err != nil {
		log.Error().Msg(err.Error())
		t.FailNow()
	}
	vEvents, err := timetable.ClassToVEvents(fm, cl, nil)
	if err != nil {
		log.Error().Msg(err.Error())
		t.FailNow()
	}
	iCal := ical.VCalendar{
		vEvents,
	}

	f, err := os.OpenFile("ical.ics", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Error().Msg(err.Error())
		t.FailNow()
	}

	f.WriteString(iCal.GetICS(""))
	f.Close()
}
