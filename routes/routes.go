package routes

import (
	"gacha/handlers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all API routes
func SetupRoutes(r *gin.Engine, gachaHandler *handlers.GachaHandler, userHandler *handlers.UserHandler, wsHandler *handlers.WebSocketHandler) {
	// WebSocket endpoint
	r.GET("/ws", wsHandler.HandleWebSocket)

	// HTTP API endpoints (kept for backward compatibility)
	api := r.Group("/api")
	{
		// Gacha routes
		gacha := api.Group("/gacha")
		{
			gacha.POST("/pull", gachaHandler.HandleSinglePull)
			gacha.POST("/pull-ten", gachaHandler.HandleTenPull)
			gacha.GET("/pool", gachaHandler.HandleGetPool)
		}

		// User routes
		user := api.Group("/user")
		{
			user.GET("/info", userHandler.HandleGetUserInfo)
			user.GET("/inventory", userHandler.HandleGetInventory)
			user.POST("/add-currency", userHandler.HandleAddCurrency)
		}
	}
}
