package service

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	v1 "toko-emas/api/v1"
	"toko-emas/internal/data"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/xuri/excelize/v2"
)

type BarangService struct {
	repo *data.BarangRepo
	auth *AuthService
}

func NewBarangService(repo *data.BarangRepo, auth *AuthService) *BarangService {
	return &BarangService{repo: repo, auth: auth}
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

// LatestBarcode godoc
// @Summary      Get latest barcode number
// @Tags         Barang
// @Produce      json
// @Success      200  {object}  v1.Response
// @Router       /api/v1/barang/latest-barcode [get]
func (s *BarangService) LatestBarcode(w http.ResponseWriter, r *http.Request) {
	maxBarcode, err := s.repo.GetLatestBarcode()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, v1.Response{Code: 500, Message: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, v1.Response{Code: 200, Message: "success", Data: maxBarcode})
}

// ImportBarang godoc
// @Summary      Import barang dari file Excel/CSV
// @Description  Import data barang massal dari file .xlsx atau .csv dengan kolom barcode, baki, nama, kadar, berat, harga jual, kondisi. Khusus admin.
// @Tags         Barang
// @Accept       multipart/form-data
// @Produce      json
// @Param        file  formData  file  true  "File Excel (.xlsx) atau CSV"
// @Success      200   {object}  v1.ImportBarangResult
// @Failure      400   {object}  v1.Response
// @Failure      403   {object}  v1.Response
// @Router       /api/v1/barang/import [post]
func (s *BarangService) Import(w http.ResponseWriter, r *http.Request) {
	role, err := s.auth.RoleFromRequest(r)
	if err != nil || role != "admin" {
		writeJSON(w, http.StatusForbidden, v1.Response{Code: 403, Message: "hanya admin yang dapat mengimport barang"})
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		writeJSON(w, http.StatusBadRequest, v1.Response{Code: 400, Message: "file terlalu besar atau form tidak valid"})
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, v1.Response{Code: 400, Message: "file wajib diisi"})
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	rows, err := parseImportFile(ext, file)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, v1.Response{Code: 400, Message: err.Error()})
		return
	}
	if len(rows) < 2 {
		writeJSON(w, http.StatusBadRequest, v1.Response{Code: 400, Message: "file tidak berisi data (butuh baris header + minimal 1 baris data)"})
		return
	}

	// Map kolom header (case-insensitive, abaikan spasi)
	col := map[string]int{}
	for i, h := range rows[0] {
		key := normalizeHeader(h)
		if key != "" {
			col[key] = i
		}
	}
	barcodeCol, barcodeOK := col["barcode"]
	namaCol, namaOK := col["nama"]
	beratCol, beratOK := col["berat"]
	if !barcodeOK || !namaOK || !beratOK {
		writeJSON(w, http.StatusBadRequest, v1.Response{Code: 400, Message: "kolom wajib 'barcode', 'nama', dan 'berat' tidak ditemukan di baris header"})
		return
	}
	bakiCol := -1
	if idx, ok := col["baki"]; ok {
		bakiCol = idx
	}
	kadarCol := -1
	if idx, ok := col["kadar"]; ok {
		kadarCol = idx
	}
	hargaCol := -1
	if idx, ok := col["harga"]; ok {
		hargaCol = idx
	}
	kondisiCol := -1
	if idx, ok := col["kondisi"]; ok {
		kondisiCol = idx
	}

	// Preload data existing agar cepat
	existingBarcodes, err := s.repo.FindAllBarcodes()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, v1.Response{Code: 500, Message: err.Error()})
		return
	}
	bakis, err := s.repo.FindAllBaki()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, v1.Response{Code: 500, Message: err.Error()})
		return
	}
	bakiMap := make(map[string]uuid.UUID, len(bakis))
	for _, b := range bakis {
		bakiMap[strings.TrimSpace(b.Nama)] = b.ID
	}
	karats, err := s.repo.FindAllKarat()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, v1.Response{Code: 500, Message: err.Error()})
		return
	}
	karatMap := make(map[string]uuid.UUID, len(karats))
	for _, k := range karats {
		karatMap[normalizeKadar(k.Name)] = k.ID
	}

	result := &v1.ImportBarangResult{}
	for idx := 1; idx < len(rows); idx++ {
		row := rows[idx]
		result.Total++

		barcode := strings.TrimSpace(cellAt(row, barcodeCol))
		nama := strings.TrimSpace(cellAt(row, namaCol))
		beratStr := strings.TrimSpace(cellAt(row, beratCol))
		bakiName := strings.TrimSpace(cellAt(row, bakiCol))
		kadarName := strings.TrimSpace(cellAt(row, kadarCol))
		hargaStr := strings.TrimSpace(cellAt(row, hargaCol))
		kondisiStr := strings.TrimSpace(cellAt(row, kondisiCol))

		if barcode == "" || nama == "" || beratStr == "" {
			appendImportError(result, idx, barcode, "kolom barcode, nama, atau berat kosong")
			continue
		}
		if existingBarcodes[barcode] {
			appendImportError(result, idx, barcode, "barcode sudah terdaftar")
			continue
		}
		berat, err := parseBerat(beratStr)
		if err != nil || berat < 0 {
			appendImportError(result, idx, barcode, "berat tidak valid: \""+beratStr+"\"")
			continue
		}

		var karatID *uuid.UUID
		if kadarName != "" {
			id, ok := karatMap[normalizeKadar(kadarName)]
			if !ok {
				appendImportError(result, idx, barcode, "kadar tidak terdaftar: \""+kadarName+"\"")
				continue
			}
			karatID = &id
		}

		harga := int64(0)
		if hargaStr != "" {
			parsed, err := parseHarga(hargaStr)
			if err != nil || parsed < 0 {
				appendImportError(result, idx, barcode, "harga jual tidak valid: \""+hargaStr+"\"")
				continue
			}
			harga = parsed
		}

		kondisi := strings.ToLower(kondisiStr)
		if kondisi == "" {
			kondisi = "baru"
		}
		switch kondisi {
		case "baru", "bekas", "rusak":
		default:
			appendImportError(result, idx, barcode, "kondisi tidak valid: \""+kondisiStr+"\" (gunakan baru, bekas, atau rusak)")
			continue
		}

		var bakiID *uuid.UUID
		if bakiName != "" {
			id, ok := bakiMap[strings.TrimSpace(bakiName)]
			if !ok {
				appendImportError(result, idx, barcode, "baki tidak terdaftar: \""+bakiName+"\"")
				continue
			}
			bakiID = &id
		}

		item := &data.Barang{
			Barcode: barcode,
			Nama:    nama,
			KaratID: karatID,
			Berat:   berat,
			Harga:   harga,
			Kondisi: kondisi,
			BakiID:  bakiID,
		}
		if err := s.repo.Create(item); err != nil {
			appendImportError(result, idx, barcode, err.Error())
			continue
		}
		existingBarcodes[barcode] = true
		result.Imported++
	}

	result.Skipped = result.Total - result.Imported
	writeJSON(w, http.StatusOK, v1.Response{Code: 200, Message: "import selesai", Data: result})
}

