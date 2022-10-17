package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/idna"
)

var DefaultSessionManager = &SessionManager{
	counter:  0,
	sessions: map[string]*Session{},
	slock:    &sync.RWMutex{},
}

type SessionManager struct {
	counter  int
	sessions map[string]*Session
	slock    *sync.RWMutex
}

func (sm *SessionManager) Del(k string) error {
	k, err := idna.ToASCII(k)
	if err != nil {
		return err
	}
	sm.slock.Lock()
	delete(sm.sessions, k)
	sm.slock.Unlock()
	return nil
}

func (sm *SessionManager) DelSession(ssn *Session) {
	sm.slock.Lock()
	for k, v := range sm.sessions {
		if v == ssn {
			delete(sm.sessions, k)
			emsg := fmt.Sprintf("Recycled %s", k)
			log.Println(emsg)
		}
	}
	sm.slock.Unlock()
}

func (sm *SessionManager) Get(k string) (*Session, bool) {
	k, _ = idna.ToASCII(k)
	host, _, _ := strings.Cut(k, ":")
	sm.slock.RLock()
	ssn, ok := sm.sessions[host]
	sm.slock.RUnlock()
	return ssn, ok
}

func (sm *SessionManager) Add(k string, ssn *Session) error {
	k, err := idna.ToASCII(k)
	if err != nil {
		return err
	}
	sm.slock.Lock()
	sm.counter += 1
	sm.sessions[k] = ssn
	sm.slock.Unlock()
	return nil
}

func (sm *SessionManager) Lease(ssn *Session, candidates []string) error {
	var err error
	allowRandom := len(candidates) == 0
	var lease string
	for _, pfx := range candidates {
		k := fmt.Sprintf("%s.%s", pfx, HOST)
		if _, exist := sm.Get(k); !exist {
			lease = k
			break
		}
	}
	if (lease == "") && !allowRandom {
		emsg := fmt.Sprintf("ERR %s: %v\n", "none of your requested subdomains are currently available", candidates)
		_, err = io.WriteString(ssn.Controller, emsg)
		if err != nil {
			return err
		}
		return nil
	}
	if lease == "" {
		lease = fmt.Sprintf("%d.%s", sm.counter, HOST)
	}
	_, err = io.WriteString(ssn.Controller, fmt.Sprintf("HOST %s\n", lease))
	if err != nil {
		return err
	}
	log.Println("leasing", lease)
	err = sm.Add(lease, ssn)
	if err != nil {
		return err
	}
	log.Println("leasing", "x"+lease)
	sm.Add("x"+lease, ssn)
	return nil
}

func (sm *SessionManager) Ping(ssn *Session) {
	for {
		_, err := io.WriteString(ssn.Controller, fmt.Sprintf("%s\n", "PING"))
		if err != nil {
			break
		}
		time.Sleep(5 * time.Second)
	}
	sm.DelSession(ssn)
}

func (sm *SessionManager) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	NotFoundHandler().ServeHTTP(w, r)
}

func (sm *SessionManager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Alt-Svc", ALT_SVC)
	// for HTTP_PROXY r.Method = GET && r.Host = google.com
	// for HTTPs_PROXY r.Method = GET && r.Host = google.com:443
	// they are currently not supported and will be handled by the 404 handler
	ssn, ok := sm.Get(r.Host)
	if !ok {
		sm.NotFoundHandler(w, r)
		return
	}

	dr := func(req *http.Request) {
		// log.Println("director: rewriting Host", r.URL, r.Host)
		req.Host = r.Host
		req.URL.Host = r.Host
		req.URL.Scheme = "http"
		// for webtransport, Proto is "webtransport" instead of "HTTP/1.1"
		// However, reverseproxy doesn't support webtransport yet
		// so setting this field currently doesn't have any effect
		req.Proto = r.Proto
	}
	tr := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return ssn.OpenConn(ctx)
		},
	}
	rp := &httputil.ReverseProxy{
		Director:  dr,
		Transport: tr,
	}
	rp.ServeHTTP(w, r)
}
