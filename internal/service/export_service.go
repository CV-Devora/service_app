package service

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	v1 "toko-emas/api/v1"

	"github.com/xuri/excelize/v2"
)

// parseDateRange membaca query param from/to dengan format YYYY-MM-DD.
// Kedua parameter wajib diisi dan 'to' tidak boleh sebelum 'from'.
func parseDateRange(r *http.Request) (time.Time, time.Time, error) {
	fromStr := strings.TrimSpace(r.URL.Query().Get("from"))
	toStr := strings.TrimSpace(r.URL.Query().Get("to"))
	if fromStr == "" || toStr == "" {
		return time.Time{}, time.Time{}, fmt.Errorf("rentang tanggal wajib diisi (param from dan to)")
	}
	from, err := time.ParseInLocation("2006-01-02", fromStr, time.Local)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("format tanggal 'from' tidak valid (gunakan YYYY-MM-DD)")
	}
	to, err := time.ParseInLocation("2006-01-02", toStr, time.Local)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("format tanggal 'to' tidak valid (gunakan YYYY-MM-DD)")
	}
	if to.Before(from) {
		return time.Time{}, time.Time{}, fmt.Errorf("tanggal 'to' tidak boleh sebelum 'from'")
	}
	return from, to, nil
}

// writeXLSXResponse membuat file Excel (xlsx) dan mengirimkannya sebagai
// attachment download.
func writeXLSXResponse(w http.ResponseWriter, filename string, headers []string, rows [][]interface{}) error {
	f := excelize.NewFile()
	defer f.Close()
	sheet := f.GetSheetName(0)

	headerStyle, err := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "1F4E78"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"DDEBF7"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	if err != nil {
		return err
	}
	bodyStyle, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	if err != nil {
		return err
	}

	widths := make([]float64, len(headers))
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		if err := f.SetCellValue(sheet, cell, h); err != nil {
			return err
		}
		widths[i] = float64(len([]rune(h)))
	}
	start, _ := excelize.CoordinatesToCellName(1, 1)
	end, _ := excelize.CoordinatesToCellName(len(headers), 1)
	if err := f.SetCellStyle(sheet, start, end, headerStyle); err != nil {
		return err
	}

	for ri, row := range rows {
		for ci, val := range row {
			cell, _ := excelize.CoordinatesToCellName(ci+1, ri+2)
			if err := f.SetCellValue(sheet, cell, val); err != nil {
				return err
			}
			if err := f.SetCellStyle(sheet, cell, cell, bodyStyle); err != nil {
				return err
			}
			if w := float64(len([]rune(fmt.Sprint(val)))); w > widths[ci] {
				widths[ci] = w
			}
		}
	}
	for i, width := range widths {
		col, _ := excelize.ColumnNumberToName(i + 1)
		w := width + 2
		if w > 50 {
			w = 50
		}
		if err := f.SetColWidth(sheet, col, col, w); err != nil {
			return err
		}
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	return f.Write(w)
}

// ExportBarang godoc
// @Summary      Export barang ke Excel
// @Description  Mengunduh data barang inventaris (yang belum terjual) sebagai file .xlsx
// @Tags         Barang
// @Produce      octet-stream
// @Success      200  {file}  binary
// @Router       /api/v1/barang/export [get]
func (s *BarangService) Export(w http.ResponseWriter, r *http.Request) {
	items, err := s.repo.FindAll()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, v1.Response{Code: 500, Message: err.Error()})
		return
	}
	soldIDs, err := s.repo.FindSoldBarangIDs()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, v1.Response{Code: 500, Message: err.Error()})
		return
	}

	headers := []string{"Barcode", "Nama Barang", "Kadar", "Berat (gr)", "Berat Atribut (gr)", "Group", "Harga Jual", "Kondisi"}
	rows := make([][]interface{}, 0, len(items))
	for _, b := range items {
		if soldIDs[b.ID] {
			continue
		}
		kadar := ""
		if b.Karat != nil {
			kadar = b.Karat.Name
		}
		grup := b.Grup
		if grup == "" {
			grup = "-"
		}
		rows = append(rows, []interface{}{
			b.Barcode,
			b.Nama,
			kadar,
			b.Berat,
			b.BeratAtribut,
			grup,
			b.Harga,
			b.Kondisi,
		})
	}

	filename := fmt.Sprintf("daftar_barang.xlsx")
	if err := writeXLSXResponse(w, filename, headers, rows); err != nil {
		writeJSON(w, http.StatusInternalServerError, v1.Response{Code: 500, Message: err.Error()})
		return
	}
}

