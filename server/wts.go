package server

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/lucas-clemente/quic-go/http3"
	"github.com/marten-seemann/webtransport-go"
)

const WebTransportProto = "webtransport"

func WebtransportServer(next http.Handler) *webtransport.Server {
	s := &webtransport.Server{
		CheckOrigin: func(*http.Request) bool { return true },
	}
	s.H3 = http3.Server{
		Addr:            PORT,
		Handler:         &WTS{s, next},
		EnableDatagrams: true,
	}
	return s
}

// WTS is a HTTP/3 server that handles:
// - UFO client registration (CONNECT HOST)
// - requests over HTTP/3 (others)
type WTS struct {
	*webtransport.Server
	Next http.Handler
}

func (s *WTS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// passthrough normal requests to next:
	// 1. simple http / websockets (Host: x.localhost)
	// 2. webtransport (Host: x.localhost:300, not yet supported by reverseproxy)
	if !IsUFORequest(r) {
		s.Next.ServeHTTP(w, r)
		return
	}
	log.Println("🛸", r.RemoteAddr, r.Proto, r.Method, r.Host, r.URL.Path)
	// handle ufo client registration
	// Host: ufo.k0s.io:300
	ssn, err := s.Upgrade(w, r)
	if err != nil {
		log.Printf("upgrading failed: %s", err)
		w.WriteHeader(500)
		return
	}
	session := &Session{
		Session: ssn,
	}
	err = session.InitController(context.Background())
	if err != nil {
		log.Printf("session init failed: %s", err)
		return
	}
	candidates := ParseDomainCandidates(r.URL.Path)
	err = DefaultSessionManager.Lease(session, candidates)
	if err != nil {
		log.Printf("leasing failed: %s", err)
		return
	}
	go DefaultSessionManager.Ping(session)
}

// IsUFORequest tells if the incoming request should be treated as UFO request
//
// if true, it will be upgraded into a webtransport session
// otherwise the request will be handled by DefaultSessionManager
func IsUFORequest(r *http.Request) bool {
	host, _, _ := strings.Cut(r.Host, ":")
	isWebtransport := r.Proto == WebTransportProto
	isConnect := r.Method == http.MethodConnect
	isRoot := host == HOST
	return isWebtransport && isRoot && isConnect
}

// ParseDomainCandidates splits a path string like /a/b/cd/😏
// into a list of subdomains: [a, b, cd, 😏]
//
// when result is empty, a random subdomain will be assigned by the server
func ParseDomainCandidates(p string) []string {
	var list []string
	parts := strings.Split(p, "/")
	for _, part := range parts {
		dom := strings.Trim(part, " ")
		if dom == "" {
			continue
		}
		list = append(list, dom)
	}
	return list
}
