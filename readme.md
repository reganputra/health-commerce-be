# Health Store API Documentation

**Version:** 1.0
**Base URL:** `http://localhost:8080` (development)
**Authentication:** JWT Bearer Token
**Content-Type:** `application/json`

---

## Table of Contents

1. [Quick Start](#quick-start)
2. [Authentication](#authentication)
3. [API Endpoints](#api-endpoints)
   - [Health Check](#health-check)
   - [Authentication](#authentication-endpoints)
   - [Products](#products)
   - [Categories](#categories)
   - [Shopping Cart](#shopping-cart)
   - [Orders](#orders)
   - [Feedback](#feedback)
   - [User Management (Admin)](#user-management)
   - [Reports (Admin)](#reports)
4. [Data Models](#data-models)
5. [Error Handling](#error-handling)
6. [Rate Limiting & Best Practices](#best-practices)

---

## Quick Start

### 1. Register a New User

```javascript
const response = await fetch("http://localhost:8080/auth/register", {
  method: "POST",
  headers: { "Content-Type": "application/json" },
  body: JSON.stringify({
    username: "john_doe",
    password: "securepass123",
    email: "john@example.com",
    dob: "1990-01-15",
    gender: "M",
    address: "123 Main Street",
    city: "New York",
    contact_number: "1234567890",
  }),
});
```

### 2. Login & Get Token

```javascript
const response = await fetch("http://localhost:8080/auth/login", {
  method: "POST",
  headers: { "Content-Type": "application/json" },
  body: JSON.stringify({
    username: "john_doe",
    password: "securepass123",
  }),
});
const { token, user } = await response.json();
// Store token for subsequent requests
localStorage.setItem("authToken", token);
```

### 3. Make Authenticated Requests

```javascript
const response = await fetch("http://localhost:8080/cart/", {
  method: "GET",
  headers: {
    Authorization: `Bearer ${localStorage.getItem("authToken")}`,
    "Content-Type": "application/json",
  },
});
```

---

## Authentication

### Authentication Flow

1. **Register** → Get user account
2. **Login** → Receive JWT token
3. **Use Token** → Include in `Authorization` header for protected routes
4. **Token Expiry** → Tokens expire after 24 hours, user must re-login

### Token Structure

```javascript
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 123,
    "username": "john_doe",
    "email": "john@example.com",
    "role": "customer" // or "admin"
  }
}
```

### Role-Based Access Control

| Role         | Permissions                                                                 |
| ------------ | --------------------------------------------------------------------------- |
| **Customer** | Browse products, manage cart, place orders, submit feedback                 |
| **Admin**    | All customer permissions + manage users, products, categories, view reports |

---

## API Endpoints

### Health Check

#### Check API Status

```http
GET /ping
```

**Response:**

```json
{
  "message": "pong"
}
```

---

## Authentication Endpoints

### Register User

```http
POST /auth/register
```

**Request Body:**

```json
{
  "username": "john_doe", // 3-50 characters
  "password": "securepass123", // min 6 characters
  "email": "john@example.com", // valid email
  "dob": "1990-01-15", // date format: YYYY-MM-DD
  "gender": "M", // "M" or "F"
  "address": "123 Main Street", // 10-255 characters
  "city": "New York", // 2-100 characters
  "contact_number": "1234567890" // 10-15 characters
}
```

**Success Response (200):**

```json
{
  "message": "User registered successfully",
  "user_id": 123
}
```

**Error Responses:**

- `400` - Validation error (missing/invalid fields)
- `409` - Username or email already exists

---

### Login

```http
POST /auth/login
```

**Request Body:**

```json
{
  "username": "john_doe",
  "password": "securepass123"
}
```

**Success Response (200):**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 123,
    "username": "john_doe",
    "email": "john@example.com",
    "role": "customer",
    "dob": "1990-01-15",
    "gender": "M",
    "address": "123 Main Street",
    "city": "New York",
    "contact_number": "1234567890"
  }
}
```

**Error Responses:**

- `400` - Invalid request format
- `401` - Invalid credentials

**Frontend Example:**

```javascript
async function login(username, password) {
  try {
    const response = await fetch("http://localhost:8080/auth/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username, password }),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error);
    }

    const { token, user } = await response.json();
    localStorage.setItem("authToken", token);
    localStorage.setItem("user", JSON.stringify(user));
    return { token, user };
  } catch (error) {
    console.error("Login failed:", error);
    throw error;
  }
}
```

---

## Products

### Get All Products

```http
GET /api/products
```

**Authentication:** Not required

**Success Response (200):**

```json
[
  {
    "id": 1,
    "category_id": 2,
    "name": "Vitamin C Tablets",
    "description": "High-quality vitamin C supplement for immune support",
    "price": 19.99,
    "stock": 150,
    "image_url": "https://example.com/images/vitamin-c.jpg",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
]
```

**Frontend Example:**

```javascript
async function getAllProducts() {
  const response = await fetch("http://localhost:8080/api/products");
  return await response.json();
}
```

---

### Get Product by ID

```http
GET /api/products/:id
```

**Authentication:** Not required

**Path Parameters:**

- `id` (integer) - Product ID

**Success Response (200):**

```json
{
  "id": 1,
  "category_id": 2,
  "name": "Vitamin C Tablets",
  "description": "High-quality vitamin C supplement for immune support",
  "price": 19.99,
  "stock": 150,
  "image_url": "https://example.com/images/vitamin-c.jpg",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

**Error Responses:**

- `400` - Invalid product ID
- `404` - Product not found

---

### Create Product (Admin Only)

```http
POST /admin/products
```

**Authentication:** Required (Admin role)

**Request Body:**

```json
{
  "category_id": 2,
  "name": "Omega-3 Fish Oil",
  "description": "Premium omega-3 supplement for heart health",
  "price": 29.99,
  "stock": 100,
  "image_url": "https://example.com/images/omega3.jpg"
}
```

**Validation Rules:**

- `name`: 2-255 characters
- `description`: 10-1000 characters
- `price`: Must be greater than 0
- `stock`: Must be >= 0
- `image_url`: Valid URL format

**Success Response (200):**

```json
{
  "id": 15,
  "category_id": 2,
  "name": "Omega-3 Fish Oil",
  "description": "Premium omega-3 supplement for heart health",
  "price": 29.99,
  "stock": 100,
  "image_url": "https://example.com/images/omega3.jpg",
  "created_at": "2024-01-20T14:30:00Z",
  "updated_at": "2024-01-20T14:30:00Z"
}
```

**Error Responses:**

- `400` - Validation error
- `401` - Not authenticated
- `403` - Insufficient permissions (not admin)

---

### Update Product (Admin Only)

```http
PUT /admin/products/:id
```

**Authentication:** Required (Admin role)

**Request Body:** (All fields optional - partial update supported)

```json
{
  "name": "Updated Product Name",
  "price": 24.99,
  "stock": 200
}
```

**Success Response (200):**

```json
{
  "id": 15,
  "category_id": 2,
  "name": "Updated Product Name",
  "description": "Premium omega-3 supplement for heart health",
  "price": 24.99,
  "stock": 200,
  "image_url": "https://example.com/images/omega3.jpg",
  "created_at": "2024-01-20T14:30:00Z",
  "updated_at": "2024-01-21T09:15:00Z"
}
```

---

### Delete Product (Admin Only)

```http
DELETE /admin/products/:id
```

**Authentication:** Required (Admin role)

**Success Response (200):**

```json
{
  "message": "Product deleted successfully"
}
```

**Frontend Example:**

```javascript
async function deleteProduct(productId) {
  const token = localStorage.getItem("authToken");
  const response = await fetch(
    `http://localhost:8080/admin/products/${productId}`,
    {
      method: "DELETE",
      headers: {
        Authorization: `Bearer ${token}`,
      },
    }
  );
  return await response.json();
}
```

---

## Categories

### Get All Categories

```http
GET /api/categories
```

**Authentication:** Not required

**Success Response (200):**

```json
[
  {
    "id": 1,
    "name": "Vitamins",
    "description": "Essential vitamins and supplements",
    "created_at": "2024-01-10T08:00:00Z",
    "updated_at": "2024-01-10T08:00:00Z"
  },
  {
    "id": 2,
    "name": "Supplements",
    "description": "Health and wellness supplements",
    "created_at": "2024-01-10T08:05:00Z",
    "updated_at": "2024-01-10T08:05:00Z"
  }
]
```

---

### Get Category by ID

```http
GET /api/categories/:id
```

**Authentication:** Not required

**Success Response (200):**

```json
{
  "id": 1,
  "name": "Vitamins",
  "description": "Essential vitamins and supplements",
  "created_at": "2024-01-10T08:00:00Z",
  "updated_at": "2024-01-10T08:00:00Z"
}
```

---

### Create Category (Admin Only)

```http
POST /admin/categories
```

**Authentication:** Required (Admin role)

**Request Body:**

```json
{
  "name": "Protein Powders",
  "description": "High-quality protein supplements for fitness"
}
```

**Success Response (200):**

```json
{
  "id": 5,
  "name": "Protein Powders",
  "description": "High-quality protein supplements for fitness",
  "created_at": "2024-01-22T11:00:00Z",
  "updated_at": "2024-01-22T11:00:00Z"
}
```

---

### Update Category (Admin Only)

```http
PUT /admin/categories/:id
```

**Authentication:** Required (Admin role)

**Request Body:**

```json
{
  "name": "Updated Category Name",
  "description": "Updated description (10-500 chars)"
}
```

**Validation Rules:**

- `name`: 2-100 characters
- `description`: 10-500 characters

---

### Delete Category (Admin Only)

```http
DELETE /admin/categories/:id
```

**Authentication:** Required (Admin role)

**Success Response (200):**

```json
{
  "message": "Category deleted successfully"
}
```

---

## Shopping Cart

### Get User's Cart

```http
GET /cart/
```

**Authentication:** Required (Customer or Admin)

**Success Response (200):**

```json
{
  "cart": {
    "id": 1,
    "user_id": 123,
    "items": [
      {
        "id": 1,
        "cart_id": 1,
        "product_id": 5,
        "quantity": 2,
        "product": {
          "id": 5,
          "name": "Vitamin C Tablets",
          "description": "High-quality vitamin C supplement",
          "price": 19.99,
          "stock": 150,
          "image_url": "https://example.com/images/vitamin-c.jpg"
        }
      },
      {
        "id": 2,
        "cart_id": 1,
        "product_id": 8,
        "quantity": 1,
        "product": {
          "id": 8,
          "name": "Omega-3 Fish Oil",
          "description": "Premium omega-3 supplement",
          "price": 29.99,
          "stock": 100,
          "image_url": "https://example.com/images/omega3.jpg"
        }
      }
    ],
    "created_at": "2024-01-15T10:00:00Z",
    "updated_at": "2024-01-20T14:30:00Z"
  },
  "total": 69.97
}
```

**Frontend Example:**

```javascript
async function getCart() {
  const token = localStorage.getItem("authToken");
  const response = await fetch("http://localhost:8080/cart/", {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
  return await response.json();
}
```

---

### Add Item to Cart

```http
POST /cart/
```

**Authentication:** Required (Customer or Admin)

**Request Body:**

```json
{
  "product_id": 5,
  "quantity": 2
}
```

**Success Response (200):**

```json
{
  "message": "Item added to cart successfully"
}
```

**Error Responses:**

- `400` - Invalid product ID or quantity
- `401` - Not authenticated
- `403` - Insufficient permissions
- `500` - Server error (e.g., insufficient stock)

**Frontend Example:**

```javascript
async function addToCart(productId, quantity) {
  const token = localStorage.getItem("authToken");
  const response = await fetch("http://localhost:8080/cart/", {
    method: "POST",
    headers: {
      Authorization: `Bearer ${token}`,
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ product_id: productId, quantity }),
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error);
  }

  return await response.json();
}
```

---

### Remove Item from Cart

```http
DELETE /cart/:id
```

**Authentication:** Required (Customer or Admin)

**Path Parameters:**

- `id` (integer) - Cart item ID (not product ID)

**Success Response (200):**

```json
{
  "message": "Item removed from cart successfully"
}
```

**Frontend Example:**

```javascript
async function removeFromCart(cartItemId) {
  const token = localStorage.getItem("authToken");
  const response = await fetch(`http://localhost:8080/cart/${cartItemId}`, {
    method: "DELETE",
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
  return await response.json();
}
```

---

## Orders

### Place Order

```http
POST /orders/
```

**Authentication:** Required (Customer or Admin)

**Request Body:**

```json
{
  "payment_method": "paypal", // Options: "paypal", "debit", "cc", "cod"
  "bank_name": "Chase Bank" // Optional, only for debit/cc
}
```

**Payment Methods:**

- `paypal` - PayPal payment
- `debit` - Debit card
- `cc` - Credit card
- `cod` - Cash on delivery

**Success Response (200):**

```json
{
  "message": "Order placed successfully",
  "order": {
    "id": 42,
    "user_id": 123,
    "status": "pending",
    "total_price": 69.97,
    "payment_method": "paypal",
    "bank_name": "",
    "created_at": "2024-01-22T15:00:00Z",
    "updated_at": "2024-01-22T15:00:00Z",
    "items": [
      {
        "id": 85,
        "order_id": 42,
        "product_id": 5,
        "quantity": 2,
        "price": 19.99
      }
    ]
  }
}
```

**Notes:**

- Creates an order from the user's current cart items
- Cart is cleared after successful order placement
- Order status is initially set to "pending"

**Frontend Example:**

```javascript
async function placeOrder(paymentMethod, bankName = "") {
  const token = localStorage.getItem("authToken");
  const response = await fetch("http://localhost:8080/orders/", {
    method: "POST",
    headers: {
      Authorization: `Bearer ${token}`,
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      payment_method: paymentMethod,
      bank_name: bankName,
    }),
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error);
  }

  return await response.json();
}
```

---

### Get User's Orders

```http
GET /orders/
```

**Authentication:** Required (Customer or Admin)

**Success Response (200):**

```json
{
  "orders": [
    {
      "id": 42,
      "user_id": 123,
      "status": "paid",
      "total_price": 69.97,
      "payment_method": "paypal",
      "bank_name": "",
      "created_at": "2024-01-22T15:00:00Z",
      "updated_at": "2024-01-22T16:00:00Z"
    }
  ],
  "count": 1
}
```

---

### Get Order by ID

```http
GET /orders/:id
```

**Authentication:** Required (Customer or Admin)

**Success Response (200):**

```json
{
  "id": 42,
  "user_id": 123,
  "status": "paid",
  "total_price": 69.97,
  "payment_method": "paypal",
  "bank_name": "",
  "created_at": "2024-01-22T15:00:00Z",
  "updated_at": "2024-01-22T16:00:00Z",
  "items": [
    {
      "id": 85,
      "order_id": 42,
      "product_id": 5,
      "quantity": 2,
      "price": 19.99,
      "product": {
        "id": 5,
        "name": "Vitamin C Tablets",
        "description": "High-quality vitamin C supplement",
        "price": 19.99,
        "stock": 148,
        "image_url": "https://example.com/images/vitamin-c.jpg"
      }
    }
  ]
}
```

---

### Cancel Order

```http
PUT /orders/:id/cancel
```

**Authentication:** Required (Customer or Admin)

**Success Response (200):**

```json
{
  "message": "Order cancelled successfully"
}
```

**Notes:**

- Changes order status to "cancelled"
- Product stock is restored

---

### Get All Orders (Admin Only)

```http
GET /admin/orders/
```

**Authentication:** Required (Admin role)

**Success Response (200):**

```json
[
  {
    "id": 42,
    "user_id": 123,
    "status": "paid",
    "total_price": 69.97,
    "payment_method": "paypal",
    "bank_name": "",
    "created_at": "2024-01-22T15:00:00Z",
    "updated_at": "2024-01-22T16:00:00Z"
  }
]
```

---

### Update Order Status (Admin Only)

```http
PUT /admin/orders/:id/status
```

**Authentication:** Required (Admin role)

**Request Body:**

```json
{
  "status": "shipped" // Options: "pending", "paid", "shipped", "cancelled"
}
```

**Order Status Values:**

- `pending` - Order placed, awaiting payment
- `paid` - Payment confirmed
- `shipped` - Order shipped to customer
- `cancelled` - Order cancelled

**Success Response (200):**

```json
{
  "message": "Order status updated successfully"
}
```

**Frontend Example:**

```javascript
async function updateOrderStatus(orderId, newStatus) {
  const token = localStorage.getItem("authToken");
  const response = await fetch(
    `http://localhost:8080/admin/orders/${orderId}/status`,
    {
      method: "PUT",
      headers: {
        Authorization: `Bearer ${token}`,
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ status: newStatus }),
    }
  );
  return await response.json();
}
```

---

## Feedback

### Submit Product Feedback

```http
POST /feedback/
```

**Authentication:** Required (Customer or Admin)

**Request Body:**

```json
{
  "product_id": 5,
  "comment": "Great product! Highly recommend.",
  "rating": 5 // Integer from 1 to 5
}
```

**Validation Rules:**

- `rating`: Must be between 1 and 5 (inclusive)
- `comment`: Required string

**Success Response (200):**

```json
{
  "id": 15,
  "userId": 123,
  "productId": 5,
  "user": {
    "id": 123,
    "username": "john_doe"
  },
  "product": {
    "id": 5,
    "name": "Vitamin C Tablets",
    "description": "High-quality vitamin C supplement",
    "image": "https://example.com/images/vitamin-c.jpg"
  },
  "comment": "Great product! Highly recommend.",
  "rating": 5,
  "createdAt": "2024-01-23T10:30:00Z"
}
```

**Frontend Example:**

```javascript
async function submitFeedback(productId, comment, rating) {
  const token = localStorage.getItem("authToken");
  const response = await fetch("http://localhost:8080/feedback/", {
    method: "POST",
    headers: {
      Authorization: `Bearer ${token}`,
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ product_id: productId, comment, rating }),
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error);
  }

  return await response.json();
}
```

---

## User Management

### Get All Users (Admin Only)

```http
GET /admin/users
```

**Authentication:** Required (Admin role)

**Success Response (200):**

```json
[
  {
    "id": 123,
    "username": "john_doe",
    "email": "john@example.com",
    "role": "customer",
    "dob": "1990-01-15",
    "gender": "M",
    "address": "123 Main Street",
    "city": "New York",
    "contact_number": "1234567890",
    "created_at": "2024-01-15T10:00:00Z",
    "updated_at": "2024-01-15T10:00:00Z"
  }
]
```

---

### Get User by ID (Admin Only)

```http
GET /admin/users/:id
```

**Authentication:** Required (Admin role)

**Success Response (200):**

```json
{
  "id": 123,
  "username": "john_doe",
  "email": "john@example.com",
  "role": "customer",
  "dob": "1990-01-15",
  "gender": "M",
  "address": "123 Main Street",
  "city": "New York",
  "contact_number": "1234567890",
  "paypal_id": "",
  "created_at": "2024-01-15T10:00:00Z",
  "updated_at": "2024-01-15T10:00:00Z"
}
```

---

### Update User (Admin Only)

```http
PUT /admin/users/:id
```

**Authentication:** Required (Admin role)

**Request Body:** (All fields optional - partial update)

```json
{
  "username": "new_username",
  "email": "newemail@example.com",
  "role": "admin",
  "address": "456 New Address",
  "city": "Los Angeles",
  "contact_number": "9876543210"
}
```

**Success Response (200):**

```json
{
  "id": 123,
  "username": "new_username",
  "email": "newemail@example.com",
  "role": "admin",
  "dob": "1990-01-15",
  "gender": "M",
  "address": "456 New Address",
  "city": "Los Angeles",
  "contact_number": "9876543210",
  "created_at": "2024-01-15T10:00:00Z",
  "updated_at": "2024-01-23T14:00:00Z"
}
```

---

### Delete User (Admin Only)

```http
DELETE /admin/users/:id
```

**Authentication:** Required (Admin role)

**Success Response (200):**

```json
{
  "message": "User deleted successfully"
}
```

---

## Reports

### Generate Report (Admin Only)

```http
GET /admin/report
```

**Authentication:** Required (Admin role)

**Query Parameters:**

- `type` (string, optional) - Report type: `summary`, `detailed`, or `financial` (default: `summary`)
- `format` (string, optional) - Output format: `pdf` or `csv` (default: `pdf`)
- `start_date` (string, optional) - Start date in `YYYY-MM-DD` format
- `end_date` (string, optional) - End date in `YYYY-MM-DD` format
- `limit` (integer, optional) - Limit number of records (default: 10)

**Report Types:**

1. **Summary Report** (`type=summary`)

   - Total revenue
   - Total orders count
   - Top 10 products by sales
   - Recent customer activity

2. **Detailed Report** (`type=detailed`)

   - Complete order details with items
   - Customer information
   - Product details

3. **Financial Report** (`type=financial`)
   - Revenue breakdown by product
   - Payment method distribution
   - Time-based revenue analysis

**Response:**

- Binary file download (PDF or CSV)
- Content-Type: `application/pdf` or `text/csv`
- Content-Disposition: `attachment; filename=report.pdf` or `report.csv`

**Frontend Example (PDF Download):**

```javascript
async function downloadReport(
  type = "summary",
  format = "pdf",
  startDate = null,
  endDate = null
) {
  const token = localStorage.getItem("authToken");

  // Build query parameters
  const params = new URLSearchParams();
  params.append("type", type);
  params.append("format", format);
  if (startDate) params.append("start_date", startDate);
  if (endDate) params.append("end_date", endDate);

  const response = await fetch(
    `http://localhost:8080/admin/report?${params.toString()}`,
    {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    }
  );

  if (!response.ok) {
    throw new Error("Failed to generate report");
  }

  // Create download link
  const blob = await response.blob();
  const url = window.URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = `report_${type}_${Date.now()}.${format}`;
  document.body.appendChild(a);
  a.click();
  window.URL.revokeObjectURL(url);
  document.body.removeChild(a);
}

// Usage examples:
// downloadReport('summary', 'pdf');
// downloadReport('financial', 'csv', '2024-01-01', '2024-01-31');
```

---

## Data Models

### User Model

```typescript
interface User {
  id: number;
  username: string;
  email: string;
  role: "admin" | "customer";
  dob: string; // YYYY-MM-DD format
  gender: "M" | "F";
  address: string;
  city: string;
  contact_number: string;
  paypal_id?: string;
  created_at: string; // ISO 8601 timestamp
  updated_at: string; // ISO 8601 timestamp
}
```

### Product Model

```typescript
interface Product {
  id: number;
  category_id: number;
  name: string;
  description: string;
  price: number;
  stock: number;
  image_url: string;
  created_at: string;
  updated_at: string;
}
```

### Category Model

```typescript
interface Category {
  id: number;
  name: string;
  description: string;
  created_at: string;
  updated_at: string;
}
```

### Cart Model

```typescript
interface Cart {
  id: number;
  user_id: number;
  items: CartItem[];
  created_at: string;
  updated_at: string;
}

interface CartItem {
  id: number;
  cart_id: number;
  product_id: number;
  quantity: number;
  product: Product;
}
```

### Order Model

```typescript
interface Order {
  id: number;
  user_id: number;
  status: "pending" | "paid" | "shipped" | "cancelled";
  total_price: number;
  payment_method: "paypal" | "debit" | "cc" | "cod";
  bank_name?: string;
  created_at: string;
  updated_at: string;
  items?: OrderItem[];
}

interface OrderItem {
  id: number;
  order_id: number;
  product_id: number;
  quantity: number;
  price: number; // Price at time of order
  product?: Product;
}
```

### Feedback Model

```typescript
interface Feedback {
  id: number;
  userId: number;
  productId: number;
  user: {
    id: number;
    username: string;
  };
  product: {
    id: number;
    name: string;
    description: string;
    image: string;
  };
  comment: string;
  rating: number; // 1-5
  createdAt: string;
}
```

---

## Error Handling

### Error Response Format

```typescript
interface ErrorResponse {
  error: string;
  code?: string;
  required?: string; // For permission errors
  role?: string; // For permission errors
}
```

### HTTP Status Codes

| Code  | Meaning               | Common Causes                            |
| ----- | --------------------- | ---------------------------------------- |
| `200` | Success               | Request completed successfully           |
| `400` | Bad Request           | Invalid input, validation errors         |
| `401` | Unauthorized          | Missing or invalid authentication token  |
| `403` | Forbidden             | Insufficient permissions (wrong role)    |
| `404` | Not Found             | Resource doesn't exist                   |
| `409` | Conflict              | Duplicate username/email on registration |
| `500` | Internal Server Error | Database errors, server issues           |

### Common Error Examples

**400 Bad Request:**

```json
{
  "error": "Invalid request format: missing required field 'username'"
}
```

**401 Unauthorized:**

```json
{
  "error": "Authorization header required"
}
```

or

```json
{
  "error": "Invalid token"
}
```

**403 Forbidden:**

```json
{
  "error": "Insufficient permissions",
  "required": "create:product",
  "role": "customer"
}
```

**404 Not Found:**

```json
{
  "code": "PAGE_NOT_FOUND",
  "message": "Page not found"
}
```

or

```json
{
  "error": "Product not found"
}
```

**409 Conflict:**

```json
{
  "error": "Username or email already exists"
}
```

**500 Internal Server Error:**

```json
{
  "error": "Failed to fetch products: database connection error"
}
```

### Frontend Error Handling Example

```javascript
async function handleApiCall(url, options = {}) {
  try {
    const response = await fetch(url, options);

    // Check if response is ok
    if (!response.ok) {
      const error = await response.json();

      switch (response.status) {
        case 400:
          throw new Error(`Validation Error: ${error.error}`);
        case 401:
          // Redirect to login
          localStorage.removeItem("authToken");
          window.location.href = "/login";
          throw new Error("Session expired. Please login again.");
        case 403:
          throw new Error(`Access Denied: ${error.error}`);
        case 404:
          throw new Error("Resource not found");
        case 409:
          throw new Error(`Conflict: ${error.error}`);
        case 500:
          throw new Error("Server error. Please try again later.");
        default:
          throw new Error(error.error || "An error occurred");
      }
    }

    return await response.json();
  } catch (error) {
    console.error("API call failed:", error);
    throw error;
  }
}
```

---

## Best Practices

### 1. Token Management

```javascript
// Store token after login
function storeAuthToken(token, user) {
  localStorage.setItem("authToken", token);
  localStorage.setItem("user", JSON.stringify(user));
}

// Get token for requests
function getAuthToken() {
  return localStorage.getItem("authToken");
}

// Clear token on logout
function logout() {
  localStorage.removeItem("authToken");
  localStorage.removeItem("user");
  window.location.href = "/login";
}

// Check if user is authenticated
function isAuthenticated() {
  return !!localStorage.getItem("authToken");
}

// Check user role
function isAdmin() {
  const user = JSON.parse(localStorage.getItem("user") || "{}");
  return user.role === "admin";
}
```

### 2. API Service Class

```javascript
class HealthStoreAPI {
  constructor(baseURL = "http://localhost:8080") {
    this.baseURL = baseURL;
  }

  getHeaders(includeAuth = true) {
    const headers = {
      "Content-Type": "application/json",
    };

    if (includeAuth) {
      const token = localStorage.getItem("authToken");
      if (token) {
        headers["Authorization"] = `Bearer ${token}`;
      }
    }

    return headers;
  }

  async request(endpoint, options = {}) {
    const url = `${this.baseURL}${endpoint}`;
    const response = await fetch(url, {
      ...options,
      headers: {
        ...this.getHeaders(options.auth !== false),
        ...options.headers,
      },
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || "Request failed");
    }

    return await response.json();
  }

  // Auth
  async register(userData) {
    return this.request("/auth/register", {
      method: "POST",
      body: JSON.stringify(userData),
      auth: false,
    });
  }

  async login(username, password) {
    return this.request("/auth/login", {
      method: "POST",
      body: JSON.stringify({ username, password }),
      auth: false,
    });
  }

  // Products
  async getProducts() {
    return this.request("/api/products", { auth: false });
  }

  async getProduct(id) {
    return this.request(`/api/products/${id}`, { auth: false });
  }

  // Cart
  async getCart() {
    return this.request("/cart/");
  }

  async addToCart(productId, quantity) {
    return this.request("/cart/", {
      method: "POST",
      body: JSON.stringify({ product_id: productId, quantity }),
    });
  }

  async removeFromCart(cartItemId) {
    return this.request(`/cart/${cartItemId}`, { method: "DELETE" });
  }

  // Orders
  async placeOrder(paymentMethod, bankName = "") {
    return this.request("/orders/", {
      method: "POST",
      body: JSON.stringify({
        payment_method: paymentMethod,
        bank_name: bankName,
      }),
    });
  }

  async getOrders() {
    return this.request("/orders/");
  }

  async getOrder(id) {
    return this.request(`/orders/${id}`);
  }

  async cancelOrder(id) {
    return this.request(`/orders/${id}/cancel`, { method: "PUT" });
  }

  // Feedback
  async submitFeedback(productId, comment, rating) {
    return this.request("/feedback/", {
      method: "POST",
      body: JSON.stringify({ product_id: productId, comment, rating }),
    });
  }
}

// Usage
const api = new HealthStoreAPI();

// Login
const { token, user } = await api.login("john_doe", "password123");
localStorage.setItem("authToken", token);

// Get products
const products = await api.getProducts();

// Add to cart
await api.addToCart(5, 2);
```

### 3. Request Optimization

- Cache public data (products, categories) for better performance
- Implement debouncing for search/filter operations
- Use pagination for large datasets (though not currently supported by API)
- Batch updates when possible

### 4. Security Best Practices

- **Never** store sensitive data in localStorage beyond the auth token
- Always validate user input before sending to API
- Implement HTTPS in production
- Handle token expiration gracefully (24-hour expiry)
- Clear auth data on logout
- Sanitize user-generated content (feedback, comments)

### 5. User Experience

- Show loading states during API calls
- Display clear error messages to users
- Implement retry logic for failed requests
- Cache non-sensitive data to reduce API calls
- Validate forms client-side before submission

---

## Testing the API

### Using cURL

**Register a user:**

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "test123",
    "email": "test@example.com",
    "dob": "1990-01-01",
    "gender": "M",
    "address": "123 Test St",
    "city": "Test City",
    "contact_number": "1234567890"
  }'
```

**Login:**

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "testuser", "password": "test123"}'
```

**Get products (no auth):**

```bash
curl http://localhost:8080/api/products
```

**Get cart (with auth):**

```bash
curl http://localhost:8080/cart/ \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Using Postman

1. Create a new collection named "Health Store API"
2. Add environment variables:
   - `base_url`: `http://localhost:8080`
   - `token`: (will be set after login)
3. Set up authentication in collection settings:
   - Type: Bearer Token
   - Token: `{{token}}`
4. Import endpoints from this documentation

---

## Changelog

### Version 1.0 (Current)

- Initial API release
- User authentication with JWT
- Product and category management
- Shopping cart functionality
- Order placement and tracking
- Feedback system
- Admin user management
- Report generation (PDF/CSV)

---

## Support

For issues, questions, or feature requests, please contact the development team or check the project repository.

**Key Files:**

- Routes: [routes/routes.go](routes/routes.go)
- Auth Handlers: [handlers/auth.go](handlers/auth.go)
- Product Handlers: [handlers/product_handlers.go](handlers/product_handlers.go)
- Models: [models/](models/)
- Middleware: [middleware/](middleware/)
