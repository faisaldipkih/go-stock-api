package models

import (
	"encoding/json"
	"time"
)

type Stock struct {
    ID       int    `json:"id"`
    Nama_Barang     string `json:"nama_barang"`
    Jumlah int    `json:"jumlah"`
    Nomor_Seri int    `json:"nomor_seri"`
	AdditionalInfo json.RawMessage `json:"additional_info,omitempty"` // JSON field
	GambarBarang  string          `json:"gambar_barang,omitempty"`
	CreatedAt     time.Time       `json:"created_at"`
    UpdatedAt     time.Time       `json:"updated_at"`
}

func (s *Stock) TableName() string {
    return "stocks"
}
