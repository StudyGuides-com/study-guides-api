package search

// convertToStringSlice converts an interface{} to []string
func convertToStringSlice(v interface{}) []string {
	if v == nil {
		return nil
	}
	
	slice, ok := v.([]interface{})
	if !ok {
		return nil
	}

	result := make([]string, 0, len(slice))
	for _, item := range slice {
		if str, ok := item.(string); ok {
			result = append(result, str)
		}
	}
	return result
}