package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nick92/solarnest/weather"
)

func Ping(c *gin.Context) {
	StatusOkResponse(c, gin.H{
		"message": "pong",
	})
}

func UpdateStatus(c *gin.Context) {
	status := c.Request.Body

	if err := c.ShouldBindJSON(&status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	StatusOkResponse(c, status)
}

func GetWeather(c *gin.Context) {
	weather := weather.GetMockWeather()
	StatusOkResponse(c, weather)
}
