package utils

import "net/mail"

func ExistsWithin(arr []interface{}, element interface{}) interface{} {
	for _, v := range arr {
		if v == element {
			return v
		}
	}

	return nil
}

func ValidateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return err
	}
	return nil
}