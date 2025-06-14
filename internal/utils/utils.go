package utils

import (
	"encoding/json"
	"fmt"
	"sort"
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

// HasID interface for structs with an ID field
type HasID interface {
	GetID() int
}

// SortByID sorts any slice of HasID by ID (ascending or descending)
func SortByID[T HasID](items []T, ascending bool) {
	sort.Slice(items, func(i, j int) bool {
		if ascending {
			return items[i].GetID() < items[j].GetID()
		}
		return items[i].GetID() > items[j].GetID()
	})
}
