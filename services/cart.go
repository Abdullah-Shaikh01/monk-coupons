package services

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Abdullah-Shaikh01/monk-coupons/models"
	"github.com/Abdullah-Shaikh01/monk-coupons/utils"
)



func ApplyCartWiseCoupon(db *sql.DB, couponID int, cart models.Cart) (models.Cart, float64, error) {
    var discountValue float64
    var discountType string
    var repetitionThreshold int
    var expirationDate time.Time

    err := db.QueryRow(`SELECT discount_value, discount_type, repetition_threshold, expiration_date FROM coupons WHERE id = ?`, couponID).
        Scan(&discountValue, &discountType, &repetitionThreshold, &expirationDate)
    if err != nil {
        return cart, 0, err
    }

    if time.Now().After(expirationDate) {
        return cart, 0, errors.New("coupon expired")
    }

    var totalPrice float64
    for _, item := range cart.Items {
        totalPrice += float64(item.Quantity) * item.Price
    }
    if int(totalPrice) < repetitionThreshold {
        return cart, 0, fmt.Errorf("cart total (%v) is less than required threshold (%d)", totalPrice, repetitionThreshold)
    }

    var discount float64
    if discountType == "percentage" {
        discount = totalPrice * discountValue / 100
    } else {
        discount = discountValue
    }

    // No per item discount breakdown here; total discount only
    return cart, discount, nil
}

func ApplyProductWiseCoupon(db *sql.DB, couponID int, cart models.Cart) (models.Cart, float64, error) {
    var discountValue float64
    var discountType string
    var productID int
    var expirationDate time.Time

    err := db.QueryRow(`SELECT discount_value, discount_type, product_id, expiration_date FROM coupons WHERE id = ?`, couponID).
        Scan(&discountValue, &discountType, &productID, &expirationDate)
    if err != nil {
        return cart, 0, err
    }

    if time.Now().After(expirationDate) {
        return cart, 0, errors.New("coupon expired")
    }

    var totalDiscount float64
    for i, item := range cart.Items {
        if item.ProductID == productID {	
            var itemDiscount float64
            if discountType == "percentage" {
                itemDiscount = item.Price * float64(item.Quantity) * discountValue / 100
            } else {
                itemDiscount = discountValue * float64(item.Quantity)
            }
            totalDiscount += itemDiscount
            cart.Items[i].TotalDiscount = itemDiscount
        } else {
            cart.Items[i].TotalDiscount = 0
        }
    }
    return cart, totalDiscount, nil
}

func ApplyBxGyCoupon(db *sql.DB, couponID int, cart models.Cart) (models.Cart, float64, error) {
    var buyQuantity, getQuantity, repetitionLimit int
    var expirationDate time.Time

    err := db.QueryRow(`SELECT buyQuantity, getQuantity, repetition_threshold, expiration_date FROM coupons WHERE id = ?`, couponID).
        Scan(&buyQuantity, &getQuantity, &repetitionLimit, &expirationDate)
    if err != nil {
        return cart, 0, err
    }
    if time.Now().After(expirationDate) {
        return cart, 0, errors.New("coupon expired")
    }

    buyProducts, getProducts, err := utils.GetBxGyProducts(db, couponID)
    if err != nil {
        return cart, 0, err
    }

    // Map productID to quantity for quick lookup
    quantityMap := make(map[int]int)
    for _, item := range cart.Items {
        quantityMap[item.ProductID] += item.Quantity
    }

    // Calculate total buy quantity in cart
    var buyCount int
    for _, pid := range buyProducts {
        buyCount += quantityMap[pid]
    }
    if buyCount < buyQuantity {
        return cart, 0, errors.New("insufficient buy products quantity for coupon")
    }

    eligibleTimes := buyCount / buyQuantity
    if eligibleTimes > repetitionLimit {
        eligibleTimes = repetitionLimit
    }
    getEligible := eligibleTimes * getQuantity

    totalDiscount := 0.0
    // Add discount for get products in cart
    for i := range cart.Items {
        for _, gpid := range getProducts {
            if cart.Items[i].ProductID == gpid && getEligible > 0 {
                discountUnits := cart.Items[i].Quantity
                if discountUnits > getEligible {
                    discountUnits = getEligible
                }
                discount := cart.Items[i].Price * float64(discountUnits)
                cart.Items[i].TotalDiscount = discount
                totalDiscount += discount
                getEligible -= discountUnits
            }
        }
    }

    return cart, totalDiscount, nil
}
