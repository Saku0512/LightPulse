# LightPulse Backend

LightPulseのバックエンドAPIサーバー

## 技術スタック

- Go 1.23+
- Gin (Webフレームワーク)
- PostgreSQL (データベース)
- アーキテクチャ: Repository-Service-Handler パターン

## プロジェクト構造

```
.
├── main.go                 # アプリケーションエントリーポイント
├── migrations/             # データベースマイグレーションファイル
│   ├── 001_create_scans_table.up.sql
│   ├── 001_create_scans_table.down.sql
│   ├── 002_create_vulnerabilities_table.up.sql
│   └── 002_create_vulnerabilities_table.down.sql
├── models/                 # データモデル
│   ├── scan.go
│   ├── vulnerability.go
│   ├── scan_request.go
│   └── response.go
├── repository/             # データアクセス層
│   ├── scan_repository.go
│   └── vulnerability_repository.go
├── service/                # ビジネスロジック層
│   ├── scan_service.go
│   └── scanner_service.go
└── handler/                # HTTPハンドラー層
    ├── scan_handler.go
    └── health_handler.go
```

## 環境変数

- `PORT`: サーバーのポート番号 (デフォルト: 8080)
- `DB_HOST`: データベースホスト (デフォルト: localhost)
- `DB_PORT`: データベースポート (デフォルト: 5432)
- `DB_USER`: データベースユーザー名 (デフォルト: postgres)
- `DB_PASSWORD`: データベースパスワード (デフォルト: postgres)
- `DB_NAME`: データベース名 (デフォルト: lightpulse)

## データベースマイグレーション

マイグレーションファイルは `migrations/` ディレクトリに配置されています。
マイグレーションツール（例: migrate, golang-migrate）を使用してデータベースを初期化してください。

例（golang-migrateを使用する場合）:
```bash
migrate -path ./migrations -database "postgres://user:password@localhost/lightpulse?sslmode=disable" up
```

## API エンドポイント

### ヘルスチェック
- `GET /api/health` - サーバーのヘルスチェック

### 検査（Scan）
- `POST /api/scans` - 新しい検査を作成
- `GET /api/scans` - 全ての検査を取得
- `GET /api/scans/:id` - 指定されたIDの検査を取得（脆弱性情報を含む）
- `DELETE /api/scans/:id` - 検査を削除

## 実行方法

```bash
# 依存関係のインストール
go mod download

# サーバーの起動
go run main.go
```

## ビルド

```bash
go build -o lightpulse-backend main.go
```

