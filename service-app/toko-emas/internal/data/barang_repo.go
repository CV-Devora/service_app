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
	err := r.db.Where("deleted_at IS NULL").Find(&items).Error
	return items, err
}

func (r *BarangRepo) FindByID(id uuid.UUID) (*Barang, error) {
	var item Barang
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&item).Error
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

func (r *BarangRepo) Delete(id uuid.UUID) error {
	return r.db.Model(&Barang{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NOW()")).Error
}
