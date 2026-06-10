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
	photostore "github.com/qirs-mezgeb/api/internal/photos"
)

const (
	pdfBrandR = 45
	pdfBrandG = 82
	pdfBrandB = 130

	pdfMarginL   = 15.0
	pdfMarginR   = 15.0
	pdfMarginTop = 34.0
	pdfMarginBot = 20.0
)

type pdfRecord struct {
	Title      string
	RecordID   string
	Status     string
	ApprovedAt *time.Time
	Fields     [][2]string
	PhotoPath  string
}

func buildImmovablePDF(record *models.ImmovableRecord, recordPhotos []models.RecordPhoto, mediaPath string) ([]byte, error) {
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
		PhotoPath:  firstPhotoPath(recordPhotos, mediaPath),
	})
}

func buildMovablePDF(record *models.MovableRecord, recordPhotos []models.RecordPhoto, mediaPath string) ([]byte, error) {
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
		PhotoPath:  firstPhotoPath(recordPhotos, mediaPath),
	})
}

func buildRecordPDF(record pdfRecord) ([]byte, error) {
	printedAt := time.Now().UTC()

	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(pdfMarginL, pdfMarginTop, pdfMarginR)
	pdf.SetAutoPageBreak(true, pdfMarginBot)

	if pdf.RegisterImageOptionsReader(
		"logo",
		fpdf.ImageOptions{ImageType: "JPG", ReadDpi: true},
		bytes.NewReader(logoJPEG),
	) == nil {
		return nil, fmt.Errorf("register logo: invalid image data")
	}

	pdf.SetHeaderFunc(func() {
		drawPDFHeader(pdf, record, pdf.PageNo() == 1)
	})
	pdf.SetFooterFunc(func() {
		drawPDFFooter(pdf, printedAt)
	})
	pdf.AliasNbPages("")

	pdf.AddPage()
	drawRecordTitleBlock(pdf, record)
	drawRecordDetails(pdf, record)

	if record.PhotoPath != "" {
		drawRecordPhoto(pdf, record.PhotoPath)
	}

	if record.Status == string(models.StatusApproved) {
		drawApprovedWatermark(pdf)
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("render pdf: %w", err)
	}

	return buf.Bytes(), nil
}

func drawPDFHeader(pdf *fpdf.Fpdf, record pdfRecord, firstPage bool) {
	pdf.ImageOptions("logo", pdfMarginL, 8, 18, 0, false, fpdf.ImageOptions{}, 0, "")

	pdf.SetXY(36, 9)
	pdf.SetFont("Helvetica", "B", 14)
	pdf.SetTextColor(pdfBrandR, pdfBrandG, pdfBrandB)
	pdf.Cell(0, 6, "Qirs Mezgeb")

	pdf.SetXY(36, 15)
	pdf.SetFont("Helvetica", "", 9)
	pdf.SetTextColor(70, 70, 70)
	pdf.Cell(0, 5, "Heritage Registration System")

	pdf.SetXY(36, 20)
	pdf.SetFont("Helvetica", "I", 8)
	pdf.Cell(0, 4, "Harari Culture, Heritage & Tourism Bureau")

	lineY := 28.0
	if !firstPage {
		pdf.SetXY(36, 14)
		pdf.SetFont("Helvetica", "", 8)
		pdf.SetTextColor(100, 100, 100)
		pdf.Cell(0, 4, sanitizePDFText(record.RecordID))
		lineY = 22
	}

	pdf.SetDrawColor(pdfBrandR, pdfBrandG, pdfBrandB)
	pdf.SetLineWidth(0.5)
	pdf.Line(pdfMarginL, lineY, 210-pdfMarginR, lineY)
	pdf.SetTextColor(0, 0, 0)
}

func drawPDFFooter(pdf *fpdf.Fpdf, printedAt time.Time) {
	pageW, pageH := pdf.GetPageSize()
	footerY := pageH - 14

	pdf.SetDrawColor(200, 200, 200)
	pdf.SetLineWidth(0.2)
	pdf.Line(pdfMarginL, footerY, pageW-pdfMarginR, footerY)

	pdf.SetY(footerY + 2)
	pdf.SetFont("Helvetica", "", 7)
	pdf.SetTextColor(90, 90, 90)
	pdf.CellFormat(55, 4, "Qirs Mezgeb Heritage Registry", "", 0, "L", false, 0, "")
	pdf.CellFormat(
		pageW-pdfMarginL-pdfMarginR-110,
		4,
		fmt.Sprintf("Page %d of {nb}", pdf.PageNo()),
		"",
		0,
		"C",
		false,
		0,
		"",
	)
	pdf.CellFormat(55, 4, printedAt.Format("2006-01-02 15:04 UTC"), "", 0, "R", false, 0, "")

	pdf.Ln(4)
	pdf.SetFont("Helvetica", "I", 6)
	pdf.CellFormat(
		0,
		3,
		"Official heritage registration document — generated electronically",
		"",
		0,
		"C",
		false,
		0,
		"",
	)
	pdf.SetTextColor(0, 0, 0)
}

