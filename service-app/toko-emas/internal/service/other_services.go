package service

import (
	"encoding/json"
	"net/http"
	"strconv"

	v1 "toko-emas/api/v1"
	"toko-emas/internal/data"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// ---- Pembelian ----

type PembelianService struct {
	repo *data.PembelianRepo
}

func NewPembelianService(repo *data.PembelianRepo) *PembelianService {
	return &PembelianService{repo: repo}
}

// ListPembelian godoc
// @Summary      List semua pembelian
// @Tags         Pembelian
// @Produce      json
// @Success      200  {object}  v1.Response
// @Router       /api/v1/pembelian [get]
func (s *PembelianService) List(w http.ResponseWriter, r *http.Request) {
	items, err := s.repo.FindAll()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, v1.Response{Code: 500, Message: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, v1.Response{Code: 200, Message: "success", Data: items})
}

// GetPembelian godoc
// @Summary      Get pembelian by ID
// @Tags         Pembelian
// @Produce      json
// @Param        id   path      int  true  "Pembelian ID"
// @Success      200  {object}  v1.Response
// @Router       /api/v1/pembelian/{id} [get]
func (s *PembelianService) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, v1.Response{Code: 400, Message: "invalid id"})
		return
	}
	item, err := s.repo.FindByID(uint(id))
	if err != nil || item == nil {
		writeJSON(w, http.StatusNotFound, v1.Response{Code: 404, Message: "not found"})
		return
	}
	writeJSON(w, http.StatusOK, v1.Response{Code: 200, Message: "success", Data: item})
}

// CreatePembelian godoc
// @Summary      Tambah pembelian
// @Tags         Pembelian
// @Accept       json
// @Produce      json
// @Param        body  body      v1.CreatePembelianRequest  true  "Data pembelian"
// @Success      201   {object}  v1.Response
// @Router       /api/v1/pembelian [post]
func (s *PembelianService) Create(w http.ResponseWriter, r *http.Request) {
	var req v1.CreatePembelianRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, v1.Response{Code: 400, Message: "invalid body"})
		return
	}
	p := &data.Pembelian{
		NoFaktur:    req.NoFaktur,
		Nama:        req.Nama,
		TipePemasok: req.TipePemasok,
		HargaDeal:   req.HargaDeal,
	}
	if err := s.repo.Create(p); err != nil {
		writeJSON(w, http.StatusInternalServerError, v1.Response{Code: 500, Message: err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, v1.Response{Code: 201, Message: "created", Data: p})
}

// UpdatePembelian godoc
// @Summary      Update pembelian
// @Tags         Pembelian
// @Accept       json
// @Produce      json
// @Param        id    path      int                        true  "Pembelian ID"
// @Param        body  body      v1.UpdatePembelianRequest  true  "Data pembelian"
// @Success      200   {object}  v1.Response
// @Router       /api/v1/pembelian/{id} [put]
func (s *PembelianService) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, v1.Response{Code: 400, Message: "invalid id"})
		return
	}
	existing, err := s.repo.FindByID(uint(id))
	if err != nil || existing == nil {
		writeJSON(w, http.StatusNotFound, v1.Response{Code: 404, Message: "not found"})
		return
	}
	var req v1.UpdatePembelianRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, v1.Response{Code: 400, Message: "invalid body"})
		return
	}
	existing.NoFaktur = req.NoFaktur
	existing.Nama = req.Nama
	existing.TipePemasok = req.TipePemasok
	existing.HargaDeal = req.HargaDeal
	if err := s.repo.Update(existing); err != nil {
		writeJSON(w, http.StatusInternalServerError, v1.Response{Code: 500, Message: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, v1.Response{Code: 200, Message: "updated", Data: existing})
}

// DeletePembelian godoc
// @Summary      Hapus pembelian
// @Tags         Pembelian
// @Produce      json
// @Param        id   path      int  true  "Pembelian ID"
// @Success      200  {object}  v1.Response
// @Router       /api/v1/pembelian/{id} [delete]
func (s *PembelianService) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, v1.Response{Code: 400, Message: "invalid id"})
		return
	}
	if err := s.repo.Delete(uint(id)); err != nil {
		writeJSON(w, http.StatusInternalServerError, v1.Response{Code: 500, Message: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, v1.Response{Code: 200, Message: "deleted"})
}

// ---- Karat ----

type KaratService struct {
	repo *data.KaratRepo
}

func NewKaratService(repo *data.KaratRepo) *KaratService {
	return &KaratService{repo: repo}
}

// ListKarat godoc
// @Summary      List semua karat
// @Tags         Karat
// @Produce      json
// @Success      200  {object}  v1.Response
// @Router       /api/v1/karat [get]
func (s *KaratService) List(w http.ResponseWriter, r *http.Request) {
	items, err := s.repo.FindAll()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, v1.Response{Code: 500, Message: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, v1.Response{Code: 200, Message: "success", Data: items})
}

// GetKarat godoc
// @Summary      Get karat by ID
// @Tags         Karat
// @Produce      json
// @Param        id   path      string  true  "Karat ID"
// @Success      200  {object}  v1.Response
// @Router       /api/v1/karat/{id} [get]
func (s *KaratService) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		writeJSON(w, http.StatusBadRequest, v1.Response{Code: 400, Message: "invalid id"})
		return
	}
	item, err := s.repo.FindByID(id)
	if err != nil || item == nil {
		writeJSON(w, http.StatusNotFound, v1.Response{Code: 404, Message: "not found"})
		return
	}
	writeJSON(w, http.StatusOK, v1.Response{Code: 200, Message: "success", Data: item})
}

