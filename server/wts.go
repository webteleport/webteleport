package main

import (
	"context"
	_ "embed"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"time"

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
		ssn, err := s.Upgrade(w, r)
		if err != nil {
			log.Printf("upgrading failed: %s", err)
			w.WriteHeader(500)
			return
		}
		// go echoConn(conn)
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

type streamConn struct {
	webtransport.Stream
	lad net.Addr
	rad net.Addr
}

func (sc streamConn) LocalAddr() net.Addr  { return sc.lad }
func (sc streamConn) RemoteAddr() net.Addr { return sc.rad }

func (sm *sessionManager) Add(ssn *webtransport.Session) error {
	stm0, err := ssn.OpenStreamSync(context.Background())
	if err != nil {
		return err
	}
	host := fmt.Sprintf("%d.quichost.k0s.io", sm.counter)
	_, err = io.WriteString(stm0, fmt.Sprintf("HOST %d\n", sm.counter))
	if err != nil {
		return err
	}
	sm.counter += 1
	sm.sessions[host] = ssn
	go func() {
		for {
			io.WriteString(stm0, fmt.Sprintf("%s\n", "LINE"))
			time.Sleep(5 * time.Second)
		}
		delete(sm.sessions, host)
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
			return streamConn{stream, ssn.LocalAddr(), ssn.RemoteAddr()}, err
		},
	}
	rp := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			log.Println("director: rewriting Host", r.URL, r.Host)
			req.Host = r.Host
			req.URL.Host = r.Host
		},
		Transport: tr,
	}
	rp.ServeHTTP(w, r)
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
