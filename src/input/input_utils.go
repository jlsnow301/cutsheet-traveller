package input

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jlsnow301/cutsheet-traveller/utils"
)

func PromptUserForNumber(maxNumber int) int {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("Enter a number between 1 and %d: ", maxNumber)

		scanner.Scan()
		userInput := scanner.Text()

		enteredNumber, err := strconv.Atoi(userInput)
		if err != nil || enteredNumber < 1 || enteredNumber > maxNumber {
			fmt.Println("Invalid input. Please enter a number within the valid range.")
			continue
		}

		return enteredNumber
	}
}

func PromptForEventTime() string {
	for {
		fmt.Print("Please enter the event time (HH:MM AM/PM): ")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			eventTime := scanner.Text()

			// Convert AM/PM to uppercase for parsing
			upperEventTime := strings.ToUpper(eventTime)

			// Try parsing with both "03:04 PM" and "3:04 PM" formats
			_, err1 := time.Parse("03:04 PM", upperEventTime)
			_, err2 := time.Parse("3:04 PM", upperEventTime)

			if err1 == nil || err2 == nil {
				return eventTime // Return the original input
			}
			utils.PrintRed("Invalid format. Please use HH:MM AM/PM.")
		}
	}
}
