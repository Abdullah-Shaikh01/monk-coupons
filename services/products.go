package services

import (
    "database/sql"
    "github.com/Abdullah-Shaikh01/monk-coupons/models"
)

func GetAllProducts(db *sql.DB) ([]models.Product, error) {
    rows, err := db.Query("SELECT id, name, price FROM products")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var products []models.Product
    for rows.Next() {
        var p models.Product
        if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
            return nil, err
        }
        products = append(products, p)
    }
    return products, nil
}
