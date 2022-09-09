package timetable

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"reflect"
	"testing"
)

var classList = &[]Class{
	{Name: "电力电子技术在电力系统中的应用", Semester: Autumn, Repeat: EveryWeek, Duration: ClassDuration{Starts: 1, Ends: 2}, Teacher: "徐政", Location: "(场地详见学院通知)", DayOfWeek: 1, RawDuration: "第一节--第二节"},
	{Name: "电力电子技术在电力系统中的应用", Semester: Autumn, Repeat: EveryWeek, Duration: ClassDuration{Starts: 1, Ends: 2}, Teacher: "徐政", Location: "(场地详见学院通知)", DayOfWeek: 2, RawDuration: "第一节--第二节"},
	{Name: "现代控制理论", Semester: AutumnWinter, Repeat: EveryWeek, Duration: ClassDuration{Starts: 1, Ends: 3}, Teacher: "王建全", Location: "玉泉教7-408(录播.4)", DayOfWeek: 3, RawDuration: "第一节--第三节"},
	{Name: "电气工程学科最新发展综述", Semester: Autumn, Repeat: EveryWeek, Duration: ClassDuration{Starts: 1, Ends: 4}, Teacher: "沈建新", Location: "玉泉教7-404(录播.4)", DayOfWeek: 5, RawDuration: "第一节--第四节"},
	{Name: "电力系统规划", Semester: Autumn, Repeat: EveryWeek, Duration: ClassDuration{Starts: 6, Ends: 9}, Teacher: "林振智(杨莉)", Location: "玉泉教7-406(录播.4)", DayOfWeek: 1, RawDuration: "第六节--第九节"},
	{Name: "DSP在机电控制中的应用", Semester: Autumn, Repeat: EveryWeek, Duration: ClassDuration{Starts: 6, Ends: 9}, Teacher: "章玮", Location: "玉泉第2教学大楼-115", DayOfWeek: 2, RawDuration: "第六节--第九节"},
	{Name: "电力系统运行与控制", Semester: Autumn, Repeat: EveryWeek, Duration: ClassDuration{Starts: 6, Ends: 9}, Teacher: "郭瑞鹏", Location: "玉泉教7-402(录播.4)", DayOfWeek: 4, RawDuration: "第六节--第九节"},
	{Name: "电力市场与电力经济", Semester: Autumn, Repeat: EveryWeek, Duration: ClassDuration{Starts: 6, Ends: 9}, Teacher: "杨莉(甘德强)", Location: "玉泉教3-340(录播)", DayOfWeek: 5, RawDuration: "第六节--第九节"},
	{Name: "运动素质课", Semester: Autumn, Repeat: EveryWeek, Duration: ClassDuration{Starts: 9, Ends: 10}, Teacher: "潘雯雯", Location: "玉泉教7-402(录播.4)", DayOfWeek: 3, RawDuration: "第九节--第十节"},
	{Name: "中国马克思主义与当代", Semester: Autumn, Repeat: EveryWeek, Duration: ClassDuration{Starts: 11, Ends: 14}, Teacher: "董扣艳", Location: "玉泉教7-406(录播.4)", DayOfWeek: 2, RawDuration: "第十一节--第十四节"},
	{Name: "中国马克思主义与当代", Semester: Autumn, Repeat: EveryWeek, Duration: ClassDuration{Starts: 11, Ends: 14}, Teacher: "桑建泉", Location: "玉泉教7-406(录播.4)", DayOfWeek: 3, RawDuration: "第十一节--第十四节"}}

