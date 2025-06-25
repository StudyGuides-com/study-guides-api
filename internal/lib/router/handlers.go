package router

import (
	"context"
	"encoding/json"
	"fmt"

	sharedpb "github.com/studyguides-com/study-guides-api/api/v1/shared"
	"github.com/studyguides-com/study-guides-api/internal/store"
)

// FormatType represents the different output formats available
type FormatType string

const (
	FormatList FormatType = "list"
	FormatJSON FormatType = "json"
	FormatCSV  FormatType = "csv"
	FormatTable FormatType = "table"
)

// TagsAsNumberedList formats a slice of tags as a numbered list with descriptions
func TagsAsNumberedList(tags []*sharedpb.Tag) string {
	if len(tags) == 0 {
		return ""
	}
	
	var response string
	for i, tag := range tags {
		response += fmt.Sprintf("%d. %s", i+1, tag.Name)
		if tag.Description != nil && *tag.Description != "" {
			response += fmt.Sprintf(" - %s", *tag.Description)
		}
		response += "\n"
	}
	return response
}

// TagsAsJSON formats a slice of tags as JSON
func TagsAsJSON(tags []*sharedpb.Tag) string {
	if len(tags) == 0 {
		return "[]"
	}
	
	// Create a simplified structure for JSON output
	type TagOutput struct {
		Name        string `json:"name"`
		Description string `json:"description,omitempty"`
		Type        string `json:"type"`
		ID          string `json:"id"`
	}
	
	var output []TagOutput
	for _, tag := range tags {
		tagOutput := TagOutput{
			Name: tag.Name,
			Type: tag.Type.String(),
			ID:   tag.Id,
		}
		if tag.Description != nil && *tag.Description != "" {
			tagOutput.Description = *tag.Description
		}
		output = append(output, tagOutput)
	}
	
	jsonBytes, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error formatting as JSON: %v", err)
	}
	return string(jsonBytes)
}

// TagsAsCSV formats a slice of tags as CSV
func TagsAsCSV(tags []*sharedpb.Tag) string {
	if len(tags) == 0 {
		return "name,description,type,id\n"
	}
	
	response := "name,description,type,id\n"
	for _, tag := range tags {
		description := ""
		if tag.Description != nil && *tag.Description != "" {
			description = *tag.Description
		}
		response += fmt.Sprintf("\"%s\",\"%s\",\"%s\",\"%s\"\n", 
			tag.Name, description, tag.Type.String(), tag.Id)
	}
	return response
}

// TagsAsTable formats a slice of tags as a simple table
func TagsAsTable(tags []*sharedpb.Tag) string {
	if len(tags) == 0 {
		return "No tags found."
	}
	
	response := "| Name | Description | Type | ID |\n"
	response += "|------|-------------|------|----|\n"
	
	for _, tag := range tags {
		description := ""
		if tag.Description != nil && *tag.Description != "" {
			description = *tag.Description
		}
		response += fmt.Sprintf("| %s | %s | %s | %s |\n", 
			tag.Name, description, tag.Type.String(), tag.Id)
	}
	return response
}

// FormatTags formats tags according to the specified format
func FormatTags(tags []*sharedpb.Tag, format FormatType) string {
	switch format {
	case FormatJSON:
		return TagsAsJSON(tags)
	case FormatCSV:
		return TagsAsCSV(tags)
	case FormatTable:
		return TagsAsTable(tags)
	case FormatList:
		fallthrough
	default:
		return TagsAsNumberedList(tags)
	}
}

// getFormatFromParams extracts the format parameter from the params map
func getFormatFromParams(params map[string]string) FormatType {
	if format, ok := params["format"]; ok && format != "" {
		return FormatType(format)
	}
	return FormatList // default format
}

// getPublicDescription converts a boolean string to a human-readable description
func getPublicDescription(publicStr string) string {
	if publicStr == "true" {
		return "public"
	} else if publicStr == "false" {
		return "private"
	}
	return "unknown status"
}

