package encrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"os"

	"github.com/rs/zerolog/log"
)

type Encrypter interface {
	Encrypt(b []byte) ([]byte, error)
	Decrypt(b []byte) ([]byte, error)
}

type encrypt struct {
	key any
}

func New(path string) *encrypt {
	key := readKeyFile(path)
	return &encrypt{
		key: key,
	}
}

func (e encrypt) Encrypt(b []byte) ([]byte, error) {
	var cipherbytes []byte

	hash := sha512.New()
	label := []byte("metrics")
	pk := e.key.(*rsa.PublicKey)

	msgLen := len(b)
	step := pk.Size() - 2*hash.Size() - 2

	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}

		cipherBlock, err := rsa.EncryptOAEP(hash, rand.Reader, pk, b[start:finish], label)
		if err != nil {
			log.Debug().AnErr("EncryptOAEP", err).Msg("Encryption Failed")
			return nil, err
		}

		cipherbytes = append(cipherbytes, cipherBlock...)
	}

	return cipherbytes, nil
}

func (e encrypt) Decrypt(b []byte) ([]byte, error) {
	var plainBytes []byte

	hash := sha512.New()
	label := []byte("metrics")
	pk := e.key.(*rsa.PrivateKey)
	step := pk.PublicKey.Size()
	msgLen := len(b)

	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}

		plainBlock, err := rsa.DecryptOAEP(hash, rand.Reader, pk, b[start:finish], label)
		if err != nil {
			log.Debug().AnErr("DecryptOAEP", err).Msg("Dencryption Failed")
			return nil, err
		}

		plainBytes = append(plainBytes, plainBlock...)
	}

	return plainBytes, nil
}

func readKeyFile(path string) any {
	var key any
	pemBytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	block, _ := pem.Decode(pemBytes)
	if block == nil {
		log.Fatal().Msg("failed to decode PEM block")
	}

	if block != nil {
		switch block.Type {
		case "RSA PRIVATE KEY":
			key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
			if err != nil {
				log.Fatal().Msg(err.Error())
			}
		case "PUBLIC KEY":
			key, err = x509.ParsePKIXPublicKey(block.Bytes)
			if err != nil {
				log.Fatal().Msg(err.Error())
			}
		default:
			log.Fatal().Msg("not supported PEM block")
		}
	}

	return key
}
