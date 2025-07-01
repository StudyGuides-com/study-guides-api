package formatting

import (
	"encoding/json"
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

// ResponseType represents the type of data being returned
type ResponseType string

const (
	ResponseTypeCount  ResponseType = "count"
	ResponseTypeList   ResponseType = "list"
	ResponseTypeSingle ResponseType = "single"
	ResponseTypeCSV    ResponseType = "csv"
	ResponseTypeTable  ResponseType = "table"
)

// APIResponse is the universal wrapper for all API responses
type APIResponse struct {
	Type        ResponseType         `json:"type"`
	Data        interface{}          `json:"data"`
	Message     string               `json:"message,omitempty"`
	ContentType string               `json:"content_type"`
	Filters     map[string]string    `json:"filters,omitempty"`
	Pagination  *PaginationInfo      `json:"pagination,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// PaginationInfo contains pagination details
type PaginationInfo struct {
	Page     int `json:"page"`
	Limit    int `json:"limit"`
	Total    int `json:"total"`
	Pages    int `json:"pages"`
	HasNext  bool `json:"has_next"`
	HasPrev  bool `json:"has_prev"`
}

// Formatter interface for different data types
type Formatter interface {
	Format(format FormatType) interface{}
}

// NewCountResponse creates a new count response
func NewCountResponse(count int64, message string, filters map[string]string) *APIResponse {
	return &APIResponse{
		Type:        ResponseTypeCount,
		Data:        count,
		Message:     message,
		ContentType: "application/json",
		Filters:     filters,
	}
}

// NewCountResponseInt creates a new count response for int values
func NewCountResponseInt(count int, message string, filters map[string]string) *APIResponse {
	return &APIResponse{
		Type:        ResponseTypeCount,
		Data:        count,
		Message:     message,
		ContentType: "application/json",
		Filters:     filters,
	}
}

// NewListResponse creates a new list response
func NewListResponse(data interface{}, message string, filters map[string]string, pagination *PaginationInfo) *APIResponse {
	contentType := "application/json"
	if str, ok := data.(string); ok {
		if len(str) > 0 {
			contentType = "text/plain"
		} else {
			// Empty string should still be JSON
			contentType = "application/json"
		}
	}
	
	return &APIResponse{
		Type:        ResponseTypeList,
		Data:        data,
		Message:     message,
		ContentType: contentType,
		Filters:     filters,
		Pagination:  pagination,
	}
}

// NewSingleResponse creates a new single item response
func NewSingleResponse(item interface{}, message string) *APIResponse {
	return &APIResponse{
		Type:        ResponseTypeSingle,
		Data:        item,
		Message:     message,
		ContentType: "application/json",
	}
}

// NewCSVResponse creates a new CSV response
func NewCSVResponse(csvData string, message string, filters map[string]string) *APIResponse {
	return &APIResponse{
		Type:        ResponseTypeCSV,
		Data:        csvData,
		Message:     message,
		ContentType: "text/csv",
		Filters:     filters,
	}
}

// NewTableResponse creates a new table response
func NewTableResponse(tableData string, message string, filters map[string]string) *APIResponse {
	return &APIResponse{
		Type:        ResponseTypeTable,
		Data:        tableData,
		Message:     message,
		ContentType: "text/plain",
		Filters:     filters,
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

// BuildLimitMessage creates a message indicating when results are limited
func BuildLimitMessage(params map[string]string) string {
	if limitStr, ok := params["limit"]; ok && limitStr != "" {
		return fmt.Sprintf(" (limited to first %s results)", limitStr)
	}
	return ""
}

// ToJSON converts an APIResponse to a JSON string
func (r *APIResponse) ToJSON() string {
	jsonBytes, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error formatting as JSON: %v", err)
	}
	return string(jsonBytes)
}

// BuildCountMessage creates a consistent message for count responses
func BuildCountMessage(count int64, itemType string, params map[string]string, hasTypeFilter, hasContextFilter, hasNameFilter, hasPublicFilter bool) string {
	if hasTypeFilter && hasContextFilter && hasNameFilter && hasPublicFilter {
		return fmt.Sprintf("Found %d %s for type '%s', context '%s', name containing '%s', and %s.", count, itemType, params["type"], params["contextType"], params["name"], GetPublicDescription(params["public"]))
	} else if hasTypeFilter && hasContextFilter && hasNameFilter {
		return fmt.Sprintf("Found %d %s for type '%s', context '%s', and name containing '%s'.", count, itemType, params["type"], params["contextType"], params["name"])
	} else if hasTypeFilter && hasContextFilter && hasPublicFilter {
		return fmt.Sprintf("Found %d %s for type '%s', context '%s', and %s.", count, itemType, params["type"], params["contextType"], GetPublicDescription(params["public"]))
	} else if hasTypeFilter && hasNameFilter && hasPublicFilter {
		return fmt.Sprintf("Found %d %s for type '%s', name containing '%s', and %s.", count, itemType, params["type"], params["name"], GetPublicDescription(params["public"]))
	} else if hasContextFilter && hasNameFilter && hasPublicFilter {
		return fmt.Sprintf("Found %d %s for context '%s', name containing '%s', and %s.", count, itemType, params["contextType"], params["name"], GetPublicDescription(params["public"]))
	} else if hasTypeFilter && hasContextFilter {
		return fmt.Sprintf("Found %d %s for type '%s' and context '%s'.", count, itemType, params["type"], params["contextType"])
	} else if hasTypeFilter && hasNameFilter {
		return fmt.Sprintf("Found %d %s for type '%s' and name containing '%s'.", count, itemType, params["type"], params["name"])
	} else if hasTypeFilter && hasPublicFilter {
		return fmt.Sprintf("Found %d %s for type '%s' and %s.", count, itemType, params["type"], GetPublicDescription(params["public"]))
	} else if hasContextFilter && hasNameFilter {
		return fmt.Sprintf("Found %d %s for context '%s' and name containing '%s'.", count, itemType, params["contextType"], params["name"])
	} else if hasContextFilter && hasPublicFilter {
		return fmt.Sprintf("Found %d %s for context '%s' and %s.", count, itemType, params["contextType"], GetPublicDescription(params["public"]))
	} else if hasNameFilter && hasPublicFilter {
		return fmt.Sprintf("Found %d %s with name containing '%s' and %s.", count, itemType, params["name"], GetPublicDescription(params["public"]))
	} else if hasTypeFilter {
		return fmt.Sprintf("Found %d %s for type '%s'.", count, itemType, params["type"])
	} else if hasContextFilter {
		return fmt.Sprintf("Found %d %s for context '%s'.", count, itemType, params["contextType"])
	} else if hasNameFilter {
		return fmt.Sprintf("Found %d %s with name containing '%s'.", count, itemType, params["name"])
	} else if hasPublicFilter {
		return fmt.Sprintf("Found %d %s that are %s.", count, itemType, GetPublicDescription(params["public"]))
	}
	return fmt.Sprintf("Found %d %s total.", count, itemType)
}

// BuildCountMessageInt creates a consistent message for count responses with int values
func BuildCountMessageInt(count int, itemType string, params map[string]string, hasTypeFilter, hasContextFilter, hasNameFilter, hasPublicFilter bool) string {
	if hasTypeFilter && hasContextFilter && hasNameFilter && hasPublicFilter {
		return fmt.Sprintf("Found %d %s for type '%s', context '%s', name containing '%s', and %s.", count, itemType, params["type"], params["contextType"], params["name"], GetPublicDescription(params["public"]))
	} else if hasTypeFilter && hasContextFilter && hasNameFilter {
		return fmt.Sprintf("Found %d %s for type '%s', context '%s', and name containing '%s'.", count, itemType, params["type"], params["contextType"], params["name"])
	} else if hasTypeFilter && hasContextFilter && hasPublicFilter {
		return fmt.Sprintf("Found %d %s for type '%s', context '%s', and %s.", count, itemType, params["type"], params["contextType"], GetPublicDescription(params["public"]))
	} else if hasTypeFilter && hasNameFilter && hasPublicFilter {
		return fmt.Sprintf("Found %d %s for type '%s', name containing '%s', and %s.", count, itemType, params["type"], params["name"], GetPublicDescription(params["public"]))
	} else if hasContextFilter && hasNameFilter && hasPublicFilter {
		return fmt.Sprintf("Found %d %s for context '%s', name containing '%s', and %s.", count, itemType, params["contextType"], params["name"], GetPublicDescription(params["public"]))
	} else if hasTypeFilter && hasContextFilter {
		return fmt.Sprintf("Found %d %s for type '%s' and context '%s'.", count, itemType, params["type"], params["contextType"])
	} else if hasTypeFilter && hasNameFilter {
		return fmt.Sprintf("Found %d %s for type '%s' and name containing '%s'.", count, itemType, params["type"], params["name"])
	} else if hasTypeFilter && hasPublicFilter {
		return fmt.Sprintf("Found %d %s for type '%s' and %s.", count, itemType, params["type"], GetPublicDescription(params["public"]))
	} else if hasContextFilter && hasNameFilter {
		return fmt.Sprintf("Found %d %s for context '%s' and name containing '%s'.", count, itemType, params["contextType"], params["name"])
	} else if hasContextFilter && hasPublicFilter {
		return fmt.Sprintf("Found %d %s for context '%s' and %s.", count, itemType, params["contextType"], GetPublicDescription(params["public"]))
	} else if hasNameFilter && hasPublicFilter {
		return fmt.Sprintf("Found %d %s with name containing '%s' and %s.", count, itemType, params["name"], GetPublicDescription(params["public"]))
	} else if hasTypeFilter {
		return fmt.Sprintf("Found %d %s for type '%s'.", count, itemType, params["type"])
	} else if hasContextFilter {
		return fmt.Sprintf("Found %d %s for context '%s'.", count, itemType, params["contextType"])
	} else if hasNameFilter {
		return fmt.Sprintf("Found %d %s with name containing '%s'.", count, itemType, params["name"])
	} else if hasPublicFilter {
		return fmt.Sprintf("Found %d %s that are %s.", count, itemType, GetPublicDescription(params["public"]))
	}
	return fmt.Sprintf("Found %d %s total.", count, itemType)
}
