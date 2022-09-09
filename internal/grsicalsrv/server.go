package grsicalsrv

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"github.com/gofiber/fiber/v2"
	cmn "grs-ical/internal/common"
	"io"
	"os"
	"text/template"
)

type SetupData struct {
	Year     int
	YearN    int
	Semester string
	Link     string
}

var EncKey, Host string
var setupTpl *template.Template
var cfg cmn.Config
var tweak cmn.TweakConfig
var sd = SetupData{
	Year:     0,
	YearN:    0,
	Semester: "",
	Link:     "",
}
var gcm cipher.AEAD

func ListenAndServe(address string) error {
	EncKey = os.Getenv("GRSICALSRV_ENCKEY")
	if EncKey == "" {
		return errors.New("encryption key not set")
	}
	c, err := aes.NewCipher([]byte(EncKey))
	if err != nil {
		return err
	}
	gcm, err = cipher.NewGCM(c)
	if err != nil {
		return err
	}

	Host = os.Getenv("GRSICALSRV_HOST")
	if Host == "" {
		return errors.New("host not set")
	}
	configFile := os.Getenv("GRSICALSRV_CFG")
	if configFile == "" {
		return errors.New("config file not set")
	}
	tweaksFile := os.Getenv("GRSICALSRV_TWEAKS")
	if tweaksFile == "" {
		return errors.New("config file not set")
	}
	// read template
	f, err := os.Open("./web/template/setup.html")
	if err != nil {
		return err
	}
	fc, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	setupTpl, err = template.New("setup").Parse(string(fc))
	if err != nil {
		return err
	}
	// read config
	cfg, err = cmn.GetConfig(configFile)
	if err != nil {
		return err
	}
	sd.Year = cfg.FetchConfig[0].Year
	sd.YearN = sd.Year + 1
	if cfg.FetchConfig[0].Semester > 12 {
		sd.Semester = "秋冬"
	} else {
		sd.Semester = "春夏"
	}
	tweak, err = cmn.GetTweakConfig(tweaksFile)
	if err != nil {
		return err
	}

	app := fiber.New()
	setRoutes(app)
	return app.Listen(address)
}

func setRoutes(app *fiber.App) {
	app.Static("/", "./web/app")
	app.Post("/", SetupPage)
	app.Get("/ical", FetchCal)
}
