package test

import (
	"TikTok/apps/app/api/utils"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestJwtAuthCreateToken(t *testing.T) {
	jwtAuth := utils.JwtAuth{
		AccessSecret: []byte("secret_key"),
		AccessExpire: 3600, // 1 hour in seconds
	}

	tokenID := int64(123)

	// Call the function being tested
	tokenString, err := jwtAuth.CreateToken(tokenID)

	// Assertions
	assert.NoError(t, err, "Error creating token")
	assert.NotEmpty(t, tokenString, "Token string should not be empty")
}

func TestJwtAuthParseToken(t *testing.T) {
	jwtAuth := utils.JwtAuth{
		AccessSecret: []byte("secret_key"),
		AccessExpire: 3600, // 1 hour in seconds
	}

	tokenID := int64(123)

	// Create a token for testing
	tokenString, _ := jwtAuth.CreateToken(tokenID)

	// Call the function being tested
	parsedTokenID, err := jwtAuth.ParseToken(tokenString)

	// Assertions
	assert.NoError(t, err, "Error parsing token")
	assert.Equal(t, tokenID, parsedTokenID, "Parsed token ID should be equal to the original token ID")
}

func TestJwtAuthExpiredToken(t *testing.T) {
	jwtAuth := utils.JwtAuth{
		AccessSecret: []byte("secret_key"),
		AccessExpire: 1, // 1 second (expired token)
	}

	tokenID := int64(456)

	// Create a token that will expire in 1 second
	tokenString, _ := jwtAuth.CreateToken(tokenID)

	// Wait for the token to expire
	time.Sleep(2 * time.Second)

	// Call the function being tested
	parsedTokenID, err := jwtAuth.ParseToken(tokenString)

	// Assertions
	assert.Error(t, err, "Token should have expired")
	assert.Equal(t, int64(-1), parsedTokenID, "Parsed token ID should be -1 (invalid)")
}
