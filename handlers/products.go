package handlers

import (
    "database/sql"
    "net/http"

    "github.com/Abdullah-Shaikh01/monk-coupons/services"
    "github.com/Abdullah-Shaikh01/monk-coupons/utils"
    "github.com/gin-gonic/gin"
)

func GetAllProducts(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        products, err := services.GetAllProducts(db)
        if err != nil {
            utils.Error(c, http.StatusInternalServerError, "Failed to fetch products", err)
            return
        }
        utils.Success(c, http.StatusOK, "Products retrieved successfully", products)
    }
}
