package repository

import (
	"database/sql"
	"time"

	"github.com/saku0512/lightpules-backend/models"
)

// ScanRepository 検査リポジトリインターフェース
type ScanRepository interface {
	Create(scan *models.Scan) error
	GetByID(id int) (*models.Scan, error)
	GetAll() ([]*models.Scan, error)
	Update(scan *models.Scan) error
	UpdateStatus(id int, status models.ScanStatus) error
	UpdateStartedAt(id int, startedAt time.Time) error
	UpdateCompletedAt(id int, completedAt time.Time, errorMessage *string) error
	Delete(id int) error
}

type scanRepository struct {
	db *sql.DB
}

// NewScanRepository 検査リポジトリのコンストラクタ
func NewScanRepository(db *sql.DB) ScanRepository {
	return &scanRepository{db: db}
}

// Create 検査を作成
func (r *scanRepository) Create(scan *models.Scan) error {
	query := `
		INSERT INTO scans (url, status, created_at)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	err := r.db.QueryRow(query, scan.URL, scan.Status, scan.CreatedAt).Scan(&scan.ID)
	if err != nil {
		return err
	}
	return nil
}

// GetByID IDで検査を取得
func (r *scanRepository) GetByID(id int) (*models.Scan, error) {
	scan := &models.Scan{}
	query := `
		SELECT id, url, status, created_at, started_at, completed_at, error_message
		FROM scans
		WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&scan.ID,
		&scan.URL,
		&scan.Status,
		&scan.CreatedAt,
		&scan.StartedAt,
		&scan.CompletedAt,
		&scan.ErrorMessage,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return scan, nil
}

// GetAll 全ての検査を取得
func (r *scanRepository) GetAll() ([]*models.Scan, error) {
	query := `
		SELECT id, url, status, created_at, started_at, completed_at, error_message
		FROM scans
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	scans := []*models.Scan{}
	for rows.Next() {
		scan := &models.Scan{}
		err := rows.Scan(
			&scan.ID,
			&scan.URL,
			&scan.Status,
			&scan.CreatedAt,
			&scan.StartedAt,
			&scan.CompletedAt,
			&scan.ErrorMessage,
		)
		if err != nil {
			return nil, err
		}
		scans = append(scans, scan)
	}
	return scans, nil
}

// Update 検査を更新
func (r *scanRepository) Update(scan *models.Scan) error {
	query := `
		UPDATE scans
		SET url = $1, status = $2, started_at = $3, completed_at = $4, error_message = $5
		WHERE id = $6
	`
	_, err := r.db.Exec(query, scan.URL, scan.Status, scan.StartedAt, scan.CompletedAt, scan.ErrorMessage, scan.ID)
	return err
}

// UpdateStatus ステータスのみ更新
func (r *scanRepository) UpdateStatus(id int, status models.ScanStatus) error {
	query := `UPDATE scans SET status = $1 WHERE id = $2`
	_, err := r.db.Exec(query, status, id)
	return err
}

// UpdateStartedAt 開始時刻を更新
func (r *scanRepository) UpdateStartedAt(id int, startedAt time.Time) error {
	query := `UPDATE scans SET started_at = $1, status = $2 WHERE id = $3`
	_, err := r.db.Exec(query, startedAt, models.ScanStatusRunning, id)
	return err
}

// UpdateCompletedAt 完了時刻を更新
func (r *scanRepository) UpdateCompletedAt(id int, completedAt time.Time, errorMessage *string) error {
	status := models.ScanStatusCompleted
	if errorMessage != nil {
		status = models.ScanStatusFailed
	}
	query := `UPDATE scans SET completed_at = $1, status = $2, error_message = $3 WHERE id = $4`
	_, err := r.db.Exec(query, completedAt, status, errorMessage, id)
	return err
}

// Delete 検査を削除
func (r *scanRepository) Delete(id int) error {
	query := `DELETE FROM scans WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

