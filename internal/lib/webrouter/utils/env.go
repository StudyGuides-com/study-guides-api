package utils

import (
	"net/http"
	"os"
	"strings"
)

// EnvironmentData contains environment information for templates
type EnvironmentData struct {
	Environment string
	Version     string
	IsDev       bool
	IsProd      bool
	IsTest      bool
}

// detectEnvironmentFromHostname determines environment based on hostname
func detectEnvironmentFromHostname(hostname string) string {
	hostname = strings.ToLower(hostname)
	
	// Check for environment subdomains
	if strings.HasPrefix(hostname, "dev.") {
		return "dev"
	}
	if strings.HasPrefix(hostname, "test.") {
		return "test"
	}
	
	// Production is the main domain without subdomain (api.studyguides.com)
	if hostname == "api.studyguides.com" {
		return "prod"
	}
	
	// For any other domain that doesn't have a known subdomain, assume production
	return "prod"
}

// GetEnvironmentData returns environment information for templates
func GetEnvironmentData(r *http.Request) EnvironmentData {
	// Environment detection - prioritize explicit ENVIRONMENT setting
	env := os.Getenv("ENVIRONMENT")
	if env == "" && r != nil {
		// Detect from hostname if request is available
		env = detectEnvironmentFromHostname(r.Host)
	} else if env == "" {
		// Default to development for safety
		env = "development"
	}
	
	// Version detection
	version := os.Getenv("VERSION")
	if version == "" {
		version = "dev"
	}
	

	


	return EnvironmentData{
		Environment: env,
		Version:     version,
		IsDev:       env == "dev",
		IsProd:      env == "prod",
		IsTest:      env == "test",
	}
}

// MergeWithEnvData merges environment data with existing template data
func MergeWithEnvData(data map[string]interface{}, r *http.Request) map[string]interface{} {
	envData := GetEnvironmentData(r)
	
	// Merge environment data with existing data
	for key, value := range map[string]interface{}{
		"Environment": envData.Environment,
		"Version":     envData.Version,
		"IsDev":       envData.IsDev,
		"IsProd":      envData.IsProd,
		"IsTest":      envData.IsTest,
	} {
		data[key] = value
	}
	
	return data
} 