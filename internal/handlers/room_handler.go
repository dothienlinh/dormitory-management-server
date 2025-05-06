package handlers

import (
	"dormitory_management/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateRoom(c *gin.Context) {
	c.JSON(http.StatusCreated, utils.SuccessResponse("Create room success", nil))
}
