package handlers

import (
	"database/sql"
	"net/http"

	"github.com/Abdullah-Shaikh01/monk-coupons/models"
	"github.com/Abdullah-Shaikh01/monk-coupons/utils"
	"github.com/gin-gonic/gin"
)

func GetAllProducts(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query("SELECT id, name, price FROM products")
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "Failed to fetch products", err)
			return
		}
		defer rows.Close()

		var products []models.Product
		for rows.Next() {
			var p models.Product
			if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
				utils.Error(c, http.StatusInternalServerError, "Failed to parse product row", err)
				return
			}
			products = append(products, p)
		}

		utils.Success(c, http.StatusOK, "Products retrieved successfully", products)
	}
}
