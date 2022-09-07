package zjuapi

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"net/url"
)

const GrsLoginUrl = "http://grs.zju.edu.cn/ssohome"
const changeLocaleUrl = "http://grs.zju.edu.cn/py/page/student/grkcb.htm?pageAction=changeLocale"

type GrsTerm int

const (
	Spring GrsTerm = 11
	Summer GrsTerm = 12
	Autumn GrsTerm = 13
	Winter GrsTerm = 14
)

func (c *ZJUAPIClient) FetchTimetable(year int, term GrsTerm) (io.Reader, error) {
	// year - year+1 学年度
	_, err := c.HttpClient.PostForm(changeLocaleUrl, url.Values{
		"locale": {"zh_CN"},
	})
	r, err := c.HttpClient.Get(fmt.Sprintf("http://grs.zju.edu.cn/py/page/student/grkcb.htm?xj=%d&xn=%d", term, year))
	if err != nil {
		e := fmt.Sprintf("failed to fetch timetable for %d-%d, error: %s", year, term, err)
		log.Error().Msg(e)
		return nil, errors.New(e)
	}
	rb, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		e := fmt.Sprintf("failed to read timetable for %d-%d, error: %s", year, term, err)
		log.Error().Msg(e)
		return nil, errors.New(e)
	}
	return bytes.NewBuffer(rb), nil
}
