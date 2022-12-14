package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	_ "github.com/go-sql-driver/mysql"
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

	//バリデーション対象のデータをセット
	var user User
	validate := validator.New()  //インスタンス生成
	err = validate.Struct(&user) //バリデーションを実行し、NGの場合、ここでエラーが返る。
	if err != nil {
		log.Fatal(err)
	}
	// c.ShouldBindJSON(&user)

	insert, err := db.Prepare("INSERT INTO users(first, last) VALUES (?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	insert.Exec(user.First, user.Last)
}
