package tools

// AppNameMapping maps friendly names to DigitalOcean App Platform app IDs
var AppNameMapping = map[string]string{
	// Slackbot apps
	"test slackbot": "2d5a77bc-380d-49c6-bf9d-3d476fcee9dd",
	"dev slackbot":  "ffbeda27-f974-443d-a3d9-15e6079c0065",
	"prod slackbot": "a7679e54-612c-4a0d-9416-ebfac674063c",
	"slackbot test": "2d5a77bc-380d-49c6-bf9d-3d476fcee9dd",
	"slackbot dev":  "ffbeda27-f974-443d-a3d9-15e6079c0065",
	"slackbot prod": "a7679e54-612c-4a0d-9416-ebfac674063c",

	// API apps
	"test api": "7090c048-e409-43d0-965d-ded47c6eb63c",
	"dev api":  "9ec95f78-85c7-41ce-9b7e-972ad2314bf6",
	"prod api": "d3f9c599-6a21-47b4-be7d-f1db89717142",
	"api test": "7090c048-e409-43d0-965d-ded47c6eb63c",
	"api dev":  "9ec95f78-85c7-41ce-9b7e-972ad2314bf6",
	"api prod": "d3f9c599-6a21-47b4-be7d-f1db89717142",

	// Web apps
	"test web": "f335086c-e4a1-4cc8-96d0-5e2602989f7e",
	"dev web":  "ce8ec6a4-1e53-4fe0-b344-1fa85e52c90f",
	"web test": "f335086c-e4a1-4cc8-96d0-5e2602989f7e",
	"web dev":  "ce8ec6a4-1e53-4fe0-b344-1fa85e52c90f",
}

// ResolveAppName attempts to resolve a friendly name to an app ID
func ResolveAppName(name string) string {
	if appID, exists := AppNameMapping[name]; exists {
		return appID
	}
	return name // Return original if no mapping found
}

// GetAvailableAppNames returns a list of available app names for help text
func GetAvailableAppNames() []string {
	names := make([]string, 0, len(AppNameMapping))
	for name := range AppNameMapping {
		names = append(names, name)
	}
	return names
} 