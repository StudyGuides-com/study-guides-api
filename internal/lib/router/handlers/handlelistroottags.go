package handlers

import (
	"context"
	"fmt"

	sharedpb "github.com/studyguides-com/study-guides-api/api/v1/shared"
	"github.com/studyguides-com/study-guides-api/internal/store"
)

func HandleListRootTags(ctx context.Context, store store.Store, params map[string]string) (string, error) {
	// Get the format specified by the AI
	format := GetFormatFromParams(params)
	
	// Check if we have any filters
	hasNameFilter := false
	hasPublicFilter := false
	if name, ok := params["name"]; ok && name != "" {
		hasNameFilter = true
	}
	if publicStr, ok := params["public"]; ok && publicStr != "" {
		hasPublicFilter = true
	}
	
	var tags []*sharedpb.Tag
	var err error
	
	if hasNameFilter || hasPublicFilter {
		// If any filter is specified, use ListTagsWithFilters with parentTagId IS NULL
		// We need to add the parentTagId filter to the params
		filterParams := make(map[string]string)
		for k, v := range params {
			filterParams[k] = v
		}
		// Add a special marker for root tags
		filterParams["rootOnly"] = "true"
		
		tags, err = store.TagStore().ListTagsWithFilters(ctx, filterParams)
	} else {
		// Get root tags without any filters
		tags, err = store.TagStore().ListRootTags(ctx, params)
	}
	
	if err != nil {
		return "", err
	}
	
	if len(tags) == 0 {
		if hasNameFilter && hasPublicFilter {
			return fmt.Sprintf("No root tags found with name containing '%s' that are %s.", params["name"], GetPublicDescription(params["public"])), nil
		} else if hasNameFilter {
			return fmt.Sprintf("No root tags found with name containing '%s'.", params["name"]), nil
		} else if hasPublicFilter {
			return fmt.Sprintf("No root tags found that are %s.", GetPublicDescription(params["public"])), nil
		}
		return "No root tags found.", nil
	}
	
	// Format the response according to the AI-specified format
	if format == FormatList {
		limitMsg := BuildLimitMessage(params)
		var response string
		if hasNameFilter && hasPublicFilter {
			response = fmt.Sprintf("Found %d root tags with name containing '%s' that are %s%s:\n", len(tags), params["name"], GetPublicDescription(params["public"]), limitMsg)
		} else if hasNameFilter {
			response = fmt.Sprintf("Found %d root tags with name containing '%s'%s:\n", len(tags), params["name"], limitMsg)
		} else if hasPublicFilter {
			response = fmt.Sprintf("Found %d root tags that are %s%s:\n", len(tags), GetPublicDescription(params["public"]), limitMsg)
		} else {
			response = fmt.Sprintf("Found %d root tags%s:\n", len(tags), limitMsg)
		}
		response += FormatTags(tags, format)
		return response, nil
	} else {
		// For other formats, just return the formatted data
		return FormatTags(tags, format), nil
	}
} 