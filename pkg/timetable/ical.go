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

		switch t.Type {
		case Remove:
			for _, mc := range matchedClass {
				for i := len(cc) - 1; i >= 0; i-- {
					c := &cc[i]
					if c == mc {
						cc = append(cc[:i], cc[i+1:]...)
						break
					}
				}
			}
		case Modify:
			for _, mc := range matchedClass {
				for k, v := range t.Op {
					switch k {
					case "name":
						mc.Name = v.(string)
					case "repeat":
						mc.Repeat = v.(Repeat)
					case "starts":
						mc.Duration.Starts = int(v.(float64))
					case "ends":
						mc.Duration.Ends = int(v.(float64))
					case "dayOfWeek":
						mc.DayOfWeek = int(v.(float64))
					default:
						log.Ctx(ctx).Warn().Msgf("unsupported op '%s' in MatchAll-Modify mode", k)
					}
				}
				mc.tweakDesc = t.Description
			}
		case Duplicate:
			for _, mc := range matchedClass {
				for k, v := range t.Op {
					c := *mc
					switch k {
					case "repeat":
						c.Repeat = v.(Repeat)
					case "dayOfWeek":
						c.DayOfWeek = int(v.(float64))
					case "starts":
						c.Duration.Starts = int(v.(float64))
					case "ends":
						c.Duration.Ends = int(v.(float64))
					}
					cc = append(cc, c)
				}
				mc.tweakDesc = t.Description
			}
		}
	}

	cd := firstMonday
	var vEvents []ical.VEvent
	for i := 0; i < 8; i++ {
	CCFOR:
		for _, c := range cc {
			d := cd.AddDate(0, 0, c.DayOfWeek-1)

			cc := c
			ccLst := []*Class{&cc}
			for _, t := range moTweaks {
				date, ok := t.MatchRule["date"]
				if !ok {
					log.Ctx(ctx).Warn().Msgf("tweak has invalid MatchOnce rule")
					continue
				}
				if _, ok := date.(string); !ok {
					log.Ctx(ctx).Warn().Msgf("tweak has invalid MatchOnce rule")
					continue
				}
				currentDate := fmt.Sprintf("%02d%02d", d.Month(), d.Day())
				if currentDate != date.(string) {
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
					case "date":
						// dummy
						continue
					default:
						log.Ctx(ctx).Warn().Msgf("unsupported match rule '%s'", k)
					}
					if cv != v {
						matched = false
						break
					}
				}
				if matched {
					switch t.Type {
					case Remove:
						continue CCFOR
					case Modify:
						for k, v := range t.Op {
							switch k {
							case "name":
								cc.Name = v.(string)
							case "starts":
								cc.Duration.Starts = int(v.(float64))
							case "ends":
								cc.Duration.Ends = int(v.(float64))
							case "date":
								d := fmt.Sprintf("%04d%s", d.Year(), v)
								date, err := time.ParseInLocation("20060102", d, CSTLocation)
								if err != nil {
									log.Warn().Msgf("invalid date in tweak")
									continue CCFOR
								}
								cc.date = date
							default:
								log.Ctx(ctx).Warn().Msgf("unsupported op '%s'", k)
							}
						}
						cc.tweakDesc = t.Description
					case Duplicate:
						ccDup := cc
						for k, v := range t.Op {
							switch k {
							case "starts":
								ccDup.Duration.Starts = int(v.(float64))
							case "ends":
								ccDup.Duration.Ends = int(v.(float64))
							case "dayOfWeek":
								ccDup.DayOfWeek = int(v.(float64))
							case "date":
								d := fmt.Sprintf("%04d%s", d.Year(), date.(string))
								date, err := time.ParseInLocation("20060102", d, CSTLocation)
								if err != nil {
									log.Warn().Msgf("invalid date in tweak")
									continue CCFOR
								}
								ccDup.date = date
							default:
								log.Ctx(ctx).Warn().Msgf("unsupported op '%s'", k)
							}
						}
						ccDup.tweakDesc = t.Description
						ccLst = append(ccLst, &ccDup)
					}
				}
			}

			for _, cc := range ccLst {
				if cc.Repeat == SingleWeek && i%2 == 1 {
					continue
				} else if cc.Repeat == DoubleWeek && i%2 == 0 {
					continue
				}
				if cc.date.IsZero() {
					cc.date = d
				}
				v := ical.VEvent{
					Summary:     cc.Name,
					Location:    cc.Location,
					Description: cc.ToDesc(i),
					StartTime:   cc.date,
					EndTime:     cc.date,
				}
				v.StartTime = cc.date.Add(time.Duration(ClassStart[cc.Duration.Starts]) * time.Minute)
				v.EndTime = cc.date.Add(time.Duration(ClassStart[cc.Duration.Ends]) * time.Minute).Add(time.Minute * 45)
				vEvents = append(vEvents, v)
			}
		}
		cd = cd.AddDate(0, 0, 7)
	}
	return &vEvents, nil
}

func ExamToVEvents(ctx context.Context, exams *[]Exam) (*[]ical.VEvent, error) {
	var vEvents []ical.VEvent
	for _, exam := range *exams {
		v := ical.VEvent{
			Summary:     fmt.Sprintf("考试：%s", exam.Name),
			Location:    fmt.Sprintf("%s %s", exam.Region, exam.Location),
			Description: exam.ToDesc(),
			StartTime:   exam.StartTime,
			EndTime:     exam.EndTime,
		}
		vEvents = append(vEvents, v)
	}
	return &vEvents, nil
}
