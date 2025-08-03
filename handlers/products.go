package handlers

import (
	"database/sql"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/Abdullah-Shaikh01/monk-coupons/models"
)

func GetAllProducts(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query("SELECT id, name, price FROM products")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
			return
		}
		defer rows.Close()

		var products []models.Product
		for rows.Next() {
			var p models.Product
			if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse products"})
				return
			}
			products = append(products, p)
		}

		c.JSON(http.StatusOK, products)
	}
}
