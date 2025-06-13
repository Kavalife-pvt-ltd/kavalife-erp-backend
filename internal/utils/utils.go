package utils

import (
	"encoding/json"
	"fmt"
	"strings"
)

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

func PrettyPrint(x any) {
	b, err := json.MarshalIndent(x, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Print(string(b))
}
