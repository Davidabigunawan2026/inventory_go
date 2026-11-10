
# README BACKEND — INVENTORY GO

## 1. Mengapa Backend Menggunakan 3 File Go?

Pada project Inventory Go, satu fitur seperti Stock In biasanya menggunakan 3 bagian:

1. Models
2. Repositories
3. Handlers

Tujuannya adalah memisahkan tugas agar program lebih terstruktur, mudah diperbaiki, dan mudah dikembangkan.

Ketiga file tersebut bukan berarti membuat 3 tabel database.

Satu fitur dapat menggunakan beberapa file Go, tetapi tetap mengakses tabel PostgreSQL yang sama.

---

## 2. Struktur Folder Backend

```text
backend/
│
├── cmd/
│   └── main.go
│
├── internal/
│   ├── database/
│   └── handlers/
│       ├── stock_in_handler.go
│       └── stock_out_handler.go
│
├── models/
│   ├── stock_in.go
│   └── stock_out.go
│
└── repositories/
    ├── stock_in_repository.go
    └── stock_out_repository.go
```

---

## 3. MODELS — Struktur Data

Lokasi:

`backend/models/stock_in.go`

Models digunakan untuk mendefinisikan struktur data yang akan digunakan oleh program Go.

Contoh:

```go
package models

type StockInTransaction struct {
    ID            int64   `json:"id"`
    TransactionNo string  `json:"transaction_no"`
    TotalQuantity float64 `json:"total_quantity"`
    TotalAmount   float64 `json:"total_amount"`
}
```

Penjelasan:

- ID: identitas transaksi.
- TransactionNo: nomor transaksi.
- TotalQuantity: total jumlah barang.
- TotalAmount: total nilai transaksi.
- json: nama field ketika dikirim ke frontend.

Models tidak menjalankan SQL.

Models hanya mendefinisikan bentuk dan tipe data.

Analogi:

**Models = Formulir Data**

---

## 4. REPOSITORIES — Akses Database

Lokasi:

`backend/repositories/stock_in_repository.go`

Repositories bertugas melakukan komunikasi dengan PostgreSQL.

Contoh:

```go
type StockInRepository struct {
    DB *pgxpool.Pool
}

func NewStockInRepository(
    db *pgxpool.Pool,
) *StockInRepository {
    return &StockInRepository{DB: db}
}
```

Contoh query:

```go
query := `
    SELECT
        si.id,
        si.transaction_no,
        COALESCE(SUM(d.quantity), 0)
    FROM stock_in si
    LEFT JOIN stock_in_details d
        ON d.stock_in_id = si.id
    GROUP BY si.id
`
```

Fungsi Repositories:

1. Membaca data dengan SELECT.
2. Menambah data dengan INSERT.
3. Mengubah data dengan UPDATE.
4. Menghapus data dengan DELETE.
5. Mengirim hasil query ke bagian program yang memanggilnya.
6. Mengembalikan error jika terjadi masalah database.

Repositories tidak mengatur tampilan Vue.

Analogi:

**Repositories = Petugas Database**

---

## 5. HANDLERS — Komunikasi API

Lokasi:

`backend/internal/handlers/stock_in_handler.go`

Handlers bertugas menerima permintaan HTTP dari frontend dan memberikan response.

Contoh:

```go
func (h *StockInHandler) GetAllStockIn(
    w http.ResponseWriter,
    r *http.Request,
) {
    transactions, err :=
        h.Repo.GetAllStockIn(r.Context())

    if err != nil {
        http.Error(
            w,
            "Failed to retrieve stock in",
            http.StatusInternalServerError,
        )
        return
    }

    w.Header().Set(
        "Content-Type",
        "application/json",
    )

    json.NewEncoder(w).Encode(
        map[string]interface{}{
            "data": transactions,
            "status": "success",
        },
    )
}
```

Penjelasan:

- Menerima HTTP request.
- Memanggil Repository.
- Memeriksa error.
- Mengirim HTTP status.
- Mengubah hasil menjadi JSON.
- Mengirim response ke frontend.

Analogi:

**Handlers = Petugas Loket API**

---

## 6. MAIN.GO — Penghubung Semua Bagian

Lokasi:

`backend/cmd/main.go`

Main.go bertugas menghubungkan Repository, Handler, dan HTTP Router.

Contoh:

```go
stockInRepo :=
    repositories.NewStockInRepository(db)

stockInHandler :=
    handlers.NewStockInHandler(stockInRepo)

mux.HandleFunc(
    "GET /api/stock-in",
    stockInHandler.GetAllStockIn,
)
```

Artinya:

Ketika frontend meminta:

