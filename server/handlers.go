package main

import (
	"chatting/repository/room"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
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
