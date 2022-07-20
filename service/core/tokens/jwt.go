package tokens

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

const _minimumSecretKeySize = 32

type JwtBuilder struct {
	secretKey string
}

func NewJwtBuilder(secretKey string) (Builder, error) {
	if len(secretKey) < _minimumSecretKeySize {
		return nil, fmt.Errorf("invalid key size: key must be at least %d characters long", _minimumSecretKeySize)
	}
	return &JwtBuilder{secretKey}, nil
}

func (builder *JwtBuilder) CreateToken(user string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(user, duration)
	if err != nil {
		return "", payload, err
	}

	jwt_token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwt_token.SignedString([]byte(builder.secretKey))
	return token, payload, err
}

func (builder *JwtBuilder) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(builder.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
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
