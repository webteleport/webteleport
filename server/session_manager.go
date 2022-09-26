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
	"golang.org/x/net/idna"
)

var DefaultSessionManager = &SessionManager{
	counter:  0,
	sessions: map[string]*webtransport.Session{},
	slock:    &sync.RWMutex{},
}

type SessionManager struct {
	counter  int
	sessions map[string]*webtransport.Session
	slock    *sync.RWMutex
}

func (sm *SessionManager) Del(k string) {
	k, _ = idna.ToASCII(k)
	sm.slock.Lock()
	delete(sm.sessions, k)
	sm.slock.Unlock()
}

func (sm *SessionManager) Get(k string) (*webtransport.Session, bool) {
	k, _ = idna.ToASCII(k)
	host, _, _ := strings.Cut(k, ":")
	sm.slock.RLock()
	ssn, ok := sm.sessions[host]
	sm.slock.RUnlock()
	return ssn, ok
}

func (sm *SessionManager) Add(k string, ssn *webtransport.Session) error {
	k, _ = idna.ToASCII(k)
	sm.slock.Lock()
	sm.counter += 1
	sm.sessions[k] = ssn
	sm.slock.Unlock()
	return nil
}

func (sm *SessionManager) Lease(ssn *webtransport.Session, domainList []string) error {
	stm0, err := ssn.OpenStreamSync(context.Background())
	if err != nil {
		return err
	}
	allowRandom := len(domainList) == 0
	var lease string
	for _, pfx := range domainList {
		k := fmt.Sprintf("%s.%s", pfx, HOST)
		if _, exist := sm.Get(k); !exist {
			lease = k
			break
		}
	}
	if (lease == "") && !allowRandom {
		emsg := fmt.Sprintf("ERR %s: %v\n", "none of your requested subdomains are currently available", domainList)
		_, err = io.WriteString(stm0, emsg)
		if err != nil {
			return err
		}
		return nil
	}
	if lease == "" {
		lease = fmt.Sprintf("%d.%s", sm.counter, HOST)
	}
	_, err = io.WriteString(stm0, fmt.Sprintf("HOST %s\n", lease))
	if err != nil {
		return err
	}
	err = sm.Add(lease, ssn)
	if err != nil {
		return err
	}
	go func() {
		var err error
		for {
			_, err = io.WriteString(stm0, fmt.Sprintf("%s\n", "PING"))
			if err != nil {
				break
			}
			time.Sleep(5 * time.Second)
		}
		sm.Del(lease)
		emsg := fmt.Sprintf("%s. Recycled %s", err, lease)
		log.Println(emsg)
	}()
	return nil
}

func (sm *SessionManager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Alt-Svc", ALT_SVC)
	ssn, ok := sm.Get(r.Host)
	if !ok {
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
	tr := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			// when there is a timeout, it still panics before MARK
			//
			// ctx, _ = context.WithTimeout(ctx, 3*time.Second)
			//
			// turns out the stream is empty so need to check stream == nil
			stream, err := ssn.OpenStreamSync(ctx)
			if err != nil {
				return nil, err
			}
			// once ctx got cancelled, err is nil but stream is empty too
			// add the check to avoid returning empty stream
			if stream == nil {
				return nil, fmt.Errorf("stream is empty")
			}
			// log.Println(`MARK`, stream)
			// MARK
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
