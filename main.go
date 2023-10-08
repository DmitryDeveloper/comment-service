package main

import (
	"net/http"
	"strconv"
	"encoding/json"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	db "comment-service/dbConnector"
	"github.com/jinzhu/gorm"
)

func main() {
	e := echo.New()

	//TODO VALIDATIONS

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/comments", getAllComments)
	e.POST("/comments", createComment)
	e.GET("/comments/:id", getComment)
	e.PUT("/comments/:id", updateComment)
	e.DELETE("/comments/:id", deleteComment)
	e.GET("/events/:id/comments", getEventComments)

	//Run migration
	db.GetDB().AutoMigrate(&Comment{})

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

type Comment struct {
	gorm.Model
	Text string `json:"text"`
	IsEdited bool `json:"is_edited"`
    EventId int `json:"event_id"`
}

func NewComment() *Comment {
    return &Comment{
        IsEdited: false,
    }
}

//----------
// Handlers
//----------

func createComment(c echo.Context) error {
	
	com := NewComment()
	if err := c.Bind(com); err != nil {
		return err
	}
    
    //STORE in DB
	db.GetDB().Create(com)

	return c.JSON(http.StatusCreated, com)
}

func getComment(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

    com := &Comment{}

	//GET IN DB by id
	db.GetDB().First(&com, id)

	return c.JSON(http.StatusOK, com)
}

func updateComment(c echo.Context) error {
	com := new(Comment)
	if err := c.Bind(com); err != nil {
		return err
	}
	id, _ := strconv.Atoi(c.Param("id"))

	// Размаршалировать JSON-данные в структуру
    var requestData map[string]interface{}
    if err := json.NewDecoder(c.Request().Body).Decode(&requestData); err != nil {
        return err
    }

    // Извлечь значение Text
    text, ok := requestData["Text"].(string)
    if !ok {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid Text format")
    }

    //UPDATE IN DB
	db.GetDB().Model(&com).Where("id = ?", id).Updates(map[string]interface{}{"Text": text, "IsEdited": true})

	return c.JSON(http.StatusOK, com)
}

func deleteComment(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	com := new(Comment)
	db.GetDB().Delete(&com, id)

	return c.NoContent(http.StatusNoContent)
}

func getAllComments(c echo.Context) error {
	var comments []Comment
	//FETCH ALL
	db.GetDB().Find(&comments)
	return c.JSON(http.StatusOK, comments)
}

func getEventComments(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var comments []Comment
	db.GetDB().Where("event_id <> ?", id).Find(&comments)
	return c.JSON(http.StatusOK, comments)
}
