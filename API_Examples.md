
# API Examples

## GET /products

**Request:**  
`GET http://localhost:8080/products`

**Response:**
```json
{
  "data": [
    { "id": 1, "name": "Laptop", "price": 75000 },
    { "id": 2, "name": "Smartphone", "price": 30000 },
    ...
  ],
  "message": "Products retrieved successfully",
  "status": 200
}
```

## GET /coupons

**Request:**  
`GET http://localhost:8080/coupons`

**Response:**
```json
{
  "data": [
    {
      "id": 1,
      "type": "cart-wise",
      "discount_value": 15,
      ...
    },
    ...
  ],
  "message": "Coupons retrieved successfully",
  "status": 200
}
```

## GET /coupons/:id

**Valid Request:**  
`GET http://localhost:8080/coupons/1`

**Response:**
```json
{
  "data": {
    "id": 1,
    "type": "cart-wise",
    "discount_value": 15,
    ...
  },
  "message": "Coupon retrieved successfully",
  "status": 200
}
```

**Invalid Request:**  
`GET http://localhost:8080/coupons/100`

**Response:**
```json
{
  "message": "Coupon not found",
  "status": 404
}
```

## POST /coupons

**Valid Request:**
```json
{
  "type": "cart-wise",
  "discount_value": 10,
  "discount_type": "percentage",
  "repetition_threshold": 100,
  "expiration_date": "2025-09-01T00:00:00Z"
}
```

**Response:**
```json
{
  "data": { "coupon_id": 12 },
  "message": "Coupon created successfully",
  "status": 201
}
```

**Invalid Request:**
```json
{
  "discount_value": 5,
  "discount_type": "percentage",
  "threshold": 500,
  "expiration_date": "2025-09-01T00:00:00Z"
}
```

**Response:**
```json
{
  "message": "Coupon type is required",
  "status": 400
}
```

**BxGy Coupon Request:**
```json
{
  "type": "bxgy",
  "expiration_date": "2025-09-01T00:00:00Z",
  "details": {
    "buy_products": [
      { "product_id": 1, "quantity": 3 },
      { "product_id": 2, "quantity": 3 }
    ],
    "get_products": [
      { "product_id": 3, "quantity": 1 }
    ],
    "repition_limit": 2
  }
}
```

**Response:**
```json
{
  "data": { "coupon_id": 15 },
  "message": "Coupon created successfully",
  "status": 201
}
```

## PUT /coupons/:id

**Update Cart-wise Coupon:**
```json
{
  "discount_value": 15.0,
  "repetition_threshold": 500
}
```

**Response:**
```json
{
  "data": null,
  "message": "Coupon updated successfully",
  "status": 200
}
```

**Update BxGy Coupon:**
```json
{
  "buyQuantity": 3,
  "getQuantity": 1
}
```

**Response:**
```json
{
  "data": null,
  "message": "Coupon updated successfully",
  "status": 200
}
```

## DELETE /coupons/:id

**Request:**  
`DELETE http://localhost:8080/coupons/5`

**Response:**
```json
{
  "data": null,
  "message": "Coupon deleted successfully",
  "status": 200
}
```

## POST /apply-coupon/:id

**Valid Product-wise Coupon Request:**
```json
{
  "cart": {
    "items": [
      { "product_id": 1, "quantity": 6, "price": 50 },
      { "product_id": 2, "quantity": 3, "price": 30 },
      { "product_id": 5, "quantity": 2, "price": 25 }
    ]
  }
}
```

**Response:**
```json
{
  "updated_cart": {
    "items": [
      { "product_id": 1, ..., "total_discount": 0 },
      { "product_id": 2, ..., "total_discount": 0 },
      { "product_id": 5, ..., "total_discount": 40 }
    ],
    "total_price": 440,
    "total_discount": 40,
    "final_price": 400
  }
}
```

**BxGy Coupon Request:**
```json
{
  "cart": {
    "items": [
      { "product_id": 1, "quantity": 6, "price": 500 },
      { "product_id": 2, "quantity": 3, "price": 300 },
      { "product_id": 3, "quantity": 2, "price": 250 }
    ]
  }
}
```

**Response:**
```json
{
  "updated_cart": {
    "items": [
      { "product_id": 1, ..., "total_discount": 0 },
      { "product_id": 2, ..., "total_discount": 0 },
      { "product_id": 3, ..., "total_discount": 500 }
    ],
    "total_price": 4400,
    "total_discount": 500,
    "final_price": 3900
  }
}
```

## POST /applicable-coupons

**Request:**
```json
{
  "cart": {
    "items": [
      { "product_id": 1, "quantity": 6, "price": 50 },
      { "product_id": 2, "quantity": 3, "price": 30 },
      { "product_id": 5, "quantity": 2, "price": 25 },
      { "product_id": 6, "quantity": 1, "price": 2000 }
    ]
  }
}
```

**Response:**
```json
{
  "applicable_coupons": [
    { "coupon_id": 1, "type": "cart-wise", "discount": 366 },
    { "coupon_id": 4, "type": "product-wise", "discount": 60 },
    ...
  ]
}
```
