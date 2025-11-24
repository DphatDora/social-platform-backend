package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeString_BasicLowercase(t *testing.T) {
	input := "Hello World"
	expected := "hello world"
	result := NormalizeString(input)

	assert.Equal(t, expected, result)
}

func TestNormalizeString_VietnameseAccents(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"Hà Nội", "ha noi"},
		{"Sài Gòn", "sai gon"},
		{"Đà Nẵng", "đa nang"},
		{"Phở", "pho"},
		{"Bánh mì", "banh mi"},
		{"Cà phê", "ca phe"},
		{"Tiếng Việt", "tieng viet"},
	}

	for _, tc := range testCases {
		result := NormalizeString(tc.input)
		assert.Equal(t, tc.expected, result)
	}
}

func TestNormalizeString_SpecialCharacters(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"Hello@World!", "hello@world!"},
		{"Test#123", "test#123"},
		{"user_name", "user_name"},
		{"email@example.com", "email@example.com"},
	}

	for _, tc := range testCases {
		result := NormalizeString(tc.input)
		assert.Equal(t, tc.expected, result)
	}
}

func TestNormalizeString_WhitespaceHandling(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"  Hello  ", "hello"},
		{"Hello   World", "hello   world"},
		{"\tHello\n", "hello"},
		{"   ", ""},
	}

	for _, tc := range testCases {
		result := NormalizeString(tc.input)
		assert.Equal(t, tc.expected, result)
	}
}

func TestNormalizeString_EmptyString(t *testing.T) {
	result := NormalizeString("")
	assert.Equal(t, "", result)
}

func TestNormalizeString_MixedCase(t *testing.T) {
	input := "HeLLo WoRLd"
	expected := "hello world"
	result := NormalizeString(input)

	assert.Equal(t, expected, result)
}

func TestBuildSearchPattern_SingleWord(t *testing.T) {
	keyword := "hello"
	result := BuildSearchPattern(keyword)

	assert.Len(t, result, 1)
	assert.Equal(t, "%hello%", result[0])
}

func TestBuildSearchPattern_MultipleWords(t *testing.T) {
	keyword := "hello world"
	result := BuildSearchPattern(keyword)

	assert.Len(t, result, 2)
	assert.Equal(t, "%hello%", result[0])
	assert.Equal(t, "%world%", result[1])
}

func TestBuildSearchPattern_VietnameseKeywords(t *testing.T) {
	keyword := "Hà Nội"
	result := BuildSearchPattern(keyword)

	assert.Len(t, result, 2)
	assert.Equal(t, "%ha%", result[0])
	assert.Equal(t, "%noi%", result[1])
}

func TestBuildSearchPattern_ExtraSpaces(t *testing.T) {
	keyword := "  hello   world  "
	result := BuildSearchPattern(keyword)

	assert.Len(t, result, 2)
	assert.Equal(t, "%hello%", result[0])
	assert.Equal(t, "%world%", result[1])
}

func TestBuildSearchPattern_EmptyString(t *testing.T) {
	keyword := ""
	result := BuildSearchPattern(keyword)

	assert.Len(t, result, 0)
}

func TestBuildSearchPattern_SpecialCharacters(t *testing.T) {
	keyword := "test@123"
	result := BuildSearchPattern(keyword)

	assert.Len(t, result, 1)
	assert.Equal(t, "%test@123%", result[0])
}

func TestBuildSearchPattern_ComplexVietnamese(t *testing.T) {
	keyword := "Phở Bò Tái"
	result := BuildSearchPattern(keyword)

	assert.Len(t, result, 3)
	assert.Equal(t, "%pho%", result[0])
	assert.Equal(t, "%bo%", result[1])
	assert.Equal(t, "%tai%", result[2])
}
