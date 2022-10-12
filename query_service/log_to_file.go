package query_service

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"
)

func CreateLogCSVFile(logs []map[string]string, keys []string) string {
	fileName := fmt.Sprintf("room_query_%d.csv", time.Now().Unix())
	file, _ := os.Create(fileName)
	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Write(keys)
	for _, item := range logs {
		itemArray := make([]string, len(keys))
		for ki, k := range keys {
			itemArray[ki] = item[k]
		}
		writer.Write(itemArray[:])
	}
	return fileName
}