// buildFilterDescription creates a consistent filter description for tag listings
func buildFilterDescription(params map[string]string, hasTypeFilter, hasContextFilter, hasPublicFilter bool) string {
	if hasTypeFilter && hasContextFilter && hasPublicFilter {
		return fmt.Sprintf(" for type '%s', context '%s', and %s", params["type"], params["contextType"], getPublicDescription(params["public"]))
	} else if hasTypeFilter && hasContextFilter {
		return fmt.Sprintf(" for type '%s' and context '%s'", params["type"], params["contextType"])
	} else if hasTypeFilter && hasPublicFilter {
		return fmt.Sprintf(" for type '%s' and %s", params["type"], getPublicDescription(params["public"]))
	} else if hasContextFilter && hasPublicFilter {
		return fmt.Sprintf(" for context '%s' and %s", params["contextType"], getPublicDescription(params["public"]))
	} else if hasTypeFilter {
		return fmt.Sprintf(" for type '%s'", params["type"])
	} else if hasContextFilter {
		return fmt.Sprintf(" for context '%s'", params["contextType"])
	} else if hasPublicFilter {
		return fmt.Sprintf(" that are %s", getPublicDescription(params["public"]))
	}
	return ""
}

func handleTagCount(ctx context.Context, store store.Store, params map[string]string) (string, error) {
	count, err := store.TagStore().CountTags(ctx, params)
	if err != nil {
		return "", err
	}

	// Build a descriptive message based on the filters used
	var filterDesc string
	
	hasTypeFilter := false
	hasContextFilter := false
	hasPublicFilter := false
	
	if tagType, ok := params["type"]; ok && tagType != "" {
		hasTypeFilter = true
	}
	
	if contextType, ok := params["contextType"]; ok && contextType != "" {
		hasContextFilter = true
	}

	if publicStr, ok := params["public"]; ok && publicStr != "" {
		hasPublicFilter = true
	}
	
	if hasTypeFilter && hasContextFilter && hasPublicFilter {
		filterDesc = fmt.Sprintf(" of type '%s', context '%s', and %s", params["type"], params["contextType"], getPublicDescription(params["public"]))
	} else if hasTypeFilter && hasContextFilter {
		filterDesc = fmt.Sprintf(" of type '%s' and context '%s'", params["type"], params["contextType"])
	} else if hasTypeFilter && hasPublicFilter {
		filterDesc = fmt.Sprintf(" of type '%s' and %s", params["type"], getPublicDescription(params["public"]))
	} else if hasContextFilter && hasPublicFilter {
		filterDesc = fmt.Sprintf(" with context '%s' and %s", params["contextType"], getPublicDescription(params["public"]))
	} else if hasTypeFilter {
		filterDesc = fmt.Sprintf(" of type '%s'", params["type"])
	} else if hasContextFilter {
		filterDesc = fmt.Sprintf(" with context '%s'", params["contextType"])
	} else if hasPublicFilter {
		filterDesc = fmt.Sprintf(" that are %s", getPublicDescription(params["public"]))
	} else {
		filterDesc = " in total"
	}

	return fmt.Sprintf("You have %d tags%s.", count, filterDesc), nil
}

