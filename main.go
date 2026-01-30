package main

import (
	"todo-app/auth"
	"todo-app/todo"

	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&todo.Todo{})

	r := gin.Default()
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/token", auth.AccessToken)

	protected := r.Group("", auth.Protect([]byte("==signature==")))

	handler := todo.NewTodoHandler(db)
	protected.POST("/todos", handler.NewTask)

	r.Run()

}
