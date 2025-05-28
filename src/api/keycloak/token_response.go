package api_keycloak



// TokenResponse representa la estructura del JSON de respuesta del token
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