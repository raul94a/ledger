package api_keycloak

import (
	"crypto/rsa"
	"fmt"
	app_errors "src/errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var CacheJwkSet = &KeycloakJwkSet{Keys: make([]KeycloakJwk, 0)}

func (k KeycloakClient) GetRsaPublicKey() (*rsa.PublicKey, app_errors.AppError) {
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

		jwkSet, err := k.GetJwkCerts()
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

func (k KeycloakClient) getJwtToken(token string) (*jwt.Token, error) {
	rsaKey, err := k.GetRsaPublicKey()
	if err != nil {
		return nil, err
	}
	jwk, err := CacheJwkSet.GetSigJwk()
	if err != nil {
		return nil, err
	}
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
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
}

func (k KeycloakClient) VerifyToken(token string) (*jwt.Token, app_errors.AppError) {

	parsedToken, jwtError := k.getJwtToken(token)

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

func (k KeycloakClient) HasTokenExpired(token string) bool {
	jwtToken, jwtError := k.getJwtToken(token)
	if jwtError != nil {
		fmt.Println("HasTokenExpired() " + jwtError.Error())
		return true
	}
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return true
	}
	exp, ok := claims["exp"].(float64)
	if !ok {
		return true
	}
	return exp < float64(time.Now().Unix())
}

func (k KeycloakClient) VerifyIssuer(token string, preferred_username string) error {
	jwtToken, jwtError := k.getJwtToken(token)
	if jwtError != nil {
		fmt.Println("VerifyIssuer() " + jwtError.Error())
		return jwtError
	}
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return &app_errors.ErrVerifyToken{Message: "error obtaining claims"}
	}
	issuerKey := "preferred_username"
	fmt.Printf("Claims %v\n", claims)
	fmt.Printf("Claimed key: " + claims[issuerKey].(string) + "\n")
	fmt.Printf("Preferred username %s\n", preferred_username)
	issuer, ok := claims[issuerKey].(string)
	sameIssuer := strings.EqualFold(strings.ToLower(issuer), strings.ToLower(preferred_username))
	if !ok || !sameIssuer {
		return fmt.Errorf("not valid credentials")
	}
	return nil

}
