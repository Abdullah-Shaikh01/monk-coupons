package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/Abdullah-Shaikh01/monk-coupons/models"
	"github.com/Abdullah-Shaikh01/monk-coupons/services"
	"github.com/Abdullah-Shaikh01/monk-coupons/utils"
	"github.com/gin-gonic/gin"
)

func ApplyCouponByID(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        idParam := c.Param("id")
        couponID, err := strconv.Atoi(idParam)
        if err != nil {
            utils.ErrorWithoutErr(c, http.StatusBadRequest, "Invalid coupon ID")
            return
        }

		couponType, err := utils.GetCouponTypeByID(db, couponID)
        if err != nil {
            if err == sql.ErrNoRows {
                utils.ErrorWithoutErr(c, http.StatusNotFound, "Coupon not found")
            } else {
                utils.Error(c, http.StatusInternalServerError, "Failed to fetch coupon", err)
            }
            return
        }

        var input struct {
            Cart models.Cart `json:"cart"`
        }
        if err := c.ShouldBindJSON(&input); err != nil {
            utils.Error(c, http.StatusBadRequest, "Invalid JSON payload", err)
            return
        }

        var updatedCart models.Cart
        var totalDiscount float64

        switch couponType {
        case "cart-wise":
            updatedCart, totalDiscount, err = services.ApplyCartWiseCoupon(db, couponID, input.Cart)
        case "product-wise":
            updatedCart, totalDiscount, err = services.ApplyProductWiseCoupon(db, couponID, input.Cart)
        case "bxgy":
            updatedCart, totalDiscount, err = services.ApplyBxGyCoupon(db, couponID, input.Cart)
        default:
            utils.ErrorWithoutErr(c, http.StatusBadRequest, "Unsupported coupon type")
            return
        }

        if err != nil {
            utils.ErrorWithoutErr(c, http.StatusBadRequest, err.Error())
            return
        }

        // Calculate totals
        var totalPrice float64
        for _, item := range updatedCart.Items {
            totalPrice += float64(item.Quantity)*item.Price
        }
        finalPrice := totalPrice - totalDiscount

        c.JSON(http.StatusOK, gin.H{
            "updated_cart": gin.H{
                "items":         updatedCart.Items,
                "total_price":   totalPrice,
                "total_discount": totalDiscount,
                "final_price":   finalPrice,
            },
        })
    }
}

func GetApplicableCouponsHandler(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var input struct {
            Cart models.Cart `json:"cart"`
        }
        if err := c.ShouldBindJSON(&input); err != nil {
            utils.Error(c, http.StatusBadRequest, "Invalid JSON payload", err)
            return
        }

        applicableCoupons, err := services.GetApplicableCoupons(db, input.Cart)
        if err != nil {
            utils.Error(c, http.StatusInternalServerError, "Failed to fetch applicable coupons", err)
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "applicable_coupons": applicableCoupons,
        })
    }
}
