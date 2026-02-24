package word

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/xuri/excelize/v2"
)

func (w *WordService) ImportFromExcel(ctx context.Context, file io.Reader, category string, addedBy string) (int, error) {
	// reading the imported file
	buff := new(bytes.Buffer)
	if _, err := io.Copy(buff, file); err != nil {
		return 0, fmt.Errorf("Error reading file to buff, %w", err)
	}

	// opens the excel file

	xlx, err := excelize.OpenReader(buff)
	if err != nil {
		return 0, fmt.Errorf("error opening the excel file, %w", err)
	}

	sheets := xlx.GetSheetList()
	if len(sheets) == 0 {
		return 0, fmt.Errorf("No sheet in file, %w", err)
	}

	rows, err := xlx.GetRows(sheets[0])
	if err != nil {
		return 0, fmt.Errorf("Error reading rows, %w", err)
	}

	count := 0
	// loops whrough each row, skipping headers and empty cells
	for i, row := range rows {
		//skipps hearde
		if i == 0 {
			continue
		}
		// skips empty cells
		if len(row) == 0 {
			continue
		}

		word := strings.TrimSpace(row[0])
		if word == "" {
			continue
		}

		err = w.AddWord(ctx, AddWordInput{
			Word:     word,
			Category: category,
			AddedBy:  addedBy,
		})

		if err != nil {
			continue
		}
		count++

	}
	return count, nil
}
