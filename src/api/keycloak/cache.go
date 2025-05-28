package api_keycloak

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var CacheJwkSet = &KeycloakJwkSet{Keys: make([]KeycloakJwk, 0)}

func GetRsaPublicKey() (*rsa.PublicKey, error) {
	if len(CacheJwkSet.Keys) != 0 {
		fmt.Println("GET JWK FROM CACHE ", CacheJwkSet)
		signingJwk, err := CacheJwkSet.GetSigJwk()
		if err != nil {
			return nil, fmt.Errorf("could not get JWK %s", err.Error())
		}
		key, err := signingJwk.ComputePublicRsaKey()
		if err != nil {
			return nil, fmt.Errorf("could not get RSA public key")
		}
		return &key, nil
	} else {
		fmt.Println("GET JWK FROM REMOTE")

		jwkSet, err := GetJwkCerts()
		CacheJwkSet = &jwkSet
		if err != nil {
			return nil, fmt.Errorf("could not get JWK %s", err.Error())
		}
		signingJwk, err := jwkSet.GetSigJwk()
		if err != nil {
			return nil, fmt.Errorf("could not get JWK")
		}
		key, err := signingJwk.ComputePublicRsaKey()
		if err != nil {
			return nil, fmt.Errorf("could not get RSA public key")
		}
		return &key, nil
	}
}

func VerifyToken(token string) (*jwt.Token, error) {
	rsaKey, err := GetRsaPublicKey()
	if err != nil {
		return nil,err
	}
	jwk, err := CacheJwkSet.GetSigJwk()
	if err != nil {
		return nil,err
	}
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
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

	if err != nil {
		return nil,err
	}

	if !parsedToken.Valid {
		return nil,fmt.Errorf("token is not valid")
	}

	return parsedToken,nil
}

func VerifyClaims(token *jwt.Token) error {

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("error obtaining claims")
	}
	
	exp, ok := claims["exp"].(float64)
	if !ok {
		return fmt.Errorf("invalid expiration claim")
	}
	if exp < float64(time.Now().Unix()) {
		return fmt.Errorf("token is expired")
	}

	// // subject
	// if sub, ok := claims["sub"].(string); !ok || sub == "" {
	// 	t.Error("Invalid or missing sub claim")
	// }
	return nil
}
