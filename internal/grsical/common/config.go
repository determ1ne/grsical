package common

import (
	"encoding/json"
	"grs-ical/pkg/timetable"
	"io"
	"net/http"
	"os"
	"strings"
)

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

func GetConfig(path string) (Config, error) {
	var r io.Reader
	c := Config{}
	if strings.HasPrefix(path, "http") {
		res, err := http.DefaultClient.Get(path)
		if err != nil {
			return c, err
		}
		r = res.Body
		defer res.Body.Close()
	} else {
		f, err := os.Open(path)
		defer f.Close()
		r = f
		if err != nil {
			return c, err
		}
	}
	cfd, err := io.ReadAll(r)
	if err != nil {
		return c, err
	}
	err = json.Unmarshal(cfd, &c)
	return c, err
}

func GetTweakConfig(path string) (TweakConfig, error) {
	tc := TweakConfig{}
	if path == "" {
		return tc, nil
	}

	var r io.Reader
	if strings.HasPrefix(path, "http") {
		res, err := http.DefaultClient.Get(path)
		if err != nil {
			return tc, err
		}
		r = res.Body
		defer res.Body.Close()
	} else {
		f, err := os.Open(path)
		defer f.Close()
		r = f
		if err != nil {
			return tc, err
		}
	}

	tfd, err := io.ReadAll(r)
	if err != nil {
		return tc, err
	}
	err = json.Unmarshal(tfd, &tc)
	return tc, err
}
