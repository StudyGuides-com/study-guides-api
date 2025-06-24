package tools

// ParameterSchema defines a cleaner way to specify function parameters
type ParameterSchema struct {
	Type       string                   `json:"type"`
	Properties map[string]PropertyValue `json:"properties"`
	Required   []string                 `json:"required"`
}

// BuildParameterSchema creates a ParameterSchema from properties and required fields
func BuildParameterSchema(properties Properties, required []string) ParameterSchema {
	return ParameterSchema{
		Type:       "object",
		Properties: properties,
		Required:   required,
	}
}

// BuildParameterSchemaFromProps creates a ParameterSchema directly from Property structs and required fields
func BuildParameterSchemaFromProps(required []string, props ...Property) ParameterSchema {
	return BuildParameterSchema(BuildProperties(props...), required)
} 