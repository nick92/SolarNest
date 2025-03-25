package api

import (
	"github.com/gin-gonic/gin"
	"github.com/nick92/solarnest/sensors"
	"github.com/nick92/solarnest/weather"
)

func Ping(c *gin.Context) {
	StatusOkResponse(c, gin.H{
		"message": "pong",
	})
}

func GetStatus(c *gin.Context) {
	status := sensors.GetStatus()
	StatusOkResponse(c, status)
}

func GetWeather(c *gin.Context) {
	weather := weather.GetMockWeather()
	StatusOkResponse(c, weather)
}
