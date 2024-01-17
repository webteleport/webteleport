package webteleport

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"github.com/quic-go/webtransport-go"
	"github.com/webteleport/utils"
)

// 2^60 == 1152921504606846976
const MaxIncomingStreams int64 = 1 << 60

func altsvcLines(txts []string) []string {
	const prefix = "Alt-Svc: "

	altsvcs := []string{}

	for _, txt := range txts {
		// Case insensitive prefix match. See Issue 22736.
		if len(txt) < len(prefix) || !strings.EqualFold(txt[:len(prefix)], prefix) {
			continue
		}
		altsvcs = append(altsvcs, txt[len(prefix):])
	}

	return altsvcs
}

// Dial is a wrapper around webtransport.Dial with automatic HTTP/3 service discovery
func Dial(ctx context.Context, u *url.URL, hdr http.Header) (*webtransport.Session, error) {
	resp, err := http.Head(u.String())
	if err != nil {
		return nil, err
	}
	txts, err := utils.LookupHostTXT(u.Host, "1.1.1.1:53")
	if err != nil {
		slog.Warn(fmt.Sprintf("dns lookup error: %v", err))
	}
	altsvcHeader := resp.Header.Get("Alt-Svc")
	altsvcs := append(altsvcLines(txts), altsvcHeader)
	// log.Println(altsvcs)
	endpoints := []string{}
	for _, altsvc := range altsvcs {
		endpoints = append(endpoints, utils.ExtractAltSvcEndpoints(altsvc, "webteleport")...)
	}
	if len(endpoints) == 0 {
		return nil, errors.New("service discovery failed: no webteleport endpoints found in Alt-Svc records / headers")
	}
	alt := endpoints[0]
	addr := utils.Graft(u.Host, alt)
	d := &webtransport.Dialer{
		RoundTripper: &http3.RoundTripper{
			QuicConfig: &quic.Config{
				MaxIncomingStreams: MaxIncomingStreams,
			},
		},
	}
	// we are dialing an HTTP/3 address, so it is guaranteed to be https
	un, _ := url.Parse(u.String())
	un.Scheme = "https"
	un.Host = addr
	params := un.Query()
	params.Add("x-webteleport-upgrade", "1")
	un.RawQuery = params.Encode()
	_, session, err := d.Dial(ctx, un.String(), hdr)
	if err != nil {
		return nil, fmt.Errorf("error dialing %s (UDP): %w", un.String(), err)
	}
	return session, nil
}
