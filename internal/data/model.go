package data

import (
	"time"

	"github.com/google/uuid"
)

// User model
type User struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	Nama      string     `gorm:"type:varchar(255)" json:"nama"`
	Username  string     `gorm:"type:varchar(255);uniqueIndex" json:"username"`
	Password  string     `gorm:"type:varchar(255)" json:"-"`
	Role      string     `gorm:"type:varchar(50)" json:"role"`
	KodeSales *int64     `gorm:"type:int;column:kode_sales" json:"kode_sales,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(tx interface{}) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// Pembelian model
type Pembelian struct {
	ID               uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	NoFaktur         string     `gorm:"type:varchar(100);uniqueIndex" json:"no_faktur"`
	Nama             string     `gorm:"type:varchar(255)" json:"nama"`
	TipePemasok      string     `gorm:"type:varchar(100)" json:"tipe_pemasok"`
	HargaDeal        int64      `json:"harga_deal"`
	BeratNota        float64    `gorm:"type:decimal(10,3);default:0" json:"berat_nota"`
	HargaNota        int64      `gorm:"default:0" json:"harga_nota"`
	HargaRata        int64      `gorm:"default:0" json:"harga_rata"`
	TipePembayaran   string     `gorm:"type:varchar(50);default:''" json:"tipe_pembayaran"`
	JumlahPembayaran int64      `gorm:"default:0" json:"jumlah_pembayaran"`
	Barang           []Barang   `gorm:"foreignKey:PembelianID" json:"barang,omitempty"`
	IsApprove        bool       `gorm:"default:false" json:"is_approve"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (Pembelian) TableName() string {
	return "pembelian"
}

// Karat model
type Karat struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	Name      string     `gorm:"type:varchar(100)" json:"name"`
	Harga     int64      `json:"harga"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (Karat) TableName() string {
	return "karat"
}

// Baki model
type Baki struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	Nama      string     `gorm:"type:varchar(255)" json:"nama"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (Baki) TableName() string {
	return "baki"
}

// Penjualan model
type Penjualan struct {
	ID         uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	NoFaktur   string     `gorm:"type:varchar(100);uniqueIndex" json:"no_faktur"`
	Nama       string     `gorm:"type:varchar(255)" json:"nama"`
	TotalHarga int64      `json:"total_harga"`
	KodeSales  *int64     `gorm:"type:int;column:kode_sales" json:"kode_sales,omitempty"`
	Barang     *Barang    `gorm:"foreignKey:PenjualanSaleID" json:"barang,omitempty"`
	HargaGram  int64      `gorm:"default:0" json:"harga_gram"`
	HargaJual  int64      `gorm:"default:0" json:"harga_jual"`
	Ongkos     int64      `gorm:"default:0" json:"ongkos"`
	Cash       int64      `gorm:"default:0" json:"cash"`
	Transfer   int64      `gorm:"default:0" json:"transfer"`
	Debet      int64      `gorm:"default:0" json:"debet"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (Penjualan) TableName() string {
	return "penjualan"
}

// Barang model
type Barang struct {
	ID               uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	Barcode          string     `gorm:"type:varchar(100);uniqueIndex" json:"barcode"`
	Nama             string     `gorm:"type:varchar(255)" json:"nama"`
	KaratID          *uuid.UUID `gorm:"type:uuid" json:"karat_id,omitempty"`
	Karat            *Karat     `gorm:"foreignKey:KaratID" json:"karat,omitempty"`
	Berat            float64    `gorm:"type:decimal(10,3)" json:"berat"`
	BeratAtribut     float64    `gorm:"type:decimal(10,3);default:0" json:"berat_atribut"`
	Harga            int64      `json:"harga"`
	Photo            string     `gorm:"type:text" json:"photo"`
	Kondisi          string     `gorm:"type:varchar(100)" json:"kondisi"`
	PembelianID      *uuid.UUID `gorm:"type:uuid" json:"pembelian_id,omitempty"`
	Pembelian        *Pembelian `gorm:"foreignKey:PembelianID" json:"pembelian,omitempty"`
	BakiID           *uuid.UUID `gorm:"type:uuid" json:"baki_id,omitempty"`
	Baki             *Baki      `gorm:"foreignKey:BakiID" json:"baki,omitempty"`
	Grup             string     `gorm:"type:varchar(255);column:grup" json:"grup,omitempty"`
	PenjualanSaleID  *uuid.UUID `gorm:"type:uuid;column:penjualan_id" json:"penjualan_id,omitempty"`
	Penjualan        *Penjualan `gorm:"foreignKey:PenjualanSaleID" json:"penjualan,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (Barang) TableName() string {
	return "barang"
}

// BarangLanding model — barang yang ditampilkan di landing page
type BarangLanding struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	Nama      string     `gorm:"type:varchar(255)" json:"nama"`
	Karat     string     `gorm:"type:varchar(100)" json:"karat"`
	Berat     float64    `gorm:"type:decimal(10,3)" json:"berat"`
	Harga     int64      `json:"harga"`
	Photo     string     `gorm:"type:text" json:"photo"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (BarangLanding) TableName() string {
	return "barang_landing"
}
