package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"github.com/saku0512/lightpules-backend/handler"
	"github.com/saku0512/lightpules-backend/repository"
	"github.com/saku0512/lightpules-backend/service"
)

func main() {
	// データベース接続
	db, err := connectDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// リポジトリの初期化
	scanRepo := repository.NewScanRepository(db)
	vulnerabilityRepo := repository.NewVulnerabilityRepository(db)

	// サービスの初期化
	scanService := service.NewScanService(scanRepo, vulnerabilityRepo)
	scannerService := service.NewScannerService()

	// ハンドラーの初期化
	scanHandler := handler.NewScanHandler(scanService, scannerService)
	healthHandler := handler.NewHealthHandler()

	// ルーターの設定
	r := gin.Default()

	// CORS設定
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"} // 開発環境用、本番環境では適切に設定
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

	// APIルート
	api := r.Group("/api")
	{
		api.GET("/health", healthHandler.Health)
		api.POST("/scans", scanHandler.CreateScan)
		api.GET("/scans", scanHandler.GetAllScans)
		api.GET("/scans/:id", scanHandler.GetScanByID)
		api.DELETE("/scans/:id", scanHandler.DeleteScan)
	}

	// サーバー起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// connectDB データベースに接続
func connectDB() (*sql.DB, error) {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "lightpulse")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// getEnv 環境変数を取得、デフォルト値を返す
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
