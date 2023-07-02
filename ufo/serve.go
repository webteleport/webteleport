package ufo

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/webteleport/auth"
	"github.com/webteleport/webteleport"
)

var DefaultTimeout = 5 * time.Second

func Serve(stationURL string, handler http.Handler) error {
	ctx, _ := context.WithTimeout(context.Background(), DefaultTimeout)

	if handler == nil {
		handler = http.DefaultServeMux
	}

	u, err := url.Parse(stationURL)
	if err != nil {
		return nil
	}
	lm := &auth.LoginMiddleware{
		Password: u.Fragment,
	}

	ln, err := webteleport.Listen(ctx, stationURL)
	if err != nil {
		return err
	}

	log.Println("🛸 listening on", ln.ClickableURL())
	if lm.IsPasswordRequired() {
		handler = lm.Wrap(handler)
		log.Println("🔒 secured by password authentication")
	} else {
		log.Println("🔓 publicly accessible without a password")
	}
	return http.Serve(ln, handler)
}
