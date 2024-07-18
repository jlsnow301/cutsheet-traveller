package header

import (
	"regexp"
	"strings"
	"time"
)

type HeaderInfo struct {
	OrderID     string
	Origin      string
	Destination string
	Size        string
	EventTime   string
	SuiteInfo   string
	EventDate   time.Time
}

func hasDatePrefix(line string) bool {
	days := []string{"monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"}
	lowerLine := strings.ToLower(line)

	for _, day := range days {
		if strings.HasPrefix(lowerLine, day) {
			return true
		}
	}

	return false

}

// Catches order IDs like "S21271"
// Line starts with S and is followed by at least 5 digits
func hasOrderID(line string) bool {
	re := regexp.MustCompile(`^S\d{5}$`)
	return re.MatchString(line)
}

// Split after the first colon
func splitAfterColon(line string) string {
	split := strings.SplitN(line, ":", 2)
	if len(split) > 1 {
		return strings.TrimSpace(split[1])
	}
	return ""
}

func getSuiteInfo(line string) string {
	re := regexp.MustCompile(`(?i)(suite\s*\S+|#\d{3})`)
	match := re.FindString(line)
	if match == "" {
		return ""
	}

	return strings.TrimSpace(line)
}

func normalizeAddress(address string) string {
	// Remove any occurrence of "Headcount" and everything after it
	headcountIndex := strings.Index(strings.ToLower(address), "headcount")
	if headcountIndex != -1 {
		address = address[:headcountIndex]
	}

	// Trim any trailing commas and spaces
	address = strings.TrimRight(address, ", ")

	// Add a space before any capitalized letter unless one already exists
	re := regexp.MustCompile(`([a-z])([A-Z])`)
	address = re.ReplaceAllString(address, "$1 $2")

	// Ensure there's a space before the ZIP code if it exists
	zipRe := regexp.MustCompile(`(\D)(\d{5})$`)
	address = zipRe.ReplaceAllString(address, "$1 $2")

	// Check if the address contains a ZIP code
	hasZip := regexp.MustCompile(`\d{5}$`).MatchString(address)

	// Check if the address contains "Seattle" (case-insensitive)
	hasSeattle := strings.Contains(strings.ToLower(address), "seattle")

	// If there's no ZIP code and no "Seattle", add "Seattle"
	if !hasZip && !hasSeattle {
		address += " Seattle"
	}

	return address
}

func ParseHeaderInfo(content []string) HeaderInfo {
	info := HeaderInfo{}
	addressParts := []string{}

	matchers := map[string]func(string){
		"Fremont":     func(s string) { info.Origin = s },
		"Eastlake":    func(s string) { info.Origin = s },
		"Start Time:": func(s string) { info.EventTime = splitAfterColon(s) },
		"Site Address:": func(s string) {
			siteAddress := splitAfterColon(s)
			addressParts = append(addressParts, siteAddress)
			if suite := getSuiteInfo(siteAddress); suite != "" {
				info.SuiteInfo = strings.TrimSpace(suite)
			}
		},
		"Site Name:": func(s string) {
			siteName := splitAfterColon(s)
			if suite := getSuiteInfo(siteName); suite != "" {
				info.SuiteInfo = strings.TrimSpace(suite)
			}
		},
		"Headcount:": func(s string) {
			info.Size = splitAfterColon(s)
			if len(addressParts) > 0 {
				info.Destination = normalizeAddress(strings.Join(addressParts, ", "))
			}
			addressParts = nil // Clear address parts after setting destination
		},
	}

	for _, line := range content {
		line = strings.TrimSpace(line)

		if info.EventDate.IsZero() && hasDatePrefix(line) {
			// Check for date in format "Day, MM/DD/YYYY"
			if date, err := time.Parse("Monday, 1/2/2006", line); err == nil {
				info.EventDate = date
				continue
			}
		}

		if info.OrderID == "" && hasOrderID(line) {
			info.OrderID = line
		}

		matched := false
		for prefix, handler := range matchers {
			if strings.HasPrefix(line, prefix) {
				handler(line)
				matched = true
				break
			}
		}

		if !matched && info.Destination == "" && len(addressParts) > 0 && line != "" {
			addressParts = append(addressParts, line)
		}
	}

	// In case the address collection wasn't terminated by a Headcount line
	if info.Destination == "" && len(addressParts) > 0 {
		info.Destination = normalizeAddress(strings.Join(addressParts, ", "))
	}

	return info
}
