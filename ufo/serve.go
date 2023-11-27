package ufo

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/webteleport/auth"
	"github.com/webteleport/webteleport"
)

var DefaultTimeout = 10 * time.Second

func Serve(stationURL string, handler http.Handler) error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	if handler == nil {
		handler = http.DefaultServeMux
	}

	u, err := url.Parse(stationURL)
	if err != nil {
		return err
	}
	lm := &auth.LoginMiddleware{
		Password: u.Fragment,
	}

	// attach extra info to the query string
	q := u.Query()
	q.Add("client", "ufo")
	for _, arg := range os.Args {
		q.Add("args", arg)
	}
	u.RawQuery = q.Encode()

	ln, err := webteleport.Listen(ctx, u.String())
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
