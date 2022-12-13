package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
)

type JsonRequest struct {
	FieldStr  string `json:"field_str"`
	FieldInt  int    `json:"field_int" validate:"gte=0,lt=130"`
	FieldBool bool   `json:"field_bool"`
}

func main() {
	r := gin.Default()
	validate := validator.New() //インスタンス生成
	r.POST("/postjson", func(c *gin.Context) {
		var json JsonRequest
		errors := validate.Struct(json) //バリデーションを実行し、NGの場合、ここでエラーが返る。
		if errors != nil {
			log.Fatal(errors)
		}
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"str": json.FieldStr, "int": json.FieldInt, "bool": json.FieldBool})
	})

	r.Run()
}
