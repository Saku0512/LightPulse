package models

// ScanRequest 検査リクエストモデル
type ScanRequest struct {
	URL string `json:"url" binding:"required,url"`
}

