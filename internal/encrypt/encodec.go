package encrypt

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

// Name is the name registered for the proto compressor.
const Name = "encodec"

type Encrypter interface {
	Encrypt(b []byte) ([]byte, error)
}

type Decrypter interface {
	Decrypt(b []byte) ([]byte, error)
}

// encodec is a Codec with encryption functionality implementation with protobuf. It is the default codec for gRPC.
type Encodec struct {
	Enc Encrypter
	Dec Decrypter
}

// Marshal encodes mesage and encrypts the encoded result. If Enc is nil, encryption is omitted.
func (e Encodec) Marshal(v interface{}) (b []byte, err error) {
	log.Debug().Str("call", "Marshal").Msg("Encodec")
	vv, ok := v.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("failed to marshal, message is %T, want proto.Message", v)
	}
	b, err = proto.Marshal(vv)
	if err != nil {
		log.Log().AnErr("Marshal", err).Msg("Encodec")
		return nil, err
	}
	if e.Enc == nil {
		return
	}
	return e.Enc.Encrypt(b)
}

// Unmarshal decrypts message and decodes result. If Dec is nil, decryption is omitted.
func (e Encodec) Unmarshal(data []byte, v interface{}) (err error) {
	log.Debug().Str("call", "Unmarshal").Msg("Encodec")
	if e.Dec != nil {
		data, err = e.Dec.Decrypt(data)
		if err != nil {
			return
		}
	}
	vv, ok := v.(proto.Message)
	if !ok {
		return fmt.Errorf("failed to unmarshal, message is %T, want proto.Message", v)
	}
	return proto.Unmarshal(data, vv)
}

func (Encodec) Name() string {
	return Name
}
