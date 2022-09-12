package common

import (
	"bytes"
	"github.com/pierrec/lz4"
	"io"
	"strings"
	"time"
)

const ContextReqId = "grsical-reqid"
const DurationOneDay = time.Hour * 24

func Lz4Encode(s string) ([]byte, error) {
	var zout bytes.Buffer
	zw := lz4.NewWriter(&zout)
	_, err := io.Copy(zw, strings.NewReader(s))
	if err != nil {
		return []byte(""), err
	}
	err = zw.Close()
	if err != nil {
		return []byte(""), err
	}
	return zout.Bytes(), nil
}

func Lz4Decode(b []byte) (string, error) {
	var out bytes.Buffer
	zin := bytes.NewBuffer(b)
	zr := lz4.NewReader(zin)
	_, err := io.Copy(&out, zr)
	return string(out.Bytes()), err
}
