package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Serve HTML at /
	r.GET("/", func(c *gin.Context) {
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
	})

	// create API endpoints group
	api := r.Group("/api")
	api.GET("/ping", Ping)
	api.POST("/status", UpdateStatus)
}
