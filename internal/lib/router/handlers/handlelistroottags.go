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
		errorMessage := fmt.Sprintf("Error retrieving root tags: %v", err)
		response := formatting.NewSingleResponse(nil, errorMessage)
		return response.ToJSON(), nil
	}

	if len(tags) == 0 {
		response := formatting.NewListResponse([]interface{}{}, "No root tags found.", nil, nil)
		return response.ToJSON(), nil
	}

	// Build filters map if any filters are present
	var filters map[string]string
	hasTypeFilter := params["type"] != ""
	hasContextFilter := params["contextType"] != ""
	hasNameFilter := params["name"] != ""
	hasPublicFilter := params["public"] != ""

	if hasTypeFilter || hasContextFilter || hasNameFilter || hasPublicFilter {
		filters = make(map[string]string)
		if hasTypeFilter {
			filters["type"] = params["type"]
		}
		if hasContextFilter {
			filters["context_type"] = params["contextType"]
		}
		if hasNameFilter {
			filters["name"] = params["name"]
		}
		if hasPublicFilter {
			filters["public"] = params["public"]
		}
	}

	// Format the data using the formatter
	formatter := formatting.NewTagFormatter(tags)
	data := formatter.Format(format)

	// Build message
	message := fmt.Sprintf("Found %d root tags", len(tags))
	if hasTypeFilter || hasContextFilter || hasNameFilter || hasPublicFilter {
		message += formatting.BuildFilterDescription(params, hasTypeFilter, hasContextFilter, hasNameFilter, hasPublicFilter)
	}
	message += formatting.BuildLimitMessage(params)

	response := formatting.NewListResponse(data, message, filters, nil)
	return response.ToJSON(), nil
}
