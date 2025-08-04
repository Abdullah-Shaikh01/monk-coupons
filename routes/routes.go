package routes

import (
	"database/sql"

	"github.com/Abdullah-Shaikh01/monk-coupons/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, db *sql.DB) {
	router.GET("/products", handlers.GetAllProducts(db))
	router.GET("/coupons", handlers.GetAllCoupons(db))
	router.GET("/coupons/:id", handlers.GetCouponByID(db))
	router.POST("/coupons", handlers.CreateCoupon(db))
	router.PUT("/coupons/:id", handlers.UpdateCoupon(db))
	router.DELETE("/coupons/:id", handlers.DeleteCoupon(db))
	router.POST("/apply-coupon/:id", handlers.ApplyCouponByID(db))
	router.POST("/applicable-coupons", handlers.GetApplicableCouponsHandler(db))

}
