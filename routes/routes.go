package routes

import (
	"net/http"

	"ticket-system/handlers"
	"ticket-system/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter registers all routes for the ticket system
func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Public routes
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", handlers.Register)
		authGroup.POST("/login", handlers.Login)
	}

	// Protected routes (require JWT authentication middleware)
	ticketGroup := router.Group("/tickets")
	ticketGroup.Use(middleware.AuthMiddleware())
	{
		ticketGroup.POST("", handlers.CreateTicket)
		ticketGroup.GET("", handlers.ListTickets)
		ticketGroup.GET("/:id", handlers.GetTicket)
		ticketGroup.PATCH("/:id/status", handlers.UpdateTicketStatus)
	}

	return router
}
