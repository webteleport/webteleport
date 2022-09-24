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
	"time"

	"github.com/btwiuse/h3/utils"
	"github.com/btwiuse/skynet"
	"github.com/lucas-clemente/quic-go/http3"
	"github.com/marten-seemann/webtransport-go"
)

var HOST = utils.EnvHost("skynet.k0s.io")

func webtransportServer(port string, next http.Handler) *webtransport.Server {
	s := &webtransport.Server{
		CheckOrigin: func(*http.Request) bool { return true },
	}
	s.H3 = http3.Server{
		Addr:            port,
		Handler:         webtransportHandler(s, next),
		EnableDatagrams: true,
	}
	return s
}

func webtransportHandler(s *webtransport.Server, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// passthrough normal requests:
		// 1. simple http
		// 2. websockets
		// 3. webtransport (not yet supported by reverseproxy)
		host, _, ok := strings.Cut(r.Host, ":")
		isSimple := !ok
		isWebtransport := ok && host != HOST
		isSkynetClient := ok && host == HOST
		switch {
		// passthrough requests made by webtransport-go, i.e.
		// strip the port:
		//
		// xxx.skynet.k0s.io:300
		// =>
		// xxx.skynet.k0s.io
		case isWebtransport:
			r.Host = host
			fallthrough
		case isSimple:
			next.ServeHTTP(w, r)
			return
		case isSkynetClient:
			break
		}
		log.Println("[01]", r.Proto, r.Method, r.Host, r.URL.Path)
		// handle skynet client registration
		// Host: skynet.k0s.io:300
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
	subhost := fmt.Sprintf("%d.%s", sm.counter, HOST)
	_, err = io.WriteString(stm0, fmt.Sprintf("HOST %s\n", subhost))
	if err != nil {
		return err
	}
	sm.counter += 1
	sm.sessions[host] = ssn
	go func() {
		var err error
		for {
			_, err = io.WriteString(stm0, fmt.Sprintf("%s\n", "PING"))
			if err != nil {
				break
			}
			time.Sleep(5 * time.Second)
		}
		delete(sm.sessions, host)
		log.Println(err, "deleted", host)
	}()
	return nil
}

func (sm *sessionManager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ssn, ok := sm.sessions[r.Host]
	if !ok {
		http.NotFoundHandler().ServeHTTP(w, r)
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
			stream, err := ssn.OpenStreamSync(ctx)
			return skynet.StreamConn{stream, ssn.LocalAddr(), ssn.RemoteAddr()}, err
		},
	}
	rp := &httputil.ReverseProxy{
		Director:  dr,
		Transport: tr,
	}
	rp.ServeHTTP(w, r)
}
