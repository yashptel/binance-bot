package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func RunHttpServer() {
	router := gin.New()
	fmt.Println("Starting server on port 8080")

	router.Run(":8080")
}