// CreateKarat godoc
// @Summary      Tambah karat
// @Tags         Karat
// @Accept       json
// @Produce      json
// @Param        body  body      v1.CreateKaratRequest  true  "Data karat"
// @Success      201   {object}  v1.Response
// @Router       /api/v1/karat [post]
func (s *KaratService) Create(w http.ResponseWriter, r *http.Request) {
	var req v1.CreateKaratRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, v1.Response{Code: 400, Message: "invalid body"})
		return
	}
	k := &data.Karat{Name: req.Name, Harga: req.Harga}
	if err := s.repo.Create(k); err != nil {
		writeJSON(w, http.StatusInternalServerError, v1.Response{Code: 500, Message: err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, v1.Response{Code: 201, Message: "created", Data: k})
}

// UpdateKarat godoc
// @Summary      Update karat
// @Tags         Karat
// @Accept       json
// @Produce      json
// @Param        id    path      string                 true  "Karat ID"
// @Param        body  body      v1.UpdateKaratRequest  true  "Data karat"
// @Success      200   {object}  v1.Response
// @Router       /api/v1/karat/{id} [put]
func (s *KaratService) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		writeJSON(w, http.StatusBadRequest, v1.Response{Code: 400, Message: "invalid id"})
		return
	}
	existing, err := s.repo.FindByID(id)
	if err != nil || existing == nil {
		writeJSON(w, http.StatusNotFound, v1.Response{Code: 404, Message: "not found"})
		return
	}
	var req v1.UpdateKaratRequest
	json.NewDecoder(r.Body).Decode(&req)
	existing.Name = req.Name
	existing.Harga = req.Harga
	s.repo.Update(existing)
	writeJSON(w, http.StatusOK, v1.Response{Code: 200, Message: "updated", Data: existing})
}

// DeleteKarat godoc
// @Summary      Hapus karat
// @Tags         Karat
// @Produce      json
// @Param        id   path      string  true  "Karat ID"
// @Success      200  {object}  v1.Response
// @Router       /api/v1/karat/{id} [delete]
func (s *KaratService) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := uuid.Parse(vars["id"])
	s.repo.Delete(id)
	writeJSON(w, http.StatusOK, v1.Response{Code: 200, Message: "deleted"})
}

// ---- Baki ----

type BakiService struct {
	repo *data.BakiRepo
}

func NewBakiService(repo *data.BakiRepo) *BakiService {
	return &BakiService{repo: repo}
}

