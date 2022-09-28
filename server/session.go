package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/marten-seemann/webtransport-go"
)

type Session struct {
	*webtransport.Session
	Candidates []string
	SecretCode string
}

// PrecheckAccessToken returns a bool that indicates whether the caller should continue
func (ssn *Session) PrecheckAccessToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hostonly, _, _ := strings.Cut(r.URL.Host, ":")
		if ssn.SecretCode == "" || strings.HasSuffix(hostonly, "localhost") {
			next.ServeHTTP(w, r)
			return
		}
		if r.URL.Path == fmt.Sprintf("/secretCode=%s", ssn.SecretCode) {
			cookies := fmt.Sprintf(`WebTeleportSecretCode="%s"; Path=/; Max-Age=2592000; HttpOnly; Domain=%s`, ssn.SecretCode, r.Host)
			w.Header().Set("Set-Cookie", cookies)
			http.Redirect(w, r, "/", 302)
			return
		}
		wtat, err := r.Cookie("WebTeleportSecretCode")
		if err != nil {
			http.Error(w, "🛸"+http.StatusText(401)+" "+err.Error(), 401)
			return
		}
		if wtat.Value != ssn.SecretCode {
			http.Error(w, "🛸"+http.StatusText(401), 401)
			return
		}
		next.ServeHTTP(w, r)
	})
}
