package api

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")

	api.GET("/ping", Ping)
	api.GET("/status", GetStatus)
	api.GET("/weather_info", GetWeather)
}
