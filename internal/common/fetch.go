package common

import (
	"context"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/html"
	"grs-ical/pkg/ical"
	"grs-ical/pkg/timetable"
	"grs-ical/pkg/zjuapi"
	"strings"
	"time"
)

func FetchToMemory(ctx context.Context, username, password string, config Config, tweaks TweakConfig) (string, error) {
	c := zjuapi.NewClient()
	log.Ctx(ctx).Info().Msgf("logging in for %s", username)
	err := c.Login(ctx, zjuapi.GrsLoginUrl, username, password)
	if err != nil {
		return "", err
	}

	var ve []ical.VEvent
	for _, fc := range config.FetchConfig {
		log.Ctx(ctx).Info().Msgf("fetching %d-%d", fc.Year, fc.Semester)
		r, err := c.FetchTimetable(ctx, fc.Year, zjuapi.GrsSemester(fc.Semester))
		if err != nil {
			return "", err
		}

		table, err := timetable.GetTable(r)
		if err != nil {
			return "", err
		}

		log.Ctx(ctx).Info().Msgf("parsing %d-%d", fc.Year, fc.Semester)
		cl, err := timetable.ParseTable(ctx, table)
		if err != nil {
			// dump table
			var b strings.Builder
			_ = html.Render(&b, table)
			return b.String(), err
		}

		fm, err := time.ParseInLocation("20060102", fc.FirstDay, time.Local)
		if err != nil {
			return "", err
		}

		log.Ctx(ctx).Info().Msgf("generating vevents %d-%d", fc.Year, fc.Semester)
		vEvents, err := timetable.ClassToVEvents(ctx, fm, cl, &tweaks.Tweaks)
		if err != nil {
			return "", err
		}

		ve = append(ve, *vEvents...)
	}

	log.Ctx(ctx).Info().Msgf("generating iCalendar file")
	iCal := ical.VCalendar{VEvents: &ve}
	return iCal.GetICS(""), nil
}
