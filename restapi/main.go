//GORM接続からの工程はこっちのrestapiでぃれくとりに書いてる
package main

import (   //必要なパッケージをインポートしている．
  "fmt"    //フォーマットされたI/Oを提供するパッケージ．fmt.Sprintfなどの関数を使える．
  "log"    //ログ出力のためのパッケージ．エラーメッセージを出力する際に使用．
  "os"     //OSとのインタラクションを行うためのパッケージ．環境変数を取得する際に使用．

  "github.com/gin-gonic/gin" //Ginのフレームワークを使用するためのパッケージ．WEBアプリケーションを構築
  "github.com/joho/godotenv" //環境変数を.envから読み込むためのパッケージ
  "gorm.io/driver/postgres"  //PostgreSQLデータベースドライバー
  "gorm.io/gorm"             //GORM=Go用のORM(Object Relational Mapping)ライブラリで，データベースとのやり取りを簡素化する．
)

//{3}
type Task struct {
//フィールド名  データ型　　　　制約
	ID          uint           `gorm:"primary_key"`  //uint(符号なし整数)で，タスクの一意な識別子として機能する．''タグは，このフィールドがテーブルの主キーであることを示している．
	Task        string         `gorm:"size:255"`     //タスクの内容や説明を格納．このフィールドは最大255文字の長さを持つ．
	IsCompleted bool           `gorm:"default:false"`//タスクが完了したかどうかを示す．新しいレコードが作成されたときにそのフィールドのデフォルト値がfaultであることを示す．
	CreatedAt   time.Time      `gorm:"default:CURRENT_TIMESTAMP"`//タスクが作成された日時を格納する．新しいレコードが作成されたときに現在のタイムスランプを自動的にこのフィールドに設定する．
	UpdatedAt   time.Time      `gorm:"default:CURRENT_TIMESTAMP"`//タスクが最後に更新された日時を格納する．レコードが更新されるたびに．．．．．．．．．．．．．．．．．．．．．．．．．．．．
  }

  //{1}
func main() {
  // .envファイルから環境変数を読み込む
  err := godotenv.Load()          //envファイルを読み込み，その内容を現在のプロセスの環境変数として設定．このおかげで.envファイル内のキー＝値のペアが環境変数として利用可能となる．
  if err != nil {         //上の環境変数の読み込みが成功すると，nilを返す．
	log.Fatal("Error loading .env file")  //環境変数の読み込みが失敗したときに出力されるエラーメッセージ
  }
   //{2}
   // 環境変数からSQLデータベースの接続情報(ユーザー名，パスワード，データベース名，ホスト，ポートとか)を取得
  dbUser := os.Getenv("POSTGRES_USER")
  dbPassword := os.Getenv("POSTGRES_PASSWORD")
  dbName := os.Getenv("POSTGRES_DB")
  dbHost := "localhost" // または環境変数から取得
  dbPort := "5432"      // または環境変数から取得

  //{2}
  //GORMと接続していく
  // PostgreSQL接続情報 
  // DSN(Data Source Name)を構築
  dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tokyo", dbHost, dbUser, dbPassword, dbName, dbPort)   //Print系関数の書き方はCとかと似ている．文字列を%sに代入するとことか

  // GORMでデータベースに接続
  db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
  if err != nil {   //nil=接続に成功
	log.Fatal("Failed to connect to database:", err)   //接続に失敗したときんおエラーコード
  }

  //{2}{3}
  // データベースに定義されたモデルについてテーブルを作成．この行は具体的なモデルを指定していないので，実際にはテーブルはさくせいされない．必要に応じてモデルを追加する
  // テーブルの自動作成：構造体に基づいて，対応するテーブルがデータメース内に存在しない場合，GORMはそのテーブルを自動的に作成する．
  // スキーマの自動更新：既存のテーブルのスキーマ（構造）を自動的に更新してくれる． 
  // 安全なスキーマ変更：AutoMigrateはデータの損失を引き起こす可能性のある変更を行わない．ー＞比較的安全に使える．
  // つまり．下の関数は，構造体を作成し，そこから，DBのスキーマを自動で生成してくれる．
  db.AutoMigrate(&Task{}) //tasksテーブルを作成するようにしている．上の構造体を用いて

  //{1}
  // Ginエンジンのインスタンスを作成
  r := gin.Default()

  //{4}
// タスクを取得するエンドポイント
r.GET("/tasks", func(c *gin.Context) {   //URLパスが/tasksであるGETリクエストをハンドルするためのルートを定義している．クライアントがこのエンドポイントにリクエストを送ると，指定された無名関数が実行される．
	var tasks []Task                     //Task型のスライスを宣言する．データベースから取得されるタスクを格納するために使用される．
	db.Find(&tasks)                      //GORMのFindメソッドを使って，データベース内の全てのTaskレコードを取得し，それらをtasksスライスに格納する．
	c.JSON(http.StatusOK, tasks)         //ステータスコード200 OKとともに取得したタスクのリストをJSON形式でクライアントに返す．
  })

  //{5}
// 新しいタスクを(Create)作成するエンドポイント
r.POST("/tasks", func(c *gin.Context) {
	var task Task
	if err := c.ShouldBindJSON(&task); err != nil {
	  c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	  return
	}
	db.Create(&task)                 //Gormの機能. taskをpostしている
	c.JSON(http.StatusOK, task)
  })

  //{6}
  // タスクを(Update)更新するエンドポイント
r.PUT("/tasks/:id", func(c *gin.Context) {
	var task Task                //Task型の変数を宣言．データベースから取得される特定のタスクを格納するために使用される．
	id := c.Param("id")          //リクエストURLからタスクのIDを取得する．
  
	if err := db.First(&task, id).Error; err != nil {     //GORMのFirstメソッドにより，指定されたIDを持つ最初のTaskレコードを取得し，それをtask変数に格納．
	  c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
	  return
	}
  
	if err := c.ShouldBindJSON(&task); err != nil {
	  c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	  return
	}
  
	db.Save(&task)        //GORMのSaveメソッドを使用して，更新されたtaskをデータベースに保存する．
	c.JSON(http.StatusOK, task)      //更新されたタスクをJSON形式でクライアントに返す．
  })

//{7}
// タスクを(Delete)削除するエンドポイント
r.DELETE("/tasks/:id", func(c *gin.Context) {
	var task Task
	id := c.Param("id")  //idを取得
  
	if err := db.First(&task, id).Error; err != nil {
	  c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
	  return
	}
  
	db.Delete(&task)  //taskを削除
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
  })

  //{1}
  // ルートURL ("/") に対するGETリクエストをハンドル
  r.GET("/", func(c *gin.Context) {
       // JSONレスポンスを返す
       c.JSON(200, gin.H{
       "message": "Hello World",
     })
  })

  //{1}
  // 8080ポートでサーバーを起動
  r.Run(":8080")
}
