package utils

import "strings"

func ToLowerNoSpaces(input string) string {
	// Remove all spaces
	noSpaces := strings.ReplaceAll(input, " ", "")
	// Convert to lowercase
	return strings.ToLower(noSpaces)
}

func AreStringsEqualIgnoreCase(a, b string) bool {
	return strings.EqualFold(a, b)
}

func StringInSlice(target string, list []string) bool {
	for _, item := range list {
		if item == target {
			return true
		}
	}
	return false
}
