package helper

func RemoveDuplicateUint64(slice []uint64) []uint64 {
	seen := make(map[uint64]bool)
	result := make([]uint64, 0, len(slice))

	for _, value := range slice {
		if !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}

	return result
}
