package query_service

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"
)

func CreateLogCSVFile(logs []map[string]string, keys []string) string {
	filePath := fmt.Sprintf("room_query_%d.csv", time.Now().Unix())
	file, _ := os.Create(filePath)
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
	return filePath
}

func findIndex(in []string, val string) int {
	for i, v := range in {
		if v == val {
			return i
		}
	}
	return -1
}

func CreateLogXLSCFile(logs []map[string]string, keys []string) string {
	filePath := fmt.Sprintf("room_query_%d.xlsx", time.Now().Unix())
	file := excelize.NewFile()
	streamWriter, _ := file.NewStreamWriter("Sheet1")

	dateIndex := findIndex(keys, "createdat")
	streamWriter.SetColWidth(dateIndex+1, dateIndex+1, 15)
	messageIndex := findIndex(keys, "message")
	streamWriter.SetColWidth(messageIndex+1, messageIndex+1, 44)

	dateStyleId, _ := file.NewStyle(&excelize.Style{NumFmt: 22})

	messageStyle, _ := file.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Indent:     1,
			Vertical:   "bottom",
			Horizontal: "left",
			WrapText:   true,
		},
		Font: &excelize.Font{
			Size: 12,
		},
	})

	headerStyleId, _ := file.NewStyle(&excelize.Style{
		Font: &excelize.Font{Size: 14},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
		},
	})
	var header []interface{}
	for _, key := range keys {
		header = append(header, excelize.Cell{StyleID: headerStyleId, Value: key})
	}
	streamWriter.SetRow("A1", header, excelize.RowOpts{Height: 44})
	keyLength := len(keys)

	for i, log := range logs {
		row := make([]interface{}, keyLength)
		for j, k := range keys {
			if j == dateIndex {
				num, numErr := strconv.ParseInt(log[k], 10, 64)
				if numErr == nil {
					t := time.Unix(num/1000, 0)
					row[j] = excelize.Cell{Value: t, StyleID: dateStyleId}
					continue
				}
			}
			row[j] = excelize.Cell{Value: log[k], StyleID: messageStyle}
		}
		cell, _ := excelize.CoordinatesToCellName(1, i+2)
		streamWriter.SetRow(cell, row)
	}

	streamWriter.Flush()
	file.SaveAs(filePath)
	return filePath
}
