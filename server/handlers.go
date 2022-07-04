package main

import (
	"chatting/model"
	"chatting/repository/room"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"time"
)

func EnterRoom(c echo.Context) error {

	param := c.Param("id")
	roomId, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "room id not appropriate",
		})
	}

	repository := room.GetMySqlRepository()
	_, err = repository.FindById(roomId)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"message": "room does not exist",
		})
	}

	http.ServeFile(c.Response().Writer, c.Request(), "home.html")
	return nil
}

type PostCreateRoomDto struct {
	Name    string `json:"name"`
	Private bool   `json:"private"`
}

func CreateRoom(c echo.Context) error {

	var requestBody PostCreateRoomDto
	err := c.Bind(requestBody)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "invalid request body",
		})
	}

	if len(requestBody.Name) > 50 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "room name is too long",
		})
	}

	newRoom := model.Room{
		Name:     requestBody.Name,
		Private:  requestBody.Private,
		CreateAt: time.Now().Truncate(time.Second),
	}
	repository := room.GetMySqlRepository()
	save, err := repository.Save(newRoom)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "invalid request body",
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"newRoomId": save.Id,
	})
}
