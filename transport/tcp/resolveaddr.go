package tcp

import (
	"net/url"

	"github.com/webteleport/utils"
)

func ResolveAddr(addr string, relayURL *url.URL) (string, error) {
	u, err := url.Parse(utils.AsURL(addr))
	if err != nil {
		return "", err
	}
	if u.Hostname() == "" {
		u.Host = relayURL.Host
	}
	return u.Host, nil
}
