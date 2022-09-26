package server

import (
	"log"
	"net/http"
	"strings"

	"github.com/lucas-clemente/quic-go/http3"
	"github.com/marten-seemann/webtransport-go"
)

const WebTransportProto = "webtransport"

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
	domainList, isUFO := parseUFO(r)
	if !isUFO {
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
	// log.Println(domainList)
	err = defaultSessionManager.Lease(ssn, domainList)
	if err != nil {
		log.Printf("leasing failed: %s", err)
	}
}

func parseUFO(r *http.Request) ([]string, bool) {
	host, _, _ := strings.Cut(r.Host, ":")
	isWebtransport := r.Proto == WebTransportProto
	isConnect := r.Method == http.MethodConnect
	isRoot := host == HOST
	isUFO := isWebtransport && isRoot && isConnect
	if !isUFO {
		return nil, false
	}
	list := domainList(r.URL.Path)
	return list, true
}

func domainList(p string) []string {
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
