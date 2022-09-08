package common

import "grs-ical/pkg/timetable"

type FetchConfig struct {
	Year     int    `json:"year"`
	Semester int    `json:"semester"`
	FirstDay string `json:"firstDay"`
}

type Config struct {
	FetchConfig []FetchConfig `json:"fetch"`
}

type TweakConfig struct {
	Tweaks []timetable.Tweak `json:"tweaks"`
}
