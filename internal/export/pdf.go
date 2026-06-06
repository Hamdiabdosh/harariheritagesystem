package export

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-pdf/fpdf"

	"github.com/qirs-mezgeb/api/internal/models"
)

type pdfRecord struct {
	Title      string
	RecordID   string
	Status     string
	ApprovedAt *time.Time
	Fields     [][2]string
	PhotoPath  string
}

func buildImmovablePDF(record *models.ImmovableRecord, photos []models.RecordPhoto, mediaPath string) ([]byte, error) {
	fields := [][2]string{
		{"Name (Amharic)", record.NameAmharic},
		{"Name (Local)", stringValue(record.NameLocal)},
		{"Category", strings.Join(record.Category, ", ")},
		{"Woreda", record.Woreda},
		{"Kebele", record.Kebele},
		{"House Number", stringValue(record.HouseNumber)},
		{"Owner", stringValue(record.OwnerName)},
		{"GPS East", floatValue(record.GPSEast)},
		{"GPS North", floatValue(record.GPSNorth)},
		{"Built By", stringValue(record.BuiltBy)},
		{"Construction Period", stringValue(record.ConstructionPeriod)},
		{"Description", stringValue(record.Description)},
		{"Overall Condition", stringValueEnum(record.OverallCondition)},
		{"Notes", stringValue(record.Notes)},
	}

	return buildRecordPDF(pdfRecord{
		Title:      "Immovable Heritage Record (Form 02)",
		RecordID:   record.RecordID,
		Status:     string(record.Status),
		ApprovedAt: record.ApprovedAt,
		Fields:     fields,
		PhotoPath:  firstPhotoPath(photos, mediaPath),
	})
}

func buildMovablePDF(record *models.MovableRecord, photos []models.RecordPhoto, mediaPath string) ([]byte, error) {
	fields := [][2]string{
		{"Name (Amharic)", record.NameAmharic},
		{"Name (Local)", stringValue(record.NameLocal)},
		{"Category", stringValue(record.Category)},
		{"Location", stringValue(record.LocationName)},
		{"Woreda", stringValue(record.Woreda)},
		{"Kebele", stringValue(record.Kebele)},
		{"Owner", stringValue(record.OwnerName)},
		{"Storage Location", stringValueEnum(record.StorageLocation)},
		{"Made By", stringValue(record.MadeBy)},
		{"Period Made", stringValue(record.PeriodMade)},
		{"Description", stringValue(record.Description)},
		{"Condition", stringValueEnum(record.Condition)},
		{"Notes", stringValue(record.Notes)},
	}

	return buildRecordPDF(pdfRecord{
		Title:      "Movable Heritage Record (Form 01)",
		RecordID:   record.RecordID,
		Status:     string(record.Status),
		ApprovedAt: record.ApprovedAt,
		Fields:     fields,
		PhotoPath:  firstPhotoPath(photos, mediaPath),
	})
}

func buildRecordPDF(record pdfRecord) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)
	pdf.AddPage()
	pdf.SetFont("Helvetica", "B", 16)
	pdf.Cell(0, 10, "Qirs Mezgeb Heritage Registry")
	pdf.Ln(12)

	pdf.SetFont("Helvetica", "B", 13)
	pdf.MultiCell(0, 7, record.Title, "", "L", false)
	pdf.Ln(2)

	pdf.SetFont("Helvetica", "", 11)
	pdf.Cell(45, 7, "Record ID:")
	pdf.SetFont("Helvetica", "B", 11)
	pdf.Cell(0, 7, record.RecordID)
	pdf.Ln(7)

	pdf.SetFont("Helvetica", "", 11)
	pdf.Cell(45, 7, "Status:")
	pdf.Cell(0, 7, record.Status)
	pdf.Ln(7)

	if record.ApprovedAt != nil {
		pdf.Cell(45, 7, "Approved At:")
		pdf.Cell(0, 7, record.ApprovedAt.UTC().Format("2006-01-02 15:04 MST"))
		pdf.Ln(10)
	}

	pdf.SetFont("Helvetica", "B", 12)
	pdf.Cell(0, 8, "Record Details")
	pdf.Ln(8)

	pdf.SetFont("Helvetica", "", 10)
	for _, field := range record.Fields {
		if strings.TrimSpace(field[1]) == "" {
			continue
		}
		pdf.SetFont("Helvetica", "B", 10)
		pdf.Cell(45, 6, field[0]+":")
		pdf.SetFont("Helvetica", "", 10)
		pdf.MultiCell(0, 6, sanitizePDFText(field[1]), "", "L", false)
		pdf.Ln(1)
	}

	if record.PhotoPath != "" {
		pdf.Ln(4)
		pdf.SetFont("Helvetica", "B", 12)
		pdf.Cell(0, 8, "Primary Photo")
		pdf.Ln(6)

		options := fpdf.ImageOptions{ImageType: imageTypeForPath(record.PhotoPath), ReadDpi: true}
		pdf.ImageOptions(record.PhotoPath, 15, pdf.GetY(), 80, 0, false, options, 0, "")
	}

	if record.Status == string(models.StatusApproved) {
		pdf.SetAlpha(0.2, "Normal")
		pdf.SetFont("Helvetica", "B", 36)
		pdf.SetTextColor(200, 0, 0)
		pdf.Text(50, 180, "APPROVED")
		pdf.SetAlpha(1, "Normal")
		pdf.SetTextColor(0, 0, 0)
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("render pdf: %w", err)
	}

	return buf.Bytes(), nil
}

func firstPhotoPath(photos []models.RecordPhoto, mediaPath string) string {
	if len(photos) == 0 {
		return ""
	}
	path := photos[0].FilePath
	if filepath.IsAbs(path) {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	fullPath := filepath.Join(mediaPath, path)
	if _, err := os.Stat(fullPath); err == nil {
		return fullPath
	}
	return ""
}

func imageTypeForPath(path string) string {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".png":
		return "PNG"
	default:
		return "JPG"
	}
}

func sanitizePDFText(value string) string {
	return strings.Map(func(r rune) rune {
		if r < 256 {
			return r
		}
		return '?'
	}, value)
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func stringValueEnum[T ~string](value *T) string {
	if value == nil {
		return ""
	}
	return string(*value)
}

func floatValue(value *float64) string {
	if value == nil {
		return ""
	}
	return fmt.Sprintf("%.6f", *value)
}
