package router

import (
	"context"
	"fmt"
	"math"
	"time"

	sharedpb "github.com/studyguides-com/study-guides-api/api/v1/shared"
	"github.com/studyguides-com/study-guides-api/internal/store"
)

func handleTagCount(ctx context.Context, store store.Store, params map[string]string) (string, error) {
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

func handleListTags(ctx context.Context, store store.Store, params map[string]string) (string, error) {
	// Debug: Print all parameters
	fmt.Printf("DEBUG: handleListTags called with params: %+v\n", params)
	
	// Get the format specified by the AI
	format := GetFormatFromParams(params)
	
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
			filterDesc := BuildFilterDescription(params, hasTypeFilter, hasContextFilter, hasNameFilter, hasPublicFilter)
			return fmt.Sprintf("No tags found%s.", filterDesc), nil
		}
		
		// Format the response according to the AI-specified format
		if format == FormatList {
			filterDesc := BuildFilterDescription(params, hasTypeFilter, hasContextFilter, hasNameFilter, hasPublicFilter)
			limitMsg := BuildLimitMessage(params)
			response := fmt.Sprintf("Found %d tags%s%s:\n", len(tags), filterDesc, limitMsg)
			response += FormatTags(tags, format)
			return response, nil
		} else {
			// For other formats, just return the formatted data
			return FormatTags(tags, format), nil
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
	if format == FormatList {
		limitMsg := BuildLimitMessage(params)
		response := fmt.Sprintf("Found %d root tags%s:\n", len(tags), limitMsg)
		response += FormatTags(tags, format)
		return response, nil
	} else {
		// For other formats, just return the formatted data
		return FormatTags(tags, format), nil
	}
}

func handleListRootTags(ctx context.Context, store store.Store, params map[string]string) (string, error) {
	// Get the format specified by the AI
	format := GetFormatFromParams(params)
	
	// Check if we have any filters
	hasNameFilter := false
	hasPublicFilter := false
	if name, ok := params["name"]; ok && name != "" {
		hasNameFilter = true
	}
	if publicStr, ok := params["public"]; ok && publicStr != "" {
		hasPublicFilter = true
	}
	
	var tags []*sharedpb.Tag
	var err error
	
	if hasNameFilter || hasPublicFilter {
		// If any filter is specified, use ListTagsWithFilters with parentTagId IS NULL
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
		tags, err = store.TagStore().ListRootTags(ctx, params)
	}
	
	if err != nil {
		return "", err
	}
	
	if len(tags) == 0 {
		if hasNameFilter && hasPublicFilter {
			return fmt.Sprintf("No root tags found with name containing '%s' that are %s.", params["name"], GetPublicDescription(params["public"])), nil
		} else if hasNameFilter {
			return fmt.Sprintf("No root tags found with name containing '%s'.", params["name"]), nil
		} else if hasPublicFilter {
			return fmt.Sprintf("No root tags found that are %s.", GetPublicDescription(params["public"])), nil
		}
		return "No root tags found.", nil
	}
	
	// Format the response according to the AI-specified format
	if format == FormatList {
		limitMsg := BuildLimitMessage(params)
		var response string
		if hasNameFilter && hasPublicFilter {
			response = fmt.Sprintf("Found %d root tags with name containing '%s' that are %s%s:\n", len(tags), params["name"], GetPublicDescription(params["public"]), limitMsg)
		} else if hasNameFilter {
			response = fmt.Sprintf("Found %d root tags with name containing '%s'%s:\n", len(tags), params["name"], limitMsg)
		} else if hasPublicFilter {
			response = fmt.Sprintf("Found %d root tags that are %s%s:\n", len(tags), GetPublicDescription(params["public"]), limitMsg)
		} else {
			response = fmt.Sprintf("Found %d root tags%s:\n", len(tags), limitMsg)
		}
		response += FormatTags(tags, format)
		return response, nil
	} else {
		// For other formats, just return the formatted data
		return FormatTags(tags, format), nil
	}
}

func handleGetTag(ctx context.Context, store store.Store, params map[string]string) (string, error) {
	tagID, ok := params["tagId"]
	if !ok || tagID == "" {
		return "Please provide a tag ID to retrieve.", nil
	}
	
	// Get the tag by ID
	tag, err := store.TagStore().GetTagByID(ctx, tagID)
	if err != nil {
		return fmt.Sprintf("Error retrieving tag: %v", err), nil
	}
	
	if tag == nil {
		return fmt.Sprintf("Tag with ID '%s' not found.", tagID), nil
	}
	
	// Get the format specified by the AI
	format := GetFormatFromParams(params)
	
	// Format the tag details using the formatting function
	return TagAsFormatted(tag, format), nil
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
		return "No context types found.", nil
	}

	response := "Available context types:\n"
	for _, contextType := range contextTypes {
		response += fmt.Sprintf("- %s\n", contextType)
	}

	return response, nil
}

