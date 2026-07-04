package handlers

import (
	"net/http"
	"strconv"

	"ticket-system/config"
	"ticket-system/models"

	"github.com/gin-gonic/gin"
)

type CreateTicketInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type UpdateTicketStatusInput struct {
	Status string `json:"status" binding:"required"`
}

// CreateTicket handles ticket creation
func CreateTicket(c *gin.Context) {
	// Extract userID from context (set by middleware)
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, ok := userIDVal.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var input CreateTicketInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ticket := models.Ticket{
		Title:       input.Title,
		Description: input.Description,
		Status:      "open",
		UserID:      userID,
	}

	if err := config.DB.Create(&ticket).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create ticket"})
		return
	}

	c.JSON(http.StatusCreated, ticket)
}

// ListTickets retrieves all tickets belonging to the logged-in user
func ListTickets(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, ok := userIDVal.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var tickets []models.Ticket
	if err := config.DB.Where("user_id = ?", userID).Find(&tickets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tickets"})
		return
	}

	c.JSON(http.StatusOK, tickets)
}

// GetTicket retrieves a single ticket by ID, validating ownership
func GetTicket(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, ok := userIDVal.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ticketIDStr := c.Param("id")
	ticketID, err := strconv.ParseUint(ticketIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	var ticket models.Ticket
	if err := config.DB.First(&ticket, ticketID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	// Validate ownership
	if ticket.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to view this ticket"})
		return
	}

	c.JSON(http.StatusOK, ticket)
}

// UpdateTicketStatus updates the status of a ticket, enforcing valid state flows and ownership
func UpdateTicketStatus(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, ok := userIDVal.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ticketIDStr := c.Param("id")
	ticketID, err := strconv.ParseUint(ticketIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	var input UpdateTicketStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate allowed status names
	newStatus := input.Status
	if newStatus != "open" && newStatus != "in_progress" && newStatus != "closed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status. Supported statuses: open, in_progress, closed"})
		return
	}

	var ticket models.Ticket
	if err := config.DB.First(&ticket, ticketID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	// Validate ownership
	if ticket.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to update this ticket"})
		return
	}

	// Validate transition
	if !isValidTransition(ticket.Status, newStatus) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid status transition. Flow must follow: open -> in_progress -> closed, and closed tickets cannot be reopened.",
		})
		return
	}

	ticket.Status = newStatus
	if err := config.DB.Save(&ticket).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update ticket status"})
		return
	}

	c.JSON(http.StatusOK, ticket)
}

// isValidTransition helper checks if a status change is legal
func isValidTransition(current, next string) bool {
	if current == next {
		return true
	}
	switch current {
	case "open":
		return next == "in_progress" || next == "closed"
	case "in_progress":
		return next == "closed"
	case "closed":
		// Closed tickets cannot transition to any other status
		return false
	}
	return false
}
