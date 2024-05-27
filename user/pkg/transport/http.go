package transport

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	httpapi "github.com/swarit-pandey/e-commerce/user/api/http/server"
	"github.com/swarit-pandey/e-commerce/user/pkg/service"
)

type Handler struct {
	userService service.UserService
}

func NewHandler(userService service.UserService) *Handler {
	return &Handler{
		userService: userService,
	}
}

// PostUsers creates a new user
func (h *Handler) PostUsers(c *gin.Context) {
	var req httpapi.UserRegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.userService.CreateUser(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) PostUsersLogin(c *gin.Context) {
	var req httpapi.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.userService.LoginUser(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// PostUsersPasswordReset initiates a password reset
func (h *Handler) PostUsersPasswordReset(c *gin.Context) {
	var req httpapi.PasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.userService.InitiatePasswordReset(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset initiated"})
}

// PutUsersPasswordReset updates the user's password
func (h *Handler) PutUsersPasswordReset(c *gin.Context) {
	var req httpapi.PasswordUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.userService.UpdatePassword(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password updated"})
}

// GetUsersUserId gets a user profile
func (h *Handler) GetUsersUserId(c *gin.Context, userID int) {
	resp, err := h.userService.GetUserProfile(context.Background(), uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// PostUsersUserId adds a user profile
func (h *Handler) PostUsersUserId(c *gin.Context, userID int) {
	var reqProfile httpapi.UserProfile
	var reqAddress httpapi.Address
	if err := c.ShouldBindJSON(&reqProfile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&reqAddress); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.userService.AddUserProfile(context.Background(), &reqProfile, &reqAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User profile added"})
}

// DeleteUsersUserIdAddressesAddressId deletes a user address
func (h *Handler) DeleteUsersUserIdAddressesAddressId(c *gin.Context, userID int, addressID int) {
	err := h.userService.DeleteUserAddress(context.Background(), uint(userID), uint(addressID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User address deleted"})
}

// PutUsersUserIdAddressesAddressId updates a user address
func (h *Handler) PutUsersUserIdAddressesAddressId(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	_, err = strconv.Atoi(c.Param("addressId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID"})
		return
	}

	var req httpapi.Address
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.userService.UpdateUserAddress(context.Background(), uint(userId), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User address updated"})
}
