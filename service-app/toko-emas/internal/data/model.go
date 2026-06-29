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
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	NoFaktur    string     `gorm:"type:varchar(100);uniqueIndex" json:"no_faktur"`
	Nama        string     `gorm:"type:varchar(255)" json:"nama"`
	TipePemasok string     `gorm:"type:varchar(100)" json:"tipe_pemasok"`
	HargaDeal   int64      `json:"harga_deal"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`
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
	KodeSales  *uuid.UUID `gorm:"type:uuid" json:"kode_sales,omitempty"`
	Sales      *User      `gorm:"foreignKey:KodeSales" json:"sales,omitempty"`
	Barang     []Barang   `gorm:"many2many:penjualan_barang;" json:"barang,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (Penjualan) TableName() string {
	return "penjualan"
}

// Barang model
type Barang struct {
	ID          uuid.UUID   `gorm:"type:uuid;primary_key" json:"id"`
	Barcode     string      `gorm:"type:varchar(100);uniqueIndex" json:"barcode"`
	Nama        string      `gorm:"type:varchar(255)" json:"nama"`
	Karat       int         `json:"karat"`
	Berat       float64     `gorm:"type:decimal(10,3)" json:"berat"`
	Harga       int64       `json:"harga"`
	Photo       string      `gorm:"type:varchar(500)" json:"photo"`
	Kondisi     string      `gorm:"type:varchar(100)" json:"kondisi"`
	PembelianID *uuid.UUID  `gorm:"type:uuid" json:"pembelian_id,omitempty"`
	Pembelian   *Pembelian  `gorm:"foreignKey:PembelianID" json:"pembelian,omitempty"`
	BakiID      *uuid.UUID  `gorm:"type:uuid" json:"baki_id,omitempty"`
	Baki        *Baki       `gorm:"foreignKey:BakiID" json:"baki,omitempty"`
	GrupID      *uuid.UUID  `gorm:"type:uuid" json:"grup_id,omitempty"`
	Penjualan   []Penjualan `gorm:"many2many:penjualan_barang;" json:"penjualan,omitempty"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	DeletedAt   *time.Time  `gorm:"index" json:"deleted_at,omitempty"`
}

func (Barang) TableName() string {
	return "barang"
}

// PenjualanBarang model
type PenjualanBarang struct {
	BarangID    uuid.UUID `gorm:"type:uuid;primaryKey" json:"barang_id"`
	PenjualanID uuid.UUID `gorm:"type:uuid;primaryKey" json:"penjualan_id"`
}

func (PenjualanBarang) TableName() string {
	return "penjualan_barang"
}
