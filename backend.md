
# Medical Equipment Online Store 

## Requirements Backend

### Tech Stack
- Golang
- GORM + MySQL
- Gin
- JWT
- Unidoc to generate PDF
- Stripe for simulate payment
- Cloudinary for image storage

### Functional Requirements
1. Auth
- Register (Visitor → Customer)
- Login (Admin & Customer)
- JWT Authentication
- Logout
2. User Management (Admin)
- CRUD Customer (view/update/delete)
- Manage Admin account
3. Products Management (CRUD)
- CRUD Product (Admin)
- CRUD Category (Admin)
- Browse Product (Visitor/Customer)
4. Add Cart
- Add/Remove Product from cart (Customer)
- View Cart (Customer)
5. Checkout / Order
- Place Order (Customer)
- Cancel before shipment (Customer)
- Shipping order (Admin update status)
- Payment (simulate CC/Debit/Paypal/Mock Stripe)
6. Feedback
- Customer give feedback
7. Generate Report
- Generate PDF transaction report (Admin/Customer allow to download)


### ERD (Entity Relationship Diagram)
Users

| Field           | Type                  | Explanation    |
| --------------- | --------------------- |----------------|
| id (PK)         | UUID / INT            | Primary key    |
| username        | VARCHAR               | Unique         |
| password        | VARCHAR               | Hashed         |
| email           | VARCHAR               | Unique         |
| dob             | DATE                  | Date of birth  |
| gender          | ENUM(M/F)             | Gender         |
| address         | TEXT                  | Address        |
| city            | VARCHAR               | City           |
| contact\_number | VARCHAR               | Contact Number |
| paypal\_id      | VARCHAR NULL          | Optional       |
| role            | ENUM(admin, customer) | User role      |
| created\_at     | TIMESTAMP             | time created   |
| updated\_at     | TIMESTAMP             | time updated   |


Categories


| Field       | Type      | Explanation   |
| ----------- | --------- |---------------|
| id (PK)     | INT       | Primary key   |
| name        | VARCHAR   | category name |
| description | TEXT      | description   |
| created\_at | TIMESTAMP |               |


Products

| Field        | Type      | Explanation         |
|--------------| --------- |---------------------|
| id (PK)      | INT       | Primary key         |
| category\_id | INT (FK)  | related to category |
| name         | VARCHAR   | product name        |
| description  | TEXT      | product description |
| price        | DECIMAL   | price               |
| stock        | INT       | stock               |
| created\_at  | TIMESTAMP |                     |
| updated\_at  | TIMESTAMP |                     |
| images_url   | VARCHAR   | images of product   |

Carts

| Field       | Type      | Explanation      |
| ----------- | --------- |------------------|
| id (PK)     | INT       | Primary key      |
| user\_id FK | INT       | related to users |
| created\_at | TIMESTAMP |                  |


Cart Items

| Field          | Type | Explanation         |
| -------------- | ---- |---------------------|
| id (PK)        | INT  | Primary key         |
| cart\_id FK    | INT  | related to carts    |
| product\_id FK | INT  | related to products |
| quantity       | INT  | item quantity       |


Orders

| Field           | Type                                    | Explanation      |
| --------------- | --------------------------------------- |------------------|
| id (PK)         | INT                                     | Primary key      |
| user\_id FK     | INT                                     | related to users |
| status          | ENUM(pending, paid, shipped, cancelled) | order status     |
| total\_price    | DECIMAL                                 | total price      |
| payment\_method | ENUM(paypal, debit, cc, cod)            | payment method   |
| bank\_name      | VARCHAR NULL                            | optional         |
| created\_at     | TIMESTAMP                               |                  |


Oder Items

| Field          | Type    | Explanation        |
| -------------- | ------- |--------------------|
| id (PK)        | INT     | Primary key        |
| order\_id FK   | INT     | related to orders  |
| product\_id FK | INT     | related to product |
| quantity       | INT     | buy quantity       |
| price          | DECIMAL | price while order  |


Feedbacks

| Field          | Type      | Explanation        |
| -------------- | --------- |--------------------|
| id (PK)        | INT       | Primary key        |
| user\_id FK    | INT       | related to users   |
| product\_id FK | INT       | related to product |
| comment        | TEXT      | comment            |
| rating         | INT (1–5) | Rating             |
| created\_at    | TIMESTAMP |                    |


### Table relationship
- User (1 → N) Orders
- User (1 → 1) Cart
- Cart (1 → N) CartItems → Products
- Category (1 → N) Products
- Order (1 → N) OrderItems → Products
- User (1 → N) Feedback → Products







