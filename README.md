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

## Implemented Use Cases

1. ✅ **Cart-wise Coupons**  
   - **Condition:** Discount applies to entire cart if total exceeds a threshold  
   - **Discount:** Fixed or percentage  
   - **Fields:** `discount_value`, `discount_type`, `repetition_threshold` (as threshold)  
   - **Examples:**  
     - 10% off if cart total > ₹500  
     - Flat ₹50 off on cart total > ₹1000  

2. ✅ **Product-wise Coupons**  
   - **Condition:** Discount only if specific product exists in cart  
   - **Discount:** Fixed or percentage  
   - **Fields:** `product_id`, `discount_value`, `discount_type`  
   - **Examples:**  
     - ₹100 off on Laptop  
     - 25% off on Mouse  

3. ✅ **BxGy (Buy X, Get Y Free) Coupons**  
   - **Condition:** Buy `buyQuantity` of products from one list, get `getQuantity` of products from another list  
   - **Fields:** `buyQuantity`, `getQuantity`, `repetition_threshold`  
   - Buy and Get product arrays (stored in separate tables)  
   - **Examples:**  
     - Buy 3 of [X, Y], Get 1 of [Z] free  
     - Repeat up to 2 times if enough items in cart  

4. ✅ **Expiry Check on Coupons**  
   - All coupons include an `expiration_date`, checked before application  

5. ✅ **Selective Field Updates**  
   - Coupon update is type-safe:  
     - Cart-wise: Only `discount_value` and `repetition_threshold`  
     - Product-wise: Only `discount_value` and `product_id`  
     - BxGy: Only `buyQuantity` and `getQuantity`  

6. ✅ **Coupon Application & Cart Update**  
   - Applies discount and returns:  
     - `updated_cart` with:  
       - `total_price`  
       - `total_discount`  
       - `final_price`  
       - Item-wise `total_discount`  
   - Validation includes:  
     - Threshold met for cart-wise  
     - Product presence for product-wise  
     - Quantity match and repetition logic for BxGy  
     - Coupon not expired  


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

## Sample API Requests and Responses

For detailed examples of API requests and responses, please refer to the [API Examples](https://github.com/Abdullah-Shaikh01/monk-coupons/blob/main/API_Examples.md) file. It contains comprehensive input/output samples to help you understand how to interact with the API effectively.

---
