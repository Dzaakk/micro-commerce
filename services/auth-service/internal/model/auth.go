package model

import "time"

type UserRole string

const (
	RoleCustomer UserRole = "customer"
	RoleSeller   UserRole = "seller"
	RoleAdmin    UserRole = "admin"
)

type User struct {
	ID           int64     `db:"id"`
	Email        string    `db:"email"`
	Username     string    `db:"username"`
	PasswordHash string    `db:"password_hash"`
	Role         UserRole  `db:"role"`
	IsActive     bool      `db:"is_active"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`

	CustomerData *CustomerData `db:"-"`
	SellerData   *SellerData   `db:"-"`
}

type CustomerData struct {
	PhoneNumber string
	Gender      int
	DateOfBirth time.Time
}

type SellerData struct {
	StoreName   string
	Address     string
	PhoneNumber string
}
