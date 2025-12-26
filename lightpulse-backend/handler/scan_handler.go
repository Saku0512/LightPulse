package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/saku0512/lightpules-backend/models"
	"github.com/saku0512/lightpules-backend/service"
)

// ScanHandler 検査ハンドラー
type ScanHandler struct {
	scanService    service.ScanService
	scannerService service.ScannerService
}

// NewScanHandler 検査ハンドラーのコンストラクタ
func NewScanHandler(scanService service.ScanService, scannerService service.ScannerService) *ScanHandler {
	return &ScanHandler{
		scanService:    scanService,
		scannerService: scannerService,
	}
}

// CreateScan 新しい検査を作成
// POST /api/scans
func (h *ScanHandler) CreateScan(c *gin.Context) {
	var req models.ScanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid request: "+err.Error()))
		return
	}

	scan, err := h.scanService.CreateScan(req.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("Failed to create scan: "+err.Error()))
		return
	}

	// 非同期でスキャンを実行
	go h.runScan(scan.ID, req.URL)

	c.JSON(http.StatusCreated, models.SuccessResponse(scan))
}

// runScan 非同期でスキャンを実行
func (h *ScanHandler) runScan(scanID int, targetURL string) {
	// スキャンを開始
	err := h.scanService.StartScan(scanID)
	if err != nil {
		errMsg := err.Error()
		h.scanService.CompleteScan(scanID, nil, &errMsg)
		return
	}

	// 脆弱性をスキャン
	vulnerabilities, err := h.scannerService.ScanURL(targetURL)
	if err != nil {
		errMsg := err.Error()
		h.scanService.CompleteScan(scanID, nil, &errMsg)
		return
	}

	// スキャンを完了
	err = h.scanService.CompleteScan(scanID, vulnerabilities, nil)
	if err != nil {
		errMsg := err.Error()
		h.scanService.CompleteScan(scanID, nil, &errMsg)
		return
	}
}

// GetScanByID IDで検査を取得
// GET /api/scans/:id
func (h *ScanHandler) GetScanByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid ID"))
		return
	}

	scan, err := h.scanService.GetScanByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse("Scan not found: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(scan))
}

// GetAllScans 全ての検査を取得
// GET /api/scans
func (h *ScanHandler) GetAllScans(c *gin.Context) {
	scans, err := h.scanService.GetAllScans()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("Failed to get scans: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(scans))
}

// DeleteScan 検査を削除
// DELETE /api/scans/:id
func (h *ScanHandler) DeleteScan(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid ID"))
		return
	}

	err = h.scanService.DeleteScan(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("Failed to delete scan: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(nil))
}
