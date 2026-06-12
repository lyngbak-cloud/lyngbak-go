package miget

import (
	"context"
	"net/http"
)

// AuthService handles the Authentication endpoints (OpenAPI tag
// "Authentication").
type AuthService struct {
	client *Client
}

// SignIn exchanges email/password credentials for a pair of JWT tokens.
// Operation: signIn (POST /api/v1/auth/sign_in)
func (s *AuthService) SignIn(ctx context.Context, email, password string) (*AuthTokens, error) {
	body := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{Email: email, Password: password}
	out := new(AuthTokens)
	if err := s.client.do(ctx, http.MethodPost, "/api/v1/auth/sign_in", nil, body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// RefreshToken exchanges a refresh token for a new access token.
// Operation: refreshToken (POST /api/v1/auth/refresh_token)
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*AccessToken, error) {
	body := struct {
		RefreshToken string `json:"refresh_token"`
	}{RefreshToken: refreshToken}
	out := new(AccessToken)
	if err := s.client.do(ctx, http.MethodPost, "/api/v1/auth/refresh_token", nil, body, out); err != nil {
		return nil, err
	}
	return out, nil
}
