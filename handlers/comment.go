package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	db "comment-service/dbConnector"
	m "comment-service/models"

	"github.com/labstack/echo/v4"
)

func CreateComment(c echo.Context) (err error) {

	com := m.NewComment()
	//Bind request data to Comment structure
	if err := c.Bind(com); err != nil {
		return err
	}

	//Validate data
	if err = c.Validate(com); err != nil {
		return err
	}

	//STORE in DB
	db.GetDB().Create(com)

	return c.JSON(http.StatusCreated, com)
}

func GetComment(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	com := &m.Comment{}

	//GET IN DB by id
	db.GetDB().First(&com, id)

	return c.JSON(http.StatusOK, com)
}

func UpdateComment(c echo.Context) error {
	com := new(m.Comment)
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

func DeleteComment(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	com := new(m.Comment)
	db.GetDB().Delete(&com, id)

	return c.NoContent(http.StatusNoContent)
}

func GetAllComments(c echo.Context) error {
	var comments []m.Comment
	//FETCH ALL
	db.GetDB().Find(&comments)
	return c.JSON(http.StatusOK, comments)
}

func GetEventComments(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var comments []m.Comment
	db.GetDB().Where("event_id <> ?", id).Find(&comments)
	return c.JSON(http.StatusOK, comments)
}
