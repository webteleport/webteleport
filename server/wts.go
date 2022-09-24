package server

import (
	"log"
	"net/http"
	"strings"

	"github.com/lucas-clemente/quic-go/http3"
	"github.com/marten-seemann/webtransport-go"
)

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
