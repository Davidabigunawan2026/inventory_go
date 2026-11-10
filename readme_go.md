# README_Go.md
# PANDUAN DASAR GOLANG (GO)

**Project:** Inventory Go  
**Language:** Go (Golang)  
**Backend:** Go HTTP Server  
**Database:** PostgreSQL  
**Purpose:** Basic Go Commands, Syntax, and Backend Development

---

# 1. PENGENALAN GOLANG

Go adalah bahasa pemrograman yang dikembangkan oleh Google.

Kegunaan Go:
- Backend Web Development
- REST API
- Database Application
- Microservices
- Cloud Application
- CLI Application
- Concurrent Programming

Kelebihan Go:
- Cepat dan efisien
- Syntax relatif sederhana
- Mendukung concurrency dengan goroutines
- Memiliki HTTP server bawaan
- Memiliki compiler sendiri

# 2. INSTALLASI DAN CEK GO

Cek versi Go:

```bash
go version
```

Cek lokasi instalasi:

```bash
go env GOROOT
go env GOPATH
```

Melihat seluruh konfigurasi:

```bash
go env
```

Cek lokasi compiler:

```bash
where go
```

# 3. PERINTAH DASAR GO

| Perintah | Fungsi |
|---|---|
| go version | Melihat versi Go |
| go help | Bantuan Go |
| go env | Konfigurasi Go |
| go mod init | Membuat module |
| go mod tidy | Mengatur dependencies |
| go run . | Menjalankan aplikasi |
| go run main.go | Menjalankan file |
| go build | Compile aplikasi |
| go install | Compile dan install |
| go fmt ./... | Format coding |
| go vet ./... | Memeriksa kemungkinan kesalahan |
| go test ./... | Menjalankan testing |
| go clean | Membersihkan hasil build |

# 4. MEMBUAT PROJECT GO

Contoh membuat project Inventory Go:

```bash
mkdir inventory_go
cd inventory_go
```

Membuat Go module:

```bash
go mod init inventory_go
```

Menjalankan project:

```bash
go run .
```

Compile project:

```bash
go build
```

Memperbarui dependencies:

```bash
go mod tidy
```

# 5. HELLO WORLD

Buat file `main.go`:

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello World!")
    fmt.Println("Inventory Go System")
}
```

Jalankan:

```bash
go run main.go
```

Output:

```text
Hello World!
Inventory Go System
```

Penjelasan:

- package main = program utama
- import fmt = menggunakan package format
- func main() = fungsi pertama yang dijalankan
- fmt.Println() = menampilkan teks

# 6. VARIABLE (VARIABEL)

Membuat variabel:

```go
package main

import "fmt"

func main() {

    var name string = "David"
    var age int = 60
    var price float64 = 15000.50
    var active bool = true

    fmt.Println(name)
    fmt.Println(age)
    fmt.Println(price)
    fmt.Println(active)
}
```

Menggunakan short declaration:

```go
name := "David"
age := 60
price := 15000
```

Catatan:
- string = teks
- int = bilangan bulat
- float64 = bilangan desimal
- bool = true atau false

# 7. CONSTANT

Nilai tetap menggunakan const:

```go
const company = "Gunawan Brothers"
const taxRate = 0.11
```

Nilai constant tidak dapat diubah setelah dideklarasikan.

# 8. OPERATOR

Arithmetic:

```go
a := 10
b := 5

fmt.Println(a + b)
fmt.Println(a - b)
fmt.Println(a * b)
fmt.Println(a / b)
fmt.Println(a % b)
```

Comparison:

```go
a == b
a != b
a > b
a < b
a >= b
a <= b
```

Logical:

```go
&&  // AND
||  // OR
!   // NOT
```

# 9. IF ELSE

```go
package main

import "fmt"

func main() {

    stock := 10

    if stock > 0 {
        fmt.Println("Stock Available")
    } else {
        fmt.Println("Stock Empty")
    }
}
```

# 10. SWITCH

```go
package main

import "fmt"

func main() {

    status := "active"

    switch status {

    case "active":
        fmt.Println("User Active")

    case "inactive":
        fmt.Println("User Inactive")

    default:
        fmt.Println("Unknown Status")
    }
}
```

# 11. LOOPING (FOR)

Go menggunakan `for` sebagai perulangan.

```go
for i := 1; i <= 5; i++ {
    fmt.Println(i)
}
```

Menggunakan while-style loop:

```go
i := 1

for i <= 5 {
    fmt.Println(i)
    i++
}
```

# 12. ARRAY DAN SLICE

Array:

```go
var numbers [3]int

numbers[0] = 10
numbers[1] = 20
numbers[2] = 30

