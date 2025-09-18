package dto

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role" validate:"required,oneof=customer seller"`

	CustomerInfo *CustomerInfo `json:"customer_info" validate:"required_if=Role customer"`
	SellerInfo   *SellerInfo   `json:"seller_info" validate:"required_if=Role seller"`
}

type CustomerInfo struct {
	PhoneNumber string `json:"phone_number" validate:"required,e164"`
	DateOfBirth string `json:"date_of_birth" validate:"required,datetime=02-01-2006"`
	Gender      string `json:"gender" validate:"required, oneof=male female"`
}

type SellerInfo struct {
	PhoneNumber string `json:"phone_number" validate:"required,e164"`
	StoreName   string `json:"store_name" validate:"required,min=3,max=100"`
	Address     string `json:"address" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int64     `json:"expires_in"`
	User         BasicUser `json:"user"`
}

type BasicUser struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Role     string `json:"role"`
}
