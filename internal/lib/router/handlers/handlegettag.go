package handlers

import (
	"context"
	"fmt"

	"github.com/studyguides-com/study-guides-api/internal/lib/router/formatting"
	"github.com/studyguides-com/study-guides-api/internal/store"
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
		response := formatting.NewSingleResponse(nil, fmt.Sprintf("Error retrieving tag: %v", err))
		return response.ToJSON(), nil
	}

	if tag == nil {
		response := formatting.NewSingleResponse(nil, fmt.Sprintf("Tag with ID '%s' not found.", tagID))
		return response.ToJSON(), nil
	}

	// Get the format specified by the AI
	format := formatting.GetFormatFromParams(params)

	// Format the tag details
	var data interface{}
	if format == formatting.FormatJSON {
		data = tag
	} else {
		data = formatting.TagAsFormatted(tag, format)
	}

	message := fmt.Sprintf("Found tag '%s'", tag.Name)
	response := formatting.NewSingleResponse(data, message)
	return response.ToJSON(), nil
}
