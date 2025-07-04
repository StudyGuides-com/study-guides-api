package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/studyguides-com/study-guides-api/internal/lib/router/formatting"
	"github.com/studyguides-com/study-guides-api/internal/store"
)

func HandleUniqueContextTypes(ctx context.Context, store store.Store, params map[string]string) (string, error) {
	contextTypes, err := store.TagStore().UniqueContextTypes(ctx)
	if err != nil {
		return "", err
	}

	if len(contextTypes) == 0 {
		response := formatting.NewListResponse([]interface{}{}, "No context types found.", nil, nil)
		return response.ToJSON(), nil
	}

	// Get the format specified by the AI
	format := formatting.GetFormatFromParams(params)

	// Format the data
	var data interface{}
	switch format {
	case formatting.FormatJSON:
		data = contextTypes
	case formatting.FormatCSV:
		data = "context_type\n" + strings.Join(contextTypes, "\n")

	case formatting.FormatList:
		data = strings.Join(contextTypes, "\n")
	default:
		data = strings.Join(contextTypes, "\n")
	}

	message := fmt.Sprintf("Found %d unique context types", len(contextTypes))

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
