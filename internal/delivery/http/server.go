package http

import (
	"dormitory_management/internal/config"
	"dormitory_management/internal/delivery/http/handler"
	"dormitory_management/internal/delivery/http/middleware"
	"dormitory_management/internal/delivery/http/router"
	"fmt"
	"log"
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
func NewServer(cfg *config.Config, handlers *handler.Handlers, mw *middleware.Middleware) *Server {
	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Create router
	ginRouter := gin.New()

	// Use middlewares
	ginRouter.Use(gin.CustomRecovery(mw.CustomRecovery()))
	ginRouter.Use(mw.LoggerMiddleware())
	ginRouter.Use(mw.ErrorMiddleware())

	// Configure CORS
	ginRouter.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// ping
	ginRouter.GET("/ping", func(c *gin.Context) {
		c.JSON(200, map[string]string{
			"message": "pong",
		})
	})

	// webhook
	ginRouter.GET("/answerurl", func(ctx *gin.Context) {
		log.Println("=======================================================================>answerurl<=======================================================================")
	})

	// Set up routes
	router.SetupRoutes(ginRouter, handlers, mw)

	return &Server{
		router: ginRouter,
		config: cfg,
	}
}

// Run starts the HTTP server
func (s *Server) Run() error {
	return s.router.Run(fmt.Sprintf(":%s", s.config.Server.Port))
}
