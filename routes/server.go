package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupRouter() {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	router.NoRoute(func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNotFound)
	})

	authRoute(router.Group("/auth"))
	groupRoute(router.Group("/"))

	router.Run()
}
