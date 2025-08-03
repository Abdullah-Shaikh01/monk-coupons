package services

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/Abdullah-Shaikh01/monk-coupons/models"
)

func GetAllCoupons(db *sql.DB) ([]models.Coupon, error) {
    rows, err := db.Query(`
        SELECT 
            id, 
            type, 
            discount_value, 
            discount_type, 
            buyQuantity,
            getQuantity,
            repetition_threshold, 
            expiration_date, 
            product_id 
        FROM coupons
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var coupons []models.Coupon

    for rows.Next() {
        var coupon models.Coupon
        err := rows.Scan(
            &coupon.ID,
            &coupon.Type,
            &coupon.DiscountValue,
            &coupon.DiscountType,
            &coupon.BuyQuantity,
            &coupon.GetQuantity,
            &coupon.RepetitionThreshold,
            &coupon.ExpirationDate,
            &coupon.ProductID,
        )
        if err != nil {
            return nil, err
        }
        coupons = append(coupons, coupon)
    }

    return coupons, nil
}


func CreateCouponService(db *sql.DB, coupon models.Coupon, buyProducts []int, getProducts []int) (int64, string, error) {
    tx, err := db.Begin()
    if err != nil {
        return 0, "Unable to connect to db", err
    }
    defer tx.Rollback()

    var couponID int64

    if coupon.Type == "bxgy" {
        // Insert bxgy coupon
        insertCoupon := `
            INSERT INTO coupons (type, buyQuantity, getQuantity, repetition_threshold, expiration_date)
            VALUES (?, ?, ?, ?, ?)
        `
		fmt.Printf("service: firstBuyQty: %v\n", coupon.BuyQuantity)
		fmt.Printf("service: firstBuyQty: %v\n", *coupon.BuyQuantity)

        res, err := tx.Exec(insertCoupon, coupon.Type, *coupon.BuyQuantity, *coupon.GetQuantity, coupon.RepetitionThreshold, coupon.ExpirationDate)
        if err != nil {
            return 0, "Unable to insert coupon",  err
        }

        couponID, err = res.LastInsertId()
        if err != nil {
            return 0, "Unable to insert coupon",  err
        }

        // Insert buy products
        for _, p := range buyProducts {
            _, err := tx.Exec("INSERT INTO coupon_buy_products (coupon_id, product_id) VALUES (?, ?)", couponID, p)
            if err != nil {
                return 0, "Product with id " + strconv.Itoa(p) + " not present in products table", err
            }
        }

        // Insert get products
        for _, p := range getProducts {
            _, err := tx.Exec("INSERT INTO coupon_get_products (coupon_id, product_id) VALUES (?, ?)", couponID, p)
            if err != nil {
                return 0, "Product with id " + strconv.Itoa(p) + " not present in products table", err
            }
        }

    } else {
        // Insert non-bxgy coupon
        query := `
            INSERT INTO coupons (type, discount_value, discount_type, repetition_threshold, product_id, expiration_date)
            VALUES (?, ?, ?, ?, ?, ?)
        `
        res, err := tx.Exec(query, coupon.Type, coupon.DiscountValue, coupon.DiscountType, coupon.RepetitionThreshold, coupon.ProductID, coupon.ExpirationDate)
        if err != nil {
            return 0, "Unable to insert coupon",  err
        }

        couponID, err = res.LastInsertId()
        if err != nil {
            return 0, "Unable to insert coupon",  err
        }
    }

    if err := tx.Commit(); err != nil {
            return 0, "Unable to commit insert coupon transaction",  err
    }
    return couponID, "Added Coupon Successfully", nil
}
