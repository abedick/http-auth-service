package main

import (
	"errors"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type publicToken struct {
	Debug    bool    `json:"debug"`
	ID       string  `json:"id"`
	Username string  `json:"username"`
	Exp      float64 `json:"exp"`
	Iat      float64 `json:"iat"`
	Iss      string  `json:"iss"`
	Access   int     `json:"access"`
}

type requestValidate struct {
	Token string `json:"token"`
}

type requestToken struct {
	AccessMode string      `json:"access_mode"`
	Claims     tokenClaims `json:"claims"`
}

type tokenClaims struct {
	ID     string            `json:"id"`
	ISS    string            `json:"iss"`
	Custom map[string]string `json:"custom_claims"`
}

// generateToken - generates a JWT token
func generateToken(request requestToken) (string, error) {

	tokenSecret := c.TokenConfig.Signature
	key := []byte(tokenSecret)
	token := jwt.New(jwt.SigningMethodHS256)
	header := token.Header
	claims := token.Claims.(jwt.MapClaims)

	access := c.TokenConfig.AccessMap[request.AccessMode]

	if access != "" {
		header["kid"] = access
		header["mode"] = request.AccessMode
	}

	claims["exp"] = time.Now().Add(*c.UserConfig.goDuration).Unix()
	claims["iat"] = time.Now().Unix()
	claims["iss"] = request.Claims.ISS
	claims["id"] = request.Claims.ID

	for k, v := range request.Claims.Custom {
		claims[k] = v
	}

	if c.UserConfig.Debug {
		claims["debug"] = true
	}

	tokenString, err := token.SignedString(key)
	return tokenString, err
}

// validateToken -- returns token claims if valid, else error
func validateToken(tokenString string) (map[string]interface{}, error) {
	var emptyMap map[string]interface{}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		mode := token.Header["mode"].(string)
		if token.Header["kid"] != c.TokenConfig.AccessMap[mode] {
			return "", errors.New("token KID did not match known kid")
		}

		tokenSecret := c.TokenConfig.Signature
		key := []byte(tokenSecret)
		return key, nil
	})
	if err != nil {
		return emptyMap, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		claims["access_mode"] = token.Header["mode"].(string)
		return claims, nil
	}
	return emptyMap, errors.New("could not validate claims")
}

// validateTokenRequest -- validates attributes of token requests
func validateTokenRequest(req requestToken) error {
	if c.TokenConfig.AccessMap[req.AccessMode] == "" {
		return errors.New("requested access mode is not recognized")
	}
	if req.Claims.ID == "" {
		return errors.New("must specify id")
	}
	if req.Claims.ISS == "" {
		return errors.New("must specify issuer")
	}
	return nil
}
