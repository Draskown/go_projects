package handler

import (
	"github.com/Draskown/WBL0/server/service"
	"github.com/gin-gonic/gin"
)

// Handler structure to encapsulate service
type Handler struct {
	service *service.Service
}

// Creates a new handler by passing a point to the service
func NewHandler(s *service.Service) *Handler {
	return &Handler{service: s}
}

// Initialises routes of the application
//
// Returns a gin Engine as a start of the application
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	// Route /:id implies that after the / route
	// any id can follow to display order's information
	router.GET("/:id", h.showOrder)
	
	// DEBUG
	router.GET("/test", h.testGetDB)
	router.POST("/test", h.testPostDB)

	return router
}
