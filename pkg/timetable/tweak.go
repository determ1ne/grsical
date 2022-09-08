package timetable

type MatchType int

const (
	MatchAll MatchType = iota
	MatchOnce
)

type Match struct {
	MatchType MatchType         `json:"match_type"`
	MatchRule map[string]string `json:"match_rule"`
}

type TweakType int

const (
	Remove TweakType = iota
	Modify
	Duplicate
)

type TweakOp struct {
}

type Tweak struct {
	Match []Match   `json:"match"`
	Type  TweakType `json:"type"`
	Op    []TweakOp `json:"op"`
}
