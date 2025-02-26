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

// ReduceMap removes all keys from the origin map that are present in any of the other maps with the same value
func ReduceMap[k comparable, v comparable](origin map[k]v, maps ...map[k]v) map[k]v {
	result := make(map[k]v)
	for key, value := range origin {
		keep := true
		for _, m := range maps {
			if val, exists := m[key]; exists {
				if val == value {
					keep = false
					break
				}
			}
		}
		if keep {
			result[key] = value
		}
	}
	return result
}
