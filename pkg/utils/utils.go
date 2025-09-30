package utils

import "strings"

func RemoveDuplicates[T comparable](slice []T) []T {
	seen := make(map[T]struct{})
	result := []T{}

	for _, item := range slice {
		if _, ok := seen[item]; !ok {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}

func RemoveURLDuplicates(slice []string) []string {
	formatted := make([]string, len(slice))

	for i, str := range slice {
		formatted[i] = strings.TrimSuffix(str, "/")
	}

	return RemoveDuplicates(formatted)
}
