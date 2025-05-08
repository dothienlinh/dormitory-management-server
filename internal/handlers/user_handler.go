package handlers

import (
	"dormitory_management/internal/models"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetListUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.logger.Info("GetListUser")
		filter := models.FilterUser{}
		if err := c.ShouldBindQuery(&filter); err != nil {
			BadRequest(c, err.Error())
			return
		}

		filter.Parse()

		conditions := []string{"role = ?"}
		values := []interface{}{models.UserRoleStudent}

		if filter.Status != "" {
			conditions = append(conditions, "status = ?")
			values = append(values, filter.Status)
		}

		if filter.Gender != "" {
			conditions = append(conditions, "gender = ?")
			values = append(values, filter.Gender)
		}

		if filter.Keyword != "" {
			conditions = append(conditions, "(full_name LIKE ? OR email LIKE ? OR phone LIKE ? OR student_code LIKE ?)")
			keyword := "%" + filter.Keyword + "%"
			values = append(values, keyword, keyword, keyword, keyword)
		}

		whereClause := strings.Join(conditions, " AND ")
		users := []models.User{}

		query := h.dbClient.Model(&models.User{}).Preload("Room")

		if err := query.Where(whereClause, values...).Count(&filter.Total).Error; err != nil {
			BadRequest(c, "Error querying data")
			return
		}

		if err := query.Where(whereClause, values...).
			Order("created_at DESC").
			Limit(filter.Limit).
			Offset(filter.GetOffset()).
			Find(&users).Error; err != nil {
			BadRequest(c, "Error querying data")
			return
		}

		Success(c, users, filter.Total)
	}
}
