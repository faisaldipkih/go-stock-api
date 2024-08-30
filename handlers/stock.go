package handlers

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "path/filepath"
    "strconv"
    // "strings"
	"time"

    "go-stock-api/db"
    "go-stock-api/models"

    "github.com/gorilla/mux"
)

// Response structure
type response struct {
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

func GetStocks(w http.ResponseWriter, r *http.Request) {
    rows, err := db.DB.Query("SELECT id, nama_barang, jumlah, nomor_seri, additional_info, created_at, updated_at FROM stocks")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var stocks []models.Stock
    for rows.Next() {
        var stock models.Stock
        if err := rows.Scan(&stock.ID, &stock.Nama_Barang, &stock.Jumlah, &stock.Nomor_Seri, &stock.AdditionalInfo, &stock.CreatedAt, &stock.UpdatedAt); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        stocks = append(stocks, stock)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response{
        Message: "success",
        Data:    stocks,
    })
}

func GetStock(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, _ := strconv.Atoi(params["id"])

    var stock models.Stock
    err := db.DB.QueryRow("SELECT id, nama_barang, jumlah, nomor_seri, additional_info, created_at, updated_at FROM stocks WHERE id = $1", id).Scan(&stock.ID, &stock.Nama_Barang, &stock.Jumlah, &stock.Nomor_Seri, &stock.AdditionalInfo, &stock.CreatedAt, &stock.UpdatedAt)
    if err == sql.ErrNoRows {
        w.WriteHeader(http.StatusNotFound)
        json.NewEncoder(w).Encode(response{
            Message: "Stock not found",
        })
        return
    } else if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response{
        Message: "success",
        Data:    stock,
    })
}

func CreateStock(w http.ResponseWriter, r *http.Request) {
    var stock models.Stock

    // Parse form data
    r.ParseMultipartForm(10 << 20) // max memory 10MB

    // Read form values
    stock.Nama_Barang = r.FormValue("nama_barang")
    stock.Jumlah, _ = strconv.Atoi(r.FormValue("jumlah"))

    additionalInfo := r.FormValue("additional_info")
    if additionalInfo != "" {
        stock.AdditionalInfo = json.RawMessage(additionalInfo)
    }

    // Handle file upload
    file, handler, err := r.FormFile("gambar_barang")
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    defer file.Close()

    // Create file path
    fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(handler.Filename))
    filePath := filepath.Join("uploads", fileName)

    // Create file on disk
    f, err := os.Create(filePath)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer f.Close()

    // Copy file data to disk
    _, err = io.Copy(f, file)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Store file path in database
    stock.GambarBarang = filePath

    // Insert into database
    err = db.DB.QueryRow(
        "INSERT INTO stocks(nama_barang, jumlah, nomor_seri, additional_info, gambar_barang) VALUES($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at",
        stock.Nama_Barang, stock.Jumlah, stock.Nomor_Seri, stock.AdditionalInfo, stock.GambarBarang,
    ).Scan(&stock.ID, &stock.CreatedAt, &stock.UpdatedAt)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response{
        Message: "Stock created successfully",
        Data:    stock,
    })
}

func UpdateStock(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, _ := strconv.Atoi(params["id"])

    var stock models.Stock

    // Parse form data
    r.ParseMultipartForm(10 << 20)

    // Read form values
    stock.Nama_Barang = r.FormValue("nama_barang")
    stock.Jumlah, _ = strconv.Atoi(r.FormValue("jumlah"))

    additionalInfo := r.FormValue("additional_info")
    if additionalInfo != "" {
        stock.AdditionalInfo = json.RawMessage(additionalInfo)
    }

    // Handle file upload
    file, handler, err := r.FormFile("gambar_barang")
    if err == nil {
        defer file.Close()

        // Create file path
        fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(handler.Filename))
        filePath := filepath.Join("uploads", fileName)

        // Create file on disk
        f, err := os.Create(filePath)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        defer f.Close()

        // Copy file data to disk
        _, err = io.Copy(f, file)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // Store new file path in database
        stock.GambarBarang = filePath
    }

    // Update in database
    _, err = db.DB.Exec(
        "UPDATE stocks SET nama_barang=$1, jumlah=$2, nomor_seri=$3, additional_info=$4, gambar_barang=COALESCE(NULLIF($5, ''), gambar_barang), updated_at=CURRENT_TIMESTAMP WHERE id=$6",
        stock.Nama_Barang, stock.Jumlah, stock.Nomor_Seri, stock.AdditionalInfo, stock.GambarBarang, id,
    )
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response{
        Message: "Stock updated successfully",
        Data:    stock,
    })
}

func DeleteStock(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, _ := strconv.Atoi(params["id"])

    _, err := db.DB.Exec("DELETE FROM stocks WHERE id = $1", id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response{
        Message: "Stok berhasil di hapus",
    })
}
