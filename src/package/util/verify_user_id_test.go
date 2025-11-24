package util

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetUserIDFromContext_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	expectedUserID := uint64(12345)
	c.Set("userID", expectedUserID)

	userID, err := GetUserIDFromContext(c)

	assert.NoError(t, err)
	assert.Equal(t, expectedUserID, userID)
}

func TestGetUserIDFromContext_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	userID, err := GetUserIDFromContext(c)

	assert.Error(t, err)
	assert.Equal(t, uint64(0), userID)
	assert.Contains(t, err.Error(), "userID not found in context")
}

func TestGetUserIDFromContext_InvalidType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set("userID", "not-a-uint64")

	userID, err := GetUserIDFromContext(c)

	assert.Error(t, err)
	assert.Equal(t, uint64(0), userID)
	assert.Contains(t, err.Error(), "invalid userID type")
}

func TestGetUserIDFromContext_IntType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set("userID", 12345)

	userID, err := GetUserIDFromContext(c)

	assert.Error(t, err)
	assert.Equal(t, uint64(0), userID)
}

func TestGetOptionalUserIDFromContext_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	expectedUserID := uint64(12345)
	c.Set("userID", expectedUserID)

	userID := GetOptionalUserIDFromContext(c)

	assert.NotNil(t, userID)
	assert.Equal(t, expectedUserID, *userID)
}

func TestGetOptionalUserIDFromContext_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	userID := GetOptionalUserIDFromContext(c)

	assert.Nil(t, userID)
}

func TestGetOptionalUserIDFromContext_InvalidType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set("userID", "not-a-uint64")

	userID := GetOptionalUserIDFromContext(c)

	assert.Nil(t, userID)
}

func TestGetOptionalUserIDFromContext_Zero(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set("userID", uint64(0))

	userID := GetOptionalUserIDFromContext(c)

	assert.NotNil(t, userID)
	assert.Equal(t, uint64(0), *userID)
}

func TestGetUserIDFromContext_MultipleSet(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set("userID", uint64(111))
	c.Set("userID", uint64(222))

	userID, err := GetUserIDFromContext(c)

	assert.NoError(t, err)
	assert.Equal(t, uint64(222), userID)
}
