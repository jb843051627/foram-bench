package format

import (
	"encoding/csv"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/jb843051627/foram-bench/internal/model"
)

type CSVRow struct {
	ID           string
	SectionID    string
	Taxon        string
	Count        int
	Preservation string
	Confidence   float64
	ObservedAt   time.Time
}

func ObservationsCSV(observations []model.Observation) (string, error) {
	var builder strings.Builder
	writer := csv.NewWriter(&builder)
	if err := writer.Write([]string{"id", "section_id", "taxon", "count", "preservation", "confidence", "observed_at"}); err != nil {
		return "", err
	}
	for _, item := range observations {
		if err := writer.Write([]string{item.ID, item.SectionID, item.Taxon, strconv.Itoa(item.Count), item.Preservation,
			Decimal(item.Confidence, 3), item.ObservedAt.Format("2006-01-02T15:04:05Z07:00")}); err != nil {
			return "", err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", err
	}
	return builder.String(), nil
}

func ReadObservationCSV(input io.Reader) ([]model.Observation, error) {
	reader := csv.NewReader(input)
	if _, err := reader.Read(); err != nil {
		return nil, err
	}
	result := make([]model.Observation, 0)
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(row) != 7 {
			return nil, io.ErrUnexpectedEOF
		}
		count, err := strconv.Atoi(row[3])
		if err != nil {
			return nil, err
		}
		confidence, err := strconv.ParseFloat(row[5], 64)
		if err != nil {
			return nil, err
		}
		observedAt, err := time.Parse(time.RFC3339, row[6])
		if err != nil {
			return nil, err
		}
		result = append(result, model.Observation{ID: row[0], SectionID: row[1], Taxon: row[2], Count: count,
			Preservation: row[4], Confidence: confidence, ObservedAt: observedAt})
	}
	return result, nil
}

func WriteRows(writer io.Writer, rows []CSVRow) error {
	csvWriter := csv.NewWriter(writer)
	if err := csvWriter.Write([]string{"id", "section_id", "taxon", "count", "preservation", "confidence", "observed_at"}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := csvWriter.Write([]string{row.ID, row.SectionID, row.Taxon, strconv.Itoa(row.Count), row.Preservation,
			Decimal(row.Confidence, 3), row.ObservedAt.Format(time.RFC3339)}); err != nil {
			return err
		}
	}
	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		return err
	}
	return nil
}

func RowsFromObservations(observations []model.Observation) []CSVRow {
	rows := make([]CSVRow, 0, len(observations))
	for _, observation := range observations {
		rows = append(rows, CSVRow{ID: observation.ID, SectionID: observation.SectionID, Taxon: observation.Taxon,
			Count: observation.Count, Preservation: observation.Preservation, Confidence: observation.Confidence,
			ObservedAt: observation.ObservedAt})
	}
	return rows
}
