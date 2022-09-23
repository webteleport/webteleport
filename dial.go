package quichost

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/ebi-yade/altsvc-go"
	"github.com/marten-seemann/webtransport-go"
)

func extractH3(h http.Header) ([]string, bool) {
	line := h.Get("Alt-Svc")
	if line == "" {
		return []string{}, false
	}
	svcs, err := altsvc.Parse(line)
	if err != nil {
		log.Println(err)
		return []string{}, false
	}
	results := []string{}
	log.Println(svcs)
	for _, svc := range svcs {
		if svc.ProtocolID == "h3" {
			results = append(results, svc.AltAuthority.Host+":"+svc.AltAuthority.Port)
		}
	}
	return results, len(results) > 0
}

func dial(ctx context.Context, u *url.URL) (*webtransport.Session, error) {
	resp, err := http.Head(u.String())
	if err != nil {
		return nil, err
	}
	endpoints, ok := extractH3(resp.Header)
	if !ok {
		return nil, errors.New("HTTP/3 service discovery failed: no Alt-Svc header found")
	}
	ep := endpoints[0]
	log.Println("Found", endpoints)
	// use u.Host for now
	// TODO fix this
	log.Println("switch to webtransport", ep, u.Host)

	ur := fmt.Sprintf("https://%s", "quichost.k0s.io:300")
	var d webtransport.Dialer
	log.Printf("dialing %s (UDP)", ur)
	resp, session, err := d.Dial(ctx, ur, nil)
	if err != nil {
		log.Fatalln(err)
	}
	_ = resp
	// handleConn(conn)
	// return nil, nil
	return session, err
}
