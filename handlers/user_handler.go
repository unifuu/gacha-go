package handlers

import (
	"net/http"

	"gacha/models"
	"gacha/services"

	"github.com/gin-gonic/gin"
)

// UserHandler handles user-related requests
type UserHandler struct {
	userService *services.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// HandleGetUserInfo returns user information
func (h *UserHandler) HandleGetUserInfo(c *gin.Context) {
	user := h.userService.GetDefaultUser()

	response := models.UserInfoResponse{
		Username:  user.Username,
		Currency:  user.Currency,
		PityCount: user.PityCount,
	}

	c.JSON(http.StatusOK, response)
}

// HandleGetInventory returns user inventory
func (h *UserHandler) HandleGetInventory(c *gin.Context) {
	user := h.userService.GetDefaultUser()

	response := models.InventoryResponse{
		Inventory: user.Inventory,
		Count:     len(user.Inventory),
	}

	c.JSON(http.StatusOK, response)
}

// HandleAddCurrency adds currency to user (for testing)
func (h *UserHandler) HandleAddCurrency(c *gin.Context) {
	var req models.AddCurrencyRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := h.userService.GetDefaultUser()
	user.AddCurrency(req.Amount)

	response := models.CurrencyResponse{
		Currency: user.Currency,
	}

	c.JSON(http.StatusOK, response)
}
