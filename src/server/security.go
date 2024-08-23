package server

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"io"
	"net/http"
)

type Security struct {
	publicKeysCache map[string]*rsa.PublicKey
	urlGetPublicKey string
}

func NewSecurity(url string) Security {
	return Security{make(map[string]*rsa.PublicKey), url}
}

func (s *Security) canPlay(tokenAsString string) bool {
	token, err := s.getJWTToken(tokenAsString)
	if err != nil {
		return false
	}
	if isAdmin, exist := token.Claims.(jwt.MapClaims)["is_admin"]; !exist || isAdmin != true {
		if isGuest, exist := token.Claims.(jwt.MapClaims)["guest"]; !exist || isGuest != true {
			return false
		}
	}
	return true
}

func (s *Security) checkAdmin(tokenAsString string) bool {
	token, err := s.getJWTToken(tokenAsString)
	if err != nil {
		return false
	}
	if isAdmin, exist := token.Claims.(jwt.MapClaims)["is_admin"]; !exist || isAdmin != true {
		return false
	}
	return true
}

func (s *Security) getJWTToken(tokenAsString string) (*jwt.Token, error) {
	return jwt.Parse(tokenAsString, func(t *jwt.Token) (interface{}, error) {
		kid, exist := t.Header["kid"]
		if !exist {
			return nil, errors.New("missing kid in geader")
		}
		return s.findPublicKey(kid.(string))
	})
}

func (s *Security) findPublicKey(kid string) (*rsa.PublicKey, error) {
	key, exist := s.publicKeysCache[kid]
	if !exist {
		var err error
		key, err = s.getPublicKey(kid)
		if err != nil {
			return nil, err
		}
		s.publicKeysCache[kid] = key
	}
	return key, nil
}

func (s *Security) getPublicKey(kid string) (*rsa.PublicKey, error) {
	resp, err := http.Get(fmt.Sprintf("%s?kid=%s", s.urlGetPublicKey, kid))
	if err != nil {
		return nil, err
	}
	if data, err := io.ReadAll(resp.Body); err == nil {
		return jwt.ParseRSAPublicKeyFromPEM(data)
	} else {
		return nil, err
	}
}
