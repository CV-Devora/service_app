package data

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ---- Karat ----

type KaratRepo struct {
	db *gorm.DB
}

func NewKaratRepo(db *gorm.DB) *KaratRepo {
	return &KaratRepo{db: db}
}

func (r *KaratRepo) FindAll() ([]Karat, error) {
	var items []Karat
	err := r.db.Where("deleted_at IS NULL").Find(&items).Error
	return items, err
}

func (r *KaratRepo) FindByID(id uuid.UUID) (*Karat, error) {
	var item Karat
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &item, err
}

func (r *KaratRepo) Create(k *Karat) error {
	k.ID = uuid.New()
	return r.db.Create(k).Error
}

func (r *KaratRepo) Update(k *Karat) error {
	return r.db.Save(k).Error
}

func (r *KaratRepo) Delete(id uuid.UUID) error {
	return r.db.Model(&Karat{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NOW()")).Error
}

// ---- Baki ----

type BakiRepo struct {
	db *gorm.DB
}

func NewBakiRepo(db *gorm.DB) *BakiRepo {
	return &BakiRepo{db: db}
}

func (r *BakiRepo) FindAll() ([]Baki, error) {
	var items []Baki
	err := r.db.Where("deleted_at IS NULL").Find(&items).Error
	return items, err
}

func (r *BakiRepo) FindByID(id uuid.UUID) (*Baki, error) {
	var item Baki
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &item, err
}

func (r *BakiRepo) Create(b *Baki) error {
	b.ID = uuid.New()
	return r.db.Create(b).Error
}

func (r *BakiRepo) Update(b *Baki) error {
	return r.db.Save(b).Error
}

func (r *BakiRepo) Delete(id uuid.UUID) error {
	return r.db.Model(&Baki{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NOW()")).Error
}

// ---- Penjualan ----

type PenjualanRepo struct {
	db *gorm.DB
}

func NewPenjualanRepo(db *gorm.DB) *PenjualanRepo {
	return &PenjualanRepo{db: db}
}

func (r *PenjualanRepo) FindAll() ([]Penjualan, error) {
	var items []Penjualan
	err := r.db.Preload("Barang").Where("deleted_at IS NULL").Find(&items).Error
	return items, err
}

func (r *PenjualanRepo) FindByID(id uuid.UUID) (*Penjualan, error) {
	var item Penjualan
	err := r.db.Preload("Barang").Where("id = ? AND deleted_at IS NULL", id).First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &item, err
}

// FindByDateRange mencari penjualan berdasarkan rentang created_at.
// 'to' bersifat inklusif untuk seluruh hari tersebut.
func (r *PenjualanRepo) FindByDateRange(from, to time.Time) ([]Penjualan, error) {
	var items []Penjualan
	err := r.db.Preload("Barang").
		Where("deleted_at IS NULL AND created_at >= ?", from).
		Where("created_at < ?", to.AddDate(0, 0, 1)).
		Find(&items).Error
	return items, err
}

func (r *PenjualanRepo) Create(p *Penjualan, barangIDs []uuid.UUID) error {
	p.ID = uuid.New()
	if err := r.db.Create(p).Error; err != nil {
		return err
	}
	if len(barangIDs) > 0 {
		if err := r.db.Model(&Barang{}).
			Where("id IN ?", barangIDs).
			Update("penjualan_id", p.ID).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *PenjualanRepo) Update(p *Penjualan) error {
	return r.db.Save(p).Error
}

func (r *PenjualanRepo) Delete(id uuid.UUID) error {
	return r.db.Model(&Penjualan{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NOW()")).Error
}

// SetBarangPenjualan menghubungkan barang-barang ke penjualan tertentu
// dengan mengupdate kolom penjualan_id di tabel barang.
func (r *PenjualanRepo) SetBarangPenjualan(penjualanID uuid.UUID, barangIDs []uuid.UUID) error {
	if len(barangIDs) == 0 {
		return nil
	}
	return r.db.Model(&Barang{}).
		Where("id IN ?", barangIDs).
		Update("penjualan_id", penjualanID).Error
}

