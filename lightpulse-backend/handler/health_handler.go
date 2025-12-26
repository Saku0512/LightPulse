package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saku0512/lightpules-backend/models"
)

// HealthHandler ヘルスチェックハンドラー
type HealthHandler struct{}

// NewHealthHandler ヘルスチェックハンドラーのコンストラクタ
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Health ヘルスチェック
// GET /api/health
func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, models.SuccessResponse(map[string]string{
		"status": "ok",
	}))
}
