package api_keycloak

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func AuthAdminUser() (TokenResponse, error) {
	host := os.Getenv("KC_HOST")
	port := os.Getenv("KD_PORT")
	authEndpoint := os.Getenv("GET_TOKEN_URL")
	url := fmt.Sprintf("%s:%s%s", host, port, authEndpoint)

	adminUser := os.Getenv("ADMIN_USER")
	password := os.Getenv("ADMIN_PASSWORD")
	clientId := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	grantType := "password"
	bodyStr := fmt.Sprintf("grant_type=%s&username=%s&password=%s&client_id=%s&client_secret=%s", grantType, adminUser, password, clientId, clientSecret)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(bodyStr)))
	if err != nil {
		return TokenResponse{}, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return TokenResponse{}, err
	}
	defer res.Body.Close()
	// variable with the response json
	tokenResponse := &TokenResponse{}
	derr := json.NewDecoder(res.Body).Decode(tokenResponse)
	if derr != nil {

		return TokenResponse{}, err
	}

	return *tokenResponse, nil
}

func GetJwkCerts() (KeycloakJwkSet, error) {
	host := os.Getenv("KC_HOST")
	port := os.Getenv("KC_PORT")
	jwkEndpoint := os.Getenv("JWK_URL")
	url := fmt.Sprintf("%s:%s%s", host, port, jwkEndpoint)
	fmt.Println("URL ", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return KeycloakJwkSet{}, err
	}

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return KeycloakJwkSet{}, err
	}
	if res.StatusCode != http.StatusOK {
		return KeycloakJwkSet{}, fmt.Errorf("error fetching JWK")
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if err != nil {
		return KeycloakJwkSet{}, err
	}
	// variable with the response json
	var jwkSet KeycloakJwkSet
	if err := json.Unmarshal(body, &jwkSet); err != nil {
		// Optionally log the raw JSON for debugging
		fmt.Printf("Failed to unmarshal JSON: %s\n", string(body))
		return KeycloakJwkSet{}, err
	}



	return jwkSet, nil
}


