package backtest

import (
	"archive/zip"
	"encoding/csv"
	"errors"
	"io"

	"github.com/ydm/commons"
)

type EitherRow struct {
	Error  error
	Values []string
	Origin string
}

func ReadCSV(origin string, reader io.Reader) <-chan EitherRow {
	rows := make(chan EitherRow)

	go func() {
		commons.Checker.Push()
		defer commons.Checker.Pop()

		defer close(rows)

		reader := csv.NewReader(reader)

		for {
			record, err := reader.Read()

			if errors.Is(err, io.EOF) {
				return
			}

			rows <- EitherRow{
				Error:  err,
				Values: record,
				Origin: origin,
			}
		}
	}()

	return rows
}

func ReadCSVZipFile(file *zip.File) <-chan EitherRow {
	rows := make(chan EitherRow)

	go func() {
		commons.Checker.Push()
		defer commons.Checker.Pop()

		defer close(rows)

		reader, err := file.Open()
		if err != nil {
			panic(err)
		}

		defer reader.Close()

		for row := range ReadCSV(file.Name, reader) {
			rows <- row
		}
	}()

	return rows
}

func ReadCSVZipArchive(path string) <-chan EitherRow {
	rows := make(chan EitherRow)

	go func() {
		commons.Checker.Push()
		defer commons.Checker.Pop()

		defer close(rows)

		read, err := zip.OpenReader(path)
		if err != nil {
			panic(err)
		}

		defer read.Close()

		for _, file := range read.File {
			for row := range ReadCSVZipFile(file) {
				rows <- row
			}
		}
	}()

	return rows
}
