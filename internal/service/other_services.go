package service

import (
	"encoding/json"
	"net/http"

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
// @Param        from  query     string  false  "Tanggal awal (YYYY-MM-DD)"
// @Param        to    query     string  false  "Tanggal akhir (YYYY-MM-DD)"
// @Success      200  {object}  v1.Response
// @Router       /api/v1/pembelian [get]
func (s *PembelianService) List(w http.ResponseWriter, r *http.Request) {
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")
	var (
		items []data.Pembelian
		err   error
	)
	if fromStr != "" || toStr != "" {
		from, to, perr := parseDateRange(r)
		if perr != nil {
			writeJSON(w, http.StatusBadRequest, v1.Response{Code: 400, Message: perr.Error()})
			return
		}
		items, err = s.repo.FindByDateRange(from, to)
	} else {
		items, err = s.repo.FindAll()
	}
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
// @Param        id   path      string  true  "Pembelian ID"
// @Success      200  {object}  v1.Response
// @Router       /api/v1/pembelian/{id} [get]
func (s *PembelianService) Get(w http.ResponseWriter, r *http.Request) {
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

	var barangItems []data.Barang
	for _, item := range req.Barang {
		var bakiID *uuid.UUID
		if item.BakiID != nil && *item.BakiID != "" {
			parsed, err := uuid.Parse(*item.BakiID)
			if err == nil {
				bakiID = &parsed
			}
		}
		var karatID *uuid.UUID
		if item.KaratID != nil && *item.KaratID != "" {
			parsed, err := uuid.Parse(*item.KaratID)
			if err == nil {
				karatID = &parsed
			}
		}
		barangItems = append(barangItems, data.Barang{
			Barcode: item.Barcode,
			Nama:    item.Nama,
			KaratID: karatID,
			Berat:   item.Berat,
			Harga:   item.Harga,
			Photo:   item.Photo,
			Kondisi: item.Kondisi,
			BakiID:  bakiID,
		})
	}

	p := &data.Pembelian{
		NoFaktur:         req.NoFaktur,
		Nama:             req.Nama,
		TipePemasok:      req.TipePemasok,
		HargaDeal:        req.HargaDeal,
		BeratNota:        req.BeratNota,
		HargaNota:        req.HargaNota,
		HargaRata:        req.HargaRata,
		TipePembayaran:   req.TipePembayaran,
		JumlahPembayaran: req.JumlahPembayaran,
		Barang:           barangItems,
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
// @Param        id    path      string                     true  "Pembelian ID"
// @Param        body  body      v1.UpdatePembelianRequest  true  "Data pembelian"
// @Success      200   {object}  v1.Response
// @Router       /api/v1/pembelian/{id} [put]
func (s *PembelianService) Update(w http.ResponseWriter, r *http.Request) {
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
	var req v1.UpdatePembelianRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, v1.Response{Code: 400, Message: "invalid body"})
		return
	}
	existing.NoFaktur = req.NoFaktur
	existing.Nama = req.Nama
	existing.TipePemasok = req.TipePemasok
	existing.HargaDeal = req.HargaDeal
	existing.BeratNota = req.BeratNota
	existing.HargaNota = req.HargaNota
	existing.HargaRata = req.HargaRata
	existing.TipePembayaran = req.TipePembayaran
	existing.JumlahPembayaran = req.JumlahPembayaran
	if err := s.repo.Update(existing); err != nil {
		writeJSON(w, http.StatusInternalServerError, v1.Response{Code: 500, Message: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, v1.Response{Code: 200, Message: "updated", Data: existing})
}

// ApprovePembelian godoc
// @Summary      Approve pembelian
// @Tags         Pembelian
// @Produce      json
// @Param        id   path      string  true  "Pembelian ID"
// @Success      200  {object}  v1.Response
// @Router       /api/v1/pembelian/{id}/approve [put]
func (s *PembelianService) Approve(w http.ResponseWriter, r *http.Request) {
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
	if existing.IsApprove {
		writeJSON(w, http.StatusBadRequest, v1.Response{Code: 400, Message: "already approved"})
		return
	}
	if err := s.repo.Approve(id); err != nil {
		writeJSON(w, http.StatusInternalServerError, v1.Response{Code: 500, Message: err.Error()})
		return
	}
	item, _ := s.repo.FindByID(id)
	writeJSON(w, http.StatusOK, v1.Response{Code: 200, Message: "approved", Data: item})
}

// DeletePembelian godoc
// @Summary      Hapus pembelian
// @Tags         Pembelian
// @Produce      json
// @Param        id   path      string  true  "Pembelian ID"
// @Success      200  {object}  v1.Response
// @Router       /api/v1/pembelian/{id} [delete]
func (s *PembelianService) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		writeJSON(w, http.StatusBadRequest, v1.Response{Code: 400, Message: "invalid id"})
		return
	}
	if err := s.repo.Delete(id); err != nil {
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
	repo      *data.PenjualanRepo
	authSvc   *AuthService
}

func NewPenjualanService(repo *data.PenjualanRepo, authSvc *AuthService) *PenjualanService {
	return &PenjualanService{repo: repo, authSvc: authSvc}
}

// ListPenjualan godoc
// @Summary      List semua penjualan
// @Tags         Penjualan
// @Produce      json
// @Param        from  query     string  false  "Tanggal awal (YYYY-MM-DD)"
// @Param        to    query     string  false  "Tanggal akhir (YYYY-MM-DD)"
// @Success      200  {object}  v1.Response
// @Router       /api/v1/penjualan [get]
func (s *PenjualanService) List(w http.ResponseWriter, r *http.Request) {
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")
	var (
		items []data.Penjualan
		err   error
	)
	if fromStr != "" || toStr != "" {
		from, to, perr := parseDateRange(r)
		if perr != nil {
			writeJSON(w, http.StatusBadRequest, v1.Response{Code: 400, Message: perr.Error()})
			return
		}
		items, err = s.repo.FindByDateRange(from, to)
	} else {
		items, err = s.repo.FindAll()
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, v1.Response{Code: 500, Message: err.Error()})
		return
	}
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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, v1.Response{Code: 400, Message: "invalid body"})
		return
	}

	_, err := s.authSvc.UserIDFromRequest(r)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, v1.Response{Code: 401, Message: "unauthorized"})
		return
	}

	p := &data.Penjualan{
		NoFaktur:   req.NoFaktur,
		Nama:       req.Nama,
		TotalHarga: req.TotalHarga,
		KodeSales:  req.KodeSales,
		HargaGram:  req.HargaGram,
		HargaJual:  req.HargaJual,
		Ongkos:     req.Ongkos,
		Cash:       req.Cash,
		Transfer:   req.Transfer,
		Debet:      req.Debet,
	}
	if err := s.repo.Create(p, req.BarangIDs); err != nil {
		writeJSON(w, http.StatusInternalServerError, v1.Response{Code: 500, Message: err.Error()})
		return
	}
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
	existing.KodeSales = req.KodeSales
	existing.HargaGram = req.HargaGram
	existing.HargaJual = req.HargaJual
	existing.Ongkos = req.Ongkos
	existing.Cash = req.Cash
	existing.Transfer = req.Transfer
	existing.Debet = req.Debet
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

	if err := s.repo.SetBarangPenjualan(penjualanID, req.BarangIDs); err != nil {
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
