package model

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

//type ConnInfo struct {
//	JumpHost string `json:"jump_host"`
//	Port     int    `json:"port"`
//	User     string `json:"user"`
//	Protocol string `json:"protocol"` // ssh / sftp
//	Client   string `json:"client"`   // securecrt / filezilla
//	Password string `json:"password"` // 临时凭证（一次性）
//}

var hmacKey = []byte("bastion-super-secret-key")

type SessionPayload struct {
	BastionHost string `json:"bh"` // 127.0.0.1 / bastion.example.com
	BastionPort int    `json:"bp"` // 2222
	Client      string `json:"c"`

	Secret string `json:"s"` // 随机 session secret

	IssuedAt int64 `json:"iat"`
	ExpireAt int64 `json:"exp"`
}

func ParseSession(token string) (*SessionPayload, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return nil, errors.New("invalid session format")
	}

	payloadB64 := parts[0]
	sigB64 := parts[1]

	sig, err := base64.RawURLEncoding.DecodeString(sigB64)
	if err != nil {
		return nil, err
	}

	mac := hmac.New(sha256.New, hmacKey)
	mac.Write([]byte(payloadB64))
	if !hmac.Equal(mac.Sum(nil), sig) {
		return nil, errors.New("invalid signature")
	}

	raw, err := base64.RawURLEncoding.DecodeString(payloadB64)
	if err != nil {
		return nil, err
	}

	var payload SessionPayload
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, err
	}

	if time.Now().Unix() > payload.ExpireAt {
		return nil, errors.New("session expired")
	}

	return &payload, nil
}
