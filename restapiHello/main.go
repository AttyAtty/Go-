package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	// Ginエンジンのインスタンスを作成
	//インスタンスとは，クラスや構造体から作成された具体的なオブジェクトを指す．クラスや構造体の定義を基に実際にメモリ上に生成される実体のこと．
	r := gin.Default()

	// ルートURL ("/") に対するGETリクエストをハンドル
	r.GET("/", func(c *gin.Context) {
		// JSONレスポンスを返す
		c.JSON(200, gin.H{                
			"message": "Hello World",   //レスポンスの内容
		})
		//cは*gin.context型の変数で，現在のHTTPリクエストに関する情報やレスポンスを行うためのコンテキストオブジェクト
		//c.JSON()メソッドは，HTTPレスポンスをJSON形式で返すためのメソッド．
		//引数1はステータスコード（ここでは200）.200はOKを意味し，リクエストが成功したことを示す．
		//引数2はデータ（ここではgin.H{...}）．JSONとして返したいデータを指定する．gin.HはGinで抵抗される便利な型で，マップ(キーと値のペア)を簡単に作成できる
	})

	// 8080ポートでHTTPサーバーを起動
	r.Run(":8080")
}


//go run main.go で実行した後，ブラウザでhttp://localhost:8080を入力するとjsonの形式で
{
	"message": "Hello World"
}
//が表示される．


///////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////

//次はGORMと接続していく
// PostgreSQL接続情報
dsn := "host=localhost user=yourusername password=yourpassword dbname=yourdbname port=5432 sslmode=disable TimeZone=Asia/Tokyo"

// GORMでデータベースに接続
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
if err != nil {
  log.Fatal("Failed to connect to database:", err)
}