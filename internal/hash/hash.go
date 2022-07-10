// Package hash provides hashing functionality.
package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/andrei-cloud/go-devops/internal/model"
)

// Create - creates SHA256 hash for given string using provided key.
func Create(src string, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(src))
	return hex.EncodeToString(h.Sum(nil))
}

// Validate - checks if given metric and it's hash is valid for key provided.
func Validate(m model.Metric, key []byte) (bool, error) {
	var data string
	if len(key) == 0 {
		return true, nil
	}

	h, err := hex.DecodeString(m.Hash)
	if err != nil {
		return false, err
	}

	switch m.MType {
	case "gauge":
		data = fmt.Sprintf("%s:gauge:%f", m.ID, *m.Value)
	case "counter":
		data = fmt.Sprintf("%s:counter:%d", m.ID, *m.Delta)
	}
	d, err := hex.DecodeString(Create(data, key))
	if err != nil {
		return false, err
	}
	return hmac.Equal(h, d), nil
}
