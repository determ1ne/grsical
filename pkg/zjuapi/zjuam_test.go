package zjuapi

import (
	"bytes"
	"encoding/json"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"net/url"
	"os"
	"testing"
)

func TestEncryption(t *testing.T) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	modulus := "9d2c7ce4add43b2caf5a6f49b2fef2ef38c008c0a8132133b85da9f577c39a49dd87c076c7087ab2fbc8c810bc8be5730ab81c18b87ef87004ef1fb8fc628c23"
	exponent := "10001"
	password := "aaaa0123"

	pk, err := NewPubKey(modulus, exponent)
	if err != nil {
		log.Error().Msg(err.Error())
		t.FailNow()
	}
	p := pk.Encrypt(password)

	if p != "3f5312e265def9ed1da2a5d3c5dcf5d6e0af31d1e8561af8d5c01326b66bb59fc1c85632fe9d8d70932d71a425b58f2e5e0c21bd9885d7feea383e6b21a974bf" {
		t.Fail()
	}
}

func TestLogin(t *testing.T) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	username, uSet := os.LookupEnv("GRSICAL_USERNAME")
	password, pSet := os.LookupEnv("GRSICAL_PASSWORD")
	if (uSet && pSet) == false {
		log.Info().Msg("username or password not set, skipping")
		t.SkipNow()
	}

	DefaultClient := NewClient()

	// see https://github.com/determ1ne/ejector/blob/fbc10d91b5d450cfa9f94a6ef22916463c9107f1/Ejector/Services/ZjuService.cs#L44
	//jar, _ := cookiejar.New(nil)
	//DefaultClient.HttpClient.Jar = jar

	// stage 1: get csrf key
	lpRes, err := DefaultClient.HttpClient.Get(GrsLoginUrl)
	if err != nil {
		panic(err)
	}
	//cookieJar.WriteString(extractCookies(lpRes.Header))
	pageContent, err := io.ReadAll(lpRes.Body)
	if err != nil {
		panic(err)
	}
	idxStart := bytes.Index(pageContent, []byte("execution\"")) + 18
	idxStop := bytes.Index(pageContent[idxStart:], []byte("\" />")) + idxStart
	csrf := pageContent[idxStart:idxStop]

	// stage 2: get pub key
	pkRes, err := DefaultClient.HttpClient.Get(publicKeyUrl)
	if err != nil {
		panic(err)
	}
	//cookieJar.WriteString(extractCookies(pkRes.Header))
	pkContent, err := io.ReadAll(pkRes.Body)
	if err != nil {
		panic(err)
	}
	var pkRaw pubkeyRaw
	err = json.Unmarshal(pkContent, &pkRaw)
	if err != nil {
		panic(err)
	}
	pk, err := NewPubKey(pkRaw.N, pkRaw.E)
	if err != nil {
		panic(err)
	}
	encP := pk.Encrypt(password)

	// stage 3: fire target
	lRes, err := DefaultClient.HttpClient.PostForm(loginUrl, url.Values{
		"username":  {username},
		"password":  {encP},
		"authcode":  {""},
		"execution": {string(csrf)},
		"_eventId":  {"submit"},
	})
	// TODO
	print(lRes)
	b, err := io.ReadAll(lRes.Body)
	print(string(b))
	//_, err = DefaultClient.HttpClient.Get(grsLoginUrl)
	kb, err := DefaultClient.HttpClient.Get("http://grs.zju.edu.cn/py/page/student/grkcb.htm?xj=13&xn=2022")
	a, err := io.ReadAll(kb.Body)
	print(string(a))
}
