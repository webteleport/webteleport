package server

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/btwiuse/h3/utils"
	"k0s.io/pkg/middleware"
)

var HOST = utils.EnvHost("localhost")
var CERT = utils.EnvCert("localhost.pem")
var KEY = utils.EnvKey("localhost-key.pem")
var PORT = utils.EnvPort(":3000")
var ALT_SVC = utils.EnvAltSvc(fmt.Sprintf(`webteleport="%s"`, PORT))

func Run([]string) error {
	log.Println("listening on TCP http://" + HOST + PORT)
	ln, err := net.Listen("tcp4", PORT)
	if err != nil {
		return err
	}

	handler := middleware.LoggingMiddleware(middleware.AllowAllCorsMiddleware(DefaultSessionManager))

	go func() {
		wts := WebtransportServer(handler)
		log.Println("listening on UDP https://" + HOST + PORT)
		log.Fatalln(wts.ListenAndServeTLS(CERT, KEY))
	}()

	return http.Serve(ln, handler)
}
