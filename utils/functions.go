package utils

func ExistsWithin(arr []interface{}, element interface{}) bool {
	for _, v := range arr {
		if v == element {
			return true
		}
	}

	return false
}