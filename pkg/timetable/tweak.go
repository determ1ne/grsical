package timetable

type MatchType int

const (
	MatchAll MatchType = iota
	MatchOnce
)

type TweakType int

const (
	Remove TweakType = iota
	Modify
	Duplicate
)

type Tweak struct {
	MatchType   MatchType              `json:"matchType"`
	Type        TweakType              `json:"type"`
	MatchRule   map[string]interface{} `json:"rule"`
	Op          map[string]interface{} `json:"op"`
	Description string                 `json:"description"`
}
