package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/go-playground/validator.v9"
)

type User struct {
	First string `json:"first" validate:"required"` //必須パラメータ
	Last  string `json:"last" validate:"required"`  //必須パラメータ
}

type Time struct {
	Time string `json:"time"`
}

type Date struct {
	Date string `json:"date"`
}

func main() {
	r := gin.Default()
	r.POST("/create", Insert)
	r.POST("/add", Add) // yyyymmddをデータベースに保存

	r.Run()

}

func Insert(c *gin.Context) {
	db, err := sql.Open("mysql", "root:password@(localhost:3306)/local?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var user User             // ここでは空
	validate := newFunction() //インスタンス生成

	c.ShouldBindJSON(&user) // BodyをUserという変数に入れる。

	err = validate.Struct(&user) //バリデーションを実行し、NGの場合、ここでエラーが返る。
	if err != nil {
		log.Fatal(err)
	}

	insert, err := db.Prepare("INSERT INTO users(first, last) VALUES (?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	insert.Exec(user.First, user.Last)
}

func Add(c *gin.Context) {
	db, err := sql.Open("mysql", "root:password@(localhost:3306)/local?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// JSON形式のボディをパース
	var date Date
	if err := c.ShouldBindJSON(&date); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// yyyymmdd形式の文字列をtime.Time型に変換
	layout := "20060102"
	t, err := time.Parse(layout, date.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// プリペアドステートメントを作成
	stmt, err := db.Prepare("INSERT INTO dates (date) VALUES (?)")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer stmt.Close()

	// パラメータを設定してクエリを実行
	_, err = stmt.Exec(t.Format("2006-01-02"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func newFunction() *validator.Validate {
	validate := validator.New()
	return validate
}
