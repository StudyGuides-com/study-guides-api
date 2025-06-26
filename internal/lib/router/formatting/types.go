package formatting

import (
	"fmt"
)

// FormatType represents the different output formats available
type FormatType string

const (
	FormatList  FormatType = "list"
	FormatJSON  FormatType = "json"
	FormatCSV   FormatType = "csv"
	FormatTable FormatType = "table"
)

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

// BuildLimitMessage creates a message indicating when results are limited
func BuildLimitMessage(params map[string]string) string {
	if limitStr, ok := params["limit"]; ok && limitStr != "" {
		return fmt.Sprintf(" (limited to first %s results)", limitStr)
	}
	return ""
} 