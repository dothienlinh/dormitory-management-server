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

type Server struct {
	router *gin.Engine
	config *config.Config
}

func NewServer(cfg *config.Config, handlers *handler.Handlers, mw *middleware.Middleware) *Server {
	gin.SetMode(cfg.Server.Mode)

	ginRouter := gin.New()

	ginRouter.Use(gin.CustomRecovery(mw.CustomRecovery()))
	ginRouter.Use(mw.LoggerMiddleware())
	ginRouter.Use(mw.ErrorMiddleware())

	ginRouter.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	ginRouter.GET("/ping", func(c *gin.Context) {
		c.JSON(200, map[string]string{
			"message": "pong",
		})
	})

	ginRouter.GET("/answerurl", func(ctx *gin.Context) {
		log.Println("=======================================================================>answerurl<=======================================================================")
	})

	router.SetupRoutes(ginRouter, handlers, mw)

	return &Server{
		router: ginRouter,
		config: cfg,
	}
}

func (s *Server) Run() error {
	return s.router.Run(fmt.Sprintf(":%s", s.config.Server.Port))
}
