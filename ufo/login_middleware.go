package ufo

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
)

//go:embed www/*
var WWW embed.FS

type LoginMiddleware struct {
	// set login password
	// empty value removes password check
	Password string
	// optional, fallback to "UFOSID" if empty
	SessionKey string
	// optional, fallback to "password" if empty
	PasswordKey string

	sessions map[string]struct{}
	mutex    *sync.RWMutex
	login    http.Handler
}

func (lm *LoginMiddleware) AddSessionId(id string) {
	lm.mutex.Lock()
	lm.sessions[id] = struct{}{}
	lm.mutex.Unlock()
}

func (lm *LoginMiddleware) HasSessionId(id string) bool {
	lm.mutex.RLock()
	_, ok := lm.sessions[id]
	lm.mutex.RUnlock()
	return ok
}

func (lm *LoginMiddleware) RedirectToLogin(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/login/":
		break
	case r.Referer() != "":
		break
	default:
		http.Redirect(w, r, "/login/", 302)
		return
	}
	lm.login.ServeHTTP(w, r)
}

func (lm *LoginMiddleware) SetCookiesAndRedirect(w http.ResponseWriter, r *http.Request) {
	sid := uuid.New().String()
	lm.AddSessionId(sid)
	cookies := fmt.Sprintf(`%s="%s"; Path=/; Max-Age=2592000; HttpOnly; Domain=%s`, lm.SessionKey, sid, r.Host)
	w.Header().Set("Set-Cookie", cookies)
	http.Redirect(w, r, "/", 302)
}

func (lm *LoginMiddleware) IsValidLogin(r *http.Request) bool {
	switch r.Method {
	case http.MethodGet:
		return r.URL.Path == fmt.Sprintf("/login/%s", lm.Password)
	case http.MethodPost:
		return r.PostFormValue(lm.PasswordKey) == lm.Password
	}
	return false
}

func (lm *LoginMiddleware) IsValidSession(r *http.Request) bool {
	sid, err := r.Cookie(lm.SessionKey)
	if err != nil || !lm.HasSessionId(sid.Value) {
		return false
	}
	return true
}

func (lm *LoginMiddleware) initialize() {
	if lm.mutex == nil {
		lm.mutex = &sync.RWMutex{}
	}
	if lm.sessions == nil {
		lm.sessions = map[string]struct{}{}
	}
	if lm.login == nil {
		fsys := fs.FS(WWW)
		html, _ := fs.Sub(fsys, "www")
		www := http.FileServer(http.FS(html))
		lm.login = http.StripPrefix("/login", www)
	}
	if lm.PasswordKey == "" {
		lm.PasswordKey = "password"
	}
	if lm.SessionKey == "" {
		lm.SessionKey = "UFOSID"
	}
}

func (lm *LoginMiddleware) IsPasswordRequired() bool {
	return lm.Password != ""
}

func (lm *LoginMiddleware) IsLocalhost(r *http.Request) bool {
	hostonly, _, _ := strings.Cut(r.URL.Host, ":")
	return strings.HasSuffix(hostonly, "localhost")
}

// PrecheckAccessToken returns a bool that indicates whether the caller should continue
func (lm *LoginMiddleware) Wrap(next http.Handler) http.Handler {
	lm.initialize()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		// skip login if no password
		case !lm.IsPasswordRequired():
			break
		// skip login if on localhost
		case lm.IsLocalhost(r):
			break
		// validate session id for all requests
		case !lm.IsValidSession(r):
			// validate password for login requests
			if lm.IsValidLogin(r) {
				lm.SetCookiesAndRedirect(w, r)
				return
			}
			lm.RedirectToLogin(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}
