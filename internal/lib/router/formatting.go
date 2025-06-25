package router

import (
	"encoding/json"
	"fmt"
	"strings"

	sharedpb "github.com/studyguides-com/study-guides-api/api/v1/shared"
)

// FormatType represents the different output formats available
type FormatType string

const (
	FormatList  FormatType = "list"
	FormatJSON  FormatType = "json"
	FormatCSV   FormatType = "csv"
	FormatTable FormatType = "table"
)

// TagsAsNumberedList formats a slice of tags as a numbered list with descriptions
func TagsAsNumberedList(tags []*sharedpb.Tag) string {
	if len(tags) == 0 {
		return ""
	}

	var response string
	for i, tag := range tags {
		response += fmt.Sprintf("%d. %s (ID: %s)", i+1, tag.Name, tag.Id)
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

// GetFormatFromParams extracts the format parameter from the params map
func GetFormatFromParams(params map[string]string) FormatType {
	if format, ok := params["format"]; ok && format != "" {
		return FormatType(format)
	}
	return FormatList // default format
}

// GetPublicDescription converts a boolean string to a human-readable description
func GetPublicDescription(publicStr string) string {
	switch publicStr {
	case "true":
		return "public"
	case "false":
		return "private"
	default:
		return "unknown status"
	}
}

// BuildFilterDescription creates a consistent filter description for tag listings
func BuildFilterDescription(params map[string]string, hasTypeFilter, hasContextFilter, hasNameFilter, hasPublicFilter bool) string {
	if hasTypeFilter && hasContextFilter && hasNameFilter && hasPublicFilter {
		return fmt.Sprintf(" for type '%s', context '%s', name containing '%s', and %s", params["type"], params["contextType"], params["name"], GetPublicDescription(params["public"]))
	} else if hasTypeFilter && hasContextFilter && hasNameFilter {
		return fmt.Sprintf(" for type '%s', context '%s', and name containing '%s'", params["type"], params["contextType"], params["name"])
	} else if hasTypeFilter && hasContextFilter && hasPublicFilter {
		return fmt.Sprintf(" for type '%s', context '%s', and %s", params["type"], params["contextType"], GetPublicDescription(params["public"]))
	} else if hasTypeFilter && hasNameFilter && hasPublicFilter {
		return fmt.Sprintf(" for type '%s', name containing '%s', and %s", params["type"], params["name"], GetPublicDescription(params["public"]))
	} else if hasContextFilter && hasNameFilter && hasPublicFilter {
		return fmt.Sprintf(" for context '%s', name containing '%s', and %s", params["contextType"], params["name"], GetPublicDescription(params["public"]))
	} else if hasTypeFilter && hasContextFilter {
		return fmt.Sprintf(" for type '%s' and context '%s'", params["type"], params["contextType"])
	} else if hasTypeFilter && hasNameFilter {
		return fmt.Sprintf(" for type '%s' and name containing '%s'", params["type"], params["name"])
	} else if hasTypeFilter && hasPublicFilter {
		return fmt.Sprintf(" for type '%s' and %s", params["type"], GetPublicDescription(params["public"]))
	} else if hasContextFilter && hasNameFilter {
		return fmt.Sprintf(" for context '%s' and name containing '%s'", params["contextType"], params["name"])
	} else if hasContextFilter && hasPublicFilter {
		return fmt.Sprintf(" for context '%s' and %s", params["contextType"], GetPublicDescription(params["public"]))
	} else if hasNameFilter && hasPublicFilter {
		return fmt.Sprintf(" with name containing '%s' and %s", params["name"], GetPublicDescription(params["public"]))
	} else if hasTypeFilter {
		return fmt.Sprintf(" for type '%s'", params["type"])
	} else if hasContextFilter {
		return fmt.Sprintf(" for context '%s'", params["contextType"])
	} else if hasNameFilter {
		return fmt.Sprintf(" with name containing '%s'", params["name"])
	} else if hasPublicFilter {
		return fmt.Sprintf(" that are %s", GetPublicDescription(params["public"]))
	}
	return ""
}

// TagAsFormatted formats a single tag according to the specified format
func TagAsFormatted(tag *sharedpb.Tag, format FormatType) string {
	switch format {
	case FormatJSON:
		return TagAsJSON(tag)
	case FormatCSV:
		return TagAsCSV(tag)
	case FormatTable:
		return TagAsTable(tag)
	case FormatList:
		fallthrough
	default:
		return TagAsDetailedText(tag)
	}
}

// TagAsJSON formats a single tag as JSON
func TagAsJSON(tag *sharedpb.Tag) string {
	// Create a comprehensive structure for JSON output
	type TagDetailOutput struct {
		ID                string            `json:"id"`
		Name              string            `json:"name"`
		Description       string            `json:"description,omitempty"`
		Type              string            `json:"type"`
		Context           string            `json:"context"`
		ParentTagID       string            `json:"parent_tag_id,omitempty"`
		ContentRating     string            `json:"content_rating"`
		ContentDescriptors []string         `json:"content_descriptors,omitempty"`
		MetaTags          []string          `json:"meta_tags,omitempty"`
		Public            bool              `json:"public"`
		AccessCount       int32             `json:"access_count"`
		Metadata          map[string]string `json:"metadata,omitempty"`
		BatchID           string            `json:"batch_id,omitempty"`
		Hash              string            `json:"hash"`
		HasQuestions      bool              `json:"has_questions"`
		HasChildren       bool              `json:"has_children"`
		OwnerID           string            `json:"owner_id,omitempty"`
		CreatedAt         string            `json:"created_at"`
		UpdatedAt         string            `json:"updated_at"`
	}

	output := TagDetailOutput{
		ID:            tag.Id,
		Name:          tag.Name,
		Type:          tag.Type.String(),
		Context:       tag.Context,
		ContentRating: tag.ContentRating.String(),
		Public:        tag.Public,
		AccessCount:   tag.AccessCount,
		Hash:          tag.Hash,
		HasQuestions:  tag.HasQuestions,
		HasChildren:   tag.HasChildren,
		CreatedAt:     tag.CreatedAt.AsTime().Format("2006-01-02 15:04:05"),
		UpdatedAt:     tag.UpdatedAt.AsTime().Format("2006-01-02 15:04:05"),
	}

	if tag.Description != nil && *tag.Description != "" {
		output.Description = *tag.Description
	}
	if tag.ParentTagId != nil && *tag.ParentTagId != "" {
		output.ParentTagID = *tag.ParentTagId
	}
	if len(tag.ContentDescriptors) > 0 {
		output.ContentDescriptors = tag.ContentDescriptors
	}
	if len(tag.MetaTags) > 0 {
		output.MetaTags = tag.MetaTags
	}
	if len(tag.Metadata) > 0 {
		output.Metadata = tag.Metadata
	}
	if tag.BatchId != nil && *tag.BatchId != "" {
		output.BatchID = *tag.BatchId
	}
	if tag.OwnerId != nil && *tag.OwnerId != "" {
		output.OwnerID = *tag.OwnerId
	}

	jsonBytes, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error formatting as JSON: %v", err)
	}
	return string(jsonBytes)
}

// TagAsCSV formats a single tag as CSV
func TagAsCSV(tag *sharedpb.Tag) string {
	description := ""
	if tag.Description != nil && *tag.Description != "" {
		description = *tag.Description
	}

	parentTagID := ""
	if tag.ParentTagId != nil && *tag.ParentTagId != "" {
		parentTagID = *tag.ParentTagId
	}

	contentDescriptors := ""
	if len(tag.ContentDescriptors) > 0 {
		contentDescriptors = strings.Join(tag.ContentDescriptors, ";")
	}

	metaTags := ""
	if len(tag.MetaTags) > 0 {
		metaTags = strings.Join(tag.MetaTags, ";")
	}

	metadata := ""
	if len(tag.Metadata) > 0 {
		var metadataPairs []string
		for key, value := range tag.Metadata {
			metadataPairs = append(metadataPairs, fmt.Sprintf("%s:%s", key, value))
		}
		metadata = strings.Join(metadataPairs, ";")
	}

	batchID := ""
	if tag.BatchId != nil && *tag.BatchId != "" {
		batchID = *tag.BatchId
	}

	ownerID := ""
	if tag.OwnerId != nil && *tag.OwnerId != "" {
		ownerID = *tag.OwnerId
	}

	header := "id,name,description,type,context,parent_tag_id,content_rating,content_descriptors,meta_tags,public,access_count,metadata,batch_id,hash,has_questions,has_children,owner_id,created_at,updated_at\n"
	row := fmt.Sprintf("\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",%t,%d,\"%s\",\"%s\",\"%s\",%t,%t,\"%s\",\"%s\",\"%s\"\n",
		tag.Id, tag.Name, description, tag.Type.String(), tag.Context, parentTagID, tag.ContentRating.String(),
		contentDescriptors, metaTags, tag.Public, tag.AccessCount, metadata, batchID, tag.Hash,
		tag.HasQuestions, tag.HasChildren, ownerID, tag.CreatedAt.AsTime().Format("2006-01-02 15:04:05"),
		tag.UpdatedAt.AsTime().Format("2006-01-02 15:04:05"))

	return header + row
}

// TagAsTable formats a single tag as a table
func TagAsTable(tag *sharedpb.Tag) string {
	description := ""
	if tag.Description != nil && *tag.Description != "" {
		description = *tag.Description
	}

	parentTagID := ""
	if tag.ParentTagId != nil && *tag.ParentTagId != "" {
		parentTagID = *tag.ParentTagId
	}

	contentDescriptors := ""
	if len(tag.ContentDescriptors) > 0 {
		contentDescriptors = strings.Join(tag.ContentDescriptors, ", ")
	}

	metaTags := ""
	if len(tag.MetaTags) > 0 {
		metaTags = strings.Join(tag.MetaTags, ", ")
	}

	metadata := ""
	if len(tag.Metadata) > 0 {
		var metadataPairs []string
		for key, value := range tag.Metadata {
			metadataPairs = append(metadataPairs, fmt.Sprintf("%s: %s", key, value))
		}
		metadata = strings.Join(metadataPairs, "; ")
	}

	batchID := ""
	if tag.BatchId != nil && *tag.BatchId != "" {
		batchID = *tag.BatchId
	}

	ownerID := ""
	if tag.OwnerId != nil && *tag.OwnerId != "" {
		ownerID = *tag.OwnerId
	}

	response := "| Field | Value |\n"
	response += "|-------|-------|\n"
	response += fmt.Sprintf("| ID | %s |\n", tag.Id)
	response += fmt.Sprintf("| Name | %s |\n", tag.Name)
	if description != "" {
		response += fmt.Sprintf("| Description | %s |\n", description)
	}
	response += fmt.Sprintf("| Type | %s |\n", tag.Type.String())
	response += fmt.Sprintf("| Context | %s |\n", tag.Context)
	if parentTagID != "" {
		response += fmt.Sprintf("| Parent Tag ID | %s |\n", parentTagID)
	}
	response += fmt.Sprintf("| Content Rating | %s |\n", tag.ContentRating.String())
	if contentDescriptors != "" {
		response += fmt.Sprintf("| Content Descriptors | %s |\n", contentDescriptors)
	}
	if metaTags != "" {
		response += fmt.Sprintf("| Meta Tags | %s |\n", metaTags)
	}
	response += fmt.Sprintf("| Public | %t |\n", tag.Public)
	response += fmt.Sprintf("| Access Count | %d |\n", tag.AccessCount)
	if metadata != "" {
		response += fmt.Sprintf("| Metadata | %s |\n", metadata)
	}
	if batchID != "" {
		response += fmt.Sprintf("| Batch ID | %s |\n", batchID)
	}
	response += fmt.Sprintf("| Hash | %s |\n", tag.Hash)
	response += fmt.Sprintf("| Has Questions | %t |\n", tag.HasQuestions)
	response += fmt.Sprintf("| Has Children | %t |\n", tag.HasChildren)
	if ownerID != "" {
		response += fmt.Sprintf("| Owner ID | %s |\n", ownerID)
	}
	response += fmt.Sprintf("| Created | %s |\n", tag.CreatedAt.AsTime().Format("2006-01-02 15:04:05"))
	response += fmt.Sprintf("| Updated | %s |\n", tag.UpdatedAt.AsTime().Format("2006-01-02 15:04:05"))

	return response
}

// TagAsDetailedText formats a single tag with comprehensive details
func TagAsDetailedText(tag *sharedpb.Tag) string {
	var response string
	response += fmt.Sprintf("**Tag Details for ID: %s**\n\n", tag.Id)
	response += fmt.Sprintf("**Name:** %s\n", tag.Name)
	
	if tag.Description != nil && *tag.Description != "" {
		response += fmt.Sprintf("**Description:** %s\n", *tag.Description)
	}
	
	response += fmt.Sprintf("**Type:** %s\n", tag.Type.String())
	response += fmt.Sprintf("**Context:** %s\n", tag.Context)
	
	if tag.ParentTagId != nil && *tag.ParentTagId != "" {
		response += fmt.Sprintf("**Parent Tag ID:** %s\n", *tag.ParentTagId)
	}
	
	response += fmt.Sprintf("**Content Rating:** %s\n", tag.ContentRating.String())
	
	if len(tag.ContentDescriptors) > 0 {
		response += fmt.Sprintf("**Content Descriptors:** %s\n", strings.Join(tag.ContentDescriptors, ", "))
	}
	
	if len(tag.MetaTags) > 0 {
		response += fmt.Sprintf("**Meta Tags:** %s\n", strings.Join(tag.MetaTags, ", "))
	}
	
	response += fmt.Sprintf("**Public:** %t\n", tag.Public)
	response += fmt.Sprintf("**Access Count:** %d\n", tag.AccessCount)
	
	if len(tag.Metadata) > 0 {
		response += "**Metadata:**\n"
		for key, value := range tag.Metadata {
			response += fmt.Sprintf("  - %s: %s\n", key, value)
		}
	}
	
	if tag.BatchId != nil && *tag.BatchId != "" {
		response += fmt.Sprintf("**Batch ID:** %s\n", *tag.BatchId)
	}
	
	response += fmt.Sprintf("**Hash:** %s\n", tag.Hash)
	response += fmt.Sprintf("**Has Questions:** %t\n", tag.HasQuestions)
	response += fmt.Sprintf("**Has Children:** %t\n", tag.HasChildren)
	
	if tag.OwnerId != nil && *tag.OwnerId != "" {
		response += fmt.Sprintf("**Owner ID:** %s\n", *tag.OwnerId)
	}
	
	response += fmt.Sprintf("**Created:** %s\n", tag.CreatedAt.AsTime().Format("2006-01-02 15:04:05"))
	response += fmt.Sprintf("**Updated:** %s\n", tag.UpdatedAt.AsTime().Format("2006-01-02 15:04:05"))
	
	return response
} 