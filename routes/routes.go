package routes

import (
	"database/sql"

	"github.com/Abdullah-Shaikh01/monk-coupons/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, db *sql.DB) {
	router.GET("/products", handlers.GetAllProducts(db))
	router.GET("/coupons", handlers.GetAllCoupons(db)) 
	router.POST("/coupons", handlers.CreateCoupon(db))

}
