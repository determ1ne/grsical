package common

import (
	"context"
	"fmt"
	"strings"
	"time"

	"grs-ical/pkg/ical"
	"grs-ical/pkg/timetable"
	"grs-ical/pkg/zjuapi"

	"github.com/rs/zerolog/log"
	"golang.org/x/net/html"
)

func FetchToMemory(ctx context.Context, username, password string, config Config, tweaks TweakConfig) (string, error) {
	c := zjuapi.NewClient()
	log.Ctx(ctx).Info().Msgf("logging in for %s", username)

	var isUGRS = false
	// maybe not enough
	if strings.HasPrefix(username, "3") {
		isUGRS = true
		log.Ctx(ctx).Info().Msgf("%s is using UGRS(beta)", username)
	}

	if isUGRS {
		err := c.Login(ctx, zjuapi.UgrsLoginUrl, username, password)
		if err != nil {
			return "", err
		}
		// need extra login
		err = c.UgrsExtraLogin(ctx, fmt.Sprintf("%s?gnmkdm=N253530&su=%s", zjuapi.UgrsLoginUrl2, username))
		if err != nil {
			return "", err
		}
	} else {
		err := c.Login(ctx, zjuapi.GrsLoginUrl, username, password)
		if err != nil {
			return "", err
		}
	}

	var ve []ical.VEvent
	year := map[int]struct{}{}
	for _, fc := range config.FetchConfig {
		log.Ctx(ctx).Info().Msgf("fetching %d-%d", fc.Year, fc.Semester)
		year[fc.Year] = struct{}{}
		r, err := c.FetchTimetable(ctx, fc.Year, zjuapi.GrsSemester(fc.Semester), isUGRS)
		if err != nil {
			return "", err
		}

		table, err := timetable.GetTable(r)
		if err != nil {
			return "", err
		}

		log.Ctx(ctx).Info().Msgf("parsing %d-%d", fc.Year, fc.Semester)
		cl, err := timetable.ParseTable(ctx, table, isUGRS)
		if err != nil {
			// dump table
			var b strings.Builder
			_ = html.Render(&b, table)
			return b.String(), err
		}

		fm, err := time.ParseInLocation("20060102", fc.FirstDay, CSTLocation)
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

	// fetch exams
	for y, _ := range year {
		log.Ctx(ctx).Info().Msgf("fetching exam info %d", y)
		r, err := c.FetchExamTable(ctx, y, zjuapi.AllSemester, isUGRS)
		if err != nil {
			log.Ctx(ctx).Error().Msgf("failed to fetch exam table, error: %s", err.Error())
			continue
		}
		table, err := timetable.GetExamTable(r)
		if err != nil {
			log.Ctx(ctx).Error().Msgf("failed to get exam table, error: %s", err.Error())
			continue
		}
		log.Ctx(ctx).Info().Msgf("parsing exam info %d", y)
		exams, err := timetable.ParseExamTable(ctx, table)
		if err != nil {
			log.Ctx(ctx).Error().Msgf("failed to parse exam table, error: %s", err.Error())
			continue
		}
		vEvents, err := timetable.ExamToVEvents(ctx, exams)
		ve = append(ve, *vEvents...)
		if err != nil {
			log.Ctx(ctx).Error().Msgf("failed to generate exam ical, error: %s", err.Error())
			continue
		}
	}

	log.Ctx(ctx).Info().Msgf("generating iCalendar file")
	iCal := ical.VCalendar{VEvents: ve}
	return iCal.GetICS(""), nil
}
