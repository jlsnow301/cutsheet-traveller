package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"

	fileutils "github.com/jlsnow301/cutsheet-traveller/files"
	"github.com/jlsnow301/cutsheet-traveller/input"
	"github.com/jlsnow301/cutsheet-traveller/utils"
)

func main() {
	envPath := filepath.Join(filepath.Dir(os.Args[0]), ".env")
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		utils.PrintRed("Please create a .env file in the project root.")
		os.Exit(1)
	}

	err := godotenv.Load(envPath)
	if err != nil {
		utils.PrintRed("Error loading .env file")
		os.Exit(1)
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		os.Exit(1)
	}

	employeesDir := filepath.Join(cwd, "employees")

	// Read the directory contents of the "employees" subfolder
	files, err := os.ReadDir(employeesDir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		os.Exit(1)
	}

	folders := []string{}
	for _, file := range files {
		if !file.IsDir() || !fileutils.ValidateFolder(file) {
			continue
		}

		folders = append(folders, file.Name())
	}

	utils.PrintHeader("Order Mileage")

	maxNumber := len(folders)
	if maxNumber == 0 {
		utils.PrintRed("No folders found in the 'employees' folder.")
		os.Exit(1)
	}

	fmt.Println("This tool searches for folders in the 'employees' folder.")
	fmt.Println("It is stupid, it assumes all cut sheets are valid.")
	fmt.Println("It will also skip cut sheets that it does not understand.")
	fmt.Println()
	utils.PrintStars()
	fmt.Println()
	fmt.Println("Please select the employee folder to get their mileage.")
	utils.PrintYellow("The valid employee folders are:")

	for i, folder := range folders {
		fmt.Printf("%d. %s\n", i+1, folder)
	}

	if maxNumber > 1 {
		maxNumber += 1
		fmt.Printf("%d. All", maxNumber)
	}
	fmt.Println()
	fmt.Println()

	userNumber := input.PromptUserForNumber(maxNumber)

	// Create an array with just the employee folder, or all folders if the user selected "All"
	foldersToSearch := []string{}
	if userNumber == maxNumber {
		foldersToSearch = folders
	} else {
		foldersToSearch = append(foldersToSearch, folders[userNumber-1])
	}

	employeeOrders, errors := fileutils.CollectOrdersAndErrors(foldersToSearch, employeesDir)

	err = fileutils.CreateExcelFile(employeeOrders, errors)
	if err != nil {
		utils.PrintRed(fmt.Sprintf("Error creating Excel file: %v", err))
		os.Exit(1)
	}

	fmt.Println("\nExcel file created successfully.")

}
