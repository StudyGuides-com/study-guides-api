package handlers

import (
	"context"
	"fmt"

	"github.com/studyguides-com/study-guides-api/internal/lib/router/formatting"
	"github.com/studyguides-com/study-guides-api/internal/store"
)

func HandleListRootTags(ctx context.Context, store store.Store, params map[string]string) (string, error) {
	// Get the format specified by the AI
	format := formatting.GetFormatFromParams(params)

	// Get root tags (tags with no parent)
	tags, err := store.TagStore().ListRootTags(ctx, params)
	if err != nil {
		return fmt.Sprintf("Error retrieving root tags: %v", err), nil
	}

	if len(tags) == 0 {
		return "No root tags found.", nil
	}

	// Build response message
	response := fmt.Sprintf("Found %d root tags", len(tags))

	// Add filter description if any filters are applied
	hasTypeFilter := params["type"] != ""
	hasContextFilter := params["contextType"] != ""
	hasNameFilter := params["name"] != ""
	hasPublicFilter := params["public"] != ""

	if hasTypeFilter || hasContextFilter || hasNameFilter || hasPublicFilter {
		response += formatting.BuildFilterDescription(params, hasTypeFilter, hasContextFilter, hasNameFilter, hasPublicFilter)
	}

	// Add limit message if results are limited
	response += formatting.BuildLimitMessage(params)
	response += ":\n\n"

	// Format the tags according to the specified format
	if format == formatting.FormatJSON || format == formatting.FormatCSV {
		// For JSON and CSV, return just the formatted data
		return formatting.FormatTags(tags, format), nil
	} else {
		// For other formats, include the response message
		response += formatting.FormatTags(tags, format)
		return response, nil
	}
}
