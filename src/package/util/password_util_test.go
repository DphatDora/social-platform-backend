package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword_Success(t *testing.T) {
	password := "SecurePassword123!"
	hash, err := HashPassword(password)

	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)
	assert.Greater(t, len(hash), 50)
}

func TestHashPassword_EmptyPassword(t *testing.T) {
	password := ""
	hash, err := HashPassword(password)

	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
}

func TestHashPassword_LongPassword(t *testing.T) {
	password := "ThisIsAVeryLongPasswordWithMoreThan72Characters1234567890!@#$%^&*()_+ThisIsAVeryLongPasswordWithMoreThan72Characters1234567890!@#$%^&*()_+"
	hash, err := HashPassword(password)

	assert.Error(t, err)
	assert.Empty(t, hash)
}

func TestComparePassword_Success(t *testing.T) {
	password := "MySecurePassword123!"
	hash, err := HashPassword(password)
	assert.NoError(t, err)

	err = ComparePassword(hash, password)
	assert.NoError(t, err)
}

func TestComparePassword_WrongPassword(t *testing.T) {
	password := "CorrectPassword123!"
	wrongPassword := "WrongPassword456!"
	hash, err := HashPassword(password)
	assert.NoError(t, err)

	err = ComparePassword(hash, wrongPassword)
	assert.Error(t, err)
}

func TestComparePassword_EmptyPassword(t *testing.T) {
	password := "SomePassword123!"
	hash, err := HashPassword(password)
	assert.NoError(t, err)

	err = ComparePassword(hash, "")
	assert.Error(t, err)
}

func TestComparePassword_InvalidHash(t *testing.T) {
	invalidHash := "not-a-valid-bcrypt-hash"
	password := "SomePassword123!"

	err := ComparePassword(invalidHash, password)
	assert.Error(t, err)
}

func TestHashPassword_Uniqueness(t *testing.T) {
	password := "SamePassword123!"
	hash1, err1 := HashPassword(password)
	hash2, err2 := HashPassword(password)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEqual(t, hash1, hash2)
}
