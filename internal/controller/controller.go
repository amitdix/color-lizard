package controller

import (
	"git.target.com/StoreDataMovement/color-lizard/config"
	"github.com/gin-gonic/gin"
	"net/http"
)



func GetRouter(endpoints []config.Endpoint, ready *bool) (r *gin.Engine) {
	gin.SetMode(gin.ReleaseMode)
	r = gin.Default()

	r.GET("/health", func(c *gin.Context) {
		if !*ready {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"healthy": false,
				"cause":   "not ready yet",
			})
			return
		}

		// Add other checks here as necessary
		c.JSON(http.StatusOK, gin.H{
			"healthy": true,
		})
	})
	r.GET("/ready", func(c *gin.Context) {
		status := http.StatusServiceUnavailable
		if *ready {
			status = http.StatusOK
		}
		c.JSON(status, gin.H{
			"ready": ready,
		})
	})

	r.GET("/colorlizard/*path", func(context *gin.Context) {
		first := context.Param("path")
		path := first
		for _, endpoint := range endpoints{
			if path == endpoint.Path && endpoint.Method == "GET" {
				for key, value := range endpoint.Headers{
					context.Header(key, value)
				}
				context.Data(endpoint.Status, "application/json; charset=utf-8", []byte(endpoint.Response))
				break
			}
		}
		context.JSON(http.StatusNotFound, "application/json; charset=utf-8")
	})

	r.POST("/colorlizard/*path", func(context *gin.Context) {
		first := context.Param("path")
		path := first
		for _, endpoint := range endpoints{
			if path == endpoint.Path && endpoint.Method == "POST" {
				for key, value := range endpoint.Headers{
					context.Header(key, value)
				}
				context.Data(endpoint.Status, "application/json; charset=utf-8", []byte(endpoint.Response))
				break
			}
		}
		context.JSON(http.StatusNotFound, "application/json; charset=utf-8")
	})

	r.PUT("/colorlizard/*path", func(context *gin.Context) {
		first := context.Param("path")
		path := first
		for _, endpoint := range endpoints{
			if path == endpoint.Path && endpoint.Method == "PUT" {
				for key, value := range endpoint.Headers{
					context.Header(key, value)
				}
				context.Data(endpoint.Status, "application/json; charset=utf-8", []byte(endpoint.Response))
				break
			}
		}
		context.JSON(http.StatusNotFound, "application/json; charset=utf-8")
	})


	r.DELETE("/colorlizard/*path", func(context *gin.Context) {
		first := context.Param("path")
		path := first
		for _, endpoint := range endpoints{
			if path == endpoint.Path && endpoint.Method == "DELETE" {
				for key, value := range endpoint.Headers{
					context.Header(key, value)
				}
				context.Data(endpoint.Status, "application/json; charset=utf-8", []byte(endpoint.Response))
				break
			}
		}
		context.JSON(http.StatusNotFound, "application/json; charset=utf-8")
	})

	return r
}

