package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/studyguides-com/study-guides-api/internal/lib/router/formatting"
	"github.com/studyguides-com/study-guides-api/internal/store"
)

func HandleUniqueTagTypes(ctx context.Context, store store.Store, params map[string]string) (string, error) {
	tagTypes, err := store.TagStore().UniqueTagTypes(ctx)
	if err != nil {
		return "", err
	}

	if len(tagTypes) == 0 {
		response := formatting.NewListResponse([]interface{}{}, "No tag types found in the system.", nil, nil)
		return response.ToJSON(), nil
	}

	// Convert tag types to strings
	var tagTypeStrings []string
	for _, tagType := range tagTypes {
		tagTypeStrings = append(tagTypeStrings, tagType.String())
	}

	// Get the format specified by the AI
	format := formatting.GetFormatFromParams(params)

	// Format the data
	var data interface{}
	switch format {
	case formatting.FormatJSON:
		data = tagTypeStrings
	case formatting.FormatCSV:
		data = "tag_type\n" + strings.Join(tagTypeStrings, "\n")

	case formatting.FormatList:
		// For list format, return as a string but ensure it's not empty
		if len(tagTypeStrings) > 0 {
			data = strings.Join(tagTypeStrings, "\n")
		} else {
			data = ""
		}
	default:
		// Default to list format
		if len(tagTypeStrings) > 0 {
			data = strings.Join(tagTypeStrings, "\n")
		} else {
			data = ""
		}
	}

	message := fmt.Sprintf("Found %d unique tag types", len(tagTypes))

	// Determine content type based on format
	var contentType string
	if format == formatting.FormatJSON {
		contentType = "application/json"
	} else if format == formatting.FormatCSV {
		contentType = "text/csv"
	} else {
		contentType = "text/plain"
	}

	// Create response with correct content type
	response := &formatting.APIResponse{
		Type:        formatting.ResponseTypeList,
		Data:        data,
		Message:     message,
		ContentType: contentType,
	}

	return response.ToJSON(), nil
}
