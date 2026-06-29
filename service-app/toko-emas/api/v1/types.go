package v1

import "github.com/google/uuid"

// ---- Common ----

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type AuthLoginRequest struct {
	Username string `json:"username" example:"budi"`
	Password string `json:"password" example:"secret123"`
}

type AuthRefreshRequest struct {
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type AuthTokenResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	ExpiresIn    int64       `json:"expires_in"`
	TokenType    string      `json:"token_type"`
	User         interface{} `json:"user,omitempty"`
}

// ---- User ----

type CreateUserRequest struct {
	Nama     string `json:"nama" example:"Budi Santoso"`
	Username string `json:"username" example:"budi"`
	Password string `json:"password" example:"secret123"`
	Role     string `json:"role" example:"kasir"`
}

type UpdateUserRequest struct {
	Nama     string `json:"nama" example:"Budi Santoso"`
	Username string `json:"username" example:"budi"`
	Role     string `json:"role" example:"kasir"`
}

// ---- Barang ---
type CreateBarangRequest struct {
	Barcode     string     `json:"barcode" example:"BC001"`
	Nama        string     `json:"nama" example:"Gelang Emas"`
	Karat       int        `json:"karat" example:"24"`
	Berat       float64    `json:"berat" example:"5.5"`
	Harga       int64      `json:"harga" example:"5000000"`
	Photo       string     `json:"photo" example:"https://example.com/photo.jpg"`
	Kondisi     string     `json:"kondisi" example:"baru"`
	PembelianID *uuid.UUID `json:"pembelian_id,omitempty" example:"null"`
	BakiID      *uuid.UUID `json:"baki_id,omitempty" example:"null"`
	GrupID      *uuid.UUID `json:"grup_id,omitempty" example:"null"`
}

type UpdateBarangRequest struct {
	Barcode     string     `json:"barcode" example:"BC001"`
	Nama        string     `json:"nama" example:"Gelang Emas"`
	Karat       int        `json:"karat" example:"24"`
	Berat       float64    `json:"berat" example:"5.5"`
	Harga       int64      `json:"harga" example:"5000000"`
	Photo       string     `json:"photo" example:"https://example.com/photo.jpg"`
	Kondisi     string     `json:"kondisi" example:"baru"`
	PembelianID *uuid.UUID `json:"pembelian_id,omitempty"`
	BakiID      *uuid.UUID `json:"baki_id,omitempty"`
	GrupID      *uuid.UUID `json:"grup_id,omitempty"`
}

// ---- Pembelian ----

type CreatePembelianRequest struct {
	NoFaktur    string `json:"no_faktur" example:"INV-2024-001"`
	Nama        string `json:"nama" example:"Toko Mas Jaya"`
	TipePemasok string `json:"tipe_pemasok" example:"supplier"`
	HargaDeal   int64  `json:"harga_deal" example:"10000000"`
}

type UpdatePembelianRequest struct {
	NoFaktur    string `json:"no_faktur" example:"INV-2024-001"`
	Nama        string `json:"nama" example:"Toko Mas Jaya"`
	TipePemasok string `json:"tipe_pemasok" example:"supplier"`
	HargaDeal   int64  `json:"harga_deal" example:"10000000"`
}

// ---- Karat ----

type CreateKaratRequest struct {
	Name  string `json:"name" example:"24K"`
	Harga int64  `json:"harga" example:"1000000"`
}

type UpdateKaratRequest struct {
	Name  string `json:"name" example:"24K"`
	Harga int64  `json:"harga" example:"1000000"`
}

// ---- Baki ----

type CreateBakiRequest struct {
	Nama string `json:"nama" example:"Baki 1"`
}

type UpdateBakiRequest struct {
	Nama string `json:"nama" example:"Baki 1"`
}

// ---- Penjualan ----

type CreatePenjualanRequest struct {
	NoFaktur   string    `json:"no_faktur" example:"SELL-2024-001"`
	Nama       string    `json:"nama" example:"Pelanggan A"`
	TotalHarga int64     `json:"total_harga" example:"5000000"`
	KodeSales  uuid.UUID `json:"kode_sales" example:"550e8400-e29b-41d4-a716-446655440000"`
}

type UpdatePenjualanRequest struct {
	NoFaktur   string    `json:"no_faktur" example:"SELL-2024-001"`
	Nama       string    `json:"nama" example:"Pelanggan A"`
	TotalHarga int64     `json:"total_harga" example:"5000000"`
	KodeSales  uuid.UUID `json:"kode_sales" example:"550e8400-e29b-41d4-a716-446655440000"`
}

type AttachBarangToPenjualanRequest struct {
	BarangIDs []uuid.UUID `json:"barang_ids" example:"550e8400-e29b-41d4-a716-446655440000"`
}
