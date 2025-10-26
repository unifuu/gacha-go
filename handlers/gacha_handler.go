package handlers

import (
	"net/http"
	"time"

	"gacha/models"
	"gacha/services"

	"github.com/gin-gonic/gin"
)

// GachaHandler handles gacha-related requests
type GachaHandler struct {
	gachaService *services.GachaService
	userService  *services.UserService
}

// NewGachaHandler creates a new gacha handler
func NewGachaHandler(gachaService *services.GachaService, userService *services.UserService) *GachaHandler {
	return &GachaHandler{
		gachaService: gachaService,
		userService:  userService,
	}
}

// HandleSinglePull handles single pull request
func (h *GachaHandler) HandleSinglePull(c *gin.Context) {
	user := h.userService.GetDefaultUser()

	if !user.DeductCurrency(160) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient currency"})
		return
	}

	char := h.gachaService.PerformSinglePull(user)
	isNew := user.AddCharacter(char)

	result := models.GachaResult{
		Characters: []models.Character{char},
		IsNew:      []bool{isNew},
		Timestamp:  time.Now().Unix(),
	}

	c.JSON(http.StatusOK, result)
}

// HandleTenPull handles ten pull request
func (h *GachaHandler) HandleTenPull(c *gin.Context) {
	user := h.userService.GetDefaultUser()

	if !user.DeductCurrency(1600) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient currency"})
		return
	}

	characters := h.gachaService.PerformTenPull(user)
	var isNewList []bool

	for _, char := range characters {
		isNew := user.AddCharacter(char)
		isNewList = append(isNewList, isNew)
	}

	result := models.GachaResult{
		Characters: characters,
		IsNew:      isNewList,
		Timestamp:  time.Now().Unix(),
	}

	c.JSON(http.StatusOK, result)
}

// HandleGetPool returns the gacha pool information
func (h *GachaHandler) HandleGetPool(c *gin.Context) {
	poolInfo := models.PoolInfo{
		Characters: h.gachaService.GetCharacterPool(),
		Rates: map[string]string{
			"ssr": "2%",
			"sr":  "10%",
			"r":   "88%",
		},
		PitySystem: "Guaranteed SSR at 90 pulls",
	}

	c.JSON(http.StatusOK, poolInfo)
}
