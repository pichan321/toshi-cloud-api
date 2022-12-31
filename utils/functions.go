package utils

func ExistsWithin(arr []interface{}, element interface{}) (interface{}) {
	for _, v := range arr {
		if v == element {
			return v
		}
	}

	return nil
}