package main

import (
	"context"
	"log"
	"strings"
	"time"
)

// AuthEvent is the payload of a Firestore Auth event.
type AuthEvent struct {
	Email    string `json:"email"`
	Metadata struct {
		CreatedAt time.Time `json:"createdAt"`
	} `json:"metadata"`
	UID string `json:"uid"`
}

// HelloAuth is triggered by Firestore Auth events.
func HelloAuth(ctx context.Context, e AuthEvent) error {
	log.Printf("Function triggered by creation or deletion of user: %q", e.UID)
	log.Printf("Created at: %v", e.Metadata.CreatedAt)
	if e.Email != "" {
		log.Printf("Email: %q", e.Email)
	}

	var claims = make(map[string]interface{})

	customClaimsAdmin := map[string]interface{}{
		"x-hasura-default-role": "admin",
		"x-hasura-allowed-roles": []string{"user", "admin"},
		"x-hasura-user-id": e.UID,
	}

	customClaimsUser := map[string]interface{}{
		"x-hasura-default-role": "admin",
		"x-hasura-allowed-roles": []string{"user", "admin"},
		"x-hasura-user-id": e.UID,
	}

	if strings.Contains(e.UID, "@hasura.io") {
		claims["https://hasura.io/jwt/claims"] = customClaimsAdmin
	} else {
		claims["https://hasura.io/jwt/claims"] = customClaimsUser
	}

	// Set admin privilege on the user corresponding to uid.
	err := client.SetCustomUserClaims(ctx, e.UID, claims)
	if err != nil {
		log.Fatalf("error setting custom claims %v\n", err)
	}

	return nil
}