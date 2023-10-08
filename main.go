package main

import (
	"net/http"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	db "comment-service/dbConnector"
	"github.com/go-playground/validator/v10"
	h "comment-service/handlers"
	m "comment-service/models"
)

type CustomValidator struct {
    validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
	  // Optionally, you could return the error to give each route more control over the status code
	  return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func main() {
	e := echo.New()

	//Define validator to validate data in handler
	e.Validator = &CustomValidator{validator: validator.New()}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/comments", h.GetAllComments)
	e.POST("/comments", h.CreateComment)
	e.GET("/comments/:id", h.GetComment)
	e.PUT("/comments/:id", h.UpdateComment)
	e.DELETE("/comments/:id", h.DeleteComment)
	e.GET("/events/:id/comments", h.GetEventComments)

	//Run migration
	db.GetDB().AutoMigrate(&m.Comment{})

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
