# README_SQL.md
# PostgreSQL – Basic SQL Command Reference

**Project:** Gunawan Brothers – Inventory Go  
**Database:** PostgreSQL  
**Purpose:** SQL Fundamental & Database Development

---

# 1. DATABASE COMMANDS

### 1.1 Create Database

```sql
CREATE DATABASE inventory_go;
```

### 1.2 Show Databases

Menggunakan PostgreSQL terminal (psql):

```sql
\l
```

### 1.3 Connect Database

```sql
\c inventory_go
```

### 1.4 Drop Database

```sql
DROP DATABASE inventory_go;
```

**Note:** DROP DATABASE menghapus database beserta seluruh isinya.

---

# 2. TABLE COMMANDS

### 2.1 Create Table

```sql
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    product_code VARCHAR(30) UNIQUE NOT NULL,
    product_name VARCHAR(100) NOT NULL,
    unit VARCHAR(20),
    price NUMERIC(15,2) DEFAULT 0,
    stock INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### 2.2 Show Tables

```sql
\dt
```

### 2.3 Describe Table

```sql
\d products
```

### 2.4 Rename Table

```sql
ALTER TABLE products
RENAME TO inventory;
```

### 2.5 Drop Table

```sql
DROP TABLE products;
```

### 2.6 Truncate Table

Menghapus semua records tanpa menghapus struktur tabel.

```sql
TRUNCATE TABLE products;
```

---

# 3. INSERT – CREATE DATA (POST)

Memasukkan data baru ke dalam tabel.

```sql
INSERT INTO products (
    product_code,
    product_name,
    unit,
    price,
    stock
)
VALUES (
    'PRD001',
    'Laptop',
    'PCS',
    7500000,
    10
);
```

### Insert Multiple Records

```sql
INSERT INTO products (
    product_code,
    product_name,
    unit,
    price,
    stock
)
VALUES
('PRD002', 'Keyboard', 'PCS', 250000, 20),
('PRD003', 'Mouse', 'PCS', 150000, 30),
('PRD004', 'Monitor', 'PCS', 1800000, 15);
```

### Insert and Return ID

```sql
INSERT INTO products (
    product_code,
    product_name,
    stock
)
VALUES ('PRD005', 'Printer', 5)
RETURNING id;
```

---

# 4. SELECT – READ DATA (GET)

### 4.1 Get All Records

```sql
SELECT * FROM products;
```

### 4.2 Get Specific Columns

```sql
SELECT
    id,
    product_code,
    product_name,
    stock
FROM products;
```

### 4.3 WHERE – Filter Data

```sql
SELECT *
FROM products
WHERE id = 1;
```

### 4.4 AND / OR

```sql
SELECT *
FROM products
WHERE stock > 10
AND price > 100000;
```

```sql
SELECT *
FROM products
WHERE product_name = 'Laptop'
OR product_name = 'Monitor';
```

### 4.5 LIKE – Search Data

```sql
SELECT *
FROM products
WHERE product_name LIKE '%lap%';
```

### 4.6 ORDER BY – Sorting

```sql
SELECT *
FROM products
ORDER BY product_name ASC;
```

Descending:

```sql
SELECT *
FROM products
ORDER BY id DESC;
```

### 4.7 LIMIT

```sql
SELECT *
FROM products
LIMIT 10;
```

### 4.8 OFFSET – Pagination

```sql
SELECT *
FROM products
ORDER BY id
LIMIT 10 OFFSET 10;
```

### 4.9 BETWEEN

```sql
SELECT *
FROM products
WHERE price BETWEEN 100000 AND 1000000;
```

### 4.10 IN

```sql
SELECT *
FROM products
WHERE id IN (1, 2, 3);
```

---

# 5. UPDATE – REVISE DATA (PUT / PATCH)

Mengubah data yang sudah ada.

### 5.1 Update One Column

```sql
UPDATE products
SET price = 8000000
WHERE id = 1;
```

### 5.2 Update Multiple Columns

```sql
UPDATE products
SET
    price = 8500000,
    stock = 15
WHERE id = 1;
```

### 5.3 Increase Stock

```sql
UPDATE products
SET stock = stock + 10
WHERE id = 1;
```

### 5.4 Decrease Stock

```sql
UPDATE products
SET stock = stock - 5
WHERE id = 1
AND stock >= 5;
```

### 5.5 Update and Return

```sql
UPDATE products
SET stock = stock + 5
WHERE id = 1
RETURNING *;
```

**Important:** Selalu gunakan WHERE saat melakukan UPDATE agar seluruh data tidak berubah secara tidak sengaja.

---

# 6. DELETE – DELETE DATA

### 6.1 Delete One Record

```sql
DELETE FROM products
WHERE id = 1;
```

### 6.2 Delete Multiple Records

```sql
DELETE FROM products
WHERE stock = 0;
```

### 6.3 Delete and Return

```sql
DELETE FROM products
WHERE id = 1
RETURNING *;
```

**Warning:** DELETE tanpa WHERE akan menghapus seluruh records dalam tabel.

---

# 7. JOIN – RELATIONSHIP TABLES

JOIN digunakan untuk menggabungkan data dari beberapa tabel.

Contoh tabel:

- products
- stock_in
- stock_in_details

### 7.1 INNER JOIN

Menampilkan data yang memiliki pasangan di kedua tabel.

```sql
SELECT
    p.product_name,
    d.quantity,
    d.unit_cost
