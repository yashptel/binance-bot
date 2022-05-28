package server

import "github.com/gin-gonic/gin"

var router *gin.Engine

func NewRouter() *gin.Engine {
	router = gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	

	return router
}
