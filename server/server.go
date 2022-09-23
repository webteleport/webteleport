package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	"github.com/btwiuse/h3/utils"
)

func Run([]string) error {
	port := utils.EnvPort(":3000")
	altsvc := utils.EnvAltSvc(fmt.Sprintf(`h3="%s"`, port))
	log.Println("listening on TCP http://127.0.0.1" + port)
	ln, err := net.Listen("tcp4", port)
	if err != nil {
		return err
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, ok := defaultSessionManager.sessions[r.Host]
		if !ok {
			w.Header().Set("Alt-Svc", altsvc)
			io.WriteString(w, r.Host)
			return
		}
		defaultSessionManager.ServeHTTP(w, r)
	})

	go func() {
		wts := webtransportServer(port, http.DefaultServeMux)
		cert := utils.EnvCert("localhost.pem")
		key := utils.EnvKey("localhost-key.pem")
		log.Println("listening on UDP https://127.0.0.1" + port)
		log.Fatalln(wts.ListenAndServeTLS(cert, key))
	}()

	return http.Serve(ln, http.DefaultServeMux)
}
