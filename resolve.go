package webteleport

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/webteleport/utils"
)

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

// Resolve gets all webteleport endpoints from Alt-Svc dns records / headers
func Resolve(u *url.URL) (endpoints []string) {
	resp, err := http.Head(u.String())
	if err != nil {
		slog.Warn(fmt.Sprintf("http head error: %v", err))
	}
	altsvcHeader := resp.Header.Get("Alt-Svc")

	txts, err := utils.LookupHostTXT(u.Host, "1.1.1.1:53")
	if err != nil {
		slog.Warn(fmt.Sprintf("dns lookup error: %v", err))
	}
	altsvcs := append(altsvcLines(txts), altsvcHeader)
	// log.Println(altsvcs)
	for _, altsvc := range altsvcs {
		endpoints = append(endpoints, utils.ExtractAltSvcEndpoints(altsvc, "webteleport")...)
	}
	return
}
