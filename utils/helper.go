package utils

import "database/sql"

func GetCouponTypeByID(db *sql.DB, couponID int) (string, error) {
	var couponType string
	if err := db.QueryRow("SELECT type FROM coupons WHERE id = ?", couponID).Scan(&couponType); err != nil {
		return couponType, err
	}
	return couponType, nil
}

func GetBxGyProducts(db *sql.DB, couponID int) ([]int, []int, error) {
	buyRows, err := db.Query("SELECT product_id FROM coupon_buy_products WHERE coupon_id = ?", couponID)
	if err != nil {
		return nil, nil, err
	}
	defer buyRows.Close()

	var buyProducts []int
	for buyRows.Next() {
		var pid int
		if err := buyRows.Scan(&pid); err != nil {
			return nil, nil, err
		}
		buyProducts = append(buyProducts, pid)
	}

	getRows, err := db.Query("SELECT product_id FROM coupon_get_products WHERE coupon_id = ?", couponID)
	if err != nil {
		return nil, nil, err
	}
	defer getRows.Close()

	var getProducts []int
	for getRows.Next() {
		var pid int
		if err := getRows.Scan(&pid); err != nil {
			return nil, nil, err
		}
		getProducts = append(getProducts, pid)
	}

	return buyProducts, getProducts, nil
}