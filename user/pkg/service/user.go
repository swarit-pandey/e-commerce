package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/mail"
	"net/smtp"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	aerobcrypt "github.com/aerospike/aerospike-client-go/v7/pkg/bcrypt"
	"github.com/golang-jwt/jwt"
	httpapi "github.com/swarit-pandey/e-commerce/user/api/http/server"
	"github.com/swarit-pandey/e-commerce/user/pkg/cache"
	servicecache "github.com/swarit-pandey/e-commerce/user/pkg/cache"
	"github.com/swarit-pandey/e-commerce/user/pkg/repository"
)

var (
	ErrInvalidUsername  = errors.New("username invalid")
	ErrInvalidEmail     = errors.New("email invalid")
	ErrInvalidPassword  = errors.New("password invalid")
	ErrHashingPassword  = errors.New("password hashing failed")
	ErrSettingToCache   = errors.New("caching user failed")
	ErrGettingFromCache = errors.New("error getting a user from cache")
	ErrPasswordMatch    = errors.New("password is incorrect")
	ErrResetToken       = errors.New("resetting token failed")
	ErrEmailSent        = errors.New("sending email failed")
)

type userService struct {
	cacheService *servicecache.CacheService
	inMemCache   *inMemory
}

func NewUserService(cacheService *cache.CacheService) UserService {
	return &userService{
		cacheService: cacheService,
		inMemCache:   newInMemory(),
	}
}

// CreateUser implements `CreateUser()` from UserService interface
func (us *userService) CreateUser(ctx context.Context, request *httpapi.UserRegistrationRequest) (*httpapi.UserRegistrationResponse, error) {
	response := &httpapi.UserRegistrationResponse{}
	err := validate(request.Username, "username")
	if err != nil {
		return response, errors.Join(ErrInvalidUsername, err)
	}

	err = validate(request.Password, "password")
	if err != nil {
		return response, errors.Join(ErrInvalidPassword, err)
	}

	err = validate(request.Email, "email")
	if err != nil {
		return response, errors.Join(ErrInvalidEmail, err)
	}

	passwordhash, err := hashPassword(request.Password)
	if err != nil {
		return response, errors.Join(ErrHashingPassword, err)
	}

	userToCache := &repository.User{
		Name:         *request.Name,
		Username:     request.Username,
		Email:        request.Email,
		PasswordHash: passwordhash,
	}
	err = us.cacheService.SetUser(ctx, userToCache)
	if err != nil {
		return response, errors.Join(ErrSettingToCache, err)
	}

	us.inMemCache.setInMemory(userToCache)
	return response, err
}

