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
	Nama      string `json:"nama" example:"Budi Santoso"`
	Username  string `json:"username" example:"budi"`
	Password  string `json:"password" example:"secret123"`
	Role      string `json:"role" example:"kasir"`
	KodeSales *int64 `json:"kode_sales,omitempty" example:"1"`
}

type UpdateUserRequest struct {
	Nama      string `json:"nama" example:"Budi Santoso"`
	Username  string `json:"username" example:"budi"`
	Role      string `json:"role" example:"kasir"`
	KodeSales *int64 `json:"kode_sales,omitempty" example:"1"`
}

// ---- Barang ---
type CreateBarangRequest struct {
	Barcode      string     `json:"barcode" example:"BC001"`
	Nama         string     `json:"nama" example:"Gelang Emas"`
	KaratID      *string    `json:"karat_id,omitempty" example:"0839294b-3a65-4e27-904e-63f5c49e7fd3"`
	Berat        float64    `json:"berat" example:"5.5"`
	BeratAtribut float64    `json:"berat_atribut" example:"0.5"`
	Harga        int64      `json:"harga" example:"5000000"`
	Photo        string     `json:"photo" example:"https://example.com/photo.jpg"`
	Kondisi      string     `json:"kondisi" example:"baru"`
	PembelianID  *uuid.UUID `json:"pembelian_id,omitempty" example:"null"`
	BakiID       *uuid.UUID `json:"baki_id,omitempty" example:"null"`
	Grup         *string    `json:"grup,omitempty" example:"null"`
}

type UpdateBarangRequest struct {
	Barcode      string     `json:"barcode" example:"BC001"`
	Nama         string     `json:"nama" example:"Gelang Emas"`
	KaratID      *string    `json:"karat_id,omitempty" example:"0839294b-3a65-4e27-904e-63f5c49e7fd3"`
	Berat        float64    `json:"berat" example:"5.5"`
	BeratAtribut float64    `json:"berat_atribut" example:"0.5"`
	Harga        int64      `json:"harga" example:"5000000"`
	Photo        string     `json:"photo" example:"https://example.com/photo.jpg"`
	Kondisi      string     `json:"kondisi" example:"baru"`
	PembelianID  *uuid.UUID `json:"pembelian_id,omitempty"`
	BakiID       *uuid.UUID `json:"baki_id,omitempty"`
	Grup         *string    `json:"grup,omitempty"`
}

// ---- Pembelian ----

type CreatePembelianItem struct {
	Barcode string  `json:"barcode" example:"BRG001"`
	Nama    string  `json:"nama" example:"Gelang Emas"`
	KaratID *string `json:"karat_id,omitempty" example:"0839294b-3a65-4e27-904e-63f5c49e7fd3"`
	Berat   float64 `json:"berat" example:"5.5"`
	Harga   int64   `json:"harga" example:"5000000"`
	Photo   string  `json:"photo" example:"https://example.com/photo.jpg"`
	Kondisi string  `json:"kondisi" example:"baru"`
	BakiID  *string `json:"baki_id,omitempty" example:"null"`
}

type CreatePembelianRequest struct {
	NoFaktur         string               `json:"no_faktur" example:"INV-2024-001"`
	Nama             string               `json:"nama" example:"Toko Mas Jaya"`
	TipePemasok      string               `json:"tipe_pemasok" example:"supplier"`
	HargaDeal        int64                `json:"harga_deal" example:"10000000"`
	BeratNota        float64              `json:"berat_nota" example:"25.5"`
	HargaNota        int64                `json:"harga_nota" example:"12000000"`
	HargaRata        int64                `json:"harga_rata" example:"470588"`
	TipePembayaran   string               `json:"tipe_pembayaran" example:"cash"`
	JumlahPembayaran int64                `json:"jumlah_pembayaran" example:"12000000"`
	Barang           []CreatePembelianItem `json:"barang,omitempty"`
}

