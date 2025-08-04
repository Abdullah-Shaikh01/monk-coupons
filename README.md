# Monk-Coupons
This project implements a coupon management system in Go, supporting cart-wise, product-wise, and buy-x-get-y (BxGy) discount types with RESTful APIs. It allows applying coupons to shopping carts with flexible discount rules, including expiration and repetition thresholds.

## Transparency Note

This is my first Go language project. I have learned Go while completing this task within 3 days. To assist with the implementation, I extensively used AI tools such as ChatGPT and Perplexity. Despite this, I have carefully analyzed and understood the code and take full ownership of its correctness and operation—except for the test files, which do not yet comprehensively cover all functionalities nor have I fully verified their accuracy.

---

## Assumptions

- **BxGy Coupon Quantity Interpretation:**  
  The requirements mention buying a specified number of products from one “array” (list) and getting products from another array for free. However, the provided payload example indicates each individual product must be bought in the specified quantity. To resolve this conflict, I have assumed the quantity specified in the buy_products array represents the total combined quantity required across all products in the list, not per product individually.  
  For example, a total of 3 items from any combination of the buy_products qualifies for the free items. This approach aligns more closely with the example and business interpretation, though it introduces redundancy in the payload.

- **Response Calculation Mismatch:**  
  Some response examples in the requirements appear to contain calculation errors. I have assumed those to be inadvertent mistakes and corrected calculations in implementation.

- **Default Discount Type:**  
  Unless explicitly specified in the payload, discounts are assumed to be percentages rather than fixed amounts.

- **BxGy Free Items Logic:**  
  Based on the description, free items in BxGy coupons come from products already in the cart. Although an example suggests adding extra free items, I chose to follow the description: free items come only from existing cart products, and quantities remain unchanged.  

- **Enhancements for Extensibility:**  
  - Support for coupons with either percentage or flat discount types (defaulted to percentage).  
  - Coupons have expiration dates (defaulted to one month from creation if unspecified).  
  - Added products table and related APIs for realistic pricing and product data, though price used in calculations is from the request payload.  
  - Partial updates supported in `PUT /coupons/{id}` constrained by coupon type, allowing only relevant fields to be updated.

---

## Implemented Cases

- **Coupon Types:**  
  - Cart-wise coupons that apply discounts on cart totals exceeding a threshold.  
  - Product-wise coupons targeting specific products.  
  - BxGy coupons supporting buy and get product sets with repetition limits.

- **API Coverage:**  
  - Complete CRUD for coupons and get all products.  
  - Application of coupon by ID and retrieving all applicable coupons given a cart.

- **Validation and Expiration:**  
  - Coupon expiration enforcement on application.  
  - Partial update validation based on coupon type fields.

---

## Unimplemented Cases

- **Stacking Coupons:**  
  Support for combining multiple coupons optimally per cart and returning best discount combinations is not implemented.

- **API Security:**  
  Middleware for authorized API access is absent.  
  Sensitive data encryption not implemented, as no sensitive data currently exists in use cases.

- **BxGy Free Item Optimization:**  
  Free items are currently selected without price-based optimization. Future work could select free products maximizing customer discount.

- **Comprehensive BxGy Coupon Details in GET Coupon Response:**  
  Buy/get product details for BxGy coupons are not reflected in the coupon retrieval API responses.

---

## Limitations

- **Price Integrity:**  
  Prices from cart requests are used directly, without verification against backend product prices, which may cause inconsistencies.

- **Simplified BxGy Quantity Logic:**  
  Quantity distribution across multiple buy-products assumes a combined total suffices, without complex distribution calculations.

- **Multiple Coupon Combinations:**  
  Evaluation supports multiple coupons independently but lacks logic for combined or prioritized discounts.

- **Security Not Enforced:**  
  APIs currently lack authentication and authorization mechanisms critical for production environments.

- **Incomplete Test Coverage:**  
  Automated tests do not yet cover all functionalities fully and require further validation.

---

This README outlines the project’s current scope, design assumptions, operational coverage, and the areas left for future improvements, providing a clear understanding of what is delivered and the intended direction.


# 🧾 Monk Coupons API

This is a Go-based backend service for managing and applying coupon logic, using MySQL as the database.

---

## 🚀 How to Run the Project Locally

### 1. 🛠 Install Go and MySQL

#### Go
- Download and install from: [https://golang.org/dl/](https://golang.org/dl/)
- Set up Go environment variables (`GOROOT`, `GOPATH`) as per your OS instructions.

#### MySQL
- Download and install from: [https://dev.mysql.com/downloads/mysql/](https://dev.mysql.com/downloads/mysql/)
- Start the MySQL service.
- Create a new database (e.g., `monk_coupons`) for this project.

---

### 2. 🧱 Set Up the Database

- Create the database:

```sql
CREATE DATABASE monk_coupons;
```

- Update your `.env` file:

```env
DB_USER=your_mysql_username
DB_PASS=your_mysql_password
DB_HOST=localhost
DB_PORT=3306
DB_NAME=monk_coupons
```

- Run database migrations:

```bash
go run migrations/migrations.go
```

---

### 3. 📦 Install Project Dependencies

From the project root, run:

```bash
go mod tidy
```

This will download all necessary Go modules.

---

### 4. 🧪 Run Tests (if any)

```bash
go test -v ./...
```

---

### 5. ▶️ Run the Project

Start the API server:

```bash
go run main.go
```

---

## 📬 API Documentation

Refer to the [Postman collection](https://www.postman.com/maintenance-cosmonaut-35719276/monk-coupons/collection/xxm173m/monk-coupons?action=share&creator=21298069) for sample requests and responses.
