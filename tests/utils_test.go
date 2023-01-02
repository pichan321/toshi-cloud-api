package tests

import (
	"file-api/utils"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFixEscape(t *testing.T) {
	original := `eastman ' '' '''.jpg`
	assert.Equal(t, `eastman '' '''' ''''''.jpg`, utils.FixEscape(original))
}

func TestHashPassword(t *testing.T) {
	unhashed := `123`
	assert.Equal(t, "a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3", utils.HashPassword(unhashed))
}	

func TestGenerateUuid(t *testing.T) {
	uuid := utils.GenerateUuid()
	assert.Equal(t, 5, len(strings.Split(uuid, "-")))
}

func TestGenerateToken(t *testing.T) {
	token := utils.GenerateToken()
	assert.Equal(t, 32, len(token))
} 

