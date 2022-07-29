package middlewares

import (
	"bytes"
	"io"
	"net/http"

	"github.com/andrei-cloud/go-devops/internal/encrypt"
	"github.com/rs/zerolog/log"
)

type cryptoRT struct {
	next http.RoundTripper
	encr encrypt.Encrypter
}

func NewCryptoRT(encr encrypt.Encrypter) *cryptoRT {
	return &cryptoRT{next: http.DefaultTransport, encr: encr}
}

func (e cryptoRT) RoundTrip(req *http.Request) (res *http.Response, err error) {
	if e.encr == nil {
		return e.next.RoundTrip(req)
	}

	var b bytes.Buffer

	if req.Method == http.MethodPost || req.Method == http.MethodPut || req.Method == http.MethodPatch {
		_, err = b.ReadFrom(req.Body)
		if err != nil {
			log.Log().AnErr("ReadFrom", err).Msg("CryptoRT")
			return nil, err
		}

		eb, err := e.encr.Encrypt(b.Bytes())
		if err != nil {
			log.Log().AnErr("Encrypt", err).Msg("CryptoRT")
			return nil, err
		}

		req.Body = io.NopCloser(bytes.NewReader(eb))
		req.ContentLength = int64(len(eb))
	}

	return e.next.RoundTrip(req)
}
