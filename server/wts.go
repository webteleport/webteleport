package server

import (
	"log"
	"net/http"
	"strings"

	"github.com/lucas-clemente/quic-go/http3"
	"github.com/marten-seemann/webtransport-go"
)

func webtransportServer(next http.Handler) *webtransport.Server {
	s := &webtransport.Server{
		CheckOrigin: func(*http.Request) bool { return true },
	}
	s.H3 = http3.Server{
		Addr:            PORT,
		Handler:         &wts{s, next},
		EnableDatagrams: true,
	}
	return s
}

// udp server handles
// - UFO client registration (CONNECT HOST)
// - requests over HTTP/3 (others)
type wts struct {
	*webtransport.Server
	Next http.Handler
}

func (s *wts) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// passthrough normal requests to next:
	// 1. simple http / websockets (Host: x.localhost)
	// 2. webtransport (Host: x.localhost:300, not yet supported by reverseproxy)
	if !isUFO(r) {
		s.Next.ServeHTTP(w, r)
		return
	}
	log.Println("[UFO]", r.RemoteAddr, r.Proto, r.Method, r.Host)
	// handle ufo client registration
	// Host: ufo.k0s.io:300
	ssn, err := s.Upgrade(w, r)
	if err != nil {
		log.Printf("upgrading failed: %s", err)
		w.WriteHeader(500)
		return
	}
	err = defaultSessionManager.Add(ssn)
	if err != nil {
		log.Printf("adding session failed: %s", err)
	}
}

func isUFO(r *http.Request) bool {
	host, _, _ := strings.Cut(r.Host, ":")
	isConnect := r.Method == http.MethodConnect
	isRoot := host == HOST
	return isRoot && isConnect
}
