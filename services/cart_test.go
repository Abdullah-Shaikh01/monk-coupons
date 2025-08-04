package services

import (
    "regexp"
    "testing"
    "time"

    "github.com/DATA-DOG/go-sqlmock"
    "github.com/Abdullah-Shaikh01/monk-coupons/models"
)

func TestApplyProductWiseCoupon(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("error opening stub db: %s", err)
    }
    defer db.Close()

    discount := 20.0
    discountType := "fixed"
    productID := 5
    expiration := time.Now().AddDate(1, 0, 0)

    rows := sqlmock.NewRows([]string{"discount_value", "discount_type", "product_id", "expiration_date"}).
        AddRow(discount, discountType, productID, expiration)

    mock.ExpectQuery(regexp.QuoteMeta("SELECT discount_value, discount_type, product_id, expiration_date FROM coupons WHERE id = ?")).
        WithArgs(1).
        WillReturnRows(rows)

    cart := models.Cart{
        Items: []models.CartItem{
            {ProductID: 5, Quantity: 2, Price: 25},
            {ProductID: 3, Quantity: 1, Price: 10},
        },
    }

    updatedCart, totalDiscount, err := ApplyProductWiseCoupon(db, 1, cart)
    if err != nil {
        t.Errorf("unexpected error: %s", err)
    }

    expectedDiscount := 20.0 * 2 // fixed discount * quantity
    if totalDiscount != expectedDiscount {
        t.Errorf("expected total discount %v, got %v", expectedDiscount, totalDiscount)
    }

    // Check TotalDiscount field on qualifying item
    for _, item := range updatedCart.Items {
        if item.ProductID == 5 && item.TotalDiscount != expectedDiscount {
            t.Errorf("expected total_discount %v on product 5, got %v", expectedDiscount, item.TotalDiscount)
        }
        if item.ProductID == 3 && item.TotalDiscount != 0 {
            t.Errorf("expected total_discount 0 on product 3, got %v", item.TotalDiscount)
        }
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("there were unfulfilled db expectations: %s", err)
    }
}

func TestApplyCartWiseCoupon(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("error opening stub db: %s", err)
    }
    defer db.Close()

    discountValue := 10.0
    discountType := "percentage"
    repetitionThreshold := 100
    expiration := time.Now().AddDate(1, 0, 0)

    rows := sqlmock.NewRows([]string{"discount_value", "discount_type", "repetition_threshold", "expiration_date"}).
        AddRow(discountValue, discountType, repetitionThreshold, expiration)

    mock.ExpectQuery(regexp.QuoteMeta("SELECT discount_value, discount_type, repetition_threshold, expiration_date FROM coupons WHERE id = ?")).
        WithArgs(2).
        WillReturnRows(rows)

    cart := models.Cart{
        Items: []models.CartItem{
            {ProductID: 1, Quantity: 3, Price: 50}, // total = 150
            {ProductID: 2, Quantity: 1, Price: 20}, // total = 20
        },
    }

    updatedCart, totalDiscount, err := ApplyCartWiseCoupon(db, 2, cart)
    if err != nil {
        t.Errorf("unexpected error: %s", err)
    }

    totalPrice := 150 + 20
    expectedDiscount := float64(totalPrice) * discountValue / 100

    if totalDiscount != expectedDiscount {
        t.Errorf("expected discount %v, got %v", expectedDiscount, totalDiscount)
    }

    if len(updatedCart.Items) != 2 {
        t.Errorf("expected 2 items in cart, got %d", len(updatedCart.Items))
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("unfulfilled expectations: %s", err)
    }
}

func TestApplyCartWiseCoupon_ThresholdNotMet(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("error opening stub db: %s", err)
    }
    defer db.Close()

    discountValue := 10.0
    discountType := "percentage"
    repetitionThreshold := 200 // Higher than cart total
    expiration := time.Now().AddDate(1, 0, 0)

    rows := sqlmock.NewRows([]string{"discount_value", "discount_type", "repetition_threshold", "expiration_date"}).
        AddRow(discountValue, discountType, repetitionThreshold, expiration)

    mock.ExpectQuery(regexp.QuoteMeta("SELECT discount_value, discount_type, repetition_threshold, expiration_date FROM coupons WHERE id = ?")).
        WithArgs(3).
        WillReturnRows(rows)

    cart := models.Cart{
        Items: []models.CartItem{
            {ProductID: 1, Quantity: 3, Price: 50}, // total = 150
        },
    }

    _, _, err = ApplyCartWiseCoupon(db, 3, cart)
    if err == nil {
        t.Errorf("expected error due to threshold not met, got nil")
    }
    expectedErr := "cart total (150) is less than required threshold (200)"
    if err.Error() != expectedErr {
        t.Errorf("expected error '%s', got '%s'", expectedErr, err.Error())
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("unfulfilled expectations: %s", err)
    }
}

func TestApplyCartWiseCoupon_CouponNotFound(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("error opening stub db: %s", err)
    }
    defer db.Close()

    mock.ExpectQuery(regexp.QuoteMeta("SELECT discount_value, discount_type, repetition_threshold, expiration_date FROM coupons WHERE id = ?")).
        WithArgs(4).
        WillReturnError(sqlmock.ErrCancelled) // simulate db error / coupon not found

    cart := models.Cart{
        Items: []models.CartItem{
            {ProductID: 1, Quantity: 3, Price: 50},
        },
    }

    _, _, err = ApplyCartWiseCoupon(db, 4, cart)
    if err == nil {
        t.Errorf("expected error due to coupon not found/db error, got nil")
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("unfulfilled expectations: %s", err)
    }
}