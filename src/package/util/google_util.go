package util

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/api/idtoken"
)

type GoogleUserInfo struct {
	GoogleID string
	Email    string
	Name     string
	Picture  string
	Verified bool
}

func VerifyGoogleIDToken(idToken string, clientID string) (*GoogleUserInfo, error) {
	ctx := context.Background()

	// Validate the ID token
	payload, err := idtoken.Validate(ctx, idToken, clientID)
	if err != nil {
		log.Printf("[Err] Error validating ID token: %v", err)
		return nil, fmt.Errorf("invalid ID token: %w", err)
	}

	claims := payload.Claims

	// Get user ID (subject)
	googleID, ok := claims["sub"].(string)
	if !ok || googleID == "" {
		log.Printf("[Err] Missing or invalid 'sub' claim in token")
		return nil, fmt.Errorf("missing user ID in token")
	}

	// Get email
	email, ok := claims["email"].(string)
	if !ok || email == "" {
		log.Printf("[Err] Missing or invalid 'email' claim in token")
		return nil, fmt.Errorf("missing email in token")
	}

	// Verify email is verified by Google
	emailVerified, ok := claims["email_verified"].(bool)
	if !ok || !emailVerified {
		log.Printf("[Err] Email not verified by Google for user: %s", email)
		return nil, fmt.Errorf("email not verified by Google")
	}

	name, ok := claims["name"].(string)
	if !ok || name == "" {
		name = email
	}

	picture, _ := claims["picture"].(string)

	userInfo := &GoogleUserInfo{
		GoogleID: googleID,
		Email:    email,
		Name:     name,
		Picture:  picture,
		Verified: emailVerified,
	}

	log.Printf("[Info] Successfully verified Google ID token for user: %s (GoogleID: %s)", email, googleID)
	return userInfo, nil
}
