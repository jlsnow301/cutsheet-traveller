package fileutils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jlsnow301/cutsheet-traveller/header"
	timeutils "github.com/jlsnow301/cutsheet-traveller/time"
	"github.com/jlsnow301/cutsheet-traveller/travel"
	"github.com/jlsnow301/cutsheet-traveller/utils"
)

type errorInfo struct {
	Employee string
	Filename string
}

type orderInfo struct {
	OrderID     string
	Mileage     float64
	Date        string
	Origin      string
	Destination string
}

// Collect all orders and errors
func CollectOrdersAndErrors(foldersToSearch []string, employeesDir string) (map[string][]orderInfo, []errorInfo) {
	employeeOrders := make(map[string][]orderInfo)
	var errors []errorInfo

	for _, searchFolder := range foldersToSearch {
		folderPath := filepath.Join(employeesDir, searchFolder)

		entries, err := os.ReadDir(folderPath)
		if err != nil {
			continue
		}

		cutsheets, err := getCutsheets(entries)
		if err != nil {
			continue
		}

		for _, cutsheet := range cutsheets {
			pdfPath := filepath.Join(folderPath, cutsheet.Name())

			orderInfo, err := getOrderInfo(pdfPath)
			if err != nil {
				errors = append(errors, errorInfo{
					Employee: searchFolder,
					Filename: cutsheet.Name(),
				})
				continue
			}

			employeeOrders[searchFolder] = append(employeeOrders[searchFolder], orderInfo)
		}
	}

	return employeeOrders, errors
}

// Get the order info from a PDF file
func getOrderInfo(pdfPath string) (orderInfo, error) {
	pdfText, err := utils.ExtractTextFromPDF(pdfPath)
	if err != nil {
		utils.PrintRed(fmt.Sprintf("Error extracting text from PDF: %v", err))
		return orderInfo{}, err
	}

	headerText, _ := utils.SplitTexts(pdfText)
	headerInfo := header.ParseHeaderInfo(headerText)
	if headerInfo.Destination == "" {
		utils.PrintRed("Unable to determine destination address.")
		return orderInfo{}, err
	}

	origin := headerInfo.Origin
	if origin == "" {
		utils.PrintRed("No origin specified.")
		return orderInfo{}, err
	}

	originAddress := os.Getenv(strings.ToUpper(origin) + "_ADDRESS")
	if originAddress == "" {
		utils.PrintRed(fmt.Sprintf("Unknown origin: %s", headerInfo.Origin))
		return orderInfo{}, err
	}

	eventTime, err := timeutils.GetEventTime(headerInfo.EventDate, headerInfo.EventTime)
	if err != nil {
		utils.PrintRed("Invalid event time. Please use HH:MM AM/PM.")
		return orderInfo{}, err
	}

	distance, err := travel.GetBaseTravelDistance(originAddress, headerInfo.Destination, eventTime)
	if err != nil {
		utils.PrintRed(fmt.Sprintf("Unable to calculate travel time: %v", err))
		return orderInfo{}, err
	}

	orderInfo := orderInfo{
		OrderID:     headerInfo.OrderID,
		Mileage:     distance,
		Date:        headerInfo.EventDate.Format("2006-01-02"),
		Origin:      headerInfo.Origin,
		Destination: headerInfo.Destination,
	}

	return orderInfo, nil
}
