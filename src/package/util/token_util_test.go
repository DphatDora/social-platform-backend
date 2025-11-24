package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateToken_Success(t *testing.T) {
	length := 32
	token, err := GenerateToken(length)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Equal(t, length*2, len(token))
}

func TestGenerateToken_DifferentLengths(t *testing.T) {
	testCases := []struct {
		length   int
		expected int
	}{
		{16, 32},
		{32, 64},
		{64, 128},
	}

	for _, tc := range testCases {
		token, err := GenerateToken(tc.length)
		assert.NoError(t, err)
		assert.Equal(t, tc.expected, len(token))
	}
}

func TestGenerateToken_Uniqueness(t *testing.T) {
	token1, err1 := GenerateToken(32)
	token2, err2 := GenerateToken(32)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEqual(t, token1, token2)
}

func TestGenerateToken_ZeroLength(t *testing.T) {
	token, err := GenerateToken(0)

	assert.NoError(t, err)
	assert.Equal(t, "", token)
}

func TestGenerateJWT_Success(t *testing.T) {
	userID := uint64(12345)
	expirationMinutes := 60
	secret := "test-secret-key-12345"
	now := time.Now()

	token, err := GenerateJWT(userID, expirationMinutes, secret, &now)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestGenerateJWT_WithoutPasswordChangedAt(t *testing.T) {
	userID := uint64(12345)
	expirationMinutes := 60
	secret := "test-secret-key-12345"

	token, err := GenerateJWT(userID, expirationMinutes, secret, nil)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestGenerateJWT_DifferentUserIDs(t *testing.T) {
	expirationMinutes := 60
	secret := "test-secret-key-12345"

	token1, err1 := GenerateJWT(1, expirationMinutes, secret, nil)
	token2, err2 := GenerateJWT(2, expirationMinutes, secret, nil)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEqual(t, token1, token2)
}

func TestVerifyJWT_Success(t *testing.T) {
	userID := uint64(12345)
	expirationMinutes := 60
	secret := "test-secret-key-12345"
	now := time.Now()

	token, err := GenerateJWT(userID, expirationMinutes, secret, &now)
	assert.NoError(t, err)

	claims, err := VerifyJWT(token, secret)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
	assert.NotNil(t, claims.PasswordChangedAt)
	assert.Equal(t, now.Unix(), *claims.PasswordChangedAt)
}

func TestVerifyJWT_WrongSecret(t *testing.T) {
	userID := uint64(12345)
	expirationMinutes := 60
	secret := "test-secret-key-12345"
	wrongSecret := "wrong-secret-key"

	token, err := GenerateJWT(userID, expirationMinutes, secret, nil)
	assert.NoError(t, err)

	claims, err := VerifyJWT(token, wrongSecret)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestVerifyJWT_InvalidToken(t *testing.T) {
	secret := "test-secret-key-12345"
	invalidToken := "invalid.jwt.token"

	claims, err := VerifyJWT(invalidToken, secret)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestVerifyJWT_ExpiredToken(t *testing.T) {
	userID := uint64(12345)
	expirationMinutes := -1
	secret := "test-secret-key-12345"

	token, err := GenerateJWT(userID, expirationMinutes, secret, nil)
	assert.NoError(t, err)

	time.Sleep(2 * time.Second)

	claims, err := VerifyJWT(token, secret)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestVerifyJWT_EmptyToken(t *testing.T) {
	secret := "test-secret-key-12345"

	claims, err := VerifyJWT("", secret)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestJWTClaims_PasswordChangedAt(t *testing.T) {
	userID := uint64(12345)
	expirationMinutes := 60
	secret := "test-secret-key-12345"
	passwordChangedAt := time.Now().Add(-24 * time.Hour)

	token, err := GenerateJWT(userID, expirationMinutes, secret, &passwordChangedAt)
	assert.NoError(t, err)

	claims, err := VerifyJWT(token, secret)
	assert.NoError(t, err)
	assert.NotNil(t, claims.PasswordChangedAt)
	assert.Equal(t, passwordChangedAt.Unix(), *claims.PasswordChangedAt)
}

func TestJWTClaims_NoPasswordChangedAt(t *testing.T) {
	userID := uint64(12345)
	expirationMinutes := 60
	secret := "test-secret-key-12345"

	token, err := GenerateJWT(userID, expirationMinutes, secret, nil)
	assert.NoError(t, err)

	claims, err := VerifyJWT(token, secret)
	assert.NoError(t, err)
	assert.Nil(t, claims.PasswordChangedAt)
}