// LoginUser implements `LoginUser()` from UserInterface
func (us *userService) LoginUser(ctx context.Context, request *httpapi.UserLoginRequest) (*httpapi.UserLoginResponse, error) {
	response := &httpapi.UserLoginResponse{}
	err := validate(request.Username, "username")
	if err != nil {
		return response, errors.Join(ErrInvalidUsername, err)
	}

	err = validate(request.Password, "password")
	if err != nil {
		return response, errors.Join(ErrInvalidPassword, err)
	}

	user, err := us.cacheService.GetUserByUsername(ctx, request.Username)
	if err != nil {
		return response, errors.Join(ErrGettingFromCache, err)
	}

	if !(aerobcrypt.Match(request.Password, user.PasswordHash)) {
		return response, errors.Join(ErrPasswordMatch, err)
	}

	token, err := generateJWT(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	response.AccessToken = &token
	return response, nil
}

// InitiatePasswordReset implements `InitiatePasswordReset()` from UserService interface
func (us *userService) InitiatePasswordReset(ctx context.Context, request *httpapi.PasswordResetRequest) error {
	err := validate(request.Email, "email")
	if err != nil {
		return errors.Join(ErrInvalidEmail, err)
	}

	var resetPasswordForUser *repository.User
	users := us.inMemCache.getBatchInMem()
	for _, u := range users {
		if request.Email == u.Email {
			resetPasswordForUser = u
			break
		}
	}

	token, err := generateResetToken(resetPasswordForUser)
	if err != nil {
		return errors.Join(ErrResetToken, err)
	}

	err = us.cacheService.SetPasswordResetToken(ctx, resetPasswordForUser.ID, token)
	if err != nil {
		return errors.Join(ErrResetToken, ErrSettingToCache, err)
	}

	err = sendResetToken(resetPasswordForUser.Email, token)
	if err != nil {
		return errors.Join(ErrResetToken, ErrEmailSent, err)
	}

	return nil
}

// UpdatePassword implements `UpdatePassword()` from UserService interface
func (us *userService) UpdatePassword(ctx context.Context, request *httpapi.PasswordUpdateRequest) error {
	err := validate(*request.NewPassword, "password")
	if err != nil {
		return errors.Join(ErrInvalidPassword, err)
	}

	userID, token := parseResetToken(request.Token)
	if userID == 0 || token == "" {
		return ErrResetToken
	}

	users := us.inMemCache.getBatchInMem()
	var user *repository.User
	for _, u := range users {
		if u.ID == userID {
			user = u
			break
		}
	}
	if user == nil {
		return errors.New("user not found")
	}

	cachedToken, err := generateResetToken(user)
	if err != nil || cachedToken != token {
		return ErrResetToken
	}

	newPasswordHash, err := hashPassword(*request.NewPassword)
	if err != nil {
		return errors.Join(ErrHashingPassword, err)
	}

	user.PasswordHash = newPasswordHash
	err = us.cacheService.SetUser(ctx, user)
	if err != nil {
		return errors.Join(ErrSettingToCache, err)
	}

	us.inMemCache.setInMemory(user)
	return nil
}

// GetUserProfile implements `GetUserProfile()` from UserService interface
func (us *userService) GetUserProfile(ctx context.Context, userID uint) (*httpapi.UserProfile, error) {
	response := &httpapi.UserProfile{}

	user, err := us.cacheService.GetUser(ctx, userID)
	if err != nil {
		return response, errors.Join(ErrGettingFromCache, err)
	}

	response.Username = &user.Username
	response.Email = &user.Email
	response.Name = &user.Name

	return response, nil
}

// AddUserProfile implements `AddUserProfile()` from UserService interface
func (us *userService) AddUserProfile(ctx context.Context, requestUser *httpapi.UserProfile, requestAddress *httpapi.Address) error {
	user, err := us.cacheService.GetUserByUsername(ctx, *requestUser.Username)
	if err != nil {
		return errors.Join(ErrGettingFromCache, err)
	}

	if requestUser.Email != nil {
		user.Email = *requestUser.Email
	}
	if requestUser.Name != nil {
		user.Name = *requestUser.Name
	}

	err = us.cacheService.SetUser(ctx, user)
	if err != nil {
		return errors.Join(ErrSettingToCache, err)
	}

	us.inMemCache.setInMemory(user)
	return nil
}

// AddUserAddress implements `AddUserAddress()` from UserService interface
func (us *userService) AddUserAddress(ctx context.Context, userID uint, request *httpapi.Address) error {
	user, err := us.cacheService.GetUser(ctx, userID)
	if err != nil {
		return errors.Join(ErrGettingFromCache, err)
	}

	address := &repository.UserAddress{
		AddressLine1: *request.AddressLine1,
		AddressLine2: *request.AddressLine2,
		City:         *request.City,
		State:        *request.State,
		Country:      *request.Country,
		PostalCode:   *request.Code,
	}

	user.Addresses = append(user.Addresses, *address)

	err = us.cacheService.SetUser(ctx, user)
	if err != nil {
		return errors.Join(ErrSettingToCache, err)
	}

	us.inMemCache.setInMemory(user)
	return nil
}

// DeleteUserAddress implements `DeleteUserAddress()` from UserService interface
func (us *userService) DeleteUserAddress(ctx context.Context, userID uint, addressID uint) error {
	user, err := us.cacheService.GetUser(ctx, userID)
	if err != nil {
		return errors.Join(ErrGettingFromCache, err)
	}

	newAddresses := make([]repository.UserAddress, 0)
	for _, address := range user.Addresses {
		if address.ID != addressID {
			newAddresses = append(newAddresses, address)
		}
	}
	user.Addresses = newAddresses

	err = us.cacheService.SetUser(ctx, user)
	if err != nil {
		return errors.Join(ErrSettingToCache, err)
	}

	us.inMemCache.setInMemory(user)
	return nil
}

// UpdateUserAddress implements `UpdateUserAddress()` from UserService interface
func (us *userService) UpdateUserAddress(ctx context.Context, userID uint, request *httpapi.Address) error {
	user, err := us.cacheService.GetUser(ctx, userID)
	if err != nil {
		return errors.Join(ErrGettingFromCache, err)
	}

	for _, address := range user.Addresses {
		if address.ID == uint(*request.Id) {
			if request.AddressLine1 != nil {
				address.AddressLine1 = *request.AddressLine1
			}
			if request.AddressLine2 != nil {
				address.AddressLine2 = *request.AddressLine2
			}
			if request.City != nil {
				address.City = *request.City
			}
			if request.State != nil {
				address.State = *request.State
			}
			if request.Country != nil {
				address.Country = *request.Country
			}
			if request.Code != nil {
				address.PostalCode = *request.Code
			}
		}
	}

	err = us.cacheService.SetUser(ctx, user)
	if err != nil {
		return errors.Join(ErrSettingToCache, err)
	}

	us.inMemCache.setInMemory(user)
	return nil
}

func validate(payload, payloadType string) error {
	switch payloadType {
	case "password":
		return validatePassword(payload)
	case "email":
		return validateEmail(payload)
	case "username":
		return validateUsername(payload)
	default:
		return errors.New("invalid type for validation")
	}
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be of length 8")
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char):
			hasSpecial = true
		}
	}

	if !(hasUpper && hasLower && hasDigit && hasSpecial) {
		return errors.New("password needs to have upper, lower, digit and special characters")
	}
	return nil
}

func validateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return errors.New("email is not in a valid form")
	}
	return nil
}

func validateUsername(username string) error {
	if len(username) < 4 || len(username) > 20 {
		return errors.New("username should be not be less than 4 and not greater than 20 in length")
	}

	// Username should only contain alphanumeric characters, underscores, and dots
	// Underscore and dot cannot be at the start or end
	// Underscore and dot cannot be next to each other
	re := regexp.MustCompile(`^[a-zA-Z0-9]+([._]?[a-zA-Z0-9]+)*$`)
	if !re.MatchString(username) {
		return errors.New("username can only contain alphanumeric characters, underscores, and dots")
	}

	return nil
}

func hashPassword(password string) (string, error) {
	hashed, err := aerobcrypt.Hash(password, "encode")
	if err != nil {
		return "", errors.New("failed to hash password")
	}
	return hashed, nil
}

func generateJWT(userID uint, username string) (string, error) {
	expiration := time.Now().Add(1 * time.Hour)

	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      expiration.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey := []byte("key")
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func generateResetToken(user *repository.User) (string, error) {
	randBytes := make([]byte, 32)
	_, err := rand.Read(randBytes)
	if err != nil {
		return "", err
	}

	randString := base64.URLEncoding.EncodeToString(randBytes)
	resetToken := fmt.Sprintf("%v:%v", user.ID, randString)
	return resetToken, nil
}

// sendResetToken might not really send the email, its just
// sort of "dummy"
func sendResetToken(email, token string) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	smtpUsername := "ecomm@gmail.com"
	smtpPassword := "password"

	from := "no-reply@gmail.com"
	to := []string{email}

	subject := "Password Reset Token"
	body := fmt.Sprintf("Your password reset token is: %s", token)
	message := []byte(fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body))

	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

func parseResetToken(token string) (uint, string) {
	parts := strings.Split(token, ":")
	if len(parts) != 2 {
		return 0, ""
	}

	userID, err := strconv.ParseUint(parts[0], 10, 32)
	if err != nil {
		return 0, ""
	}

	return uint(userID), parts[1]
}
