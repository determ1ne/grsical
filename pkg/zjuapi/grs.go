package zjuapi

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"net/url"
)

const GrsLoginUrl = "http://grs.zju.edu.cn/ssohome"
const changeLocaleUrl = "http://grs.zju.edu.cn/py/page/student/grkcb.htm?pageAction=changeLocale"

type GrsSemester int
type GrsExamSemester int

const (
	Spring GrsSemester = 11
	Summer GrsSemester = 12
	Autumn GrsSemester = 13
	Winter GrsSemester = 14
)

const (
	AllSemester  GrsExamSemester = -1
	AutumnWInter GrsExamSemester = 16
	SpringSummer GrsExamSemester = 15
)

func (c *ZJUAPIClient) FetchTimetable(ctx context.Context, year int, semester GrsSemester) (io.Reader, error) {
	// TODO: 考虑更换成浙大钉 API
	// year - year+1 学年度
	_, err := c.HttpClient.PostForm(changeLocaleUrl, url.Values{
		"locale": {"zh_CN"},
	})
	r, err := c.HttpClient.Get(fmt.Sprintf("http://grs.zju.edu.cn/py/page/student/grkcb.htm?xj=%d&xn=%d", semester, year))
	if err != nil {
		e := fmt.Sprintf("failed to fetch timetable for %d-%d, error: %s", year, semester, err)
		log.Ctx(ctx).Error().Msg(e)
		return nil, errors.New(e)
	}
	rb, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		e := fmt.Sprintf("failed to read timetable for %d-%d, error: %s", year, semester, err)
		log.Ctx(ctx).Error().Msg(e)
		return nil, errors.New(e)
	}
	return bytes.NewBuffer(rb), nil
}

func (c *ZJUAPIClient) FetchExamTable(ctx context.Context, year int, semester GrsExamSemester) (io.Reader, error) {
	// 这里没有做本地化，无需调整语言
	// 考试日期页面显示较慢，需要做超时准备
	// 浙大钉 API 这里显示不全，暂时不用
	url := fmt.Sprintf("http://grs.zju.edu.cn/py/page/student/grksap.htm?xj=%d&xn=%d", semester, year)
	r, err := c.HttpClient.Get(url)
	if err != nil {
		e := fmt.Sprintf("failed to fetch exam for %d-%d, error: %s", year, semester, err.Error())
		log.Ctx(ctx).Error().Msg(e)
		return nil, errors.New(e)
	}
	rb, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		e := fmt.Sprintf("failed to read exam for %d-%d, error: %s", year, semester, err.Error())
		log.Ctx(ctx).Error().Msg(e)
		return nil, errors.New(e)
	}
	return bytes.NewBuffer(rb), nil
}
