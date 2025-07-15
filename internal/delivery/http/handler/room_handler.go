package handler

import (
	"dormitory_management/internal/domain/entity"
	"dormitory_management/internal/domain/response"
	"dormitory_management/internal/domain/usecase"
	"dormitory_management/pkg/logger"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RoomHandler struct {
	useCases usecase.UseCases
	logger   logger.Logger
}

func NewRoomHandler(useCases usecase.UseCases, logger logger.Logger) *RoomHandler {
	return &RoomHandler{
		useCases: useCases,
		logger:   logger,
	}
}

func (h *RoomHandler) CreateRoom() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse

		var room entity.CreateRoom
		if err := c.ShouldBindJSON(&room); err != nil {
			h.logger.Error("Failed to bind JSON", zap.Error(err))
			c.Error(err)
			return
		}

		resp = h.useCases.Room().CreateRoom(c, &room)
		c.JSON(resp.Status, resp.Response)
	}
}

func (h *RoomHandler) GetRoomByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse

		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			h.logger.Error("Failed to parse room ID", zap.Error(err))
			resp = response.BadRequest("Invalid room ID")
			c.JSON(resp.Status, resp.Response)
			return
		}

		resp = h.useCases.Room().GetRoomByID(c, id)
		c.JSON(resp.Status, resp.Response)
	}
}

func (h *RoomHandler) GetListRooms() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse

		var filter entity.RoomFilter
		if err := c.ShouldBindQuery(&filter); err != nil {
			h.logger.Error("Failed to bind query parameters", zap.Error(err))
			c.Error(err)
			return
		}

		resp = h.useCases.Room().GetListRooms(c, &filter)
		c.JSON(resp.Status, resp.Response)
	}
}

func (h *RoomHandler) UpdateRoom() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse

		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			h.logger.Error("Failed to parse room ID", zap.Error(err))
			resp = response.BadRequest("Invalid room ID")
			c.JSON(resp.Status, resp.Response)
			return
		}

		var room entity.UpdateRoom
		if err := c.ShouldBindJSON(&room); err != nil {
			h.logger.Error("Failed to bind JSON", zap.Error(err))
			c.Error(err)
			return
		}

		resp = h.useCases.Room().UpdateRoom(c, id, &room)
		c.JSON(resp.Status, resp.Response)
	}
}

func (h *RoomHandler) DeleteRoom() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse

		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			h.logger.Error("Failed to parse room ID", zap.Error(err))
			resp = response.BadRequest("Invalid room ID")
			c.JSON(resp.Status, resp.Response)
			return
		}

		resp = h.useCases.Room().DeleteRoom(c, id)
		c.JSON(resp.Status, resp.Response)
	}
}

type RoomCategoryHandler struct {
	useCases usecase.UseCases
	logger   logger.Logger
}

func NewRoomCategoryHandler(useCases usecase.UseCases, logger logger.Logger) *RoomCategoryHandler {
	return &RoomCategoryHandler{
		useCases: useCases,
		logger:   logger,
	}
}

func (h *RoomCategoryHandler) CreateRoomCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse

		var category entity.CreateRoomCategory
		if err := c.ShouldBindJSON(&category); err != nil {
			h.logger.Error("Failed to bind JSON", zap.Error(err))
			c.Error(err)
			return
		}

		resp = h.useCases.RoomCategory().CreateRoomCategory(c, &category)
		c.JSON(resp.Status, resp.Response)
	}
}

func (h *RoomCategoryHandler) GetRoomCategoryByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse

		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			h.logger.Error("Failed to parse room category ID", zap.Error(err))
			resp = response.BadRequest("Invalid room category ID")
			c.JSON(resp.Status, resp.Response)
			return
		}

		resp = h.useCases.RoomCategory().GetRoomCategoryByID(c, id)
		c.JSON(resp.Status, resp.Response)
	}
}

func (h *RoomCategoryHandler) GetListRoomCategories() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse

		var filter entity.RoomCategoryFilter
		if err := c.ShouldBindQuery(&filter); err != nil {
			h.logger.Error("Failed to bind query parameters", zap.Error(err))
			c.Error(err)
			return
		}

		resp = h.useCases.RoomCategory().GetListRoomCategories(c, &filter)
		c.JSON(resp.Status, resp.Response)
	}
}

func (h *RoomCategoryHandler) UpdateRoomCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse

		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			h.logger.Error("Failed to parse room category ID", zap.Error(err))
			resp = response.BadRequest("Invalid room category ID")
			c.JSON(resp.Status, resp.Response)
			return
		}

		var category entity.UpdateRoomCategory
		if err := c.ShouldBindJSON(&category); err != nil {
			h.logger.Error("Failed to bind JSON", zap.Error(err))
			c.Error(err)
			return
		}

		resp = h.useCases.RoomCategory().UpdateRoomCategory(c, id, &category)
		c.JSON(resp.Status, resp.Response)
	}
}

