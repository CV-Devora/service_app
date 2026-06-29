package service

import (
	"encoding/json"
	"net/http"

	v1 "toko-emas/api/v1"
	"toko-emas/internal/data"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type BarangService struct {
	repo *data.BarangRepo
}

func NewBarangService(repo *data.BarangRepo) *BarangService {
	return &BarangService{repo: repo}
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// ListBarang godoc
// @Summary      List semua barang
// @Description  Mendapatkan semua data barang
// @Tags         Barang
// @Accept       json
// @Produce      json
// @Success      200  {object}  v1.Response
// @Router       /api/v1/barang [get]
func (s *BarangService) List(w http.ResponseWriter, r *http.Request) {
	items, err := s.repo.FindAll()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, v1.Response{Code: 500, Message: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, v1.Response{Code: 200, Message: "success", Data: items})
}

// GetBarang godoc
// @Summary      Get barang by ID
// @Description  Mendapatkan detail barang berdasarkan ID
// @Tags         Barang
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Barang ID"
// @Success      200  {object}  v1.Response
// @Failure      404  {object}  v1.Response
// @Router       /api/v1/barang/{id} [get]
func (s *BarangService) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		writeJSON(w, http.StatusBadRequest, v1.Response{Code: 400, Message: "invalid id"})
		return
	}
	item, err := s.repo.FindByID(id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, v1.Response{Code: 500, Message: err.Error()})
		return
	}
	if item == nil {
		writeJSON(w, http.StatusNotFound, v1.Response{Code: 404, Message: "not found"})
		return
	}
	writeJSON(w, http.StatusOK, v1.Response{Code: 200, Message: "success", Data: item})
}

// CreateBarang godoc
// @Summary      Tambah barang
// @Description  Menambahkan barang baru
// @Tags         Barang
// @Accept       json
// @Produce      json
// @Param        body  body      v1.CreateBarangRequest  true  "Data barang"
// @Success      201   {object}  v1.Response
// @Failure      400   {object}  v1.Response
// @Router       /api/v1/barang [post]
func (s *BarangService) Create(w http.ResponseWriter, r *http.Request) {
	var req v1.CreateBarangRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, v1.Response{Code: 400, Message: "invalid body"})
		return
	}
	b := &data.Barang{
		Barcode:     req.Barcode,
		Nama:        req.Nama,
		Karat:       req.Karat,
		Berat:       req.Berat,
		Harga:       req.Harga,
		Photo:       req.Photo,
		Kondisi:     req.Kondisi,
		PembelianID: req.PembelianID,
		BakiID:      req.BakiID,
		GrupID:      req.GrupID,
	}
	if err := s.repo.Create(b); err != nil {
		writeJSON(w, http.StatusInternalServerError, v1.Response{Code: 500, Message: err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, v1.Response{Code: 201, Message: "created", Data: b})
}

// UpdateBarang godoc
// @Summary      Update barang
// @Description  Memperbarui data barang berdasarkan ID
// @Tags         Barang
// @Accept       json
// @Produce      json
// @Param        id    path      string                  true  "Barang ID"
// @Param        body  body      v1.UpdateBarangRequest  true  "Data barang"
// @Success      200   {object}  v1.Response
// @Router       /api/v1/barang/{id} [put]
func (s *BarangService) Update(w http.ResponseWriter, r *http.Request) {
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
	var req v1.UpdateBarangRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, v1.Response{Code: 400, Message: "invalid body"})
		return
	}
	existing.Barcode = req.Barcode
	existing.Nama = req.Nama
	existing.Karat = req.Karat
	existing.Berat = req.Berat
	existing.Harga = req.Harga
	existing.Photo = req.Photo
	existing.Kondisi = req.Kondisi
	existing.PembelianID = req.PembelianID
	existing.BakiID = req.BakiID
	existing.GrupID = req.GrupID
	if err := s.repo.Update(existing); err != nil {
		writeJSON(w, http.StatusInternalServerError, v1.Response{Code: 500, Message: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, v1.Response{Code: 200, Message: "updated", Data: existing})
}

// DeleteBarang godoc
// @Summary      Hapus barang
// @Description  Menghapus barang (soft delete)
// @Tags         Barang
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Barang ID"
// @Success      200  {object}  v1.Response
// @Router       /api/v1/barang/{id} [delete]
func (s *BarangService) Delete(w http.ResponseWriter, r *http.Request) {
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
