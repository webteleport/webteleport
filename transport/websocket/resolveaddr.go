package websocket

import (
	"net/url"
)

func ResolveAddr(_ string, relayURL *url.URL) (string, error) {
	u := *relayURL
	params := u.Query()
	params.Set(UpgradeQuery, "1")
	u.RawQuery = params.Encode()
	return u.String(), nil
}
