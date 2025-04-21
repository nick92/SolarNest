package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nick92/solarnest/db/schemas"
	"github.com/nick92/solarnest/weather"
)

func (api *API) Ping(c *gin.Context) {
	StatusOkResponse(c, gin.H{
		"message": "pong",
	})
}

func (api *API) UpdateStatus(c *gin.Context) {
	status := c.Request.Body

	if err := c.ShouldBindJSON(&status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	StatusOkResponse(c, status)
}

func (api *API) GetWeather(c *gin.Context) {
	weather := weather.GetMockWeather()
	StatusOkResponse(c, weather)
}

func (api *API) ListDevices(c *gin.Context) {
	var device schemas.Device
	c.JSON(http.StatusOK, device.GetDevices(api.DB))
}
