package server

import "github.com/gin-gonic/gin"

const prefix = "/api"

var router *gin.Engine

func NewRouter() *gin.Engine {
	router = gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	api := router.Group(prefix)
	api.GET("/health", healthCheck)

	v1 := api.Group("/v1")
	v1.GET("/ping", ping)

	return router
}

func healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "OK",
	})
}

func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
