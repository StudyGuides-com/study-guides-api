package handlers

import (
	"context"
	"fmt"

	"github.com/studyguides-com/study-guides-api/internal/lib/router/formatting"
	"github.com/studyguides-com/study-guides-api/internal/store"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func HandleGetTag(ctx context.Context, store store.Store, params map[string]string) (string, error) {
	tagID, ok := params["tagId"]
	if !ok || tagID == "" {
		response := formatting.NewSingleResponse(nil, "Please provide a tag ID to retrieve.")
		return response.ToJSON(), nil
	}

	// Get the tag by ID
	tag, err := store.TagStore().GetTagByID(ctx, tagID)
	if err != nil {
		// Check if it's a NotFound error and provide a more user-friendly message
		if status.Code(err) == codes.NotFound {
			response := formatting.NewSingleResponse(nil, fmt.Sprintf("Tag with ID '%s' not found.", tagID))
			return response.ToJSON(), nil
		}
		response := formatting.NewSingleResponse(nil, fmt.Sprintf("Error retrieving tag: %v", err))
		return response.ToJSON(), nil
	}

	// Get the format specified by the AI
	format := formatting.GetFormatFromParams(params)

	// Format the tag details
	var data interface{}
	var contentType string

	switch format {
	case formatting.FormatJSON:
		data = formatting.TagAsJSONObject(tag)
		contentType = "application/json"
	case formatting.FormatCSV:
		data = formatting.TagAsFormatted(tag, format)
		contentType = "text/csv"
	default:
		data = formatting.TagAsFormatted(tag, format)
		contentType = "text/plain"
	}

	message := fmt.Sprintf("Found tag '%s'", tag.Name)

	// Create response with correct content type
	response := &formatting.APIResponse{
		Type:        formatting.ResponseTypeSingle,
		Data:        data,
		Message:     message,
		ContentType: contentType,
	}

	return response.ToJSON(), nil
}
