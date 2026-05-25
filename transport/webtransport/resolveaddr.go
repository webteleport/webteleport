package webtransport

import (
	"net/url"

	"github.com/webteleport/utils"
)

func ResolveAddr(addr string, relayURL *url.URL) (string, error) {
	u, err := url.Parse(utils.AsURL(addr))
	if err != nil {
		return "", err
	}
	u.Scheme = "https"
	u.Path = relayURL.Path
	u.RawPath = relayURL.RawPath
	params := relayURL.Query()
	params.Set(UpgradeQuery, "1")
	u.RawQuery = params.Encode()
	return u.String(), nil
}
