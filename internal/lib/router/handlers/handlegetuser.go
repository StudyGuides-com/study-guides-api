package handlers

import (
	"context"
	"fmt"

	"github.com/studyguides-com/study-guides-api/internal/lib/router/formatting"
	"github.com/studyguides-com/study-guides-api/internal/store"
)

func HandleGetUser(ctx context.Context, store store.Store, params map[string]string) (string, error) {
	userEmail, ok := params["userEmail"]
	if !ok || userEmail == "" {
		return "Please provide a user email to retrieve.", nil
	}
	
	// Get the user by ID
	user, err := store.UserStore().UserByEmail(ctx, userEmail)
	if err != nil {
		return fmt.Sprintf("Error retrieving user: %v", err), nil
	}
	
	if user == nil {
		return fmt.Sprintf("User with email '%s' not found.", userEmail), nil
	}
	
	// Get the format specified by the AI
	format := formatting.GetFormatFromParams(params)
	
	// Format the user details using the formatting function
	return formatting.UserAsFormatted(user, format), nil
} 