fmt.Println(numbers)
```

Slice:

```go
products := []string{
    "Laptop",
    "Keyboard",
    "Mouse",
}

fmt.Println(products)
```

Menambahkan data:

```go
products = append(products, "Monitor")
```

Looping slice:

```go
for index, product := range products {
    fmt.Println(index, product)
}
```

# 13. MAP

Map menyimpan data dalam pasangan key dan value.

```go
product := map[string]interface{}{

    "id": 1,
    "name": "Laptop",
    "price": 7500000,
    "stock": 10,
}

fmt.Println(product["name"])
fmt.Println(product["price"])
```

# 14. FUNCTION

Membuat function:

```go
func greeting() {
    fmt.Println("Welcome to Inventory Go")
}
```

Memanggil function:

```go
func main() {
    greeting()
}
```

Function dengan parameter:

```go
func add(a int, b int) int {

    return a + b
}
```

Memanggil:

```go
result := add(10, 20)

fmt.Println(result)
```

# 15. STRUCT

Struct digunakan untuk membuat struktur data.

Contoh Product:

```go
type Product struct {

    ID    int
    Name  string
    Price float64
    Stock int
}
```

Membuat object:

```go
product := Product{

    ID: 1,
    Name: "Laptop",
    Price: 7500000,
    Stock: 10,
}

fmt.Println(product.Name)
```

Struct dapat digunakan sebagai representasi tabel database.

# 16. POINTER

Pointer menyimpan alamat memori suatu variabel.

```go
package main

import "fmt"

func main() {

    price := 10000

    pointer := &price

    fmt.Println(price)
    fmt.Println(pointer)
    fmt.Println(*pointer)
}
```

Keterangan:

- & = mengambil alamat memori
- * = mengambil nilai dari alamat pointer

# 17. ERROR HANDLING

Go menggunakan error sebagai nilai yang dikembalikan oleh function.

```go
package main

import (
    "errors"
    "fmt"
)

func divide(a, b int) (int, error) {

    if b == 0 {
        return 0, errors.New("division by zero")
    }

    return a / b, nil
}

func main() {

    result, err := divide(10, 0)

    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    fmt.Println(result)
}
```

Pola umum:

```go
result, err := someFunction()

if err != nil {
    return err
}
```

# 18. PACKAGE DAN IMPORT

Contoh:

```go
package main

import (
    "fmt"
    "time"
    "net/http"
)
```

Package berguna untuk memisahkan fungsi dan tanggung jawab program.

Contoh package bawaan:

| Package | Kegunaan |
|---|---|
| fmt | Input dan output |
| time | Tanggal dan waktu |
| net/http | HTTP Server |
| encoding/json | JSON |
| os | Sistem operasi |
| log | Logging |
| strings | Manipulasi teks |
| strconv | Konversi tipe data |
| database/sql | Database SQL |
| context | Context management |

# 19. JSON ENCODING

Go dapat mengubah Struct menjadi JSON.

```go
package main

import (
    "encoding/json"
    "fmt"
)

type Product struct {
    ID    int     `json:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}

func main() {

    product := Product{
        ID: 1,
        Name: "Laptop",
        Price: 7500000,
    }

    data, err := json.Marshal(product)

    if err != nil {
        panic(err)
    }

    fmt.Println(string(data))
}
```

Output:

```json
{
  "id": 1,
  "name": "Laptop",
  "price": 7500000
}
```

# 20. MEMBUAT HTTP SERVER

Buat file `main.go`:

```go
package main

import (
    "fmt"
    "net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {

    fmt.Fprintln(w, "Inventory Go API is running!")
}

func main() {

    http.HandleFunc("/", homeHandler)

    fmt.Println("Server running on port 8080")

    http.ListenAndServe(":8080", nil)
}
```

Jalankan:

```bash
go run main.go
```

Buka browser:

```text
http://localhost:8080
```

# 21. MEMBUAT REST API JSON

Contoh endpoint sederhana:

```go
package main

import (
    "encoding/json"
    "net/http"
)

type Product struct {

    ID    int     `json:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price"`
    Stock int     `json:"stock"`
}

func getProducts(w http.ResponseWriter, r *http.Request) {

    w.Header().Set(
        "Content-Type",
        "application/json",
    )

    products := []Product{

        {
            ID: 1,
            Name: "Laptop",
            Price: 7500000,
            Stock: 10,
        },

        {
            ID: 2,
            Name: "Keyboard",
            Price: 250000,
            Stock: 20,
        },
    }

    json.NewEncoder(w).Encode(products)
}

