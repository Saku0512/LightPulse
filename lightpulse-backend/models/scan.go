package models

import "time"

// ScanStatus 検査のステータス
type ScanStatus string

const (
	ScanStatusPending   ScanStatus = "pending"
	ScanStatusRunning   ScanStatus = "running"
	ScanStatusCompleted ScanStatus = "completed"
	ScanStatusFailed    ScanStatus = "failed"
)

// Scan 検査情報モデル
type Scan struct {
	ID           int        `json:"id" db:"id"`
	URL          string     `json:"url" db:"url"`
	Status       ScanStatus `json:"status" db:"status"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	StartedAt    *time.Time `json:"started_at,omitempty" db:"started_at"`
	CompletedAt  *time.Time `json:"completed_at,omitempty" db:"completed_at"`
	ErrorMessage *string    `json:"error_message,omitempty" db:"error_message"`
}

// ScanWithVulnerabilities 脆弱性情報を含む検査モデル
type ScanWithVulnerabilities struct {
	Scan          *Scan           `json:"scan"`
	Vulnerabilities []*Vulnerability `json:"vulnerabilities"`
}

