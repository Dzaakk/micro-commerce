package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Dzaakk/micro-commerce/services/auth-service/internal/model"
	"github.com/lib/pq"
)

type authRepository struct {
	db *sql.DB
}

func NewauthRepository(db *sql.DB) *authRepository {
	return &authRepository{db: db}
}

// OAuth Clients
func (r *authRepository) GetClientByID(ctx context.Context, clientID string) (*model.OAuthClient, error) {
	query := `
		SELECT id, client_id, client_secret, name, redirect_uris, grant_types, scope, created_at
		FROM oauth_clients
		WHERE client_id = $1
	`

	var client model.OAuthClient
	err := r.db.QueryRowContext(ctx, query, clientID).Scan(
		&client.ID,
		&client.ClientID,
		&client.ClientSecret,
		&client.Name,
		pq.Array(&client.RedirectURIs),
		pq.Array(&client.GrantTypes),
		&client.Scope,
		&client.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("client not found")
	}
	if err != nil {
		return nil, err
	}

	return &client, nil
}

func (r *authRepository) ValidateClientCredentials(ctx context.Context, clientID, clientSecret string) (*model.OAuthClient, error) {
	client, err := r.GetClientByID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	if client.ClientSecret != clientSecret {
		return nil, fmt.Errorf("invalid client credentials")
	}

	return client, nil
}

// Authorization Codes
func (r *authRepository) CreateAuthorizationCode(code *model.AuthorizationCode) error {
	query := `
		INSERT INTO authorization_codes (code, client_id, user_id, redirect_uri, scope, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(query,
		code.Code,
		code.ClientID,
		code.UserID,
		code.RedirectURI,
		code.Scope,
		code.ExpiresAt,
	)

	return err
}

func (r *authRepository) GetAuthorizationCode(code string) (*model.AuthorizationCode, error) {
	query := `
		SELECT code, client_id, user_id, redirect_uri, scope, expires_at, created_at
		FROM authorization_codes
		WHERE code = $1 AND expires_at > NOW()
	`

	var authCode model.AuthorizationCode
	err := r.db.QueryRow(query, code).Scan(
		&authCode.Code,
		&authCode.ClientID,
		&authCode.UserID,
		&authCode.RedirectURI,
		&authCode.Scope,
		&authCode.ExpiresAt,
		&authCode.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("authorization code not found or expired")
	}
	if err != nil {
		return nil, err
	}

	return &authCode, nil
}

func (r *authRepository) DeleteAuthorizationCode(code string) error {
	query := `DELETE FROM authorization_codes WHERE code = $1`
	_, err := r.db.Exec(query, code)
	return err
}

// Access Tokens
func (r *authRepository) CreateAccessToken(token *model.AccessToken) error {
	query := `
		INSERT INTO access_tokens (id, token, client_id, user_id, scope, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(query,
		token.ID,
		token.Token,
		token.ClientID,
		token.UserID,
		token.Scope,
		token.ExpiresAt,
	)

	return err
}

func (r *authRepository) GetAccessToken(token string) (*model.AccessToken, error) {
	query := `
		SELECT id, token, client_id, user_id, scope, expires_at, created_at
		FROM access_tokens
		WHERE token = $1 AND expires_at > NOW()
	`

	var accessToken model.AccessToken
	err := r.db.QueryRow(query, token).Scan(
		&accessToken.ID,
		&accessToken.Token,
		&accessToken.ClientID,
		&accessToken.UserID,
		&accessToken.Scope,
		&accessToken.ExpiresAt,
		&accessToken.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("access token not found or expired")
	}
	if err != nil {
		return nil, err
	}

	return &accessToken, nil
}

func (r *authRepository) DeleteAccessToken(token string) error {
	query := `DELETE FROM access_tokens WHERE token = $1`
	_, err := r.db.Exec(query, token)
	return err
}

func (r *authRepository) DeleteExpiredAccessTokens() error {
	query := `DELETE FROM access_tokens WHERE expires_at < NOW()`
	_, err := r.db.Exec(query)
	return err
}

// Refresh Tokens
func (r *authRepository) CreateRefreshToken(token *model.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (id, token, client_id, user_id, access_token_id, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(query,
		token.ID,
		token.Token,
		token.ClientID,
		token.UserID,
		token.AccessTokenID,
		token.ExpiresAt,
	)

	return err
}

func (r *authRepository) GetRefreshToken(token string) (*model.RefreshToken, error) {
	query := `
		SELECT id, token, client_id, user_id, access_token_id, expires_at, created_at
		FROM refresh_tokens
		WHERE token = $1 AND expires_at > NOW()
	`

	var refreshToken model.RefreshToken
	err := r.db.QueryRow(query, token).Scan(
		&refreshToken.ID,
		&refreshToken.Token,
		&refreshToken.ClientID,
		&refreshToken.UserID,
		&refreshToken.AccessTokenID,
		&refreshToken.ExpiresAt,
		&refreshToken.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("refresh token not found or expired")
	}
	if err != nil {
		return nil, err
	}

	return &refreshToken, nil
}

func (r *authRepository) DeleteRefreshToken(token string) error {
	query := `DELETE FROM refresh_tokens WHERE token = $1`
	_, err := r.db.Exec(query, token)
	return err
}

func (r *authRepository) DeleteRefreshTokenByAccessTokenID(accessTokenID string) error {
	query := `DELETE FROM refresh_tokens WHERE access_token_id = $1`
	_, err := r.db.Exec(query, accessTokenID)
	return err
}

func (r *authRepository) CleanupExpiredTokens() error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete expired authorization codes
	_, err = tx.Exec(`DELETE FROM authorization_codes WHERE expires_at < NOW() - INTERVAL '1 day'`)
	if err != nil {
		return err
	}

	// Delete expired access tokens
	_, err = tx.Exec(`DELETE FROM access_tokens WHERE expires_at < NOW() - INTERVAL '1 day'`)
	if err != nil {
		return err
	}

	// Delete expired refresh tokens
	_, err = tx.Exec(`DELETE FROM refresh_tokens WHERE expires_at < NOW() - INTERVAL '1 day'`)
	if err != nil {
		return err
	}

	return tx.Commit()
}
