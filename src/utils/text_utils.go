package utils

import (
	"fmt"
	"strings"

	"github.com/dslipak/pdf"
	"github.com/fatih/color"
)

// Extracts all lines from the PDF file
func getAllLines(file *pdf.Reader) ([]string, error) {
	var allLines []string

	for pageIndex := 0; pageIndex < file.NumPage(); pageIndex++ {
		page := file.Page(pageIndex + 1)
		rows, err := page.GetTextByRow()
		if err != nil {
			return nil, err
		}

		for _, row := range rows {
			for _, word := range row.Content {
				words := strings.TrimSpace(word.S)
				if words != "" {
					allLines = append(allLines, words)
				}

			}
		}
	}

	return allLines, nil

}

// Prettifies the text by removing newlines and combining lines
func getProcessedLines(allLines []string) ([]string, error) {
	var processedLines []string

	for i := 0; i < len(allLines); i++ {
		line := allLines[i]

		if strings.TrimSpace(line) != "-" {
			processedLines = append(processedLines, line)
			continue
		}

		if len(processedLines) > 0 {
			// Append "-" to the last item in processedLines
			processedLines[len(processedLines)-1] += " - "
		}
		// Skip the next line and add the line after that, ensuring we don't go out of bounds
		if i+1 < len(allLines) {
			processedLines[len(processedLines)-1] += allLines[i+1]
		}
		i += 1 // Skip the next line and the line after the next in the iteration

	}

	return processedLines, nil
}

func ExtractTextFromPDF(pdfPath string) ([]string, error) {
	file, err := pdf.Open(pdfPath)
	if err != nil {
		return nil, err
	}

	allLines, err := getAllLines(file)
	if err != nil {
		return nil, err
	}

	processedLines, err := getProcessedLines(allLines)
	if err != nil {
		return nil, err
	}

	return processedLines, nil
}

// splitTexts splits the text into header and food service items.
func SplitTexts(lines []string) (headerText, remainingText []string) {
	splitIndex := -1

	for i, line := range lines {
		if strings.TrimSpace(line) == "Food/Service Item" {
			splitIndex = i
			break
		}
	}

	if splitIndex != -1 {
		headerText = lines[:splitIndex+1]
		remainingText = lines[splitIndex+1:]
	} else {
		headerText = lines
		remainingText = []string{}
	}

	return headerText, remainingText
}

// printStars prints a line of stars in yellow.
func PrintStars() {
	PrintYellow("**********************************************")
}

// printRed prints text in red.
func PrintRed(text string) {
	color.Red(text)
}

// printGreen prints text in green.
func PrintGreen(text string) {
	color.Green(text)
}

// printYellow prints text in yellow.
func PrintYellow(text string) {
	color.Yellow(text)
}

// printCyan prints text in cyan.
func PrintCyan(text string) {
	color.Cyan(text)
}

// printStats pretty prints statistics as yellow text until the colon.
func PrintStats(text string) {
	yellow := color.New(color.FgYellow).SprintFunc()

	splitText := strings.Split(text, ":")
	fmt.Printf("%s%s\n", yellow(splitText[0]+":"), splitText[1])
}

// printHeader prints a yellow star and the rest in white.
func PrintHeader(text string) {
	yellow := color.New(color.FgYellow).SprintFunc()

	PrintStars()
	PrintYellow("*")
	fmt.Printf("%s %s\n", yellow("*"), text)
	PrintYellow("*")
	fmt.Println()
}