var conflictClassList = &[]Class{
	{Name: "电力电子技术在电力系统中的应用", Semester: Autumn, Repeat: EveryWeek, Duration: ClassDuration{Starts: 1, Ends: 2}, Teacher: "徐政", Location: "(场地详见学院通知)", DayOfWeek: 1, RawDuration: "第一节--第二节"},
	{Name: "电力电子技术在电力系统中的应用", Semester: Autumn, Repeat: EveryWeek, Duration: ClassDuration{Starts: 1, Ends: 2}, Teacher: "徐政", Location: "(场地详见学院通知)", DayOfWeek: 2, RawDuration: "第一节--第二节"},
	{Name: "现代控制理论", Semester: AutumnWinter, Repeat: EveryWeek, Duration: ClassDuration{Starts: 1, Ends: 3}, Teacher: "王建全", Location: "玉泉教7-408(录播.4)", DayOfWeek: 3, RawDuration: "第一节--第三节"},
	{Name: "电气工程学科最新发展综述", Semester: Autumn, Repeat: EveryWeek, Duration: ClassDuration{Starts: 1, Ends: 4}, Teacher: "沈建新", Location: "玉泉教7-404(录播.4)", DayOfWeek: 5, RawDuration: "第一节--第四节"},
	{Name: "电力系统规划", Semester: Autumn, Repeat: EveryWeek, Duration: ClassDuration{Starts: 6, Ends: 7}, Teacher: "林振智(杨莉)", Location: "玉泉教7-406(录播.4)", DayOfWeek: 1, RawDuration: "第六节--第九节"},
	{Name: "高尔夫球", Semester: Autumn, Repeat: EveryWeek, Duration: ClassDuration{Starts: 6, Ends: 7}, Teacher: "楼恒阳", Location: "(场地详见学院通知)", DayOfWeek: 1, RawDuration: "第六节--第七节"},
	{Name: "DSP在机电控制中的应用", Semester: Autumn, Repeat: EveryWeek, Duration: ClassDuration{Starts: 6, Ends: 9}, Teacher: "章玮", Location: "玉泉第2教学大楼-115", DayOfWeek: 2, RawDuration: "第六节--第九节"},
	{Name: "电力系统运行与控制", Semester: Autumn, Repeat: EveryWeek, Duration: ClassDuration{Starts: 6, Ends: 9}, Teacher: "郭瑞鹏", Location: "玉泉教7-402(录播.4)", DayOfWeek: 4, RawDuration: "第六节--第九节"},
	{Name: "电力市场与电力经济", Semester: Autumn, Repeat: EveryWeek, Duration: ClassDuration{Starts: 6, Ends: 9}, Teacher: "杨莉(甘德强)", Location: "玉泉教3-340(录播)", DayOfWeek: 5, RawDuration: "第六节--第九节"},
	{Name: "电力系统规划", Semester: Autumn, Repeat: EveryWeek, Duration: ClassDuration{Starts: 8, Ends: 9}, Teacher: "林振智(杨莉)", Location: "玉泉教7-406(录播.4)", DayOfWeek: 1, RawDuration: "第六节--第九节"},
	{Name: "运动素质课", Semester: Autumn, Repeat: EveryWeek, Duration: ClassDuration{Starts: 9, Ends: 10}, Teacher: "潘雯雯", Location: "玉泉教7-402(录播.4)", DayOfWeek: 3, RawDuration: "第九节--第十节"},
	{Name: "中国马克思主义与当代", Semester: Autumn, Repeat: EveryWeek, Duration: ClassDuration{Starts: 11, Ends: 14}, Teacher: "董扣艳", Location: "玉泉教7-406(录播.4)", DayOfWeek: 2, RawDuration: "第十一节--第十四节"},
	{Name: "中国马克思主义与当代", Semester: Autumn, Repeat: EveryWeek, Duration: ClassDuration{Starts: 11, Ends: 14}, Teacher: "桑建泉", Location: "玉泉教7-406(录播.4)", DayOfWeek: 3, RawDuration: "第十一节--第十四节"}}

func testParser(t *testing.T, fileName string, classList *[]Class) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	ctx := context.Background()

	f, err := os.Open(fileName)
	defer f.Close()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer f.Close()

	table, err := GetTable(f)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	cl, err := ParseTable(ctx, table)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if !reflect.DeepEqual(*cl, *classList) {
		fmt.Printf("%#v\n", cl)
		t.FailNow()
	}
}

func TestNormalParser(t *testing.T) {
	testParser(t, "./test_assets/timetable.html", classList)
}

func TestConflictParser(t *testing.T) {
	testParser(t, "./test_assets/timetable-conflict.html", conflictClassList)
}