FROM stock_in_details d
INNER JOIN products p
ON d.product_id = p.id;
```

### 7.2 LEFT JOIN

Menampilkan semua produk termasuk yang belum memiliki transaksi Stock In.

```sql
SELECT
    p.product_name,
    d.quantity
FROM products p
LEFT JOIN stock_in_details d
ON p.id = d.product_id;
```

### 7.3 RIGHT JOIN

```sql
SELECT
    p.product_name,
    d.quantity
FROM products p
RIGHT JOIN stock_in_details d
ON p.id = d.product_id;
```

### 7.4 FULL OUTER JOIN

```sql
SELECT *
FROM products p
FULL OUTER JOIN stock_in_details d
ON p.id = d.product_id;
```

### 7.5 JOIN Three Tables

```sql
SELECT
    si.transaction_no,
    si.transaction_date,
    p.product_name,
    d.quantity,
    d.unit_cost
FROM stock_in si
INNER JOIN stock_in_details d
    ON si.id = d.stock_in_id
INNER JOIN products p
    ON d.product_id = p.id;
```

---

# 8. AGGREGATE FUNCTIONS

Digunakan untuk menghitung dan menganalisis data.

### COUNT

```sql
SELECT COUNT(*) FROM products;
```

### SUM

```sql
SELECT SUM(stock) FROM products;
```

### AVG

```sql
SELECT AVG(price) FROM products;
```

### MIN

```sql
SELECT MIN(price) FROM products;
```

### MAX

```sql
SELECT MAX(price) FROM products;
```

### GROUP BY

Menghitung total barang berdasarkan unit.

```sql
SELECT
    unit,
    SUM(stock) AS total_stock
FROM products
GROUP BY unit;
```

### HAVING

```sql
SELECT
    unit,
    SUM(stock) AS total_stock
FROM products
GROUP BY unit
HAVING SUM(stock) > 20;
```

---

# 9. CREATE INDEX

Index digunakan untuk mempercepat pencarian data.

### 9.1 Create Index

```sql
CREATE INDEX idx_products_name
ON products(product_name);
```

### 9.2 Unique Index

```sql
CREATE UNIQUE INDEX idx_products_code
ON products(product_code);
```

### 9.3 Multiple Column Index

```sql
CREATE INDEX idx_stock_in_date
ON stock_in(transaction_date, id);
```

### 9.4 Drop Index

```sql
DROP INDEX idx_products_name;
```

### 9.5 Check Query Performance

```sql
EXPLAIN ANALYZE
SELECT *
FROM products
WHERE product_name = 'Laptop';
```

**Note:** Jangan membuat terlalu banyak index karena index juga menambah beban saat INSERT, UPDATE dan DELETE.

---

# 10. ALTER TABLE

Mengubah struktur tabel.

### Add Column

```sql
ALTER TABLE products
ADD COLUMN category VARCHAR(50);
```

### Change Column Type

```sql
ALTER TABLE products
ALTER COLUMN price TYPE NUMERIC(18,2);
```

### Rename Column

```sql
ALTER TABLE products
RENAME COLUMN price TO unit_price;
```

### Drop Column

```sql
ALTER TABLE products
DROP COLUMN category;
```

### Set Default Value

```sql
ALTER TABLE products
ALTER COLUMN stock SET DEFAULT 0;
```

---

# 11. PRIMARY KEY & FOREIGN KEY

### Primary Key

```sql
CREATE TABLE suppliers (
    id SERIAL PRIMARY KEY,
    supplier_name VARCHAR(100) NOT NULL
);
```

### Foreign Key

```sql
CREATE TABLE purchases (
    id SERIAL PRIMARY KEY,
    supplier_id INTEGER REFERENCES suppliers(id),
    purchase_date DATE DEFAULT CURRENT_DATE
);
```

### Add Foreign Key

```sql
ALTER TABLE stock_in
ADD CONSTRAINT fk_stock_supplier
FOREIGN KEY (supplier_id)
REFERENCES suppliers(id);
```

### Check Constraints

```sql
ALTER TABLE products
ADD CONSTRAINT chk_stock_positive
CHECK (stock >= 0);
```

---

# 12. TRANSACTION – BEGIN / COMMIT / ROLLBACK

Transaction menjaga agar beberapa perintah SQL berhasil sebagai satu kesatuan.

### BEGIN

Memulai transaksi.

### COMMIT

Menyimpan perubahan.

### ROLLBACK

Membatalkan perubahan.

Contoh:

```sql
BEGIN;

