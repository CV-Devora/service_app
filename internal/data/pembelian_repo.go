package data

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PembelianRepo struct {
	db *gorm.DB
}

func NewPembelianRepo(db *gorm.DB) *PembelianRepo {
	return &PembelianRepo{db: db}
}

func (r *PembelianRepo) FindAll() ([]Pembelian, error) {
	var items []Pembelian
	err := r.db.Preload("Barang").Where("deleted_at IS NULL").Find(&items).Error
	return items, err
}

func (r *PembelianRepo) FindByID(id uuid.UUID) (*Pembelian, error) {
	var item Pembelian
	err := r.db.Preload("Barang").Where("id = ? AND deleted_at IS NULL", id).First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &item, err
}

// FindByDateRange mencari pembelian berdasarkan rentang created_at.
// 'to' bersifat inklusif untuk seluruh hari tersebut.
func (r *PembelianRepo) FindByDateRange(from, to time.Time) ([]Pembelian, error) {
	var items []Pembelian
	err := r.db.Preload("Barang").
		Where("deleted_at IS NULL AND created_at >= ?", from).
		Where("created_at < ?", to.AddDate(0, 0, 1)).
		Find(&items).Error
	return items, err
}

func (r *PembelianRepo) Create(p *Pembelian) error {
	p.ID = uuid.New()
	for i := range p.Barang {
		p.Barang[i].ID = uuid.New()
		p.Barang[i].PembelianID = &p.ID
	}
	return r.db.Create(p).Error
}

func (r *PembelianRepo) Update(p *Pembelian) error {
	return r.db.Save(p).Error
}

func (r *PembelianRepo) Delete(id uuid.UUID) error {
	return r.db.Model(&Pembelian{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NOW()")).Error
}

func (r *PembelianRepo) Approve(id uuid.UUID) error {
	return r.db.Model(&Pembelian{}).Where("id = ?", id).Update("is_approve", true).Error
}
