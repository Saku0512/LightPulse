package service

import (
	"errors"
	"time"

	"github.com/saku0512/lightpules-backend/models"
	"github.com/saku0512/lightpules-backend/repository"
)

// ScanService 検査サービスインターフェース
type ScanService interface {
	CreateScan(url string) (*models.Scan, error)
	GetScanByID(id int) (*models.ScanWithVulnerabilities, error)
	GetAllScans() ([]*models.Scan, error)
	StartScan(id int) error
	CompleteScan(id int, vulnerabilities []*models.Vulnerability, errorMessage *string) error
	DeleteScan(id int) error
}

type scanService struct {
	scanRepo          repository.ScanRepository
	vulnerabilityRepo repository.VulnerabilityRepository
}

// NewScanService 検査サービスのコンストラクタ
func NewScanService(scanRepo repository.ScanRepository, vulnerabilityRepo repository.VulnerabilityRepository) ScanService {
	return &scanService{
		scanRepo:          scanRepo,
		vulnerabilityRepo: vulnerabilityRepo,
	}
}

// CreateScan 新しい検査を作成
func (s *scanService) CreateScan(url string) (*models.Scan, error) {
	if url == "" {
		return nil, errors.New("URL is required")
	}

	scan := &models.Scan{
		URL:       url,
		Status:    models.ScanStatusPending,
		CreatedAt: time.Now(),
	}

	err := s.scanRepo.Create(scan)
	if err != nil {
		return nil, err
	}

	return scan, nil
}

// GetScanByID IDで検査を取得（脆弱性情報も含む）
func (s *scanService) GetScanByID(id int) (*models.ScanWithVulnerabilities, error) {
	scan, err := s.scanRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if scan == nil {
		return nil, errors.New("scan not found")
	}

	vulnerabilities, err := s.vulnerabilityRepo.GetByScanID(id)
	if err != nil {
		return nil, err
	}

	return &models.ScanWithVulnerabilities{
		Scan:            scan,
		Vulnerabilities: vulnerabilities,
	}, nil
}

// GetAllScans 全ての検査を取得
func (s *scanService) GetAllScans() ([]*models.Scan, error) {
	return s.scanRepo.GetAll()
}

// StartScan 検査を開始
func (s *scanService) StartScan(id int) error {
	scan, err := s.scanRepo.GetByID(id)
	if err != nil {
		return err
	}
	if scan == nil {
		return errors.New("scan not found")
	}

	if scan.Status != models.ScanStatusPending {
		return errors.New("scan is not in pending status")
	}

	return s.scanRepo.UpdateStartedAt(id, time.Now())
}

// CompleteScan 検査を完了
func (s *scanService) CompleteScan(id int, vulnerabilities []*models.Vulnerability, errorMessage *string) error {
	scan, err := s.scanRepo.GetByID(id)
	if err != nil {
		return err
	}
	if scan == nil {
		return errors.New("scan not found")
	}

	// 既存の脆弱性を削除
	err = s.vulnerabilityRepo.DeleteByScanID(id)
	if err != nil {
		return err
	}

	// 新しい脆弱性を保存
	for _, vuln := range vulnerabilities {
		vuln.ScanID = id
		vuln.CreatedAt = time.Now()
		err := s.vulnerabilityRepo.Create(vuln)
		if err != nil {
			return err
		}
	}

	// 検査を完了状態に更新
	return s.scanRepo.UpdateCompletedAt(id, time.Now(), errorMessage)
}

// DeleteScan 検査を削除
func (s *scanService) DeleteScan(id int) error {
	scan, err := s.scanRepo.GetByID(id)
	if err != nil {
		return err
	}
	if scan == nil {
		return errors.New("scan not found")
	}

	// 脆弱性も一緒に削除される（CASCADE）
	return s.scanRepo.Delete(id)
}
