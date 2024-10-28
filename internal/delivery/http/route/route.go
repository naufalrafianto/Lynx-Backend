package route

import (
	"auth-service/internal/delivery/http/handler"
	"auth-service/internal/delivery/http/middleware"
	"auth-service/internal/domain"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, authUsecase domain.UserUsecase) {
	// Create handler
	authHandler := handler.NewAuthHandler(authUsecase)
	userHandler := handler.NewUserHandler(authUsecase)

	// Public routes
	public := router.Group("/api/auth")
	{
		public.POST("/register", authHandler.Register)
		public.POST("/login", authHandler.Login)
		public.POST("/verify-otp", authHandler.VerifyOTP)
		public.POST("/resend-otp", authHandler.ResendOTP)
	}

	// Protected routes example
	protected := router.Group("/api")
	protected.Use(middleware.JWTAuth())
	{
		protected.GET("/me", userHandler.GetMe)

	}
}
