package api_keycloak

import (
	"crypto/rsa"
	"fmt"
	app_errors "src/errors"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

var CacheJwkSet = &KeycloakJwkSet{Keys: make([]KeycloakJwk, 0)}

func GetRsaPublicKey() (*rsa.PublicKey, app_errors.AppError) {
	if len(CacheJwkSet.Keys) != 0 {
		// fmt.Println("GET JWK FROM CACHE ", CacheJwkSet)
		signingJwk, err := CacheJwkSet.GetSigJwk()
		if err != nil {
			return nil, err
		}
		key, err := signingJwk.ComputePublicRsaKey()
		if err != nil {
			return nil, err
		}
		return &key, nil
	} else {
		// fmt.Println("GET JWK FROM REMOTE")

		jwkSet, err := GetJwkCerts()
		CacheJwkSet = &jwkSet
		if err != nil {
			return nil, err
		}
		signingJwk, err := jwkSet.GetSigJwk()
		if err != nil {
			return nil, err
		}
		key, err := signingJwk.ComputePublicRsaKey()
		if err != nil {
			return nil, err
		}
		return &key, nil
	}
}

func VerifyToken(token string) (*jwt.Token, app_errors.AppError) {
	rsaKey, err := GetRsaPublicKey()
	if err != nil {
		return nil, err
	}
	jwk, err := CacheJwkSet.GetSigJwk()
	if err != nil {
		return nil, err
	}
	parsedToken, jwtError := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Verificar el algoritmo de firma
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Verificar el kid
		kid, ok := token.Header["kid"].(string)
		if !ok || kid != jwk.Kid {
			return nil, fmt.Errorf("invalid or missing kid: expected %s, got %s", jwk.Kid, kid)
		}
		return rsaKey, nil
	})

	if jwtError != nil {
		return nil, &app_errors.ErrVerifyToken{Message: jwtError.Error()}
	}

	if !parsedToken.Valid {
		return nil, &app_errors.ErrVerifyToken{Message: "Token is not valid"}
	}

	return parsedToken, nil
}

func VerifyClaims(token *jwt.Token) app_errors.AppError {

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return &app_errors.ErrVerifyToken{Message: "error obtaining claims"}
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return &app_errors.ErrVerifyToken{Message: "invalid expiration claim"}
	}
	if exp < float64(time.Now().Unix()) {
		return &app_errors.ErrVerifyToken{Message: "token is expired"}
	}

	// // subject
	// if sub, ok := claims["sub"].(string); !ok || sub == "" {
	// 	t.Error("Invalid or missing sub claim")
	// }
	return nil
}
