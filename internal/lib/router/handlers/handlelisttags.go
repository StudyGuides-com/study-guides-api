package handlers

import (
	"context"
	"fmt"

	sharedpb "github.com/studyguides-com/study-guides-api/api/v1/shared"
	"github.com/studyguides-com/study-guides-api/internal/lib/router/formatting"
	"github.com/studyguides-com/study-guides-api/internal/store"
)

func HandleListTags(ctx context.Context, store store.Store, params map[string]string) (string, error) {
	// Get the format specified by the AI
	format := formatting.GetFormatFromParams(params)

	// Check if we have any filters
	hasTypeFilter := params["type"] != ""
	hasContextFilter := params["contextType"] != ""
	hasNameFilter := params["name"] != ""
	hasPublicFilter := params["public"] != ""

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

	// Get tags from store
	var tags []*sharedpb.Tag
	var err error

	if hasTypeFilter || hasContextFilter || hasNameFilter || hasPublicFilter {
		// Validate tag type if specified
		if hasTypeFilter {
			uniqueTagTypes, err := store.TagStore().UniqueTagTypes(ctx)
			if err != nil {
				return "", err
			}

			tagType := params["type"]
			var found bool
			for _, dbTagType := range uniqueTagTypes {
				if dbTagType.String() == tagType {
					found = true
					break
				}
			}

			if !found {
				var availableTypes []string
				for _, t := range uniqueTagTypes {
					availableTypes = append(availableTypes, t.String())
				}
				errorMessage := fmt.Sprintf("Invalid tag type '%s'. Available types: %v", tagType, availableTypes)
				response := formatting.NewSingleResponse(nil, errorMessage)
				return response.ToJSON(), nil
			}
		}

		tags, err = store.TagStore().ListTagsWithFilters(ctx, params)
	} else {
		tags, err = store.TagStore().ListRootTags(ctx, params)
	}

	if err != nil {
		return "", err
	}

	// Handle empty results
	if len(tags) == 0 {
		var message string
		if hasTypeFilter || hasContextFilter || hasNameFilter || hasPublicFilter {
			filterDesc := formatting.BuildFilterDescription(params, hasTypeFilter, hasContextFilter, hasNameFilter, hasPublicFilter)
			message = fmt.Sprintf("No tags found%s.", filterDesc)
		} else {
			message = "No root tags found."
		}
		response := formatting.NewListResponse([]interface{}{}, message, filters, nil)
		return response.ToJSON(), nil
	}

	// Format the data using the formatter
	formatter := formatting.NewTagFormatter(tags)
	data := formatter.Format(format)

	// Build message
	var message string
	if hasTypeFilter || hasContextFilter || hasNameFilter || hasPublicFilter {
		filterDesc := formatting.BuildFilterDescription(params, hasTypeFilter, hasContextFilter, hasNameFilter, hasPublicFilter)
		limitMsg := formatting.BuildLimitMessage(params)
		message = fmt.Sprintf("Found %d tags%s%s", len(tags), filterDesc, limitMsg)
	} else {
		limitMsg := formatting.BuildLimitMessage(params)
		message = fmt.Sprintf("Found %d root tags%s", len(tags), limitMsg)
	}

	// Create response
	response := formatting.NewListResponse(data, message, filters, nil)
	return response.ToJSON(), nil
}

