package services

import (
    "regexp"
    "testing"

    "github.com/DATA-DOG/go-sqlmock"
)

func TestGetAllProducts(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
    }
    defer db.Close()

    rows := sqlmock.NewRows([]string{"id", "name", "price"}).
        AddRow(1, "Laptop", 75000.00).
        AddRow(2, "Mouse", 800.00)

    mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, price FROM products")).
        WillReturnRows(rows)

    products, err := GetAllProducts(db)
    if err != nil {
        t.Errorf("unexpected error: %s", err)
    }

    if len(products) != 2 {
        t.Errorf("expected 2 products, got %d", len(products))
    }

    if products[0].Name != "Laptop" || products[1].Price != 800.00 {
        t.Errorf("unexpected product data")
    }

    // Ensure all expectations were met
    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("there were unfulfilled expectations: %s", err)
    }
}
