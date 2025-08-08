package main

import (
	db "herp/db/sqlc"
	_ "herp/docs/swagger" // Import generated swagger docs
	"herp/internal/auth"
	"herp/internal/config"
	"herp/internal/docs"
	"herp/internal/pos"
	"herp/internal/server"
	"herp/pkg/database"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title Hotel ERP API
// @version 1.0.0
// @description This is the Hotel ERP API server. It provides endpoints for managing hotel operations including authentication, point of sale, inventory, and more.
// @termsOfService http://swagger.io/terms/

// @contact.name Hotel ERP API Support
// @contact.email support@hotel-erp.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:9000
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT Authorization header using the Bearer scheme. Example: "Authorization: Bearer {token}"

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
	}

	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Load database
	dbs, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize sqlc
	queries := db.New(dbs)

	// Initialiaze services
	authSvc := auth.NewService(queries, cfg.JWTSecret, time.Duration(cfg.JWTExpiry)*time.Hour)

	r := gin.Default()

	// Setup API documentation
	docsConfig := docs.DefaultSwaggerConfig()
	docsConfig.Host = "localhost:" + cfg.Port
	docsConfig.Enabled = true

	// Add CORS for docs
	r.Use(docs.CORSForDocs())

	// Add API docs middleware
	r.Use(docs.APIDocsMiddleware())

	// Setup Swagger documentation
	docs.SetupSwagger(r, docsConfig)

	// Setup Redocly documentation (alternative)
	docs.SetupRedocly(r, docsConfig)

	// register routes
	v1 := r.Group("/api/v1")
	adminHandler := auth.NewAdminHandler(authSvc)
	adminHandler.RegisterAdminRoutes(v1, authSvc)
	authHandler := auth.NewHandler(authSvc)
	v1.POST("/auth/login", authHandler.Login)
	pos.RegisterRoutes(v1, authSvc)

	// Create server with graceful shutdown
	serverConfig := server.Config{
		Port:            cfg.Port,
		ReadTimeout:     15 * time.Second,
		WriteTimeout:    15 * time.Second,
		ShutdownTimeout: 30 * time.Second,
	}

	srv := server.New(r, dbs, serverConfig)

	// Add health check endpoint
	// @Summary Health check
	// @Description Check the health status of the API server
	// @Tags health
	// @Produce json
	// @Success 200 {object} map[string]string "Service is healthy"
	// @Failure 500 {object} map[string]string "Service is unhealthy"
	// @Router /health [get]
	r.GET("/health", func(c *gin.Context) {
		if err := srv.Health(); err != nil {
			c.JSON(500, gin.H{"status": "unhealthy", "error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// Start server with graceful shutdown
	log.Printf("Starting Hotel ERP server version %s...", cfg.ApiVersion)
	if err := srv.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
