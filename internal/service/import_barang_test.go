package service

import (
	"strings"
	"testing"

	"github.com/xuri/excelize/v2"
)

func TestParseCSVComma(t *testing.T) {
	csv := "barcode,baki,nama,berat\n000000001,Baki 1,Cincin Emas Polos,5.500\n000000002,Baki 2,Gelang Emas,12,25\n"
	rows, err := parseCSV(strings.NewReader(csv))
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 3 {
		t.Fatalf("expected 3 rows, got %d", len(rows))
	}
	if rows[1][0] != "000000001" || rows[1][3] != "5.500" {
		t.Fatalf("unexpected row: %v", rows[1])
	}
}

func TestParseCSVSemicolon(t *testing.T) {
	csv := "barcode;baki;nama;berat\n000000001;Baki 1;Cincin Emas Polos;5,500\n"
	rows, err := parseCSV(strings.NewReader(csv))
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 || rows[1][3] != "5,500" {
		t.Fatalf("unexpected rows: %v", rows)
	}
}

func TestParseCSVWithBOM(t *testing.T) {
	csv := "\uFEFFbarcode,baki,nama,berat\n000000001,Baki 1,Cincin Emas Polos,5.5\n"
	rows, err := parseCSV(strings.NewReader(csv))
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 || rows[0][0] != "barcode" {
		t.Fatalf("unexpected rows: %v", rows)
	}
}

func TestParseXLSX(t *testing.T) {
	f := excelize.NewFile()
	sheet := f.GetSheetName(0)
	rows := [][]interface{}{
		{"barcode", "baki", "nama", "berat"},
		{"000000001", "Baki 1", "Cincin Emas Polos", 5.5},
	}
	for i, row := range rows {
		for j, cell := range row {
			cellRef, _ := excelize.CoordinatesToCellName(j+1, i+1)
			f.SetCellValue(sheet, cellRef, cell)
		}
	}
	buf, err := f.WriteToBuffer()
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := parseXLSX(buf)
	if err != nil {
		t.Fatal(err)
	}
	if len(parsed) != 2 || parsed[1][3] != "5.5" {
		t.Fatalf("unexpected rows: %v", parsed)
	}
}

func TestParseBerat(t *testing.T) {
	cases := map[string]float64{
		"5.5":       5.5,
		"5,5":       5.5,
		"5.5 gr":    5.5,
		"5.5gr":     5.5,
		"5.5 GR":    5.5,
		"12.25 g":   12.25,
		"2,5 gram":  2.5,
		" 3.75 gr ": 3.75,
		"7 gram":    7,
	}
	for input, want := range cases {
		got, err := parseBerat(input)
		if err != nil {
			t.Errorf("parseBerat(%q) error: %v", input, err)
			continue
		}
		if got != want {
			t.Errorf("parseBerat(%q) = %v, want %v", input, got, want)
		}
	}
	for _, invalid := range []string{"abc", "5.5 kg", "gr", ""} {
		if _, err := parseBerat(invalid); err == nil {
			t.Errorf("parseBerat(%q) should return error", invalid)
		}
	}
}

func TestNormalizeKadar(t *testing.T) {
	cases := map[string]string{
		"24K":        "24",
		"24 K":       "24",
		"24k":        "24",
		"22 Karat":   "22",
		"22karat":    "22",
		"18 Karat ":  "18",
		"24":         "24",
		"750":        "750",
		" Emas 24K ": "emas 24k",
	}
	for input, want := range cases {
		if got := normalizeKadar(input); got != want {
			t.Errorf("normalizeKadar(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestParseHarga(t *testing.T) {
	cases := map[string]int64{
		"2500000":        2500000,
		"2.500.000":      2500000,
		"2,500,000":      2500000,
		"Rp 2.500.000":   2500000,
		"Rp.2.500.000":   2500000,
		"Rp2.500.000":    2500000,
		"2.5 rb":         2500,
		"2,5 rb":         2500,
		"2500rb":         2500000,
		"1.200 ribu":     1200000,
		"  Rp 3.000 rb ": 3000000,
	}
	for input, want := range cases {
		got, err := parseHarga(input)
		if err != nil {
			t.Errorf("parseHarga(%q) error: %v", input, err)
			continue
		}
		if got != want {
			t.Errorf("parseHarga(%q) = %d, want %d", input, got, want)
		}
	}
	for _, invalid := range []string{"abc", "Rp abc", "2.5.5 rb", "1.000.000 rb rb"} {
		if _, err := parseHarga(invalid); err == nil {
			t.Errorf("parseHarga(%q) should return error", invalid)
		}
	}
}

func TestNormalizeHeader(t *testing.T) {
	cases := map[string]string{
		"barcode":         "barcode",
		"Barcode":         "barcode",
		"NAMA BARANG":     "nama",
		"berat (gr)":      "berat",
		"Berat (Gram)":    "berat",
		"Baki":            "baki",
		"Nama Baki":       "baki",
		"name":            "nama",
		"weight":          "berat",
		"lokasi":          "baki",
		"  Nama  Barang ": "nama",
		"kadar":           "kadar",
		"Karat":           "kadar",
		"harga jual":      "harga",
		"harga_jual":      "harga",
		"Harga Jual (Rp)": "harga",
		"kondisi":         "kondisi",
		"Condition":       "kondisi",
	}
	for input, want := range cases {
		if got := normalizeHeader(input); got != want {
			t.Errorf("normalizeHeader(%q) = %q, want %q", input, got, want)
		}
	}
}
