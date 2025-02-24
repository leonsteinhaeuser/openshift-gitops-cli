package utils

func MapKeysToList[k comparable, v any](m map[k]v) []k {
	keys := make([]k, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
