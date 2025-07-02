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
		response := formatting.NewSingleResponse(nil, "Please provide a user email to retrieve.")
		return response.ToJSON(), nil
	}

	// Get the user by email
	user, err := store.UserStore().UserByEmail(ctx, userEmail)
	if err != nil {
		response := formatting.NewSingleResponse(nil, fmt.Sprintf("Error retrieving user: %v", err))
		return response.ToJSON(), nil
	}

	if user == nil {
		response := formatting.NewSingleResponse(nil, fmt.Sprintf("User with email '%s' not found.", userEmail))
		return response.ToJSON(), nil
	}

	// Get the format specified by the AI
	format := formatting.GetFormatFromParams(params)

	// Format the user details
	var data interface{}
	var contentType string
	
	switch format {
	case formatting.FormatJSON:
		data = user
		contentType = "application/json"
	case formatting.FormatCSV:
		data = formatting.UserAsFormatted(user, format)
		contentType = "text/csv"
	default:
		data = formatting.UserAsFormatted(user, format)
		contentType = "text/plain"
	}

	message := fmt.Sprintf("Found user '%s'", user.GetEmail())
	
	// Create response with correct content type
	response := &formatting.APIResponse{
		Type:        formatting.ResponseTypeSingle,
		Data:        data,
		Message:     message,
		ContentType: contentType,
	}
	
	return response.ToJSON(), nil
}
