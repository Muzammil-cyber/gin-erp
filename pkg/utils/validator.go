package utils

import (
	"regexp"
	"strings"
)

// ValidatePakistaniPhone validates Pakistani phone number format (+923xxxxxxxxx)
func ValidatePakistaniPhone(phone string) bool {
	// Remove any spaces or dashes
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")

	// Pakistani phone number regex: +923xxxxxxxxx (total 13 characters)
	// +92 (country code) + 3 (mobile prefix) + 9 digits
	regex := regexp.MustCompile(`^\+923[0-9]{9}$`)
	return regex.MatchString(phone)
}

// NormalizePakistaniPhone normalizes phone number to standard format
func NormalizePakistaniPhone(phone string) string {
	// Remove spaces and dashes
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")

	// If starts with 0, replace with +92
	if strings.HasPrefix(phone, "03") {
		phone = "+92" + phone[1:]
	}

	// If starts with 92, add +
	if strings.HasPrefix(phone, "92") && !strings.HasPrefix(phone, "+92") {
		phone = "+" + phone
	}

	return phone
}
