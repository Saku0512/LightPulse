-- 脆弱性情報テーブル
CREATE TABLE IF NOT EXISTS vulnerabilities (
    id SERIAL PRIMARY KEY,
    scan_id INTEGER NOT NULL REFERENCES scans(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL, -- sql_injection, xss, etc.
    severity VARCHAR(20) NOT NULL, -- low, medium, high, critical
    location TEXT NOT NULL, -- パラメータ名やURLなど
    payload TEXT, -- 使用したペイロード
    description TEXT,
    recommendation TEXT, -- 推奨される対策
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- インデックス
CREATE INDEX idx_vulnerabilities_scan_id ON vulnerabilities(scan_id);
CREATE INDEX idx_vulnerabilities_type ON vulnerabilities(type);
CREATE INDEX idx_vulnerabilities_severity ON vulnerabilities(severity);

