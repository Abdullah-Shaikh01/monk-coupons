package services

import (
    "regexp"
    "strings"
    "testing"
    "time"

    "github.com/DATA-DOG/go-sqlmock"
    "github.com/Abdullah-Shaikh01/monk-coupons/models"
)

func TestGetAllCoupons(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("failed to open mock db: %s", err)
    }
    defer db.Close()

    rows := sqlmock.NewRows([]string{
        "id", "type", "discount_value", "discount_type", "buyQuantity", "getQuantity", "repetition_threshold", "expiration_date", "product_id",
    }).
        AddRow(1, "cart-wise", 10.0, "percentage", nil, nil, 100, time.Now(), nil).
        AddRow(2, "product-wise", 20.0, "fixed", nil, nil, nil, time.Now(), 5)

    mock.ExpectQuery(regexp.QuoteMeta("SELECT id, type, discount_value, discount_type, buyQuantity, getQuantity, repetition_threshold, expiration_date, product_id FROM coupons")).
        WillReturnRows(rows)

    coupons, err := GetAllCoupons(db)
    if err != nil {
        t.Errorf("unexpected error: %s", err)
    }
    if len(coupons) != 2 {
        t.Errorf("expected 2 coupons, got %d", len(coupons))
    }
    if coupons[0].Type != "cart-wise" || coupons[1].ProductID == nil && *coupons[1].ProductID != 5 {
        t.Errorf("unexpected coupon data")
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("unfulfilled expectations: %s", err)
    }
}

func TestGetCouponByID(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("failed to open mock db: %s", err)
    }
    defer db.Close()

    couponID := "1"
    expirationDate := time.Now()

    row := sqlmock.NewRows([]string{
        "id", "type", "discount_value", "discount_type", "buyQuantity", "getQuantity", "repetition_threshold", "expiration_date", "product_id",
    }).AddRow(1, "cart-wise", 10.0, "percentage", nil, nil, 100, expirationDate, nil)

    mock.ExpectQuery(regexp.QuoteMeta("SELECT id, type, discount_value, discount_type, buyQuantity, getQuantity, repetition_threshold, expiration_date, product_id FROM coupons WHERE id = ?")).
        WithArgs(couponID).
        WillReturnRows(row)

    coupon, err := GetCouponByID(db, couponID)
    if err != nil {
        t.Errorf("unexpected error: %s", err)
    }
    if coupon.ID != 1 || coupon.Type != "cart-wise" {
        t.Errorf("unexpected coupon data")
    }
    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("unfulfilled expectations: %s", err)
    }
}

func TestCreateCouponService_NonBxGy(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("failed to open mock db: %s", err)
    }
    defer db.Close()

    coupon := models.Coupon{
        Type:               "cart-wise",
        DiscountValue:      ptrFloat64(10.0),
        DiscountType:       "percentage",
        RepetitionThreshold: ptrInt(100),
        ProductID:          nil,
        ExpirationDate:     time.Now(),
    }

    mock.ExpectBegin()
    mock.ExpectExec(regexp.QuoteMeta("INSERT INTO coupons (type, discount_value, discount_type, repetition_threshold, product_id, expiration_date) VALUES (?, ?, ?, ?, ?, ?)")).
        WithArgs(coupon.Type, coupon.DiscountValue, coupon.DiscountType, coupon.RepetitionThreshold, coupon.ProductID, coupon.ExpirationDate).
        WillReturnResult(sqlmock.NewResult(1, 1))
    mock.ExpectCommit()

    id, msg, err := CreateCouponService(db, coupon, nil, nil)
    if err != nil {
        t.Errorf("unexpected error: %s", err)
    }
    if id != 1 {
        t.Errorf("expected inserted id 1, got %d", id)
    }
    if !strings.Contains(msg, "Successfully") {
        t.Errorf("unexpected message: %s", msg)
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("unfulfilled expectations: %s", err)
    }
}

func TestCreateCouponService_BxGy(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("failed to open mock db: %s", err)
    }
    defer db.Close()

    coupon := models.Coupon{
        Type:               "bxgy",
        BuyQuantity:        ptrInt(2),
        GetQuantity:        ptrInt(1),
        RepetitionThreshold: ptrInt(3),
        ExpirationDate:     time.Now(),
    }
    buyProducts := []int{1, 2}
    getProducts := []int{3}

    mock.ExpectBegin()
    mock.ExpectExec(regexp.QuoteMeta("INSERT INTO coupons (type, buyQuantity, getQuantity, repetition_threshold, expiration_date) VALUES (?, ?, ?, ?, ?)")).
        WithArgs(coupon.Type, *coupon.BuyQuantity, *coupon.GetQuantity, coupon.RepetitionThreshold, coupon.ExpirationDate).
        WillReturnResult(sqlmock.NewResult(5,1))
    mock.ExpectExec(regexp.QuoteMeta("INSERT INTO coupon_buy_products (coupon_id, product_id) VALUES (?, ?)")).
        WithArgs(5, 1).
        WillReturnResult(sqlmock.NewResult(1,1))
    mock.ExpectExec(regexp.QuoteMeta("INSERT INTO coupon_buy_products (coupon_id, product_id) VALUES (?, ?)")).
        WithArgs(5, 2).
        WillReturnResult(sqlmock.NewResult(2,1))
    mock.ExpectExec(regexp.QuoteMeta("INSERT INTO coupon_get_products (coupon_id, product_id) VALUES (?, ?)")).
        WithArgs(5, 3).
        WillReturnResult(sqlmock.NewResult(3,1))
    mock.ExpectCommit()

    id, msg, err := CreateCouponService(db, coupon, buyProducts, getProducts)
    if err != nil {
        t.Errorf("unexpected error: %s", err)
    }
    if id != 5 {
        t.Errorf("expected inserted id 5, got %d", id)
    }
    if !strings.Contains(msg, "Successfully") {
        t.Errorf("unexpected message: %s", msg)
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("unfulfilled expectations: %s", err)
    }
}

func TestUpdateCouponService_Success(t *testing.T){
    db, mock, err:= sqlmock.New()
    if err != nil{
        t.Fatalf("failed to open mock db: %s", err)
    }
    defer db.Close()

    updates := map[string]interface{}{
        "discount_value": 15,
        "repetition_threshold": 90,
    }
    couponType := "cart-wise"
    couponID := 10

    mock.ExpectExec(regexp.QuoteMeta("UPDATE coupons SET discount_value = ?, repetition_threshold = ? WHERE id = ?")).
    WithArgs(15, 90, couponID).
    WillReturnResult(sqlmock.NewResult(0,1))

    err = UpdateCouponService(db, couponID, couponType, updates)
    if err != nil{
        t.Errorf("unexpected error: %s", err)
    }

    if err := mock.ExpectationsWereMet(); err != nil{
        t.Errorf("unfulfilled expectations: %s", err)
    }
}

func TestDeleteCouponService_Success(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("failed to open mock db: %s", err)
    }
    defer db.Close()

    couponID := 5
    mock.ExpectExec(regexp.QuoteMeta("DELETE FROM coupons WHERE id = ?")).
        WithArgs(couponID).
        WillReturnResult(sqlmock.NewResult(1, 1))

    err = DeleteCouponService(db, couponID)
    if err != nil {
        t.Errorf("unexpected error: %s", err)
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("unfulfilled expectations: %s", err)
    }
}

// Helper functions to get pointers to basic types
func ptrFloat64(f float64) *float64 { return &f }
func ptrInt(i int) *int             { return &i }
