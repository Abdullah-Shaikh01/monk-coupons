package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/Abdullah-Shaikh01/monk-coupons/models"
	"github.com/Abdullah-Shaikh01/monk-coupons/services"
	"github.com/Abdullah-Shaikh01/monk-coupons/utils"
	"github.com/gin-gonic/gin"
)

func GetAllCoupons(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		coupons, err := services.GetAllCoupons(db)
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "Failed to fetch coupons", err)
			return
		}
		utils.Success(c, http.StatusOK, "Coupons retrieved successfully", coupons)
	}
}

func CreateCoupon(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var payload map[string]interface{}

        if err := c.ShouldBindJSON(&payload); err != nil {
            utils.Error(c, http.StatusBadRequest, "Invalid JSON payload: ", err)
            return
        }

        couponType, ok := payload["type"].(string)
        if !ok || couponType == "" {
            utils.ErrorWithoutErr(c, http.StatusBadRequest, "Coupon type is required") 
            return
        }

        var coupon models.Coupon
        coupon.Type = couponType

        // Parse expiration date (fallback one month if missing)
        if expStr, ok := payload["expiration_date"].(string); ok && expStr != "" {
            exp, err := time.Parse(time.RFC3339, expStr)
            if err != nil {
                utils.Error(c, http.StatusBadRequest, "Invalid expiration_date format. Use RFC3339 (e.g., 2025-09-01T00:00:00Z)", err)
                return
            }
            coupon.ExpirationDate = exp
        } else {
            coupon.ExpirationDate = time.Now().AddDate(0, 1, 0)
        }

		var buyProducts, getProducts []int

        if couponType == "bxgy" {
			fmt.Printf("firstBuyQty: %v\n", payload)

            details, ok := payload["details"].(map[string]interface{})
            if !ok {
                utils.ErrorWithoutErr(c, http.StatusBadRequest, "Details field is required for bxgy coupons")
                return
            }
			fmt.Printf("firstBuyQty: %v\n", details)

            bRaw, ok := details["buy_products"].([]interface{})
            if !ok || len(bRaw) == 0 {
                utils.ErrorWithoutErr(c, http.StatusBadRequest, "buy_products array is required and must not be empty")
                return
            }
			fmt.Printf("buy_products array: %#v\n", bRaw)


            gRaw, ok := details["get_products"].([]interface{})
            if !ok || len(gRaw) == 0 {
                utils.ErrorWithoutErr(c, http.StatusBadRequest, "get_products array is required and must not be empty")
                return
            }

            // Extract quantities from first element
            firstBuyQty := int(bRaw[0].(map[string]interface{})["quantity"].(float64))
			fmt.Printf("firstBuyQty: %d\n", firstBuyQty)

            firstGetQty := int(gRaw[0].(map[string]interface{})["quantity"].(float64))

            // buyQty := int(firstBuy["quantity"].(float64))
            // getQty := int(firstGet["quantity"].(float64))
            coupon.BuyQuantity = &firstBuyQty
            coupon.GetQuantity = &firstGetQty

            if limit, ok := details["repition_limit"].(float64); ok {
                limitInt := int(limit)
                coupon.RepetitionThreshold = &limitInt
            } else {
				 utils.ErrorWithoutErr(c, http.StatusBadRequest, "repition_limit must be provided and be a number")
				return
			}

            // Map products for service func
            for _, item := range bRaw {
                p := item.(map[string]interface{})
                productID, ok := p["product_id"].(float64)
                if !ok {
                    utils.ErrorWithoutErr(c, http.StatusBadRequest, "Invalid product_id in buy_products")
                    return
                }
                buyProducts = append(buyProducts,  int(productID))
            }
            for _, item := range gRaw {
                p := item.(map[string]interface{})
                productID, ok := p["product_id"].(float64)
                if !ok {
                    utils.ErrorWithoutErr(c, http.StatusBadRequest, "Invalid product_id in get_products")
                    return
                }
                getProducts = append(getProducts,  int(productID))
            }	

        } else {
            // Non-bxgy fields from payload:
            if dv, ok := payload["discount_value"].(float64); ok {
                coupon.DiscountValue = &dv
            } else {
				utils.ErrorWithoutErr(c, http.StatusBadRequest, "No Discount value provided")
				return
			}
            if dt, ok := payload["discount_type"].(string); ok {
                coupon.DiscountType = dt
            } else {
                coupon.DiscountType = "percentage"
            }
			if couponType == "cart-wise" {
				if rt, ok := payload["threshold"].(float64); ok {
					rtInt := int(rt)
					coupon.RepetitionThreshold = &rtInt
				} else {
					utils.ErrorWithoutErr(c, http.StatusBadRequest, "Threshold is required for cart-wise coupons and be a number")
					return 
				}
			}
            if pid, ok := payload["product_id"].(float64); ok {
                pidInt := int(pid)
                coupon.ProductID = &pidInt
            }
        }

        // Call the service function
        couponID, message, err := services.CreateCouponService(db, coupon, buyProducts, getProducts)
        if err != nil {
            utils.Error(c, http.StatusInternalServerError, message, err)
            return
        }
        utils.Success(c, http.StatusCreated, "Coupon created successfully", map[string]int64{"coupon_id": couponID})
    }
}