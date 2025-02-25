package utils

func MapKeysToList[k comparable, v any](m map[k]v) []k {
	keys := make([]k, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func MergeMaps[k comparable, v any](maps ...map[k]v) map[k]v {
	result := map[k]v{}
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}
