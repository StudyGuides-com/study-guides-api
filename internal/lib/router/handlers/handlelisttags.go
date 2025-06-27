package handlers

import (
	"context"
	"fmt"

	"github.com/studyguides-com/study-guides-api/internal/lib/router/formatting"
	"github.com/studyguides-com/study-guides-api/internal/store"
)

func HandleListTags(ctx context.Context, store store.Store, params map[string]string) (string, error) {
	// Debug: Print all parameters
	fmt.Printf("DEBUG: handleListTags called with params: %+v\n", params)

	// Get the format specified by the AI
	format := formatting.GetFormatFromParams(params)

	// Check if we have any filters (type, contextType, name, or public)
	hasTypeFilter := false
	hasContextFilter := false
	hasNameFilter := false
	hasPublicFilter := false

	if tagType, ok := params["type"]; ok && tagType != "" {
		hasTypeFilter = true
		fmt.Printf("DEBUG: Type parameter found: '%s'\n", tagType)
	}

	if contextType, ok := params["contextType"]; ok && contextType != "" {
		hasContextFilter = true
		fmt.Printf("DEBUG: ContextType parameter found: '%s'\n", contextType)
	}

	if name, ok := params["name"]; ok && name != "" {
		hasNameFilter = true
		fmt.Printf("DEBUG: Name parameter found: '%s'\n", name)
	}

	if publicStr, ok := params["public"]; ok && publicStr != "" {
		hasPublicFilter = true
		fmt.Printf("DEBUG: Public parameter found: '%s'\n", publicStr)
	}

	// If we have any filters, use the new ListTagsWithFilters method
	if hasTypeFilter || hasContextFilter || hasNameFilter || hasPublicFilter {
		// Get the actual unique tag types from the database to validate type parameter
		if hasTypeFilter {
			uniqueTagTypes, err := store.TagStore().UniqueTagTypes(ctx)
			if err != nil {
				return "", err
			}

			// Check if the requested tag type exists in the database
			tagType := params["type"]
			var found bool
			for _, dbTagType := range uniqueTagTypes {
				if dbTagType.String() == tagType {
					found = true
					break
				}
			}

			if !found {
				// Build a list of available tag types for the error message
				var availableTypes []string
				for _, t := range uniqueTagTypes {
					availableTypes = append(availableTypes, t.String())
				}
				return fmt.Sprintf("Invalid tag type '%s'. Available types: %v", tagType, availableTypes), nil
			}
		}

		// Use ListTagsWithFilters for filtered queries
		tags, err := store.TagStore().ListTagsWithFilters(ctx, params)
		if err != nil {
			return "", err
		}

		if len(tags) == 0 {
			filterDesc := formatting.BuildFilterDescription(params, hasTypeFilter, hasContextFilter, hasNameFilter, hasPublicFilter)
			return fmt.Sprintf("No tags found%s.", filterDesc), nil
		}

		// Format the response according to the AI-specified format
		if format == formatting.FormatList {
			filterDesc := formatting.BuildFilterDescription(params, hasTypeFilter, hasContextFilter, hasNameFilter, hasPublicFilter)
			limitMsg := formatting.BuildLimitMessage(params)
			response := fmt.Sprintf("Found %d tags%s%s:\n", len(tags), filterDesc, limitMsg)
			response += formatting.FormatTags(tags, format)
			return response, nil
		} else {
			// For other formats, just return the formatted data
			return formatting.FormatTags(tags, format), nil
		}
	}

	// Default to listing root tags if no filters specified
	tags, err := store.TagStore().ListRootTags(ctx, params)
	if err != nil {
		return "", err
	}

	if len(tags) == 0 {
		return "No root tags found.", nil
	}

	// Format the response according to the AI-specified format
	if format == formatting.FormatList {
		limitMsg := formatting.BuildLimitMessage(params)
		response := fmt.Sprintf("Found %d root tags%s:\n", len(tags), limitMsg)
		response += formatting.FormatTags(tags, format)
		return response, nil
	} else {
		// For other formats, just return the formatted data
		return formatting.FormatTags(tags, format), nil
	}
}
