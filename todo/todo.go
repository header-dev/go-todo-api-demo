package todo

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"gorm.io/gorm"
)

type Todo struct {
	Title string `json:"text" validate:"required"`
	gorm.Model
}

func (Todo) TableName() string {
	return "todoList"
}

type TodoHandler struct {
	db       *gorm.DB
	validate *validator.Validate
}

func NewTodoHandler(db *gorm.DB) *TodoHandler {

	return &TodoHandler{
		db:       db,
		validate: validator.New(),
	}
}

func (t *TodoHandler) NewTask(c *gin.Context) {

	var todo Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if todo.Title == "sleep" {
		transId := c.Request.Header.Get("X-Transaction-ID")
		aud, _ := c.Get("aud")
		log.Println(transId, aud, "not allowed")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "access denied",
		})
		return
	}

	err := t.validate.Struct(todo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	r := t.db.Create(&todo)
	if err := r.Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"ID": todo.Model.ID,
	})
}
