package timetable

import (
	"grs-ical/pkg/ical"
	"time"
)

func ClassToVEvents(firstMonday time.Time, class *[]Class, tweaks *[]Tweak) (*[]ical.VEvent, error) {
	// TODO: Apply Tweak

	cd := firstMonday
	var vEvents []ical.VEvent
	for i := 0; i < 8; i++ {
		for _, c := range *class {
			d := cd.AddDate(0, 0, c.DayOfWeek-1)
			// TODO: Apply Tweak
			v := ical.VEvent{
				Summary:     c.Name,
				Location:    c.Location,
				Description: c.ToDesc(),
				StartTime:   d,
				EndTime:     d,
			}
			v.StartTime = v.StartTime.Add(time.Duration(ClassStart[c.Duration.Starts]) * time.Minute)
			v.EndTime = v.StartTime.Add(45 * time.Minute)
			vEvents = append(vEvents, v)
		}
		cd = cd.AddDate(0, 0, 7)
	}
	return &vEvents, nil
}
