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
	hp, err := setport(u.Host, ep)
	if err != nil {
		return nil, err
	}
	var d webtransport.Dialer
	up := fmt.Sprintf("https://%s", hp)
	log.Printf("dialing %s (UDP)", up)
	_, session, err := d.Dial(ctx, up, nil)
	if err != nil {
		return nil, err
	}
	return session, err
}