func (h *RoomCategoryHandler) DeleteRoomCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse

		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			h.logger.Error("Failed to parse room category ID", zap.Error(err))
			resp = response.BadRequest("Invalid room category ID")
			c.JSON(resp.Status, resp.Response)
			return
		}

		resp = h.useCases.RoomCategory().DeleteRoomCategory(c, id)
		c.JSON(resp.Status, resp.Response)
	}
}

// Student Room APIs

// GetStudentRoomDetails retrieves detailed room information for a student
func (h *RoomHandler) GetStudentRoomDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint64("user_id")
		if userID == 0 {
			resp := response.Unauthorized("Unauthorized: user_id missing")
			c.JSON(resp.Status, resp.Response)
			return
		}
		resp := h.useCases.Room().GetStudentRoomDetails(c, userID)
		c.JSON(resp.Status, resp.Response)
	}
}

// GetRoomStats retrieves room statistics for a student
func (h *RoomHandler) GetRoomStats() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint64("user_id")

		resp := h.useCases.Room().GetRoomStats(c, userID)
		c.JSON(resp.Status, resp.Response)
	}
}

// GetRoomIssues retrieves room issues with pagination
func (h *RoomHandler) GetRoomIssues() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint64("user_id")
		if userID == 0 {
			resp := response.Unauthorized("Unauthorized: user_id missing")
			c.JSON(resp.Status, resp.Response)
			return
		}
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		status := c.Query("status")
		category := c.Query("category")

		resp := h.useCases.Room().GetRoomIssues(c, userID, page, limit, status, category)
		c.JSON(resp.Status, resp.Response)
	}
}

// CreateRoomIssue creates a new room issue
func (h *RoomHandler) CreateRoomIssue() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint64("user_id")

		var req struct {
			Title       string `json:"title" binding:"required"`
			Description string `json:"description" binding:"required"`
			Category    string `json:"category" binding:"required"`
			Priority    string `json:"priority" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			h.logger.Error("Failed to bind JSON", zap.Error(err))
			resp := response.BadRequest("Invalid request payload")
			c.JSON(resp.Status, resp.Response)
			return
		}

		resp := h.useCases.Room().CreateRoomIssue(c, userID, req.Title, req.Description, req.Category, req.Priority)
		c.JSON(resp.Status, resp.Response)
	}
}

// GetRoomIssueDetails retrieves details of a specific room issue
func (h *RoomHandler) GetRoomIssueDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint64("user_id")

		issueID, err := strconv.ParseUint(c.Param("issueId"), 10, 64)
		if err != nil {
			h.logger.Error("Failed to parse issue ID", zap.Error(err))
			resp := response.BadRequest("Invalid issue ID")
			c.JSON(resp.Status, resp.Response)
			return
		}

		resp := h.useCases.Room().GetRoomIssueDetails(c, issueID, userID)
		c.JSON(resp.Status, resp.Response)
	}
}

// GetRoomBills retrieves room bills with pagination
func (h *RoomHandler) GetRoomBills() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint64("user_id")

		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		year, _ := strconv.Atoi(c.Query("year"))
		month, _ := strconv.Atoi(c.Query("month"))
		status := c.Query("status")

		resp := h.useCases.Room().GetRoomBills(c, userID, page, limit, year, month, status)
		c.JSON(resp.Status, resp.Response)
	}
}

// GetRoomRules retrieves all room rules
func (h *RoomHandler) GetRoomRules() gin.HandlerFunc {
	return func(c *gin.Context) {
		resp := h.useCases.Room().GetRoomRules(c)
		c.JSON(resp.Status, resp.Response)
	}
}

// GetCleaningSchedule retrieves cleaning schedule for user's room
func (h *RoomHandler) GetCleaningSchedule() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint64("user_id")
		if userID == 0 {
			resp := response.Unauthorized("Unauthorized: user_id missing")
			c.JSON(resp.Status, resp.Response)
			return
		}
		resp := h.useCases.Room().GetCleaningSchedule(c, userID)
		c.JSON(resp.Status, resp.Response)
	}
}

// GetRoommates retrieves roommates for a user
func (h *RoomHandler) GetRoommates() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint64("user_id")

		resp := h.useCases.Room().GetRoommates(c, userID)
		c.JSON(resp.Status, resp.Response)
	}
}
