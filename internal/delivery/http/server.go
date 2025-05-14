package http

import (
	"dormitory_management/internal/config"
	"dormitory_management/internal/delivery/http/middleware"
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Server represents the HTTP server
type Server struct {
	router *gin.Engine
	config *config.Config
}

// NewServer creates a new HTTP server
func NewServer(cfg *config.Config, handlers *Handlers, mw *middleware.Middleware) *Server {
	// Set Gin mode
	// gin.SetMode(cfg.Server.Mode)
	gin.SetMode(gin.ReleaseMode)

	// Create router
	router := gin.New()

	// Use middlewares
	router.Use(gin.Recovery())
	router.Use(mw.LoggerMiddleware())
	router.Use(mw.ErrorMiddleware())

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Set up routes
	setupRoutes(router, handlers, mw)

	return &Server{
		router: router,
		config: cfg,
	}
}

// Run starts the HTTP server
func (s *Server) Run() error {
	return s.router.Run(fmt.Sprintf(":%s", s.config.Server.Port))
}

// setupRoutes configures all the routes
func setupRoutes(router *gin.Engine, handlers *Handlers, mw *middleware.Middleware) {
	// API group
	apiV1 := router.Group("/api/v1")
	{
		// Auth routes
		auth := apiV1.Group("/auth")
		{
			auth.POST("/register", handlers.Auth.Register())
			auth.POST("/login", handlers.Auth.Login())
			auth.POST("/refresh", handlers.Auth.RefreshToken())
			auth.POST("/logout", mw.AuthMiddleware(), handlers.Auth.Logout())
		}

		// User routes (requires authentication)
		users := apiV1.Group("/users", mw.AuthMiddleware())
		{
			users.GET("", handlers.User.GetListUsers())
			users.GET("/:id", handlers.User.GetUserByID())
			users.PUT("/:id", handlers.User.UpdateUser())
			users.DELETE("/:id", mw.AdminMiddleware(), handlers.User.DeleteUser())
			users.POST("/room", handlers.User.AddUserToRoom())
			users.POST("/remove-room", handlers.User.RemoveUserFromRoom())
		}

		// Room routes
		rooms := apiV1.Group("/rooms")
		{
			rooms.GET("", handlers.Room.GetListRooms())
			rooms.GET("/:id", handlers.Room.GetRoomByID())

			// These routes require authentication and admin role
			roomsAdmin := rooms.Group("", mw.AuthMiddleware(), mw.AdminMiddleware())
			{
				roomsAdmin.POST("", handlers.Room.CreateRoom())
				roomsAdmin.PUT("/:id", handlers.Room.UpdateRoom())
				roomsAdmin.DELETE("/:id", handlers.Room.DeleteRoom())
			}
		}

		// Room category routes
		roomCategories := apiV1.Group("/room-categories")
		{
			roomCategories.GET("", handlers.RoomCategory.GetListRoomCategories())
			roomCategories.GET("/:id", handlers.RoomCategory.GetRoomCategoryByID())

			// These routes require authentication and admin role
			roomCategoryAdmin := roomCategories.Group("", mw.AuthMiddleware(), mw.AdminMiddleware())
			{
				roomCategoryAdmin.POST("", handlers.RoomCategory.CreateRoomCategory())
				roomCategoryAdmin.PUT("/:id", handlers.RoomCategory.UpdateRoomCategory())
				roomCategoryAdmin.DELETE("/:id", handlers.RoomCategory.DeleteRoomCategory())
			}
		}

		// Contract routes (requires authentication)
		contracts := apiV1.Group("/contracts", mw.AuthMiddleware())
		{
			contracts.GET("", handlers.Contract.GetListContracts())
			contracts.GET("/:id", handlers.Contract.GetContractByID())
			contracts.GET("/user/:user_id", handlers.Contract.GetContractByUserID())

			// These routes require admin role
			contractAdmin := contracts.Group("", mw.AdminMiddleware())
			{
				contractAdmin.POST("", handlers.Contract.CreateContract())
				contractAdmin.PUT("/:id", handlers.Contract.UpdateContract())
				contractAdmin.DELETE("/:id", handlers.Contract.DeleteContract())
			}
		}
	}
}
