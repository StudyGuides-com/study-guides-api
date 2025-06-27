package handlers

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/studyguides-com/study-guides-api/internal/store"
)

func HandleUserCount(ctx context.Context, store store.Store, params map[string]string) (string, error) {
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
			if math.Abs(float64(currentMonth-monthInt)) > 1 {
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
			if currentYear-yearInt > 1 {
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
