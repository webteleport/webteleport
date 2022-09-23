package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/lucas-clemente/quic-go/http3"
	"github.com/marten-seemann/webtransport-go"
	"github.com/btwiuse/quichost"
)

func webtransportServer(port string) *webtransport.Server {
	s := &webtransport.Server{
		CheckOrigin: func(*http.Request) bool { return true },
	}
	s.H3 = http3.Server{
		Addr:            port,
		Handler:         webtransportHandler(s),
		EnableDatagrams: true,
	}
	return s
}

func webtransportHandler(s *webtransport.Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ssn, err := s.Upgrade(w, r)
		if err != nil {
			log.Printf("upgrading failed: %s", err)
			w.WriteHeader(500)
			return
		}
		err = defaultSessionManager.Add(ssn)
		if err != nil {
			log.Printf("initializing session failed: %s", err)
		}
	})
}

var defaultSessionManager = &sessionManager{
	counter:  0,
	sessions: map[string]*webtransport.Session{},
}

type sessionManager struct {
	counter  int
	sessions map[string]*webtransport.Session
}

func (sm *sessionManager) Add(ssn *webtransport.Session) error {
	stm0, err := ssn.OpenStreamSync(context.Background())
	if err != nil {
		return err
	}
	host := fmt.Sprintf("%d.quichost.k0s.io", sm.counter)
	_, err = io.WriteString(stm0, fmt.Sprintf("HOST %s\n", host))
	if err != nil {
		return err
	}
	sm.counter += 1
	sm.sessions[host] = ssn
	go func() {
		for {
			io.WriteString(stm0, fmt.Sprintf("%s\n", "PING"))
			time.Sleep(5 * time.Second)
		}
		delete(sm.sessions, host)
		log.Println("deleted", host)
	}()
	return nil
}

func (sm *sessionManager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ssn, ok := sm.sessions[r.Host]
	if !ok {
		http.NotFoundHandler().ServeHTTP(w, r)
		return
	}
	tr := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			stream, err := ssn.OpenStreamSync(ctx)
			return quichost.StreamConn{stream, ssn.LocalAddr(), ssn.RemoteAddr()}, err
		},
	}
	rp := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			log.Println("director: rewriting Host", r.URL, r.Host)
			req.Host = r.Host
			req.URL.Host = r.Host
			req.URL.Scheme = "http"
		},
		Transport: tr,
	}
	rp.ServeHTTP(w, r)
}