func main() {

    http.HandleFunc("/api/products", getProducts)

    http.ListenAndServe(":8080", nil)
}
```

Endpoint:

```text
GET http://localhost:8080/api/products
```

# 22. DATABASE POSTGRESQL

Go dapat terhubung dengan PostgreSQL menggunakan driver PostgreSQL.

Contoh menggunakan pgx:

```bash
go get github.com/jackc/pgx/v5/stdlib
```

Contoh import:

```go
import (
    "database/sql"
    _ "github.com/jackc/pgx/v5/stdlib"
)
```

Membuka koneksi:

```go
db, err := sql.Open(
    "pgx",
    "postgres://username:password@localhost:5432/inventory",
)

if err != nil {
    panic(err)
}

defer db.Close()

err = db.Ping()

if err != nil {
    panic(err)
}
```

Catatan:
- Ganti username dan password sesuai PostgreSQL.
- Gunakan environment variables untuk menyimpan credential.
- Jangan menyimpan password database langsung di GitHub.

# 23. GOROUTINE

Goroutine digunakan untuk menjalankan pekerjaan secara concurrent.

```go
package main

import (
    "fmt"
    "time"
)

func process() {

    fmt.Println("Processing...")
}

func main() {

    go process()

    time.Sleep(time.Second)
}
```

Keyword:

```go
go functionName()
```

Untuk aplikasi produksi, gunakan synchronization seperti WaitGroup atau channel jika perlu menunggu pekerjaan selesai.

# 24. STRUKTUR PROJECT INVENTORY GO

Contoh struktur backend:

```text
inventory_go/
│
├── go.mod
├── go.sum
│
├── main.go
│
├── internal/
│   │
│   ├── database/
│   │   └── postgres.go
│   │
│   ├── models/
│   │   └── product.go
│   │
│   ├── repositories/
│   │   └── product_repository.go
│   │
│   ├── handlers/
│   │   └── product_handler.go
│   │
│   └── routes/
│       └── routes.go
│
└── config/
    └── config.go
```

Pembagian tugas:

| Folder | Fungsi |
|---|---|
| database | Koneksi database |
| models | Struktur data |
| repositories | SQL Query |
| handlers | HTTP Request/Response |
| routes | Pengaturan endpoint |
| config | Konfigurasi aplikasi |
| main.go | Entry point |

# 25. WORKFLOW CODING GO

Urutan kerja ketika mengembangkan aplikasi:

1. Membuat project.
2. Membuat Go module.
3. Membuat database connection.
4. Membuat model.
5. Membuat repository.
6. Membuat handler.
7. Membuat routes.
8. Menjalankan server.
9. Testing API.
10. Build aplikasi.

Perintah yang sering digunakan:

```bash
go mod tidy

go fmt ./...

go vet ./...

go test ./...

go run .

go build -o inventory.exe .
```

Menjalankan hasil build:

```bash
./inventory.exe
```

# 26. PERBANDINGAN GO DENGAN JAVASCRIPT

| Konsep | Go | JavaScript |
|---|---|---|
| Variable | var / := | let / const |
| Constant | const | const |
| Function | func | function |
| Object structure | struct | object / class |
| Array | array | array |
| Collection | slice | array |
| Key-value | map | object / Map |
| Loop | for | for / while |
| Condition | if / else | if / else |
| Error | error | Error / try-catch |
| JSON | encoding/json | JSON |
| HTTP Server | net/http | Node.js / Express |
| Database | database/sql | pg / mysql2 |

# 27. LATIHAN DASAR GO

Latihan 1:
Buat program menampilkan nama perusahaan.

Latihan 2:
Buat variabel Product Name, Price dan Stock.

Latihan 3:
Buat function untuk menghitung total harga.

Latihan 4:
Buat Struct Product.

Latihan 5:
Buat Slice berisi 5 produk.

Latihan 6:
Buat Map untuk data Customer.

Latihan 7:
Buat HTTP Server port 8080.

Latihan 8:
Buat REST API GET /api/products.

Latihan 9:
Hubungkan Go dengan PostgreSQL.

Latihan 10:
Buat CRUD Product lengkap.

---

# 28. KESIMPULAN

Golang memiliki syntax yang relatif sederhana, tetapi sangat powerful untuk backend application.

Untuk pengembangan Gunawan Brothers Inventory Go, fokus pembelajaran adalah:

1. Go Fundamentals
2. Struct and Function
3. Package and Modules
4. Error Handling
5. HTTP Server
6. REST API
7. PostgreSQL
8. CRUD
9. Authentication
10. Deployment

**Target akhir:**

Membangun Inventory Management System menggunakan Go + PostgreSQL dengan REST API yang dapat diakses oleh frontend.

---

**Documentation:** README_Go.md  
**Company:** Gunawan Brothers  
**Project:** Inventory Go  
**Version:** 1.0

cmd backend :
go run ./cmd


cmd frontend :
npm run dev
http://localhost:5173


