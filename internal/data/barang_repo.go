package data

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BarangRepo struct {
	db *gorm.DB
}

func NewBarangRepo(db *gorm.DB) *BarangRepo {
	return &BarangRepo{db: db}
}

func (r *BarangRepo) FindAll() ([]Barang, error) {
	var items []Barang
	err := r.db.Preload("Karat").Preload("Baki").
		Where("deleted_at IS NULL AND penjualan_id IS NULL").Find(&items).Error
	return items, err
}

func (r *BarangRepo) FindByID(id uuid.UUID) (*Barang, error) {
	var item Barang
	err := r.db.Preload("Karat").Preload("Baki").Where("id = ? AND deleted_at IS NULL", id).First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &item, err
}

func (r *BarangRepo) Create(b *Barang) error {
	b.ID = uuid.New()
	return r.db.Create(b).Error
}

func (r *BarangRepo) Update(b *Barang) error {
	return r.db.Save(b).Error
}

func (r *BarangRepo) GetLatestBarcode() (int64, error) {
	var maxBarcode int64
	err := r.db.Raw(`
		SELECT COALESCE(MAX(NULLIF(REGEXP_REPLACE(barcode, '[^0-9]', '', 'g'), '')::BIGINT), 0)
		FROM barang WHERE deleted_at IS NULL
	`).Scan(&maxBarcode).Error
	return maxBarcode, err
}

func (r *BarangRepo) Delete(id uuid.UUID) error {
	return r.db.Model(&Barang{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NOW()")).Error
}

func (r *BarangRepo) FindAllBarcodes() (map[string]bool, error) {
	var codes []string
	err := r.db.Model(&Barang{}).Where("deleted_at IS NULL").Pluck("barcode", &codes).Error
	if err != nil {
		return nil, err
	}
	set := make(map[string]bool, len(codes))
	for _, c := range codes {
		set[c] = true
	}
	return set, nil
}

func (r *BarangRepo) FindAllBaki() ([]Baki, error) {
	var items []Baki
	err := r.db.Where("deleted_at IS NULL").Find(&items).Error
	return items, err
}

func (r *BarangRepo) FindAllKarat() ([]Karat, error) {
	var items []Karat
	err := r.db.Where("deleted_at IS NULL").Find(&items).Error
	return items, err
}

// FindSoldBarangIDs mengembalikan ID barang yang sudah terhubung ke transaksi
// penjualan yang masih aktif.
func (r *BarangRepo) FindSoldBarangIDs() (map[uuid.UUID]bool, error) {
	var ids []uuid.UUID
	err := r.db.Raw(`
		SELECT b.id
		FROM barang b
		JOIN penjualan p ON p.id = b.penjualan_id
		WHERE b.penjualan_id IS NOT NULL
		  AND p.deleted_at IS NULL
	`).Scan(&ids).Error
	if err != nil {
		return nil, err
	}
	set := make(map[uuid.UUID]bool, len(ids))
	for _, id := range ids {
		set[id] = true
	}
	return set, nil
}