func drawRecordTitleBlock(pdf *fpdf.Fpdf, record pdfRecord) {
	pdf.SetFont("Helvetica", "B", 13)
	pdf.SetTextColor(pdfBrandR, pdfBrandG, pdfBrandB)
	pdf.MultiCell(0, 7, record.Title, "", "L", false)
	pdf.Ln(3)

	pdf.SetFillColor(240, 245, 250)
	pdf.SetDrawColor(pdfBrandR, pdfBrandG, pdfBrandB)
	pdf.SetLineWidth(0.2)

	metaY := pdf.GetY()
	pdf.Rect(pdfMarginL, metaY, 180, 18, "FD")

	pdf.SetXY(pdfMarginL+3, metaY+3)
	pdf.SetFont("Helvetica", "B", 9)
	pdf.SetTextColor(60, 60, 60)
	pdf.Cell(38, 5, "Record ID:")
	pdf.SetFont("Helvetica", "", 9)
	pdf.SetTextColor(0, 0, 0)
	pdf.Cell(52, 5, sanitizePDFText(record.RecordID))

	pdf.SetXY(pdfMarginL+3, metaY+10)
	pdf.SetFont("Helvetica", "B", 9)
	pdf.SetTextColor(60, 60, 60)
	pdf.Cell(38, 5, "Status:")
	pdf.SetFont("Helvetica", "", 9)
	pdf.SetTextColor(0, 0, 0)
	pdf.Cell(52, 5, humanizeStatus(record.Status))

	if record.ApprovedAt != nil {
		pdf.SetXY(105, metaY+3)
		pdf.SetFont("Helvetica", "B", 9)
		pdf.SetTextColor(60, 60, 60)
		pdf.Cell(38, 5, "Approved:")
		pdf.SetFont("Helvetica", "", 9)
		pdf.SetTextColor(0, 0, 0)
		pdf.Cell(0, 5, record.ApprovedAt.UTC().Format("2006-01-02 15:04 UTC"))
	}

	pdf.SetY(metaY + 22)
}

func drawRecordDetails(pdf *fpdf.Fpdf, record pdfRecord) {
	pdf.SetFont("Helvetica", "B", 11)
	pdf.SetTextColor(pdfBrandR, pdfBrandG, pdfBrandB)
	pdf.Cell(0, 7, "Record Details")
	pdf.Ln(2)

	pdf.SetDrawColor(pdfBrandR, pdfBrandG, pdfBrandB)
	pdf.SetLineWidth(0.3)
	pdf.Line(pdfMarginL, pdf.GetY(), pdfMarginL+40, pdf.GetY())
	pdf.Ln(4)

	pdf.SetFont("Helvetica", "", 10)
	pdf.SetTextColor(0, 0, 0)

	for _, field := range record.Fields {
		if strings.TrimSpace(field[1]) == "" {
			continue
		}
		pdf.SetFont("Helvetica", "B", 9)
		pdf.SetTextColor(60, 60, 60)
		pdf.Cell(48, 6, field[0]+":")
		pdf.SetFont("Helvetica", "", 9)
		pdf.SetTextColor(0, 0, 0)
		pdf.MultiCell(0, 6, sanitizePDFText(field[1]), "", "L", false)
		pdf.Ln(1)
	}
}

func drawRecordPhoto(pdf *fpdf.Fpdf, photoPath string) {
	pdf.Ln(4)
	pdf.SetFont("Helvetica", "B", 11)
	pdf.SetTextColor(pdfBrandR, pdfBrandG, pdfBrandB)
	pdf.Cell(0, 7, "Primary Photo")
	pdf.Ln(2)
	pdf.SetDrawColor(pdfBrandR, pdfBrandG, pdfBrandB)
	pdf.Line(pdfMarginL, pdf.GetY(), pdfMarginL+40, pdf.GetY())
	pdf.Ln(5)

	options := fpdf.ImageOptions{ImageType: imageTypeForPath(photoPath), ReadDpi: true}
	pdf.ImageOptions(photoPath, pdfMarginL, pdf.GetY(), 80, 0, false, options, 0, "")
	pdf.Ln(55)
}

func drawApprovedWatermark(pdf *fpdf.Fpdf) {
	pdf.SetAlpha(0.15, "Normal")
	pdf.SetFont("Helvetica", "B", 48)
	pdf.SetTextColor(200, 0, 0)
	pdf.Text(35, 180, "APPROVED")
	pdf.SetAlpha(1, "Normal")
	pdf.SetTextColor(0, 0, 0)
}

func humanizeStatus(status string) string {
	return strings.ReplaceAll(strings.ReplaceAll(status, "_", " "), "-", " ")
}

func firstPhotoPath(recordPhotos []models.RecordPhoto, mediaPath string) string {
	if len(recordPhotos) == 0 {
		return ""
	}
	path := recordPhotos[0].FilePath
	candidates := []string{
		photostore.ResolveAbsolutePath(mediaPath, path),
		filepath.Join(mediaPath, filepath.Base(path)),
	}
	for _, abs := range candidates {
		if _, err := os.Stat(abs); err == nil {
			return abs
		}
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
