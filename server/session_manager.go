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

	"github.com/btwiuse/ufo"
	"github.com/marten-seemann/webtransport-go"
)

var defaultSessionManager = &sessionManager{
	counter:  0,
	sessions: map[string]*webtransport.Session{},
	slock:    &sync.RWMutex{},
}

type sessionManager struct {
	counter  int
	sessions map[string]*webtransport.Session
	slock    *sync.RWMutex
}

func (sm *sessionManager) Del(k string) {
	sm.slock.Lock()
	delete(sm.sessions, k)
	sm.slock.Unlock()
}

func (sm *sessionManager) Get(k string) (*webtransport.Session, bool) {
	host, _, _ := strings.Cut(k, ":")
	sm.slock.RLock()
	ssn, ok := sm.sessions[host]
	sm.slock.RUnlock()
	return ssn, ok
}

func (sm *sessionManager) Add(ssn *webtransport.Session) error {
	stm0, err := ssn.OpenStreamSync(context.Background())
	if err != nil {
		return err
	}
	subhost := fmt.Sprintf("%d.%s", sm.counter, HOST)
	_, err = io.WriteString(stm0, fmt.Sprintf("HOST %s\n", subhost))
	if err != nil {
		return err
	}
	sm.slock.Lock()
	sm.counter += 1
	sm.sessions[subhost] = ssn
	sm.slock.Unlock()
	go func() {
		var err error
		for {
			_, err = io.WriteString(stm0, fmt.Sprintf("%s\n", "PING"))
			if err != nil {
				break
			}
			time.Sleep(5 * time.Second)
		}
		sm.Del(subhost)
		log.Println(err, "deleted", subhost)
	}()
	return nil
}

func (sm *sessionManager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ssn, ok := sm.Get(r.Host)
	if !ok {
		w.Header().Set("Alt-Svc", ALT_SVC)
		http.Error(w, r.Host+" not found", http.StatusNotFound)
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
		// req.Proto = r.Proto
	}
	stmN := 0
	tr := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			stream, err := ssn.OpenStreamSync(ctx)
			if err != nil {
				log.Println(err, stmN)
				return nil, err
			}
			stmN += 1
			conn := &ufo.StreamConn{stream, ssn}
			return conn, nil
		},
	}
	rp := &httputil.ReverseProxy{
		Director:  dr,
		Transport: tr,
	}
	rp.ServeHTTP(w, r)
}
