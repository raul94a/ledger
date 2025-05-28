package api_keycloak

import (
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"math/big"
	app_errors "src/errors"
)

/*
	RFC7517: JSON Web Key (JWK)

As said in the RFC 7517, there can exist a set of JWK,
each one for a certain purpose.

Keycloack offers a set of two keys to:
  - Verify the JWT `use: sig`.
  - Encrypt data `use: enc`
*/
type KeycloakJwkSet struct {
	Keys []KeycloakJwk `json:"keys"`
}

/*
KeyCloakJWK JSON attributes. The RFC identifies more attributes that can be used in other scenarios.
Keycloak makes use of the folling:

Kid: It's the ID of the key and it's used to select a key among others.

Kty: Key type. Two possible values as for RFC7517: RSA or EC

Alg: The algorithm used for the encryption. https://www.iana.org/assignments/jose/jose.xhtml#web-signature-encryption-algorithms

Use: The intended use of the public key. Two possible values: enc or sig

N: The modulus of the public key

E: The public Exponent

X5c: A list of PKIX Certificates (RFC5280). The first element of the list MUST be

	the one containing the key. This first certicate may be followed by a more certs,
	with each subsequent certificate being the one used to certify the previous one.

X5t: Base64url encoded SHA-1 thumprint —digest— of the DER encoding of an X.509 Cert.

X5t#S256: Base64url encoded SHA-256 thumprint of the DER encoding of a X.509 Cert.
*/
type KeycloakJwk struct {
	Kid     string   `json:"kid"`
	Kty     string   `json:"kty"`
	Alg     string   `json:"alg"`
	Use     string   `json:"use"`
	N       string   `json:"n"`
	E       string   `json:"e"`
	X5c     []string `json:"x5c"`
	X5t     string   `json:"x5t"`
	X5tS256 string   `json:"x5t#S256"`
}

func (r *KeycloakJwkSet) GetSigJwk() (KeycloakJwk, app_errors.AppError) {
	const signatureUse string = "sig"

	for _, element := range r.Keys {
		if element.Use == signatureUse {
			return element, nil
		}
	}
	return KeycloakJwk{}, &app_errors.ErrNotJwkFound{}
}

func (r *KeycloakJwk) ComputePublicRsaKey() (rsa.PublicKey, app_errors.AppError) {

	nBytes, err := base64.RawURLEncoding.DecodeString(r.N)
	if err != nil {
		return rsa.PublicKey{}, &app_errors.ErrRsaPublicKey{Message: fmt.Sprintf("failed to decode RSA modulus %s . %s", r.N, err.Error())}

	}

	n := new(big.Int).SetBytes(nBytes)

	eBytes, err := base64.RawURLEncoding.DecodeString(r.E)
	if err != nil {
		return rsa.PublicKey{}, &app_errors.ErrRsaPublicKey{Message: fmt.Sprintf("failed to decode RSA exponent %s . %s", r.E, err.Error())}
	}
	e := new(big.Int).SetBytes(eBytes)

	rsaKey := &rsa.PublicKey{
		N: n,
		E: int(e.Int64()),
	}
	return *rsaKey, nil
}




// Keycloack TokenResponse
type TokenResponse struct {
    AccessToken      string    `json:"access_token"`
    ExpiresIn        int       `json:"expires_in"`
    RefreshExpiresIn int       `json:"refresh_expires_in"`
    RefreshToken     string    `json:"refresh_token"`
    TokenType        string    `json:"token_type"`
    NotBeforePolicy  int64     `json:"not-before-policy"`
    SessionState     string    `json:"session_state"`
    Scope            string    `json:"scope"`
}