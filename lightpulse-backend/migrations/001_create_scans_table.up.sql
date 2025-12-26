-- 検査情報テーブル
CREATE TABLE IF NOT EXISTS scans (
    id SERIAL PRIMARY KEY,
    url TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending, running, completed, failed
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    error_message TEXT
);

-- インデックス
CREATE INDEX idx_scans_status ON scans(status);
CREATE INDEX idx_scans_created_at ON scans(created_at DESC);

