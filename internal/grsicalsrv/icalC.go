package grsicalsrv

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	common2 "grs-ical/internal/common"
)

func decrypt(b []byte) ([]byte, error) {
	ns := gcm.NonceSize()
	if len(b) < ns {
		return []byte(""), errors.New("invalid data")
	}
	nonce, ct := b[:ns], b[ns:]
	p, err := gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return []byte(""), err
	}
	return p, nil
}

func FetchCal(ctx *fiber.Ctx) error {
	p := ctx.Query("p")
	if p == "" {
		return ctx.SendString("invalid p")
	}
	b, err := base64.URLEncoding.DecodeString(p)
	if err != nil {
		return ctx.SendString("invalid p2")
	}
	unpw, err := decrypt(b)
	if err != nil {
		return ctx.SendString("invalid p3")
	}
	un := unpw[:12]
	pw := unpw[12:]
	for i := 11; i >= 0; i-- {
		if un[i] != '#' {
			un = un[:i+1]
			break
		}
	}

	reqId := uuid.NewString()
	c := log.With().Str("u", string(un)).Str("r", reqId).Logger().WithContext(ctx.UserContext())
	r, err := common2.FetchToMemory(c, string(un), string(pw), cfg, tweak)
	if err != nil {
		if r != "" && rc != nil {
			tableComp, _ := common2.Lz4Encode(r)
			_ = rc.Set(c, fmt.Sprintf("dump%s", reqId), tableComp, common2.DurationOneDay*7)
		}
		return ctx.SendString(fmt.Sprintf("%s\r\nrequest id=%s", err.Error(), reqId))
	}

	ctx.Set("Content-Type", "text/calendar")
	return ctx.SendString(r)
}
