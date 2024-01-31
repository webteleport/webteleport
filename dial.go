package webteleport

import (
	"context"
	"errors"
	// "fmt"
	"net/http"
	"net/url"

	// "github.com/quic-go/quic-go"
	// "github.com/quic-go/quic-go/http3"
	webtransportGo "github.com/quic-go/webtransport-go"
	"github.com/webteleport/utils"
	"github.com/webteleport/webteleport/webtransport"
)

// 2^60 == 1152921504606846976
const MaxIncomingStreams int64 = 1 << 60

// Dial is a wrapper around webtransport.Dial with automatic HTTP/3 service discovery
func Dial(ctx context.Context, u *url.URL, hdr http.Header) (*webtransportGo.Session, error) {
	endpoints := Resolve(u)
	if len(endpoints) == 0 {
		return nil, errors.New("service discovery failed: no webteleport endpoints found in Alt-Svc records / headers")
	}
	alt := endpoints[0]
	addr := utils.Graft(u.Host, alt)
	return webtransport.DialWebtransport(ctx, addr, hdr)
}

/*
type webtransportSession struct {
	*webtransport.Session
}

func (s *webtransportSession) AcceptStream(ctx context.Context) (Stream, error) {
	stm, err := s.Session.AcceptStream(ctx)
	return &webtransportStream{stm}, err
}

type webtransportStream struct {
	webtransport.Stream
}
*/