func handleUserCount(ctx context.Context, store store.Store, params map[string]string) (string, error) {
	// Debug: Log the incoming parameters
	fmt.Printf("DEBUG: handleUserCount called with params: %+v\n", params)
	
	// Use standard time package for date handling
	now := time.Now()
	
	// Handle relative date expressions more intelligently
	if month, hasMonth := params["month"]; hasMonth && month != "" {
		var monthInt int
		if _, err := fmt.Sscanf(month, "%d", &monthInt); err == nil {
			// Check if this looks like a relative date (current month should be close)
			currentMonth := int(now.Month())
			if math.Abs(float64(currentMonth - monthInt)) > 1 {
				// Month is significantly off, likely cached data - use current month
				params["month"] = fmt.Sprintf("%d", currentMonth)
				fmt.Printf("DEBUG: Correcting month from %s to current month %d\n", month, currentMonth)
			}
		}
	}
	
	// Handle year more intelligently
	if year, hasYear := params["year"]; hasYear {
		var yearInt int
		if _, err := fmt.Sscanf(year, "%d", &yearInt); err == nil {
			currentYear := now.Year()
			// If year is more than 1 year old, it's probably cached data
			if currentYear - yearInt > 1 {
				params["year"] = fmt.Sprintf("%d", currentYear)
				fmt.Printf("DEBUG: Overriding outdated year %d to current year %d\n", yearInt, currentYear)
			}
		}
	} else if _, hasMonth := params["month"]; hasMonth {
		// No year provided but month is, add current year
		params["year"] = fmt.Sprintf("%d", now.Year())
		fmt.Printf("DEBUG: Added current year %d\n", now.Year())
	}

	fmt.Printf("DEBUG: Final params after processing: %+v\n", params)
	
	count, err := store.UserStore().UserCount(ctx, params)
	if err != nil {
		return "", err
	}

	// Build a descriptive message based on the filters used
	var filterDesc string
	
	hasSinceFilter := false
	hasUntilFilter := false
	hasDaysFilter := false
	hasMonthsFilter := false
	hasYearsFilter := false
	hasMonthFilter := false
	hasYearFilter := false
	
	if since, ok := params["since"]; ok && since != "" {
		hasSinceFilter = true
	}
	
	if until, ok := params["until"]; ok && until != "" {
		hasUntilFilter = true
	}

	if days, ok := params["days"]; ok && days != "" {
		hasDaysFilter = true
	}

	if months, ok := params["months"]; ok && months != "" {
		hasMonthsFilter = true
	}

	if years, ok := params["years"]; ok && years != "" {
		hasYearsFilter = true
	}

	if month, ok := params["month"]; ok && month != "" {
		hasMonthFilter = true
	}

	if year, ok := params["year"]; ok && year != "" {
		hasYearFilter = true
	}
	
	// Build filter description with all possible combinations
	if hasSinceFilter && hasUntilFilter {
		filterDesc = fmt.Sprintf(" created between %s and %s", params["since"], params["until"])
	} else if hasSinceFilter {
		filterDesc = fmt.Sprintf(" created since %s", params["since"])
	} else if hasUntilFilter {
		filterDesc = fmt.Sprintf(" created until %s", params["until"])
	} else if hasDaysFilter {
		filterDesc = fmt.Sprintf(" created in the last %s days", params["days"])
	} else if hasMonthsFilter {
		filterDesc = fmt.Sprintf(" created in the last %s months", params["months"])
	} else if hasYearsFilter {
		filterDesc = fmt.Sprintf(" created in the last %s years", params["years"])
	} else if hasMonthFilter && hasYearFilter {
		monthNames := map[string]string{
			"1": "January", "2": "February", "3": "March", "4": "April",
			"5": "May", "6": "June", "7": "July", "8": "August",
			"9": "September", "10": "October", "11": "November", "12": "December",
		}
		monthName := monthNames[params["month"]]
		if monthName == "" {
			monthName = "month " + params["month"]
		}
		filterDesc = fmt.Sprintf(" created in %s %s", monthName, params["year"])
	} else if hasMonthFilter {
		monthNames := map[string]string{
			"1": "January", "2": "February", "3": "March", "4": "April",
			"5": "May", "6": "June", "7": "July", "8": "August",
			"9": "September", "10": "October", "11": "November", "12": "December",
		}
		monthName := monthNames[params["month"]]
		if monthName == "" {
			monthName = "month " + params["month"]
		}
		filterDesc = fmt.Sprintf(" created in %s", monthName)
	} else if hasYearFilter {
		filterDesc = fmt.Sprintf(" created in %s", params["year"])
	} else {
		filterDesc = " in total"
	}

	result := fmt.Sprintf("You have %d users%s.", count, filterDesc)
	fmt.Printf("DEBUG: Returning result: %s\n", result)
	return result, nil
}

func handleUnknown(ctx context.Context, store store.Store, params map[string]string) (string, error) {
	return "I'm not sure how to help with that request. Could you please rephrase or ask about something else?", nil
}