`GET http://localhost:8080/api/stock-in`

Server akan menjalankan:

`stockInHandler.GetAllStockIn`

Handler kemudian memanggil Repository untuk membaca PostgreSQL.

---

## 7. Alur Lengkap Frontend ke Database

```text
          VUE FRONTEND
               |
               |
     GET /api/stock-in
               |
               v
          GO MAIN.GO
          HTTP ROUTER
               |
               v
          HANDLERS
       GetAllStockIn()
               |
               v
         REPOSITORIES
       GetAllStockIn()
               |
               v
          POSTGRESQL
       stock_in tables
               |
               v
         REPOSITORIES
        Scan SQL Result
               |
               v
            MODELS
     StockInTransaction
               |
               v
           HANDLERS
        JSON Response
               |
               v
          VUE FRONTEND
        Display Table
```

Catatan:

Models bukan proses terpisah yang harus dijalankan. Models adalah struktur data yang digunakan oleh Repository dan Handler.

---

## 8. Mengapa Tidak Digabung dalam Satu File?

Secara teknis, semua kode bisa ditulis dalam satu file Go.

Namun, ketika aplikasi semakin besar, program akan sulit dikelola.

Keuntungan pemisahan:

### A. Mudah Diperbaiki

Jika SQL bermasalah:

Periksa Repositories.

Jika JSON response bermasalah:

Periksa Handlers.

Jika struktur field bermasalah:

Periksa Models.

Jika route API tidak ditemukan:

Periksa main.go.

### B. Mudah Dikembangkan

Misalnya menambahkan:

- Inventory Dashboard
- Purchase Order
- Sales Order
- Stock Adjustment
- Stock Transfer
- Inventory Valuation

Kita dapat mengikuti pola struktur yang sama.

### C. Mudah Diuji

Setiap bagian dapat diuji secara terpisah.

### D. Mengurangi Duplikasi

Satu model dapat digunakan oleh beberapa fungsi Repository atau Handler.

---

## 9. Hubungan dengan Bahasa Pemrograman Lain

Konsep ini juga dikenal dalam pengembangan aplikasi lain.

| Go Backend | Konsep |
|---|---|
| Models | Data Structure / Entity |
| Repositories | Data Access Layer |
| Handlers | HTTP Controller |
| main.go | Application Bootstrap / Router |

Perbandingan sederhana dengan sistem database tradisional:

- Models mirip definisi struktur record.
- Repositories mirip modul prosedur SQL.
- Handlers mirip controller yang menerima permintaan pengguna.
- main.go mirip program utama yang menghubungkan modul.

Ini adalah analogi konsep, bukan persamaan sintaks satu banding satu.

---

## 10. Contoh Debugging Inventory Go

Jika tabel Stock In kosong di frontend:

Langkah 1:

Periksa apakah server Go berjalan.

Langkah 2:

Buka:

`http://localhost:8080/api/stock-in`

Langkah 3:

Jika API mengembalikan error 500:

Periksa Handler, Repository, dan PostgreSQL.

Langkah 4:

Jika API mengembalikan data JSON dengan benar:

Periksa kode Vue, fetch API, dan tampilan tabel.

Langkah 5:

Jika terjadi error scan atau field:

Periksa kecocokan query SQL dan Models.

---

## 11. Pengembangan Backend Lebih Lanjut

Untuk aplikasi ERP yang lebih kompleks, kita dapat menambahkan:

### Services

Berisi business logic atau aturan bisnis.

Contoh:

- Menghitung Moving Weighted Average.
- Memeriksa stok sebelum Stock Out.
- Menghitung nilai persediaan.
- Memproses transaksi multi-tabel.

Contoh struktur:

```text
backend/
├── cmd/
├── internal/
│   ├── database/
│   ├── handlers/
│   └── services/
├── models/
└── repositories/
```

Services tidak wajib untuk setiap aplikasi kecil, tetapi berguna ketika aturan bisnis semakin kompleks.

---

## 12. Kesimpulan

Cara mudah mengingat:

**MODELS = Bentuk Data**

**REPOSITORIES = Akses Database**

**HANDLERS = Komunikasi API**

**MAIN.GO = Penghubung Komponen**

Alur utamanya:

Frontend -> Handler -> Repository -> PostgreSQL

Hasilnya:

PostgreSQL -> Repository -> Handler -> JSON -> Frontend

Dengan struktur ini, Inventory Go lebih mudah dikembangkan menjadi sistem ERP profesional.

---

Gunawan Brothers — Inventory Go

Backend: Golang + PostgreSQL

Frontend: Vue 3 + Vite

Dokumentasi pembelajaran backend.
