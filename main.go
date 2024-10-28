package main

import (
	"auth-service/internal/cache"
	"auth-service/internal/config"
	"auth-service/internal/delivery/http/route"
	"auth-service/internal/repository"
	"auth-service/internal/usecase"
	"auth-service/internal/utils"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize PostgreSQL connection
	dbInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
	)

	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize Redis connection
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
		Password: cfg.Redis.Password,
	})

	// Initialize cache service
	cacheService := cache.NewRedisCache(redisClient)

	// Initialize repositories
	userRepo := repository.NewCachedUserRepository(db, cacheService)
	redisRepo := repository.NewRedisRepository(redisClient)

	// Initialize email service
	emailService := utils.NewEmailService()

	// Initialize usecase
	authUsecase := usecase.NewAuthUsecase(userRepo, redisRepo, emailService)

	// Initialize Gin router
	gin.SetMode(cfg.App.GinMode)
	router := gin.Default()

	// Setup routes
	route.SetupRoutes(router, authUsecase)

	// Start server
	log.Printf("Server starting on port %s in %s mode", cfg.App.Port, os.Getenv("APP_ENV"))
	if err := router.Run(":" + cfg.App.Port); err != nil {
		log.Fatal(err)
	}
}