var beratUnitPattern = regexp.MustCompile(`(?i)^\s*([0-9]+(?:[.,][0-9]+)?)\s*(gr|gram|g)?\s*$`)

// parseBerat mem-parsing nilai berat dari file import. Label satuan
// "gr", "gram", atau "g" (dengan atau tanpa spasi) diabaikan.
func parseBerat(s string) (float64, error) {
	m := beratUnitPattern.FindStringSubmatch(strings.TrimSpace(s))
	if m == nil {
		return 0, fmt.Errorf("format berat tidak dikenal: %q", s)
	}
	return strconv.ParseFloat(strings.ReplaceAll(m[1], ",", "."), 64)
}

var kadarUnitPattern = regexp.MustCompile(`(?i)^\s*([0-9.]+)\s*(k|karat)?\s*$`)

// normalizeKadar menormalkan nama kadar untuk pencocokan: huruf kecil,
// spasi dihilangkan, dan label "k"/"karat" (dengan atau tanpa spasi)
// diabaikan. Contoh: "24K", "24 K", "24karat", "24 Karat" -> "24".
func normalizeKadar(s string) string {
	m := kadarUnitPattern.FindStringSubmatch(strings.TrimSpace(s))
	if m == nil {
		return strings.ToLower(strings.TrimSpace(s))
	}
	return m[1]
}

var hargaUnitPattern = regexp.MustCompile(`(?i)^\s*(?:rp\.?\s*)?([0-9][0-9.,]*)\s*(rb|ribu)?\s*$`)

// parseHarga mem-parsing harga jual dari file import. Awalan "rp" dan
// separator titik/koma ribuan diabaikan ("Rp 2.500.000", "2,500,000").
// Akhiran "rb"/"ribu" berarti ribuan: "2.5 rb" -> 2500, "2500rb" -> 2500000,
// "1.200 rb" -> 1200000 (kelompok akhir 3 digit dianggap separator ribuan).
func parseHarga(s string) (int64, error) {
	m := hargaUnitPattern.FindStringSubmatch(strings.TrimSpace(s))
	if m == nil {
		return 0, fmt.Errorf("format harga tidak dikenal: %q", s)
	}
	num := m[1]
	if m[2] != "" {
		lastSep := strings.LastIndexAny(num, ".,")
		if lastSep >= 0 && len(num)-lastSep-1 == 3 {
			val, err := strconv.ParseInt(stripHargaSeparators(num), 10, 64)
			if err != nil {
				return 0, err
			}
			return val * 1000, nil
		}
		val, err := strconv.ParseFloat(strings.ReplaceAll(num, ",", "."), 64)
		if err != nil {
			return 0, err
		}
		return int64(val * 1000), nil
	}
	return strconv.ParseInt(stripHargaSeparators(num), 10, 64)
}

