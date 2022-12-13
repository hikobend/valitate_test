package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	_ "github.com/go-sql-driver/mysql"
)

type User_JSON struct { // JSON
	First string `json:"first"`
	Last  string `json:"last"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

type User struct {
	Id        int
	FirstName string `validate:"required"`       //必須パラメータ
	LastName  string `validate:"required"`       //必須パラメータ
	Age       uint8  `validate:"gte=0,lt=130"`   // 0以上、130未満
	Email     string `validate:"required,email"` //必須パラメータ、かつ、emailフォーマット
}

func main() {
	// //バリデーション対象のデータをセット
	// user := &User{
	// 	FirstName: "Badger",
	// 	LastName:  "Smith",
	// 	Age:       135,
	// 	Email:     "Badger.Smithgmail.com",
	// }
	// validate := validator.New()     //インスタンス生成
	// errors := validate.Struct(user) //バリデーションを実行し、NGの場合、ここでエラーが返る。
	// log.Fatalln(errors)

	r := gin.Default()
	r.POST("/create", Insert)

	r.Run()
}

func Insert(c *gin.Context) {
	validate := validator.New() //インスタンス生成

	db, err := sql.Open("mysql", "root:password@(localhost:3306)/local?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var user User_JSON
	// user := &User{
	// 	FirstName: "Badger",
	// 	LastName:  "Smith",
	// 	Age:       135,
	// 	Email:     "Badger.Smithgmail.com",
	// }

	c.ShouldBindJSON(&user)

	insert, err := db.Prepare("INSERT INTO user(first, last, age, email) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	insert.Exec(validate.user.First, validate.user.Last, validate.user.Age, validate.user.Email)

}
