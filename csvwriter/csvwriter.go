package csvwriter

import (
	"encoding/csv"
	"esgi-archive-crypto-history/kraken"
	"fmt"
	"os"
	"time"
)

func GetCSVFileName(t time.Time) string {
	minuteWindow := (t.Minute() / 5) * 5
	return fmt.Sprintf("pair_%02d_%02d_%d_%02d_%02d.csv", t.Day(), t.Month(), t.Year(), t.Hour(), minuteWindow)
}

func WriteAssetPairsCSV(filename string, pairs map[string]kraken.AssetPair) error {
	var file *os.File
	var err error

	if _, err = os.Stat(filename); os.IsNotExist(err) {
		file, err = os.Create(filename)
		if err != nil {
			return err
		}
		writer := csv.NewWriter(file)
		header := []string{"PairKey", "Altname", "Wsname"}
		if err := writer.Write(header); err != nil {
			file.Close()
			return err
		}
		writer.Flush()
	} else {
		file, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			return err
		}
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	for key, pair := range pairs {
		record := []string{key, pair.Altname, pair.Wsname}
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error()
}
