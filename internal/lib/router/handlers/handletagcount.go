package handlers

import (
	"context"
	"fmt"

	"github.com/studyguides-com/study-guides-api/internal/store"
)

func HandleTagCount(ctx context.Context, store store.Store, params map[string]string) (string, error) {
	count, err := store.TagStore().CountTags(ctx, params)
	if err != nil {
		return "", err
	}

	// Build a descriptive message based on the filters used
	var filterDesc string
	
	hasTypeFilter := false
	hasContextFilter := false
	hasNameFilter := false
	hasPublicFilter := false
	
	if tagType, ok := params["type"]; ok && tagType != "" {
		hasTypeFilter = true
	}
	
	if contextType, ok := params["contextType"]; ok && contextType != "" {
		hasContextFilter = true
	}

	if name, ok := params["name"]; ok && name != "" {
		hasNameFilter = true
	}

	if publicStr, ok := params["public"]; ok && publicStr != "" {
		hasPublicFilter = true
	}
	
	// Build filter description with all possible combinations
	if hasTypeFilter && hasContextFilter && hasNameFilter && hasPublicFilter {
		filterDesc = fmt.Sprintf(" of type '%s', context '%s', name containing '%s', and %s", params["type"], params["contextType"], params["name"], GetPublicDescription(params["public"]))
	} else if hasTypeFilter && hasContextFilter && hasNameFilter {
		filterDesc = fmt.Sprintf(" of type '%s', context '%s', and name containing '%s'", params["type"], params["contextType"], params["name"])
	} else if hasTypeFilter && hasContextFilter && hasPublicFilter {
		filterDesc = fmt.Sprintf(" of type '%s', context '%s', and %s", params["type"], params["contextType"], GetPublicDescription(params["public"]))
	} else if hasTypeFilter && hasNameFilter && hasPublicFilter {
		filterDesc = fmt.Sprintf(" of type '%s', name containing '%s', and %s", params["type"], params["name"], GetPublicDescription(params["public"]))
	} else if hasContextFilter && hasNameFilter && hasPublicFilter {
		filterDesc = fmt.Sprintf(" with context '%s', name containing '%s', and %s", params["contextType"], params["name"], GetPublicDescription(params["public"]))
	} else if hasTypeFilter && hasContextFilter {
		filterDesc = fmt.Sprintf(" of type '%s' and context '%s'", params["type"], params["contextType"])
	} else if hasTypeFilter && hasNameFilter {
		filterDesc = fmt.Sprintf(" of type '%s' and name containing '%s'", params["type"], params["name"])
	} else if hasTypeFilter && hasPublicFilter {
		filterDesc = fmt.Sprintf(" of type '%s' and %s", params["type"], GetPublicDescription(params["public"]))
	} else if hasContextFilter && hasNameFilter {
		filterDesc = fmt.Sprintf(" with context '%s' and name containing '%s'", params["contextType"], params["name"])
	} else if hasContextFilter && hasPublicFilter {
		filterDesc = fmt.Sprintf(" with context '%s' and %s", params["contextType"], GetPublicDescription(params["public"]))
	} else if hasNameFilter && hasPublicFilter {
		filterDesc = fmt.Sprintf(" with name containing '%s' and %s", params["name"], GetPublicDescription(params["public"]))
	} else if hasTypeFilter {
		filterDesc = fmt.Sprintf(" of type '%s'", params["type"])
	} else if hasContextFilter {
		filterDesc = fmt.Sprintf(" with context '%s'", params["contextType"])
	} else if hasNameFilter {
		filterDesc = fmt.Sprintf(" with name containing '%s'", params["name"])
	} else if hasPublicFilter {
		filterDesc = fmt.Sprintf(" that are %s", GetPublicDescription(params["public"]))
	} else {
		filterDesc = " in total"
	}

	return fmt.Sprintf("You have %d tags%s.", count, filterDesc), nil
} 