UPDATE products
SET stock = stock + 10
WHERE id = 1;

UPDATE products
SET stock = stock - 5
WHERE id = 2;

COMMIT;
```

Jika terjadi kesalahan:

```sql
BEGIN;

UPDATE products
SET stock = stock + 10
WHERE id = 1;

ROLLBACK;
```

**Important:** Dalam Inventory Go, transaksi Stock In harus menggunakan database transaction agar header transaksi, detail barang, dan penambahan stok konsisten.

---

# 13. USEFUL POSTGRESQL FUNCTIONS

### Current Date

```sql
SELECT CURRENT_DATE;
```

### Current Timestamp

```sql
SELECT CURRENT_TIMESTAMP;
```

### COALESCE

Mengganti NULL dengan nilai tertentu.

```sql
SELECT
    product_name,
    COALESCE(stock, 0) AS stock
FROM products;
```

### CONCAT

```sql
SELECT CONCAT(product_code, ' - ', product_name)
FROM products;
```

### ROUND

```sql
SELECT ROUND(AVG(price), 2)
FROM products;
```

### CASE WHEN

```sql
SELECT
    product_name,
    stock,
    CASE
        WHEN stock = 0 THEN 'OUT OF STOCK'
        WHEN stock < 10 THEN 'LOW STOCK'
        ELSE 'AVAILABLE'
    END AS stock_status
FROM products;
```

---

# 14. VIEW

VIEW merupakan query yang disimpan sebagai objek database.

```sql
CREATE VIEW vw_inventory AS
SELECT
    product_code,
    product_name,
    unit,
    price,
    stock
FROM products;
```

Mengakses VIEW:

```sql
SELECT * FROM vw_inventory;
```

Menghapus VIEW:

```sql
DROP VIEW vw_inventory;
```

---

# 15. BACKUP & RESTORE

Menggunakan PostgreSQL command line.

### Backup Database

```bash
pg_dump -U postgres -d inventory_go -F c -f inventory_go.backup
```

### Restore Database

```bash
pg_restore -U postgres -d inventory_go inventory_go.backup
```

### SQL Backup

```bash
pg_dump -U postgres -d inventory_go -f inventory_go.sql
```

---

# 16. SQL COMMAND vs REST API

Hubungan SQL dengan backend Go REST API:

| SQL | HTTP Method | Function |
|---|---|---|
| SELECT | GET | Read data |
| INSERT | POST | Create data |
| UPDATE | PUT | Replace/update data |
| UPDATE | PATCH | Partial update |
| DELETE | DELETE | Delete data |

Contoh API Inventory Go:

```text
GET     /api/products
GET     /api/products/:id
POST    /api/products
PUT     /api/products/:id
PATCH   /api/products/:id
DELETE  /api/products/:id

GET     /api/stock-in
POST    /api/stock-in
GET     /api/stock-in/:id
```

Note: GET, POST, PUT, PATCH dan DELETE adalah HTTP methods, bukan perintah SQL.

---

# 17. INVENTORY GO – STOCK IN EXAMPLE

Contoh query untuk melihat transaksi Stock In.

```sql
SELECT
    si.id,
    si.transaction_no,
    si.transaction_date,
    p.product_name,
    d.quantity,
    d.unit_cost,
    (d.quantity * d.unit_cost) AS total_amount
FROM stock_in si
JOIN stock_in_details d
    ON si.id = d.stock_in_id
JOIN products p
    ON d.product_id = p.id
ORDER BY si.transaction_date DESC;
```

### Calculate Total Stock In

```sql
SELECT
    p.product_name,
    SUM(d.quantity) AS total_stock_in
FROM stock_in_details d
JOIN products p
    ON d.product_id = p.id
GROUP BY p.id, p.product_name
ORDER BY p.product_name;
```

---

# 18. SQL QUICK REFERENCE

| Command | Function |
|---|---|
| CREATE DATABASE | Create database |
| CREATE TABLE | Create table |
| INSERT INTO | Insert data |
| SELECT | Read data |
| WHERE | Filter records |
| UPDATE | Revise data |
| DELETE | Delete records |
| JOIN | Combine tables |
| GROUP BY | Group records |
| ORDER BY | Sort records |
| COUNT | Count records |
| SUM | Calculate total |
| ALTER TABLE | Change table |
| DROP TABLE | Remove table |
| TRUNCATE | Clear table |
| CREATE INDEX | Improve search performance |
| BEGIN | Start transaction |
| COMMIT | Save transaction |
| ROLLBACK | Cancel transaction |
| CREATE VIEW | Create database view |
| EXPLAIN ANALYZE | Analyze query execution |

---

**END OF README_SQL.md**

Gunawan Brothers – Inventory Go  
Database Technology: PostgreSQL  
Backend Technology: Go (Golang)
