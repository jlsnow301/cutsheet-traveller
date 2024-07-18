package fileutils

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xuri/excelize/v2"
)

func ValidateFolder(file os.DirEntry) bool {
	ignoreDirs := []string{".git", ".vscode", "src"}

	for _, ignoreDir := range ignoreDirs {
		if file.Name() == ignoreDir {
			return false
		}
	}

	return true
}

func getCutsheets(files []os.DirEntry) ([]os.DirEntry, error) {
	var cutsheets []os.DirEntry

	for _, entry := range files {
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()
		if filepath.Ext(fileName) != ".pdf" {
			continue
		}

		cutsheets = append(cutsheets, entry)
	}

	return cutsheets, nil
}

func CreateExcelFile(employeeOrders map[string][]orderInfo, errors []errorInfo) error {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Rainbow colors
	colors := []string{"#FF0000", "#FF7F00", "#FFFF00", "#00FF00", "#0000FF", "#8B00FF"}
	colorIndex := 0

	var firstSheet string
	for employee, orders := range employeeOrders {
		// Use employee name as sheet name, replacing any invalid characters
		sheetName := sanitizeSheetName(employee)
		index, err := f.NewSheet(sheetName)
		if err != nil {
			return err
		}

		// Set first sheet name
		if firstSheet == "" {
			firstSheet = sheetName
		}

		// Set employee name and color
		f.SetCellValue(sheetName, "A1", employee)
		bgColor := colors[colorIndex]
		textColor := getContrastColor(bgColor)
		style, _ := f.NewStyle(&excelize.Style{
			Fill: excelize.Fill{Type: "pattern", Color: []string{bgColor}, Pattern: 1},
			Font: &excelize.Font{Bold: true, Color: textColor},
		})
		f.SetCellStyle(sheetName, "A1", "A1", style)

		// Set headers
		headers := []string{"Order ID", "Mileage", "Date", "Origin", "Destination"}
		for col, header := range headers {
			cell := string(rune('A'+col)) + "2"
			f.SetCellValue(sheetName, cell, header)
		}

		// Set header style
		headerStyle, _ := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{Bold: true},
			Fill: excelize.Fill{Type: "pattern", Color: []string{"#E0E0E0"}, Pattern: 1},
		})
		f.SetCellStyle(sheetName, "A2", "E2", headerStyle)

		// Fill in order data
		row := 3
		totalMileage := 0.0
		for _, order := range orders {
			f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), order.OrderID)
			f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), order.Mileage)
			f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), order.Date)
			f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), order.Origin)
			f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), order.Destination)
			totalMileage += order.Mileage
			row++
		}

		// Set total mileage with the same color scheme as the employee header
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row+1), "Total Mileage:")
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row+1), totalMileage)
		f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row+1), fmt.Sprintf("B%d", row+1), style)

		// Set column widths
		f.SetColWidth(sheetName, "A", "E", 15)

		// Set active sheet
		f.SetActiveSheet(index)

		// Move to next color
		colorIndex = (colorIndex + 1) % len(colors)
	}

	defaultSheet := f.GetSheetName(0)
	if defaultSheet != firstSheet {
		f.DeleteSheet(defaultSheet)
	}

	// Set the first created sheet as active
	if firstSheet != "" {
		firstSheetIndex, _ := f.GetSheetIndex(firstSheet)
		f.SetActiveSheet(firstSheetIndex)
	}

	// Add error information (same as before)
	if len(errors) == 0 {
		return f.SaveAs("orders_report.xlsx")
	}

	errorSheetName := "Errors"
	f.NewSheet(errorSheetName)

	f.SetCellValue(errorSheetName, "A1", "Errors")
	errorHeaderStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 14},
	})
	f.SetCellStyle(errorSheetName, "A1", "A1", errorHeaderStyle)

	row := 2
	currentEmployee := ""
	for _, err := range errors {
		if err.Employee != currentEmployee {
			if row > 2 {
				row++ // Add a blank row between employees
			}
			f.SetCellValue(errorSheetName, fmt.Sprintf("A%d", row), err.Employee)
			employeeStyle, _ := f.NewStyle(&excelize.Style{
				Font: &excelize.Font{Bold: true},
			})
			f.SetCellStyle(errorSheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), employeeStyle)
			row++
			currentEmployee = err.Employee
		}
		f.SetCellValue(errorSheetName, fmt.Sprintf("A%d", row), err.Filename)
		row++
	}

	f.SetColWidth(errorSheetName, "A", "A", 30)

	// Save the Excel file
	return f.SaveAs("orders_report.xlsx")
}

func sanitizeSheetName(name string) string {
	// Replace characters that are not allowed in Excel sheet names
	invalid := []string{":", "\\", "/", "?", "*", "[", "]"}
	for _, char := range invalid {
		name = strings.ReplaceAll(name, char, "_")
	}
	// Truncate to 31 characters (Excel's limit)
	if len(name) > 31 {
		name = name[:31]
	}
	return name
}

// Helper function to determine contrasting text color
func getContrastColor(bgColor string) string {
	// Convert hex to RGB
	rgb, _ := hex.DecodeString(bgColor[1:])
	r, g, b := float64(rgb[0]), float64(rgb[1]), float64(rgb[2])

	// Calculate luminance
	luminance := (0.299*r + 0.587*g + 0.114*b) / 255

	if luminance > 0.5 {
		return "#000000" // Black text for light backgrounds
	}
	return "#FFFFFF" // White text for dark backgrounds
}