// ListBaki godoc
// @Summary      List semua baki
// @Tags         Baki
// @Produce      json
// @Success      200  {object}  v1.Response
// @Router       /api/v1/baki [get]
func (s *BakiService) List(w http.ResponseWriter, r *http.Request) {
	items, _ := s.repo.FindAll()
	writeJSON(w, http.StatusOK, v1.Response{Code: 200, Message: "success", Data: items})
}

// GetBaki godoc
// @Summary      Get baki by ID
// @Tags         Baki
// @Produce      json
// @Param        id   path      string  true  "Baki ID"
// @Success      200  {object}  v1.Response
// @Router       /api/v1/baki/{id} [get]
func (s *BakiService) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := uuid.Parse(vars["id"])
	item, _ := s.repo.FindByID(id)
	if item == nil {
		writeJSON(w, http.StatusNotFound, v1.Response{Code: 404, Message: "not found"})
		return
	}
	writeJSON(w, http.StatusOK, v1.Response{Code: 200, Message: "success", Data: item})
}

// CreateBaki godoc
// @Summary      Tambah baki
// @Tags         Baki
// @Accept       json
// @Produce      json
// @Param        body  body      v1.CreateBakiRequest  true  "Data baki"
// @Success      201   {object}  v1.Response
// @Router       /api/v1/baki [post]
func (s *BakiService) Create(w http.ResponseWriter, r *http.Request) {
	var req v1.CreateBakiRequest
	json.NewDecoder(r.Body).Decode(&req)
	b := &data.Baki{Nama: req.Nama}
	s.repo.Create(b)
	writeJSON(w, http.StatusCreated, v1.Response{Code: 201, Message: "created", Data: b})
}

// UpdateBaki godoc
// @Summary      Update baki
// @Tags         Baki
// @Accept       json
// @Produce      json
// @Param        id    path      string                true  "Baki ID"
// @Param        body  body      v1.UpdateBakiRequest  true  "Data baki"
// @Success      200   {object}  v1.Response
// @Router       /api/v1/baki/{id} [put]
func (s *BakiService) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := uuid.Parse(vars["id"])
	existing, _ := s.repo.FindByID(id)
	if existing == nil {
		writeJSON(w, http.StatusNotFound, v1.Response{Code: 404, Message: "not found"})
		return
	}
	var req v1.UpdateBakiRequest
	json.NewDecoder(r.Body).Decode(&req)
	existing.Nama = req.Nama
	s.repo.Update(existing)
	writeJSON(w, http.StatusOK, v1.Response{Code: 200, Message: "updated", Data: existing})
}

// DeleteBaki godoc
// @Summary      Hapus baki
// @Tags         Baki
// @Produce      json
// @Param        id   path      string  true  "Baki ID"
// @Success      200  {object}  v1.Response
// @Router       /api/v1/baki/{id} [delete]
func (s *BakiService) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := uuid.Parse(vars["id"])
	s.repo.Delete(id)
	writeJSON(w, http.StatusOK, v1.Response{Code: 200, Message: "deleted"})
}

// ---- Penjualan ----

type PenjualanService struct {
	repo *data.PenjualanRepo
}

func NewPenjualanService(repo *data.PenjualanRepo) *PenjualanService {
	return &PenjualanService{repo: repo}
}

// ListPenjualan godoc
// @Summary      List semua penjualan
// @Tags         Penjualan
// @Produce      json
// @Success      200  {object}  v1.Response
// @Router       /api/v1/penjualan [get]
func (s *PenjualanService) List(w http.ResponseWriter, r *http.Request) {
	items, _ := s.repo.FindAll()
	writeJSON(w, http.StatusOK, v1.Response{Code: 200, Message: "success", Data: items})
}

// GetPenjualan godoc
// @Summary      Get penjualan by ID
// @Tags         Penjualan
// @Produce      json
// @Param        id   path      string  true  "Penjualan ID"
// @Success      200  {object}  v1.Response
// @Router       /api/v1/penjualan/{id} [get]
func (s *PenjualanService) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := uuid.Parse(vars["id"])
	item, _ := s.repo.FindByID(id)
	if item == nil {
		writeJSON(w, http.StatusNotFound, v1.Response{Code: 404, Message: "not found"})
		return
	}
	writeJSON(w, http.StatusOK, v1.Response{Code: 200, Message: "success", Data: item})
}

