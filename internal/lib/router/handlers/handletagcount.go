package handlers

import (
	"context"

	"github.com/studyguides-com/study-guides-api/internal/lib/router/formatting"
	"github.com/studyguides-com/study-guides-api/internal/store"
)

func HandleTagCount(ctx context.Context, store store.Store, params map[string]string) (string, error) {
	// Check if we have any filters
	hasTypeFilter := false
	hasContextFilter := false
	hasNameFilter := false
	hasPublicFilter := false

	if typeStr, ok := params["type"]; ok && typeStr != "" {
		hasTypeFilter = true
	}
	if contextType, ok := params["contextType"]; ok && contextType != "" {
		hasContextFilter = true
	}
	if name, ok := params["name"]; ok && name != "" {
		hasNameFilter = true
	}
	if publicStr, ok := params["public"]; ok && publicStr != "" {
		hasPublicFilter = true
	}

	var count int
	var err error

	if hasTypeFilter || hasContextFilter || hasNameFilter || hasPublicFilter {
		// If any filter is specified, use CountTags with params
		count, err = store.TagStore().CountTags(ctx, params)
	} else {
		// Get total count without any filters
		count, err = store.TagStore().CountTags(ctx, params)
	}

	if err != nil {
		return "", err
	}

	// Build filters map if any filters are present
	var filters map[string]string
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

	// Build the message using the helper function
	message := formatting.BuildCountMessageInt(count, "tags", params, hasTypeFilter, hasContextFilter, hasNameFilter, hasPublicFilter)

	// Create response using the universal wrapper
	response := formatting.NewCountResponseInt(count, message, filters)
	return response.ToJSON(), nil
}