func stripHargaSeparators(s string) string {
	s = strings.ReplaceAll(s, ".", "")
	return strings.ReplaceAll(s, ",", "")
}

func appendImportError(result *v1.ImportBarangResult, row int, barcode, message string) {
	result.Errors = append(result.Errors, v1.ImportBarangError{
		Row:     row + 1,
		Barcode: barcode,
		Message: message,
	})
}

func cellAt(row []string, col int) string {
	if col < 0 || col >= len(row) {
		return ""
	}
	return row[col]
}

func normalizeHeader(h string) string {
	h = strings.ToLower(strings.TrimSpace(h))
	if h == "" {
		return ""
	}
	// Abaikan satuan di dalam tanda kurung, mis. "berat (gr)"
	if i := strings.Index(h, "("); i >= 0 {
		h = strings.TrimSpace(h[:i])
	}
	h = strings.ReplaceAll(h, "_", "")
	h = strings.ReplaceAll(h, " ", "")
	switch h {
	case "barcode", "kode", "nobarcode":
		return "barcode"
	case "nama", "name", "namabarang":
		return "nama"
	case "berat", "weight", "gram", "gr", "beratgr":
		return "berat"
	case "kadar", "karat", "kadarkarat":
		return "kadar"
	case "harga", "hargajual", "hargajualrp":
		return "harga"
	case "baki", "namabaki", "lokasi", "tray":
		return "baki"
	case "kondisi", "condition":
		return "kondisi"
	default:
		return h
	}
}

func parseImportFile(ext string, file io.Reader) ([][]string, error) {
	switch ext {
	case ".xlsx":
		return parseXLSX(file)
	case ".csv":
		return parseCSV(file)
	default:
		return nil, fmt.Errorf("format file tidak didukung: %s (gunakan .xlsx atau .csv)", ext)
	}
}

func parseXLSX(file io.Reader) ([][]string, error) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca file Excel: %v", err)
	}
	defer f.Close()
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("file Excel tidak memiliki sheet")
	}
	return f.GetRows(sheets[0])
}

func parseCSV(file io.Reader) ([][]string, error) {
	// Auto-detect delimiter: titik koma umum untuk ekspor Excel id-ID
	head := make([]byte, 4096)
	n, _ := io.ReadFull(file, head)
	if n > 0 {
		head = head[:n]
	}
	content := strings.TrimPrefix(string(head), "\uFEFF")
	semicolons := strings.Count(content, ";")
	commas := strings.Count(content, ",")
	comma := ','
	if semicolons > commas {
		comma = ';'
	}
	reader := csv.NewReader(io.MultiReader(strings.NewReader(content), file))
	reader.Comma = comma
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("gagal membaca file CSV: %v", err)
	}
	return rows, nil
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
	var karatID *uuid.UUID
	if req.KaratID != nil && *req.KaratID != "" {
		parsed, err := uuid.Parse(*req.KaratID)
		if err == nil {
			karatID = &parsed
		}
	}
	b := &data.Barang{
		Barcode:      req.Barcode,
		Nama:         req.Nama,
		KaratID:      karatID,
		Berat:        req.Berat,
		BeratAtribut: req.BeratAtribut,
		Harga:        req.Harga,
		Photo:        req.Photo,
		Kondisi:      req.Kondisi,
		PembelianID:  req.PembelianID,
		BakiID:       req.BakiID,
	}
	if req.Grup != nil {
		b.Grup = *req.Grup
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
	var karatID *uuid.UUID
	if req.KaratID != nil && *req.KaratID != "" {
		parsed, err2 := uuid.Parse(*req.KaratID)
		if err2 == nil {
			karatID = &parsed
		}
	}
	existing.Barcode = req.Barcode
	existing.Nama = req.Nama
	existing.KaratID = karatID
	existing.Berat = req.Berat
	existing.BeratAtribut = req.BeratAtribut
	existing.Harga = req.Harga
	existing.Photo = req.Photo
	existing.Kondisi = req.Kondisi
	existing.PembelianID = req.PembelianID
	existing.BakiID = req.BakiID
	if req.Grup != nil {
		existing.Grup = *req.Grup
	} else {
		existing.Grup = ""
	}
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
