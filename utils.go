package ufo

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/ebi-yade/altsvc-go"
	"github.com/mattn/go-isatty"
	"golang.org/x/net/idna"
)

// ExtractURLPort returns the :port part from URL.Host (host[:port])
//
// An empty string is returned if no port is found
func ExtractURLPort(u *url.URL) string {
	_, p, ok := strings.Cut(u.Host, ":")
	if ok {
		return ":" + p
	}
	return ""
}

// ToIdna converts a string to its idna form at best effort
// Should only be used on the hostname part without port
func ToIdna(s string) string {
	ascii, err := idna.ToASCII(s)
	if err != nil {
		log.Println(err)
		return s
	}
	return ascii
}

// ExtractAltSvcH3Endpoints reads Alt-Svc header
// returns a list of [host]:port endpoints
func ExtractAltSvcEndpoints(h http.Header, protocolId string) []string {
	line := h.Get("Alt-Svc")
	if line == "" {
		return nil
	}
	svcs, err := altsvc.Parse(line)
	if err != nil {
		log.Println(err)
		return nil
	}
	results := []string{}
	for _, svc := range svcs {
		if svc.ProtocolID != protocolId {
			continue
		}
		// host could be empty, port must not
		ep := svc.AltAuthority.Host + ":" + svc.AltAuthority.Port
		results = append(results, ep)
	}
	return results
}

// Graft returns Host(base):Port(alt)
//
// assuming
// - base is host[:port]
// - alt is [host]:port
func Graft(base, alt string) string {
	althost, altport, _ := strings.Cut(alt, ":")
	if altport == "" {
		// altport not found
		// it should never happen
		return base
	}
	if althost != "" {
		// alt is host:port
		// it is rare
		return alt
	}
	basehost, _, _ := strings.Cut(base, ":")
	return basehost + ":" + altport
}

// Print Hyperlink via OSC 8 ansi sequence.
// The syntax is: 'OSC 8 ; params ; url ST text OSC 8 ; ; ST'
// for more info see https://gist.github.com/egmontkob/eb114294efbcd5adb1944c9f3cb5feda
func Hyperlink(name, url string) string {
	return fmt.Sprintf("\u001B]8;%s;%s\u001B\\%s\u001B]8;;\u001B\\", "", url, name)
}

// MaybeHyperlink turns input into ANSI hyperlink when stdin is a tty
func MaybeHyperlink(l string) string {
	if isatty.IsTerminal(os.Stdin.Fd()) {
		return Hyperlink(l, l)
	}
	return l
}
