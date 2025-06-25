package router

import (
	"context"
	"fmt"

	sharedpb "github.com/studyguides-com/study-guides-api/api/v1/shared"
	"github.com/studyguides-com/study-guides-api/internal/store"
)

func handleTagCount(ctx context.Context, store store.Store, params map[string]string) (string, error) {
	count, err := store.TagStore().CountTags(ctx, params)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("You have %d tags matching the criteria.", count), nil
}

func handleListTags(ctx context.Context, store store.Store, params map[string]string) (string, error) {
	// Check if type parameter is provided
	if tagType, ok := params["type"]; ok && tagType != "" {
		// Get the actual unique tag types from the database
		uniqueTagTypes, err := store.TagStore().UniqueTagTypes(ctx)
		if err != nil {
			return "", err
		}
		
		// Check if the requested tag type exists in the database
		var foundTagType sharedpb.TagType
		var found bool
		for _, dbTagType := range uniqueTagTypes {
			if dbTagType.String() == tagType {
				foundTagType = dbTagType
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
		
		// Use ListTagsByType if type is specified
		tags, err := store.TagStore().ListTagsByType(ctx, foundTagType)
		if err != nil {
			return "", err
		}
		
		if len(tags) == 0 {
			return "No tags found for the specified type.", nil
		}
		
		// Build response with tag details
		response := fmt.Sprintf("Found %d tags of type %s:\n", len(tags), tagType)
		for i, tag := range tags {
			response += fmt.Sprintf("%d. %s", i+1, tag.Name)
			if tag.Description != nil && *tag.Description != "" {
				response += fmt.Sprintf(" - %s", *tag.Description)
			}
			response += "\n"
		}
		return response, nil
	}
	
	// Default to listing root tags if no type specified
	tags, err := store.TagStore().ListRootTags(ctx)
	if err != nil {
		return "", err
	}
	
	if len(tags) == 0 {
		return "No root tags found.", nil
	}
	
	// Build response with tag details
	response := fmt.Sprintf("Found %d root tags:\n", len(tags))
	for i, tag := range tags {
		response += fmt.Sprintf("%d. %s", i+1, tag.Name)
		if tag.Description != nil && *tag.Description != "" {
			response += fmt.Sprintf(" - %s", *tag.Description)
		}
		response += "\n"
	}
	return response, nil
}

func handleUniqueTagTypes(ctx context.Context, store store.Store, params map[string]string) (string, error) {
	tagTypes, err := store.TagStore().UniqueTagTypes(ctx)
	if err != nil {
		return "", err
	}
	
	if len(tagTypes) == 0 {
		return "No tag types found in the system.", nil
	}
	
	response := fmt.Sprintf("Found %d unique tag types:\n", len(tagTypes))
	for i, tagType := range tagTypes {
		response += fmt.Sprintf("%d. %s\n", i+1, tagType.String())
	}
	return response, nil
}

func handleUnknown(ctx context.Context, store store.Store, params map[string]string) (string, error) {
	return "I'm sorry, I don't understand your request. Could you please rephrase it or ask about something else?", nil
}
