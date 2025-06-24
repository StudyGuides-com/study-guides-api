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
		// Use ListTagsByType if type is specified
		tags, err := store.TagStore().ListTagsByType(ctx, sharedpb.TagType(sharedpb.TagType_value[tagType]))
		if err != nil {
			return "", err
		}
		
		if len(tags) == 0 {
			return "No tags found for the specified type.", nil
		}
		
		return fmt.Sprintf("Found %d tags of type %s.", len(tags), tagType), nil
	}
	
	// Default to listing root tags if no type specified
	tags, err := store.TagStore().ListRootTags(ctx)
	if err != nil {
		return "", err
	}
	
	if len(tags) == 0 {
		return "No root tags found.", nil
	}
	
	return fmt.Sprintf("Found %d root tags.", len(tags)), nil
}

func handleUnknown(ctx context.Context, store store.Store, params map[string]string) (string, error) {
	return "I'm sorry, I don't understand your request. Could you please rephrase it or ask about something else?", nil
}
