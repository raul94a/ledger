package api_keycloak

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	app_errors "src/errors"
)

type AdminCredentials struct {
	Username string
	Password string
}

type KcEndpoints struct {
	AuthorizationEndpoint string
	JwkEndpoint           string
	CreateUserEndpoint    string
}

type KeycloakClient struct {
	Url              string
	ClientId         string
	Secret           string
	AdminCredentials AdminCredentials
	KcEndpoints      KcEndpoints
}

// The method is supposed to be used after the .env is loaded
func BuildKeycloakClientFromEnv() KeycloakClient {
	host := os.Getenv("KC_HOST")
	adminUser := os.Getenv("ADMIN_USER")
	password := os.Getenv("ADMIN_PASSWORD")
	clientId := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	authEndpoint := os.Getenv("GET_TOKEN_URL")
	jwkEndpoint := os.Getenv("JWK_URL")
	createUserEndpoint := os.Getenv("CREATE_USER_URL")
	adminCreds := AdminCredentials{
		Username: adminUser,
		Password: password,
	}
	kcEndpoints := KcEndpoints{
		AuthorizationEndpoint: authEndpoint,
		JwkEndpoint:           jwkEndpoint,
		CreateUserEndpoint:    createUserEndpoint,
	}
	client := KeycloakClient{
		Url:              host,
		ClientId:         clientId,
		Secret:           clientSecret,
		AdminCredentials: adminCreds,
		KcEndpoints:      kcEndpoints,
	}
	return client
}

func (k KeycloakClient) AuthAdminUser() (TokenResponse, app_errors.AppError) {

	url := fmt.Sprintf("%s%s", k.Url, k.KcEndpoints.AuthorizationEndpoint)

	adminUser := k.AdminCredentials.Username
	password := k.AdminCredentials.Password
	clientId := k.ClientId
	clientSecret := k.Secret
	grantType := "password"
	bodyStr := fmt.Sprintf("grant_type=%s&username=%s&password=%s&client_id=%s&client_secret=%s", grantType, adminUser, password, clientId, clientSecret)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(bodyStr)))
	if err != nil {
		return TokenResponse{}, &app_errors.ErrBadRequest{}
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return TokenResponse{}, &app_errors.ErrBadRequest{}
	}
	defer res.Body.Close()
	// variable with the response json
	tokenResponse := &TokenResponse{}
	derr := json.NewDecoder(res.Body).Decode(tokenResponse)
	if derr != nil {

		return TokenResponse{}, &app_errors.ErrBadRequest{Message: derr.Error()}
	}

	return *tokenResponse, nil
}

func (k KeycloakClient) AuthUser(username string, password string) (TokenResponse, app_errors.AppError) {

	url := fmt.Sprintf("%s%s", k.Url, k.KcEndpoints.AuthorizationEndpoint)
	clientId := k.ClientId
	clientSecret := k.Secret
	grantType := "password"
	bodyStr := fmt.Sprintf("grant_type=%s&username=%s&password=%s&client_id=%s&client_secret=%s", grantType, username, password, clientId, clientSecret)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(bodyStr)))
	if err != nil {
		return TokenResponse{}, &app_errors.ErrBadRequest{}
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return TokenResponse{}, &app_errors.ErrBadRequest{}
	}
	statusCode := res.StatusCode
	if statusCode != 200 {
		return TokenResponse{}, &app_errors.ErrBadRequest{Message: "Bad Request"}
	}
	defer res.Body.Close()
	// variable with the response json
	tokenResponse := &TokenResponse{}

	derr := json.NewDecoder(res.Body).Decode(tokenResponse)
	if derr != nil {

		return TokenResponse{}, &app_errors.ErrBadRequest{Message: derr.Error()}
	}

	return *tokenResponse, nil
}

func (k KeycloakClient) GetJwkCerts() (KeycloakJwkSet, app_errors.AppError) {

	jwkEndpoint := k.KcEndpoints.JwkEndpoint
	url := fmt.Sprintf("%s%s", k.Url, jwkEndpoint)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return KeycloakJwkSet{}, &app_errors.ErrBadRequest{}
	}

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return KeycloakJwkSet{}, &app_errors.ErrBadRequest{}
	}
	if res.StatusCode != http.StatusOK {
		return KeycloakJwkSet{}, &app_errors.ErrBadRequest{}
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if err != nil {
		return KeycloakJwkSet{}, &app_errors.ErrBadRequest{Message: err.Error()}
	}
	// variable with the response json
	var jwkSet KeycloakJwkSet
	if err := json.Unmarshal(body, &jwkSet); err != nil {
		// Optionally log the raw JSON for debugging
		fmt.Printf("Failed to unmarshal JSON: %s\n", string(body))
		return KeycloakJwkSet{}, &app_errors.ErrBadRequest{Message: err.Error()}
	}

	return jwkSet, nil
}

func (k KeycloakClient) CreateUser(createUserReq KcCreateUserRequest) app_errors.AppError {

	url := fmt.Sprintf("%s%s", k.Url, k.KcEndpoints.CreateUserEndpoint)
	jsonBytes, err := json.Marshal(createUserReq)
	if err != nil {
		return &app_errors.ErrBadRequest{Reason: err, Message: err.Error()}
	}
	token, err := k.AuthAdminUser()
	if err != nil {
		return &app_errors.ErrInternalServer{Reason: err, Message: err.Error()}
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return &app_errors.ErrInternalServer{Reason: err}
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	client := http.Client{}
	res, err := client.Do(req)


	if err != nil {
		return &app_errors.ErrBadRequest{Reason: err, Message: err.Error()}
	}
	statusCode := res.StatusCode
	allowedStatusCode := []int{201, 409}
	if slices.Contains(allowedStatusCode, statusCode) {
		return nil
	}
	return &app_errors.ErrInternalServer{Message: fmt.Sprintf("Request finished with status code %d", statusCode)}

}


func (k KeycloakClient) IsAdminUser(user, password string) bool {
	creds := k.AdminCredentials
	adminUser := creds.Username
	adminPassword := creds.Password
	return (user == adminUser) && (password == adminPassword)
}