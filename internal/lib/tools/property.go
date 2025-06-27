package tools

// PropertyValue represents a single property value with type and description
type PropertyValue struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

// NewPropertyValue creates a new PropertyValue
func NewPropertyValue(propertyType, description string) PropertyValue {
	return PropertyValue{
		Type:        propertyType,
		Description: description,
	}
}

// Property defines a single parameter property
type Property struct {
	Name  string `json:"name"`
	Value PropertyValue
}

// NewProperty creates a new Property
func NewProperty(name, propertyType, description string) Property {
	return Property{
		Name:  name,
		Value: NewPropertyValue(propertyType, description),
	}
}

// ToMapEntry returns the property as a map entry for easy use in Properties maps
func (p Property) ToMapEntry() (string, PropertyValue) {
	return p.Name, p.Value
}

type Properties map[string]PropertyValue

// BuildProperties creates a Properties map from a slice of Property structs
func BuildProperties(props ...Property) Properties {
	result := make(Properties)
	for _, prop := range props {
		result[prop.Name] = prop.Value
	}
	return result
}
