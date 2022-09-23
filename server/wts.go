package main

import (
	"context"
	_ "embed"
	"io"
	"log"
	"net/http"

	"github.com/lucas-clemente/quic-go/http3"
	"github.com/marten-seemann/webtransport-go"
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
		conn, err := s.Upgrade(w, r)
		if err != nil {
			log.Printf("upgrading failed: %s", err)
			w.WriteHeader(500)
			return
		}
		go echoConn(conn)
	})
}

func echoConn(conn *webtransport.Session) {
	log.Println(conn.RemoteAddr(), "new session")
	ctx := context.Background()
	for {
		stream, err := conn.AcceptStream(ctx)
		if err != nil {
			log.Println(conn.RemoteAddr(), "session closed")
			break
		}
		log.Println(conn.RemoteAddr(), "new stream")
		go io.Copy(stream, stream)
	}
}