// ExportPembelian godoc
// @Summary      Export pembelian ke Excel
// @Description  Mengunduh data pembelian dalam rentang tanggal sebagai file .xlsx
// @Tags         Pembelian
// @Produce      octet-stream
// @Param        from  query     string  true  "Tanggal awal (YYYY-MM-DD)"
// @Param        to    query     string  true  "Tanggal akhir (YYYY-MM-DD)"
// @Success      200   {file}    binary
// @Router       /api/v1/pembelian/export [get]
func (s *PembelianService) Export(w http.ResponseWriter, r *http.Request) {
	from, to, err := parseDateRange(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, v1.Response{Code: 400, Message: err.Error()})
		return
	}
	items, err := s.repo.FindByDateRange(from, to)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, v1.Response{Code: 500, Message: err.Error()})
		return
	}

	headers := []string{"No. Faktur", "Nama Pemasok", "Tipe Pemasok", "Berat Nota (gr)", "Harga Nota", "Harga Rata", "Tipe Pembayaran", "Jumlah Pembayaran", "Status"}
	rows := make([][]interface{}, 0, len(items))
	for _, p := range items {
		status := "Pending"
		if p.IsApprove {
			status = "Disetujui"
		}
		rows = append(rows, []interface{}{p.NoFaktur, p.Nama, p.TipePemasok, p.BeratNota, p.HargaNota, p.HargaRata, p.TipePembayaran, p.JumlahPembayaran, status})
	}

	filename := fmt.Sprintf("pembelian_%s_%s.xlsx", from.Format("20060102"), to.Format("20060102"))
	if err := writeXLSXResponse(w, filename, headers, rows); err != nil {
		writeJSON(w, http.StatusInternalServerError, v1.Response{Code: 500, Message: err.Error()})
		return
	}
}

// ExportPenjualan godoc
// @Summary      Export penjualan ke Excel
// @Description  Mengunduh data penjualan dalam rentang tanggal sebagai file .xlsx
// @Tags         Penjualan
// @Produce      octet-stream
// @Param        from  query     string  true  "Tanggal awal (YYYY-MM-DD)"
// @Param        to    query     string  true  "Tanggal akhir (YYYY-MM-DD)"
// @Success      200   {file}    binary
// @Router       /api/v1/penjualan/export [get]
func (s *PenjualanService) Export(w http.ResponseWriter, r *http.Request) {
	from, to, err := parseDateRange(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, v1.Response{Code: 400, Message: err.Error()})
		return
	}
	items, err := s.repo.FindByDateRange(from, to)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, v1.Response{Code: 500, Message: err.Error()})
		return
	}

	headers := []string{"No. Faktur", "Jam", "Nama Pelanggan", "Nama Barang", "Harga Gram", "Harga Jual", "Ongkos", "Total Nilai", "Cash", "Transfer", "Debet", "Nama Sales"}
	rows := make([][]interface{}, 0, len(items))
	for _, p := range items {
		salesName := ""
		if p.KodeSales != nil {
			salesName = fmt.Sprintf("K%d", *p.KodeSales)
		}
		barangNama := ""
		if p.Barang != nil {
			barangNama = p.Barang.Nama
		}
		rows = append(rows, []interface{}{p.NoFaktur, p.CreatedAt.Format("2006-01-02 15:04"), p.Nama, barangNama, p.HargaGram, p.HargaJual, p.Ongkos, p.TotalHarga, p.Cash, p.Transfer, p.Debet, salesName})
	}

	filename := fmt.Sprintf("penjualan_%s_%s.xlsx", from.Format("20060102"), to.Format("20060102"))
	if err := writeXLSXResponse(w, filename, headers, rows); err != nil {
		writeJSON(w, http.StatusInternalServerError, v1.Response{Code: 500, Message: err.Error()})
		return
	}
}
