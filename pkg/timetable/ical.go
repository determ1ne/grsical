package timetable

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"grs-ical/pkg/ical"
	"time"
)

func ClassToVEvents(ctx context.Context, firstMonday time.Time, class *[]Class, tweaks *[]Tweak) (*[]ical.VEvent, error) {
	cc := *class
	var maTweaks, moTweaks []Tweak
	for _, t := range *tweaks {
		if t.MatchType == MatchAll {
			maTweaks = append(maTweaks, t)
		} else {
			moTweaks = append(moTweaks, t)
		}
	}
	for _, t := range maTweaks {
		var matchedClass []*Class
		for i, c := range cc {
			// match
			matched := true
			for k, v := range t.MatchRule {
				var cv interface{}
				switch k {
				case "name":
					cv = c.Name
				case "semester":
					cv = c.Semester
				case "repeat":
					cv = c.Repeat
				case "starts":
					cv = c.Duration.Starts
				case "ends":
					cv = c.Duration.Ends
				case "dayOfWeek":
					cv = c.DayOfWeek
				default:
					log.Ctx(ctx).Warn().Msgf("unsupported match rule '%s'", k)
				}
				if cv != v {
					matched = false
					break
				}
			}
			if matched {
				matchedClass = append(matchedClass, &cc[i])
			}
		}

		for _, mc := range matchedClass {
			for k, v := range t.Op {
				switch k {
				case "name":
					mc.Name = v.(string)
				case "semester":
					mc.Semester = v.(Semester)
				case "repeat":
					mc.Repeat = v.(Repeat)
				case "starts":
					mc.Duration.Starts = int(v.(float64))
				case "ends":
					mc.Duration.Ends = int(v.(float64))
				case "dayOfWeek":
					mc.DayOfWeek = int(v.(float64))
				default:
					log.Ctx(ctx).Warn().Msgf("unsupported op '%s'", k)
				}
			}
		}
	}

	cd := firstMonday
	var vEvents []ical.VEvent
	for i := 0; i < 8; i++ {
		for _, c := range cc {
			d := cd.AddDate(0, 0, c.DayOfWeek-1)

			cc := c
			for _, t := range moTweaks {
				if t.MatchType != MatchOnce {
					continue
				}
				date, ok := t.MatchRule["date"]
				if !ok {
					log.Ctx(ctx).Warn().Msgf("tweak has invalid MatchOnce rule")
					continue
				}
				if _, ok := date.(string); !ok {
					log.Ctx(ctx).Warn().Msgf("tweak has invalid MatchOnce rule")
					continue
				}
				if fmt.Sprintf("%02d%02d", d.Month(), d.Day()) != date.(string) {
					continue
				}

				// matched date
				matched := true
				for k, v := range t.MatchRule {
					var cv interface{}
					switch k {
					case "name":
						cv = c.Name
					case "semester":
						cv = c.Semester
					case "repeat":
						cv = c.Repeat
					case "starts":
						cv = c.Duration.Starts
					case "ends":
						cv = c.Duration.Ends
					case "dayOfWeek":
						cv = c.DayOfWeek
					default:
						log.Ctx(ctx).Warn().Msgf("unsupported match rule '%s'", k)
					}
					if cv != v {
						matched = false
						break
					}
				}
				if matched {
					for k, v := range t.Op {
						switch k {
						case "name":
							cc.Name = v.(string)
						case "semester":
							cc.Semester = v.(Semester)
						case "repeat":
							cc.Repeat = v.(Repeat)
						case "starts":
							cc.Duration.Starts = int(v.(float64))
						case "ends":
							cc.Duration.Ends = int(v.(float64))
						case "dayOfWeek":
							cc.DayOfWeek = int(v.(float64))
						default:
							log.Ctx(ctx).Warn().Msgf("unsupported op '%s'", k)
						}
					}
				}
			}

			v := ical.VEvent{
				Summary:     cc.Name,
				Location:    cc.Location,
				Description: cc.ToDesc(),
				StartTime:   d,
				EndTime:     d,
			}
			v.StartTime = v.StartTime.Add(time.Duration(ClassStart[cc.Duration.Starts]) * time.Minute)
			v.EndTime = v.StartTime.Add(45 * time.Minute)
			vEvents = append(vEvents, v)
		}
		cd = cd.AddDate(0, 0, 7)
	}
	return &vEvents, nil
}
