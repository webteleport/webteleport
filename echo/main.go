// this is an example consumer of the ufo package
// it listens on a random ufo with registered handlers for webtransport && websocket connections
// currently websocket works fine
// while webtransport is broken because reverseproxy doesn't support it yet

package echo

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/btwiuse/pretty"
	"github.com/btwiuse/ufo"
	"github.com/lucas-clemente/quic-go/http3"
	"k0s.io/pkg/wrap"
)

func Run([]string) error {
	ln, err := ufo.Listen("https://ufo.k0s.io")
	if err != nil {
		return err
	}
	addr := ln.Addr()
	location := fmt.Sprintf("%s://%s", addr.Network(), addr.String())
	log.Println("listening on", location)
	return http.Serve(ln, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isWT := r.Method == http.MethodConnect
		isWS := strings.Split(r.Header.Get("Upgrade"), ",")[0] == "websocket"
		w.Header().Set("Alt-Svc", `h3=":300"`)
		log.Println(pretty.YAMLString(r.Header))
		log.Println("Method", r.Method)
		log.Println("Proto", r.Proto)
		log.Println("Path", r.URL.Path)
		log.Println("isWS", isWS)
		log.Println("isWT", isWT)

		if isWT {
			handleWebtransport(w, r)
			return
		}

		if isWS {
			handleWebsocket(w, r)
			return
		}
	}))
}

// HOST=0.ufo.k0s.io PORT=300 h3 client
// websocat --binary wss://0.ufo.k0s.io
func handleWebsocket(w http.ResponseWriter, r *http.Request) {
	connSrc, err := wrap.Wrconn(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	go io.Copy(connSrc, connSrc)
	go io.Copy(connSrc, connSrc)
}

// HOST=0.ufo.k0s.io PORT=300 h3 client
// doesn't work yet
func handleWebtransport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodConnect {
		log.Println("warning: expected CONNECT request, got", r.Method)
	}
	if r.Proto != "webtransport" {
		log.Println("warning: expected webtransport proto, got", r.Proto)
	}
	if v, ok := r.Header["Sec-Webtransport-Http3-Draft02"]; !ok || len(v) != 1 || v[0] != "1" {
		log.Println("warning: missing or invalid header:", "Sec-Webtransport-Http3-Draft02")
	}
	w.Header().Add("Sec-Webtransport-Http3-Draft02", "draft02")
	w.WriteHeader(http.StatusOK)
	w.(http.Flusher).Flush()

	httpStreamer, ok := r.Body.(http3.HTTPStreamer)
	if !ok {
		// should never happen
		log.Println("warning: failed to take over HTTP stream")
	}
	str := httpStreamer.HTTPStream()
	sID := uint64(str.StreamID())

	hijacker, ok := w.(http3.Hijacker)
	if !ok {
		// should never happen
		log.Println("warning: failed to hijack")
	}
	_ = sID
	_ = hijacker
}

// curl3 https://7.ufo.k0s.io:300 --http3 -H "Host: 7.ufo.k0s.io"
