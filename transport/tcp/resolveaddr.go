package tcp

import (
	"net/url"
	"strings"

	"github.com/webteleport/utils"
)

func ResolveAddr(addr string, relayURL *url.URL) (string, error) {
	// Bare port like ":8080" — graft relay host, keep endpoint port
	if strings.HasPrefix(addr, ":") {
		addr = utils.Graft(relayURL.Host, addr)
	}
	u, err := url.Parse(utils.AsURL(addr))
	if err != nil {
		return "", err
	}
	return u.Host, nil
}
