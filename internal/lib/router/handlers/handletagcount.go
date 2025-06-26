package handlers

import (
	"context"
	"fmt"

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
	
	// Build the response message
	var response string
	if hasTypeFilter && hasContextFilter && hasNameFilter && hasPublicFilter {
		response = fmt.Sprintf("Found %d tags for type '%s', context '%s', name containing '%s', and %s.", count, params["type"], params["contextType"], params["name"], formatting.GetPublicDescription(params["public"]))
	} else if hasTypeFilter && hasContextFilter && hasNameFilter {
		response = fmt.Sprintf("Found %d tags for type '%s', context '%s', and name containing '%s'.", count, params["type"], params["contextType"], params["name"])
	} else if hasTypeFilter && hasContextFilter && hasPublicFilter {
		response = fmt.Sprintf("Found %d tags for type '%s', context '%s', and %s.", count, params["type"], params["contextType"], formatting.GetPublicDescription(params["public"]))
	} else if hasTypeFilter && hasNameFilter && hasPublicFilter {
		response = fmt.Sprintf("Found %d tags for type '%s', name containing '%s', and %s.", count, params["type"], params["name"], formatting.GetPublicDescription(params["public"]))
	} else if hasContextFilter && hasNameFilter && hasPublicFilter {
		response = fmt.Sprintf("Found %d tags for context '%s', name containing '%s', and %s.", count, params["contextType"], params["name"], formatting.GetPublicDescription(params["public"]))
	} else if hasTypeFilter && hasContextFilter {
		response = fmt.Sprintf("Found %d tags for type '%s' and context '%s'.", count, params["type"], params["contextType"])
	} else if hasTypeFilter && hasNameFilter {
		response = fmt.Sprintf("Found %d tags for type '%s' and name containing '%s'.", count, params["type"], params["name"])
	} else if hasTypeFilter && hasPublicFilter {
		response = fmt.Sprintf("Found %d tags for type '%s' and %s.", count, params["type"], formatting.GetPublicDescription(params["public"]))
	} else if hasContextFilter && hasNameFilter {
		response = fmt.Sprintf("Found %d tags for context '%s' and name containing '%s'.", count, params["contextType"], params["name"])
	} else if hasContextFilter && hasPublicFilter {
		response = fmt.Sprintf("Found %d tags for context '%s' and %s.", count, params["contextType"], formatting.GetPublicDescription(params["public"]))
	} else if hasNameFilter && hasPublicFilter {
		response = fmt.Sprintf("Found %d tags with name containing '%s' and %s.", count, params["name"], formatting.GetPublicDescription(params["public"]))
	} else if hasTypeFilter {
		response = fmt.Sprintf("Found %d tags for type '%s'.", count, params["type"])
	} else if hasContextFilter {
		response = fmt.Sprintf("Found %d tags for context '%s'.", count, params["contextType"])
	} else if hasNameFilter {
		response = fmt.Sprintf("Found %d tags with name containing '%s'.", count, params["name"])
	} else if hasPublicFilter {
		response = fmt.Sprintf("Found %d tags that are %s.", count, formatting.GetPublicDescription(params["public"]))
	} else {
		response = fmt.Sprintf("Found %d tags total.", count)
	}
	
	return response, nil
} 