package export

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"time"

	"github.com/qirs-mezgeb/api/internal/dashboard"
)

func buildCSV(records []dashboard.RecordSummary) ([]byte, error) {
	buf := &bytes.Buffer{}
	writer := csv.NewWriter(buf)

	if err := writer.Write([]string{
		"record_type",
		"record_id",
		"name_amharic",
		"status",
		"woreda",
		"kebele",
		"registrar_id",
		"created_at",
		"updated_at",
	}); err != nil {
		return nil, fmt.Errorf("write csv header: %w", err)
	}

	for _, record := range records {
		if err := writer.Write([]string{
			string(record.RecordType),
			record.RecordID,
			record.NameAmharic,
			string(record.Status),
			stringPtr(record.Woreda),
			stringPtr(record.Kebele),
			record.RegistrarID.String(),
			formatTime(record.CreatedAt),
			formatTime(record.UpdatedAt),
		}); err != nil {
			return nil, fmt.Errorf("write csv row: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("flush csv: %w", err)
	}

	return buf.Bytes(), nil
}

func stringPtr(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func formatTime(value time.Time) string {
	return value.UTC().Format(time.RFC3339)
}

func csvFilename() string {
	return fmt.Sprintf("qirs-mezgeb-records-%s.csv", time.Now().UTC().Format("20060102-150405"))
}

func pdfFilename(recordID string) string {
	return fmt.Sprintf("%s.pdf", recordID)
}
