package handlers

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/DATA-DOG/go-sqlmock"
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

func TestGetAllProductsHandler(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
    }
    defer db.Close()

    rows := sqlmock.NewRows([]string{"id", "name", "price"}).
        AddRow(1, "Laptop", 75000.00).
        AddRow(2, "Mouse", 800.00)

    mock.ExpectQuery("SELECT id, name, price FROM products").WillReturnRows(rows)

    router := gin.Default()
    router.GET("/products", GetAllProducts(db))

    req, _ := http.NewRequest("GET", "/products", nil)
    resp := httptest.NewRecorder()

    router.ServeHTTP(resp, req)

    assert.Equal(t, 200, resp.Code)
    assert.Contains(t, resp.Body.String(), "Laptop")
    assert.Contains(t, resp.Body.String(), "Mouse")

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("there were unfulfilled expectations: %s", err)
    }
}