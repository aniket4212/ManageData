package utils

import "regexp"

func IsValidEmail(email string) bool {
	// Basic regex: some text + "@" + some text + "." + some text
	re := regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
	return re.MatchString(email)
}

func IsValidPhone(phone string) bool {
	// Allow +, numbers, spaces, dashes, parentheses
	re := regexp.MustCompile(`^\+?[0-9\s\-\(\)]+$`)
	return re.MatchString(phone)
}