// CreatePenjualan godoc
// @Summary      Tambah penjualan
// @Tags         Penjualan
// @Accept       json
// @Produce      json
// @Param        body  body      v1.CreatePenjualanRequest  true  "Data penjualan"
// @Success      201   {object}  v1.Response
// @Router       /api/v1/penjualan [post]
func (s *PenjualanService) Create(w http.ResponseWriter, r *http.Request) {
	var req v1.CreatePenjualanRequest
	json.NewDecoder(r.Body).Decode(&req)
	p := &data.Penjualan{
		NoFaktur:   req.NoFaktur,
		Nama:       req.Nama,
		TotalHarga: req.TotalHarga,
		KodeSales:  &req.KodeSales,
	}
	s.repo.Create(p)
	writeJSON(w, http.StatusCreated, v1.Response{Code: 201, Message: "created", Data: p})
}

// UpdatePenjualan godoc
// @Summary      Update penjualan
// @Tags         Penjualan
// @Accept       json
// @Produce      json
// @Param        id    path      string                     true  "Penjualan ID"
// @Param        body  body      v1.UpdatePenjualanRequest  true  "Data penjualan"
// @Success      200   {object}  v1.Response
// @Router       /api/v1/penjualan/{id} [put]
func (s *PenjualanService) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := uuid.Parse(vars["id"])
	existing, _ := s.repo.FindByID(id)
	if existing == nil {
		writeJSON(w, http.StatusNotFound, v1.Response{Code: 404, Message: "not found"})
		return
	}
	var req v1.UpdatePenjualanRequest
	json.NewDecoder(r.Body).Decode(&req)
	existing.NoFaktur = req.NoFaktur
	existing.Nama = req.Nama
	existing.TotalHarga = req.TotalHarga
	existing.KodeSales = &req.KodeSales
	s.repo.Update(existing)
	writeJSON(w, http.StatusOK, v1.Response{Code: 200, Message: "updated", Data: existing})
}

// DeletePenjualan godoc
// @Summary      Hapus penjualan
// @Tags         Penjualan
// @Produce      json
// @Param        id   path      string  true  "Penjualan ID"
// @Success      200  {object}  v1.Response
// @Router       /api/v1/penjualan/{id} [delete]
func (s *PenjualanService) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := uuid.Parse(vars["id"])
	s.repo.Delete(id)
	writeJSON(w, http.StatusOK, v1.Response{Code: 200, Message: "deleted"})
}

// AttachBarangToPenjualan godoc
// @Summary      Tambahkan barang ke penjualan
// @Tags         Penjualan
// @Accept       json
// @Produce      json
// @Param        id    path      string                          true  "Penjualan ID"
// @Param        body  body      v1.AttachBarangToPenjualanRequest  true  "Barang IDs"
// @Success      200   {object}  v1.Response
// @Router       /api/v1/penjualan/{id}/barang [post]
func (s *PenjualanService) AttachBarang(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	penjualanID, err := uuid.Parse(vars["id"])
	if err != nil {
		writeJSON(w, http.StatusBadRequest, v1.Response{Code: 400, Message: "invalid id"})
		return
	}

	var req v1.AttachBarangToPenjualanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, v1.Response{Code: 400, Message: "invalid body"})
		return
	}

	if err := s.repo.AddBarang(penjualanID, req.BarangIDs); err != nil {
		writeJSON(w, http.StatusInternalServerError, v1.Response{Code: 500, Message: err.Error()})
		return
	}

	item, err := s.repo.FindByID(penjualanID)
	if err != nil || item == nil {
		writeJSON(w, http.StatusNotFound, v1.Response{Code: 404, Message: "not found"})
		return
	}

	writeJSON(w, http.StatusOK, v1.Response{Code: 200, Message: "success", Data: item})
}
