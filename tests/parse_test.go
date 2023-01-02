package tests

import (
	"file-api/handlers"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseText(t *testing.T) {
	file, err := os.Open("../test-files/text.txt")
	assert.NoError(t, err)
	defer file.Close()

	content, err := io.ReadAll(file)
	assert.Equal(t, ".txt", handlers.ProcessParsedText(string(content)))
}

func TestParseCSV(t *testing.T) {
	file, err := os.Open("../test-files/csv.txt")
	assert.NoError(t, err)
	defer file.Close()

	content, err := io.ReadAll(file)
	assert.Equal(t, ".csv", handlers.ProcessParsedText(string(content)))
}

func TestParseJson(t *testing.T) {
	file, err := os.Open("../test-files/json.txt")
	assert.NoError(t, err)
	defer file.Close()

	content, err := io.ReadAll(file)
	assert.Equal(t, ".json", handlers.ProcessParsedText(string(content)))
}