package ufo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"

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

func setport(base, ep string) (string, error) {
	host, port, err := net.SplitHostPort(ep)
	if err != nil {
		return "", err
	}
	if host != "" {
		return "", err
	}
	basehost, _, _ := strings.Cut(base, ":")
	return basehost + ":" + port, nil
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
	log.Println("found H3 endpoints", endpoints)
	hp, err := setport(u.Host, ep)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("getting webtransport address", u.Host, "+", ep, "=", hp)
	ur := fmt.Sprintf("https://%s", hp)
	var d webtransport.Dialer
	log.Printf("dialing %s (UDP)", ur)
	_, session, err := d.Dial(ctx, ur, nil)
	if err != nil {
		log.Fatalln(err)
	}
	return session, err
}