type UpdatePembelianRequest struct {
	NoFaktur         string `json:"no_faktur" example:"INV-2024-001"`
	Nama             string `json:"nama" example:"Toko Mas Jaya"`
	TipePemasok      string `json:"tipe_pemasok" example:"supplier"`
	HargaDeal        int64  `json:"harga_deal" example:"10000000"`
	BeratNota        float64 `json:"berat_nota" example:"25.5"`
	HargaNota        int64  `json:"harga_nota" example:"12000000"`
	HargaRata        int64  `json:"harga_rata" example:"470588"`
	TipePembayaran   string `json:"tipe_pembayaran" example:"cash"`
	JumlahPembayaran int64  `json:"jumlah_pembayaran" example:"12000000"`
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

// ---- Barang Landing ----

type CreateBarangLandingRequest struct {
	Nama  string  `json:"nama" example:"Cincin Emas"`
	Karat string  `json:"karat" example:"24"`
	Berat float64 `json:"berat" example:"5.5"`
	Harga int64   `json:"harga" example:"5000000"`
	Photo string  `json:"photo" example:"https://example.com/photo.jpg"`
}

type UpdateBarangLandingRequest struct {
	Nama  string  `json:"nama" example:"Cincin Emas"`
	Karat string  `json:"karat" example:"24"`
	Berat float64 `json:"berat" example:"5.5"`
	Harga int64   `json:"harga" example:"5000000"`
	Photo string  `json:"photo" example:"https://example.com/photo.jpg"`
}

// ---- Penjualan ----

type CreatePenjualanRequest struct {
	NoFaktur   string      `json:"no_faktur" example:"SELL-2024-001"`
	Nama       string      `json:"nama" example:"Pelanggan A"`
	TotalHarga int64       `json:"total_harga" example:"5000000"`
	KodeSales  *int64      `json:"kode_sales,omitempty" example:"1"`
	HargaGram  int64       `json:"harga_gram" example:"750000"`
	HargaJual  int64       `json:"harga_jual" example:"4500000"`
	Ongkos     int64       `json:"ongkos" example:"100000"`
	Cash       int64       `json:"cash" example:"2000000"`
	Transfer   int64       `json:"transfer" example:"2000000"`
	Debet      int64       `json:"debet" example:"600000"`
	BarangIDs  []uuid.UUID `json:"barang_ids" example:"550e8400-e29b-41d4-a716-446655440000"`
}

type UpdatePenjualanRequest struct {
	NoFaktur   string    `json:"no_faktur" example:"SELL-2024-001"`
	Nama       string    `json:"nama" example:"Pelanggan A"`
	TotalHarga int64     `json:"total_harga" example:"5000000"`
	KodeSales  *int64    `json:"kode_sales,omitempty" example:"1"`
	HargaGram  int64     `json:"harga_gram" example:"750000"`
	HargaJual  int64     `json:"harga_jual" example:"4500000"`
	Ongkos     int64     `json:"ongkos" example:"100000"`
	Cash       int64     `json:"cash" example:"2000000"`
	Transfer   int64     `json:"transfer" example:"2000000"`
	Debet      int64     `json:"debet" example:"600000"`
}

type AttachBarangToPenjualanRequest struct {
	BarangIDs []uuid.UUID `json:"barang_ids" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// ---- Dashboard ----

type ChartData struct {
	Label string `json:"label"`
	Value int64  `json:"value"`
}

type DashboardResponse struct {
	BarangChart                []ChartData `json:"barang_chart"`
	PenjualanChart             []ChartData `json:"penjualan_chart"`
	PembelianChart             []ChartData `json:"pembelian_chart"`
	PenjualanHariIni           int64       `json:"penjualan_hari_ini"`
	BarangTersedia             int64       `json:"barang_tersedia"`
	PenjualanBulanan           []ChartData `json:"penjualan_bulanan"`
	PenjualanPerGrup           []ChartData `json:"penjualan_per_grup"`
	PenjualanHariIniSales      int64       `json:"penjualan_hari_ini_sales"`
	PenjualanBulananSales      int64       `json:"penjualan_bulanan_sales"`
	PenjualanTahunanSalesChart []ChartData `json:"penjualan_tahunan_sales_chart"`
}

// ---- Import Barang ----

type ImportBarangError struct {
	Row     int    `json:"row"`
	Barcode string `json:"barcode,omitempty"`
	Message string `json:"message"`
}

type ImportBarangResult struct {
	Total       int                  `json:"total"`
	Imported    int                  `json:"imported"`
	Skipped     int                  `json:"skipped"`
	Errors      []ImportBarangError  `json:"errors,omitempty"`
}
