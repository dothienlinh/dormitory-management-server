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

// ContractHandler handles HTTP requests related to contracts
type ContractHandler struct {
	useCases usecase.UseCases
	logger   logger.Logger
}

// NewContractHandler creates a new ContractHandler
func NewContractHandler(useCases usecase.UseCases, logger logger.Logger) *ContractHandler {
	return &ContractHandler{
		useCases: useCases,
		logger:   logger,
	}
}

// CreateContract handles the request to create a contract
func (h *ContractHandler) CreateContract() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse
		h.logger.Info("CreateContract")

		var contract entity.Contract
		if err := c.ShouldBindJSON(&contract); err != nil {
			h.logger.Error("Failed to bind JSON", zap.Error(err))
			resp = response.BadRequest(err.Error())
			c.JSON(resp.Status, resp.Response)
			return
		}

		resp = h.useCases.Contract().CreateContract(c, &contract)
		c.JSON(resp.Status, resp.Response)
	}
}

// GetContractByID handles the request to get a contract by ID
func (h *ContractHandler) GetContractByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse
		h.logger.Info("GetContractByID")

		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			h.logger.Error("Failed to parse contract ID", zap.Error(err))
			resp = response.BadRequest("Invalid contract ID")
			c.JSON(resp.Status, resp.Response)
			return
		}

		resp = h.useCases.Contract().GetContractByID(c, uint(id))
		c.JSON(resp.Status, resp.Response)
	}
}

// GetContractByUserID handles the request to get a contract by user ID
func (h *ContractHandler) GetContractByUserID() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse
		h.logger.Info("GetContractByUserID")

		userID, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
		if err != nil {
			h.logger.Error("Failed to parse user ID", zap.Error(err))
			resp = response.BadRequest("Invalid user ID")
			c.JSON(resp.Status, resp.Response)
			return
		}

		resp = h.useCases.Contract().GetContractByUserID(c, uint(userID))
		c.JSON(resp.Status, resp.Response)
	}
}

// GetListContracts handles the request to get a list of contracts
func (h *ContractHandler) GetListContracts() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse
		h.logger.Info("GetListContracts")

		var filter entity.ContractFilter
		if err := c.ShouldBindQuery(&filter); err != nil {
			h.logger.Error("Failed to bind query parameters", zap.Error(err))
			resp = response.BadRequest(err.Error())
			c.JSON(resp.Status, resp.Response)
			return
		}

		resp = h.useCases.Contract().GetListContracts(c, &filter)
		c.JSON(resp.Status, resp.Response)
	}
}

// UpdateContract handles the request to update a contract
func (h *ContractHandler) UpdateContract() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse
		h.logger.Info("UpdateContract")

		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			h.logger.Error("Failed to parse contract ID", zap.Error(err))
			resp = response.BadRequest("Invalid contract ID")
			c.JSON(resp.Status, resp.Response)
			return
		}

		var updateData entity.UpdateContract
		if err := c.ShouldBindJSON(&updateData); err != nil {
			h.logger.Error("Failed to bind JSON", zap.Error(err))
			resp = response.BadRequest(err.Error())
			c.JSON(resp.Status, resp.Response)
			return
		}

		resp = h.useCases.Contract().UpdateContract(c, uint(id), &updateData)
		c.JSON(resp.Status, resp.Response)
	}
}

// DeleteContract handles the request to delete a contract
func (h *ContractHandler) DeleteContract() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response.StatusResponse
		h.logger.Info("DeleteContract")

		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			h.logger.Error("Failed to parse contract ID", zap.Error(err))
			resp = response.BadRequest("Invalid contract ID")
			c.JSON(resp.Status, resp.Response)
			return
		}

		resp = h.useCases.Contract().DeleteContract(c, uint(id))
		c.JSON(resp.Status, resp.Response)
	}
}
