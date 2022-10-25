package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const minSecretKeyLen = 32

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeyLen {
		return nil, fmt.Errorf("invalid secret key length")
	}
	return &JWTMaker{
		secretKey,
	}, nil
}

func (t *JWTMaker) CreateToken(userName string, duration time.Duration) (string, error) {
	payload, err := NewPayload(userName, duration)
	if err != nil {
		return "", err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(t.secretKey))
}

func (t *JWTMaker) VerifyToken(tString string) (*Payload, error) {
	keyFunc := func(tString *jwt.Token) (interface{}, error) {
		_, ok := tString.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(t.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(tString, &Payload{}, keyFunc)
	if err != nil {
		vErr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(vErr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}
	return payload, nil
}
