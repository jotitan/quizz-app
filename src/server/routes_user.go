package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func createUserRoutes(server *gin.Engine) {
	api := server.Group("/api/user")
	api.GET("/is_admin", isAdminRoute)
}

func isAdminRoute(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"is_admin": securityManager.checkAdmin(token),
	})
}
