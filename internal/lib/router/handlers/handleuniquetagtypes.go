package handlers

import (
	"context"
	"fmt"

	"github.com/studyguides-com/study-guides-api/internal/store"
)

func HandleUniqueTagTypes(ctx context.Context, store store.Store, params map[string]string) (string, error) {
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