package main

import (
  "fmt"
  "log"
  "os"

  "github.com/gin-gonic/gin"
  "github.com/joho/godotenv"
  "gorm.io/driver/postgres"
  "gorm.io/gorm"
)

func main() {
  // .envファイルから環境変数を読み込む
  err := godotenv.Load()
  if err != nil {
	log.Fatal("Error loading .env file")
  }

   // 環境変数から接続情報を取得
  dbUser := os.Getenv("POSTGRES_USER")
  dbPassword := os.Getenv("POSTGRES_PASSWORD")
  dbName := os.Getenv("POSTGRES_DB")
  dbHost := "localhost" // または環境変数から取得
  dbPort := "5432"      // または環境変数から取得

  // DSNを構築
  dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tokyo", dbHost, dbUser, dbPassword, dbName, dbPort)

  // GORMでデータベースに接続
  db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
  if err != nil {
	log.Fatal("Failed to connect to database:", err)
  }

  // データベースにテーブルを作成
  db.AutoMigrate()


  // Ginエンジンのインスタンスを作成
  r := gin.Default()

  // ルートURL ("/") に対するGETリクエストをハンドル
  r.GET("/", func(c *gin.Context) {
       // JSONレスポンスを返す
       c.JSON(200, gin.H{
       "message": "Hello World",
     })
  })

  // 8080ポートでサーバーを起動
  r.Run(":8080")
}
