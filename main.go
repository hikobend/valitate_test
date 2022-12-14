package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/go-playground/validator.v9"
)

type User struct {
	First string `json:"first" validate:"required"` //必須パラメータ
	Last  string `json:"last" validate:"required"`  //必須パラメータ
}

func main() {
	r := gin.Default()
	r.POST("/create", Insert)

	r.Run()

}

func Insert(c *gin.Context) {
	db, err := sql.Open("mysql", "root:password@(localhost:3306)/local?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var user User               // ここでは空
	validate := validator.New() //インスタンス生成

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
