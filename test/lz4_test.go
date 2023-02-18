package test

import (
	"bytes"
	"github.com/pierrec/lz4"
	"io"
	"strings"
	"testing"
)

func TestLz4(t *testing.T) {
	encode := func(s string) ([]byte, error) {
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
	decode := func(b []byte) (string, error) {
		var out bytes.Buffer
		zin := bytes.NewBuffer(b)
		zr := lz4.NewReader(zin)
		_, err := io.Copy(&out, zr)
		return string(out.Bytes()), err
	}

	a := "12313oi12j3io1j3io1j32oi"
	println(a)
	x, _ := encode(a)
	println(x)
	b, _ := decode(x)
	println(b)
}