func handleListTags(ctx context.Context, store store.Store, params map[string]string) (string, error) {
	// Debug: Print all parameters
	fmt.Printf("DEBUG: handleListTags called with params: %+v\n", params)
	
	// Get the format specified by the AI
	format := getFormatFromParams(params)
	
	// Check if we have any filters (type, contextType, or public)
	hasTypeFilter := false
	hasContextFilter := false
	hasPublicFilter := false
	
	if tagType, ok := params["type"]; ok && tagType != "" {
		hasTypeFilter = true
		fmt.Printf("DEBUG: Type parameter found: '%s'\n", tagType)
	}
	
	if contextType, ok := params["contextType"]; ok && contextType != "" {
		hasContextFilter = true
		fmt.Printf("DEBUG: ContextType parameter found: '%s'\n", contextType)
	}

	if publicStr, ok := params["public"]; ok && publicStr != "" {
		hasPublicFilter = true
		fmt.Printf("DEBUG: Public parameter found: '%s'\n", publicStr)
	}
	
	// If we have any filters, use the new ListTagsWithFilters method
	if hasTypeFilter || hasContextFilter || hasPublicFilter {
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
			filterDesc := buildFilterDescription(params, hasTypeFilter, hasContextFilter, hasPublicFilter)
			return fmt.Sprintf("No tags found%s.", filterDesc), nil
		}
		
		// Format the response according to the AI-specified format
		if format == FormatList {
			filterDesc := buildFilterDescription(params, hasTypeFilter, hasContextFilter, hasPublicFilter)
			response := fmt.Sprintf("Found %d tags%s:\n", len(tags), filterDesc)
			response += FormatTags(tags, format)
			return response, nil
		} else {
			// For other formats, just return the formatted data
			return FormatTags(tags, format), nil
		}
	}
	
	// Default to listing root tags if no filters specified
	tags, err := store.TagStore().ListRootTags(ctx)
	if err != nil {
		return "", err
	}
	
	if len(tags) == 0 {
		return "No root tags found.", nil
	}
	
	// Format the response according to the AI-specified format
	if format == FormatList {
		response := fmt.Sprintf("Found %d root tags:\n", len(tags))
		response += FormatTags(tags, format)
		return response, nil
	} else {
		// For other formats, just return the formatted data
		return FormatTags(tags, format), nil
	}
}

func handleListRootTags(ctx context.Context, store store.Store, params map[string]string) (string, error) {
	// Get the format specified by the AI
	format := getFormatFromParams(params)
	
	// Check if we have a public filter
	hasPublicFilter := false
	if publicStr, ok := params["public"]; ok && publicStr != "" {
		hasPublicFilter = true
	}
	
	var tags []*sharedpb.Tag
	var err error
	
	if hasPublicFilter {
		// If public filter is specified, use ListTagsWithFilters with parentTagId IS NULL
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
		tags, err = store.TagStore().ListRootTags(ctx)
	}
	
	if err != nil {
		return "", err
	}
	
	if len(tags) == 0 {
		if hasPublicFilter {
			return fmt.Sprintf("No root tags found that are %s.", getPublicDescription(params["public"])), nil
		}
		return "No root tags found.", nil
	}
	
	// Format the response according to the AI-specified format
	if format == FormatList {
		var response string
		if hasPublicFilter {
			response = fmt.Sprintf("Found %d root tags that are %s:\n", len(tags), getPublicDescription(params["public"]))
		} else {
			response = fmt.Sprintf("Found %d root tags:\n", len(tags))
		}
		response += FormatTags(tags, format)
		return response, nil
	} else {
		// For other formats, just return the formatted data
		return FormatTags(tags, format), nil
	}
}

func handleUniqueTagTypes(ctx context.Context, store store.Store, params map[string]string) (string, error) {
	tagTypes, err := store.TagStore().UniqueTagTypes(ctx)
	if err != nil {
		return "", err
	}
	
	if len(tagTypes) == 0 {
		return "No tag types found in the system.", nil
	}
	
	// For now, tag types only support list format since they're simple strings
	response := fmt.Sprintf("Found %d unique tag types:\n", len(tagTypes))
	for i, tagType := range tagTypes {
		response += fmt.Sprintf("%d. %s\n", i+1, tagType.String())
	}
	return response, nil
}

func handleUniqueContextTypes(ctx context.Context, store store.Store, params map[string]string) (string, error) {
	contextTypes, err := store.TagStore().UniqueContextTypes(ctx)
	if err != nil {
		return "", err
	}
	
	if len(contextTypes) == 0 {
		return "No context types found in the system.", nil
	}
	
	// For now, context types only support list format since they're simple strings
	response := fmt.Sprintf("Found %d unique context types:\n", len(contextTypes))
	for i, contextType := range contextTypes {
		response += fmt.Sprintf("%d. %s\n", i+1, contextType)
	}
	return response, nil
}

func handleUnknown(ctx context.Context, store store.Store, params map[string]string) (string, error) {
	return "I'm sorry, I don't understand your request. Could you please rephrase it or ask about something else?", nil
}
