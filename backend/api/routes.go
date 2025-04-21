package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type API struct {
	DB *gorm.DB
}

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	api := &API{DB: db}

	// Serve HTML at /
	r.GET("/", indexPage)

	// create API endpoints group
	apiRoutes := r.Group("/api")
	apiRoutes.GET("/ping", api.Ping)
	apiRoutes.GET("/devices", api.ListDevices)
	apiRoutes.POST("/status", api.UpdateStatus)
}

func indexPage(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>MQTT Status</title>
			<style>
				body { font-family: sans-serif; text-align: center; padding-top: 50px; background: #f4f4f4; }
				h1 { color: #2c3e50; }
			</style>
		</head>
		<body>
			<h1>✅ MQTT Server is Running</h1>
			<p>Powered by Go + Gin + Fly.io 🚀</p>
		</body>
		</html>
	`))
}
