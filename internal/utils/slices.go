package utils

import "slices"

func SortStringSlice(slice []string) []string {
	slices.Sort(slice)
	return slice
}
