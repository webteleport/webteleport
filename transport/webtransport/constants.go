package webtransport

import "net/http"

var (
	UpgradeQuery  = "x-webtransport-upgrade"
	UpgradeHeader = http.CanonicalHeaderKey(UpgradeQuery)
)
