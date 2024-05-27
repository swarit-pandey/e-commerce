package service

import (
	"context"

	httpapi "github.com/swarit-pandey/e-commerce/user/api/http/server"
)

// UserService interface exposes services layer methods to interact with user and address entity
type UserService interface {
	// CreateUser will create a new user
	CreateUser(ctx context.Context, request *httpapi.UserRegistrationRequest) (*httpapi.UserRegistrationResponse, error)

	// LoginUser will login the user
	LoginUser(ctx context.Context, request *httpapi.UserLoginRequest) (*httpapi.UserLoginResponse, error)

	// InitiatePasswordReset instantiates password reset
	InitiatePasswordReset(ctx context.Context, request *httpapi.PasswordResetRequest) error

	// UpdatePassword will update the password
	UpdatePassword(ctx context.Context, request *httpapi.PasswordUpdateRequest) error

	// GetUserProfile will fetch the user profile
	GetUserProfile(ctx context.Context, userID uint) (*httpapi.UserProfile, error)

	// AddUserProfile will add a user profile usually done when updation is needed
	AddUserProfile(ctx context.Context, requestUser *httpapi.UserProfile, requestAddress *httpapi.Address) error

	// AddUserAddress will add user address
	AddUserAddress(ctx context.Context, userID uint, request *httpapi.Address) error

	// DeleteUserAddress will delete user address
	DeleteUserAddress(ctx context.Context, userID uint, addressID uint) error

	// UpdateUserAddress will update the user address
	UpdateUserAddress(ctx context.Context, userID uint, request *httpapi.Address) error
}
