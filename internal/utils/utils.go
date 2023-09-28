package utils

import (
	"golang.org/x/exp/slices"
)

// Intersect allowed groups and token groups
func InterSectInterface(slice1 []interface{}, slice2 []string) []string {
	var intersect []string
	for _, part1 := range slice1 {
		for _, part2 := range slice2 {
			if part1 == part2 {
				intersect = append(intersect, part1.(string))
			}
		}
	}
	return intersect
}

// Check String exists in Array
func CheckStringInArray(dataarray []string, searchstring string) bool {
	if slices.Contains(dataarray, searchstring) {
		return true
	} else {
		return false
	}
}
