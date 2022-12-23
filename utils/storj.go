package utils

import "fmt"

func StorjFilename(uuid string, filename string, delimiter string) string {
	return fmt.Sprintf("%s%s%s", uuid, delimiter, filename)
}