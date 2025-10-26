package main

import (
	"math/rand"
	"time"

	"gacha/config"
	"gacha/handlers"
	"gacha/routes"
	"gacha/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize services
	userService := services.NewUserService()
	gachaService := services.NewGachaService()

	// Initialize handlers
	gachaHandler := handlers.NewGachaHandler(gachaService, userService)
	userHandler := handlers.NewUserHandler(userService)
	wsHandler := handlers.NewWebSocketHandler(gachaService, userService)

	// Setup Gin router
	r := gin.Default()

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.Server.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Upgrade", "Connection"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowWebSockets:  true,
		MaxAge:           12 * time.Hour,
	}))

	// Setup routes
	routes.SetupRoutes(r, gachaHandler, userHandler, wsHandler)

	// Start server
	r.Run(cfg.Server.Port)
}
