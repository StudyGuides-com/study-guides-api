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
		return "Please provide a tag ID to retrieve.", nil
	}

	// Get the tag by ID
	tag, err := store.TagStore().GetTagByID(ctx, tagID)
	if err != nil {
		return fmt.Sprintf("Error retrieving tag: %v", err), nil
	}

	if tag == nil {
		return fmt.Sprintf("Tag with ID '%s' not found.", tagID), nil
	}

	// Get the format specified by the AI
	format := formatting.GetFormatFromParams(params)

	// Format the tag details using the formatting function
	return formatting.TagAsFormatted(tag, format), nil
}
