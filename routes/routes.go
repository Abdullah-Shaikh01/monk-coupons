package routes

import (
	"database/sql"

	"github.com/Abdullah-Shaikh01/monk-coupons/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, db *sql.DB) {
	// Product routes
	router.GET("/products", handlers.GetAllProducts(db))

	// Add other routes here (e.g., /coupons etc.)
}