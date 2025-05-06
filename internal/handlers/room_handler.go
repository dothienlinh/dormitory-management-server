package handlers

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateRoom() gin.HandlerFunc {
	return func(c *gin.Context) {
		Success(c, nil, 0)
	}
}
