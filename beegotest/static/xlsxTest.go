// xlsx文件的读取与写入
// 注意单元格合并的情况
package main

import (
	"fmt"

	"github.com/tealeg/xlsx"
)

func main() {
	// 写xlsx文件
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")
	if err != nil {
		fmt.Println(err)
		return
	}
	row := sheet.AddRow()
	cell := row.AddCell()
	cell.Value = "I am a cell"
	cell1 := row.AddCell()
	cell1.Value = "I am a cell1"
	row2 := sheet.AddRow()
	cell2 := row2.AddCell()
	cell2.Value = "I am a cell2"
	cell22 := row2.AddCell()
	cell22.Value = "I am a cell22"
	err = file.Save("MyXLSXFile.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}

	// 读xlsx文件
	xlFile, err := xlsx.OpenFile("MyXLSXFile.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, sheet := range xlFile.Sheets {
		fmt.Printf("sheet name: %s\n", sheet.Name)
		for i, row := range sheet.Rows {
			fmt.Printf("row: %d\n", i)
			for j, cell := range row.Cells {
				fmt.Printf("cell%d: ", j)
				text, _ := cell.String()
				fmt.Printf("%s\n", text)
			}
		}
		break
	}
}
