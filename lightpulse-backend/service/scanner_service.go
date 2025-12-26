package service

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/saku0512/lightpules-backend/models"
)

// ScannerService 脆弱性スキャンサービスインターフェース
type ScannerService interface {
	ScanURL(targetURL string) ([]*models.Vulnerability, error)
}

type scannerService struct {
	httpClient *http.Client
}

// NewScannerService スキャンサービスのコンストラクタ
func NewScannerService() ScannerService {
	return &scannerService{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ScanURL URLをスキャンして脆弱性を検出
func (s *scannerService) ScanURL(targetURL string) ([]*models.Vulnerability, error) {
	var vulnerabilities []*models.Vulnerability

	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	// SQLインジェクションのチェック
	sqlVulns := s.checkSQLInjection(parsedURL)
	vulnerabilities = append(vulnerabilities, sqlVulns...)

	// XSSのチェック
	xssVulns := s.checkXSS(parsedURL)
	vulnerabilities = append(vulnerabilities, xssVulns...)

	return vulnerabilities, nil
}

// checkSQLInjection SQLインジェクションをチェック
func (s *scannerService) checkSQLInjection(parsedURL *url.URL) []*models.Vulnerability {
	var vulnerabilities []*models.Vulnerability

	// 基本的なSQLインジェクションペイロード
	payloads := []string{
		"' OR '1'='1",
		"' OR '1'='1'--",
		"1' OR '1'='1",
		"admin'--",
		"' UNION SELECT NULL--",
	}

	values := parsedURL.Query()
	for paramName := range values {
		for _, payload := range payloads {
			testValues := url.Values{}
			for k, v := range values {
				if k == paramName {
					testValues[k] = []string{payload}
				} else {
					testValues[k] = v
				}
			}

			testURL := *parsedURL
			testURL.RawQuery = testValues.Encode()

			isVulnerable, err := s.testSQLInjectionVulnerability(testURL.String(), payload)
			if err != nil {
				continue
			}

			if isVulnerable {
				description := fmt.Sprintf("SQLインジェクションの脆弱性が検出されました。パラメータ '%s' が脆弱です。", paramName)
				recommendation := "パラメータの入力値を検証し、プリペアドステートメントを使用してください。"

				vuln := &models.Vulnerability{
					Type:           models.VulnerabilityTypeSQLInjection,
					Severity:       models.SeverityHigh,
					Location:       fmt.Sprintf("パラメータ: %s", paramName),
					Payload:        &payload,
					Description:    &description,
					Recommendation: &recommendation,
				}
				vulnerabilities = append(vulnerabilities, vuln)
				break // 1つのパラメータにつき1つの脆弱性として報告
			}
		}
	}

	return vulnerabilities
}

// testSQLInjectionVulnerability SQLインジェクションの脆弱性をテスト
func (s *scannerService) testSQLInjectionVulnerability(testURL, payload string) (bool, error) {
	resp, err := s.httpClient.Get(testURL)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// レスポンスボディを読み込んでSQLエラーメッセージをチェック
	body := make([]byte, 4096)
	n, _ := resp.Body.Read(body)
	bodyStr := strings.ToLower(string(body[:n]))

	// SQLエラーを示すキーワードをチェック
	sqlErrorKeywords := []string{
		"sql syntax",
		"mysql",
		"postgresql",
		"sqlite",
		"ora-",
		"sql error",
		"sql exception",
		"database error",
	}

	for _, keyword := range sqlErrorKeywords {
		if strings.Contains(bodyStr, keyword) {
			return true, nil
		}
	}

	return false, nil
}

// checkXSS XSSをチェック
func (s *scannerService) checkXSS(parsedURL *url.URL) []*models.Vulnerability {
	var vulnerabilities []*models.Vulnerability

	// 基本的なXSSペイロード
	payloads := []string{
		"<script>alert('XSS')</script>",
		"<img src=x onerror=alert('XSS')>",
		"javascript:alert('XSS')",
		"<svg onload=alert('XSS')>",
	}

	values := parsedURL.Query()
	for paramName := range values {
		for _, payload := range payloads {
			testValues := url.Values{}
			for k, v := range values {
				if k == paramName {
					testValues[k] = []string{payload}
				} else {
					testValues[k] = v
				}
			}

			testURL := *parsedURL
			testURL.RawQuery = testValues.Encode()

			isVulnerable, err := s.testXSSVulnerability(testURL.String(), payload)
			if err != nil {
				continue
			}

			if isVulnerable {
				description := fmt.Sprintf("XSS（クロスサイトスクリプティング）の脆弱性が検出されました。パラメータ '%s' が脆弱です。", paramName)
				recommendation := "ユーザー入力値を適切にエスケープまたはサニタイズし、Content-Security-Policyヘッダーを設定してください。"

				vuln := &models.Vulnerability{
					Type:           models.VulnerabilityTypeXSS,
					Severity:       models.SeverityMedium,
					Location:       fmt.Sprintf("パラメータ: %s", paramName),
					Payload:        &payload,
					Description:    &description,
					Recommendation: &recommendation,
				}
				vulnerabilities = append(vulnerabilities, vuln)
				break // 1つのパラメータにつき1つの脆弱性として報告
			}
		}
	}

	return vulnerabilities
}

// testXSSVulnerability XSSの脆弱性をテスト
func (s *scannerService) testXSSVulnerability(testURL, payload string) (bool, error) {
	resp, err := s.httpClient.Get(testURL)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// レスポンスボディを読み込んでペイロードがリフレクトされているかチェック
	body := make([]byte, 8192)
	n, _ := resp.Body.Read(body)
	bodyStr := string(body[:n])

	// ペイロードがそのままリフレクトされているかチェック（簡易的な検出）
	if strings.Contains(bodyStr, payload) {
		return true, nil
	}

	return false, nil
}
