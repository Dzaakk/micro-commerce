package model

import (
	"time"

	"github.com/lib/pq"
)

type UserRole string

const (
	RoleCustomer UserRole = "customer"
	RoleSeller   UserRole = "seller"
	RoleAdmin    UserRole = "admin"
)

type User struct {
	ID           string    `db:"id"`
	Email        string    `db:"email"`
	Username     string    `db:"username"`
	PasswordHash string    `db:"password_hash"`
	Role         UserRole  `db:"role"`
	IsActive     bool      `db:"is_active"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type OAuthClient struct {
	ID           string         `db:"id"`
	ClientID     string         `db:"client_id"`
	ClientSecret string         `db:"client_secret"`
	Name         string         `db:"name"`
	RedirectURIs pq.StringArray `db:"redirect_uris"`
	GrantTypes   pq.StringArray `db:"grant_types"`
	Scope        string         `db:"scope"`
	CreatedAt    time.Time      `db:"created_at"`
}

type AuthorizationCode struct {
	Code        string    `db:"code"`
	ClientID    string    `db:"client_id"`
	UserID      string    `db:"user_id"`
	RedirectURI string    `db:"redirect_uri"`
	Scope       string    `db:"scope"`
	ExpiresAt   time.Time `db:"expires_at"`
	CreatedAt   time.Time `db:"created_at"`
}

type AccessToken struct {
	ID        string    `db:"id"`
	Token     string    `db:"token"`
	ClientID  string    `db:"client_id"`
	UserID    string    `db:"user_id"`
	Scope     string    `db:"scope"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
}

type RefreshToken struct {
	ID            string    `db:"id"`
	Token         string    `db:"token"`
	ClientID      string    `db:"client_id"`
	UserID        string    `db:"user_id"`
	AccessTokenID string    `db:"access_token_id"`
	ExpiresAt     time.Time `db:"expires_at"`
	CreatedAt     time.Time `db:"created_at"`
}
