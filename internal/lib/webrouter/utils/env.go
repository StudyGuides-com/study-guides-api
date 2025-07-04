package utils

import (
	"net/http"
	"os"
	"strings"
	"time"
)

// EnvironmentData contains environment information for templates
type EnvironmentData struct {
	Environment string
	Version     string
	BuildTime   string
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
	
	// Version detection - prioritize explicit VERSION setting
	version := os.Getenv("VERSION")
	if version == "" {
		// Fall back to Digital Ocean version if available
		version = os.Getenv("DIGITALOCEAN_APP_VERSION")
		if version == "" {
			version = "dev"
		}
	}
	
	// Build time detection - prioritize explicit BUILD_TIME setting
	buildTime := os.Getenv("BUILD_TIME")
	if buildTime == "" {
		// Fall back to Digital Ocean build time if available
		buildTime = os.Getenv("DIGITALOCEAN_APP_BUILD_TIME")
		if buildTime == "" {
			buildTime = time.Now().Format("2006-01-02T15:04:05Z")
		}
	}
	


	return EnvironmentData{
		Environment: env,
		Version:     version,
		BuildTime:   buildTime,
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
		"BuildTime":   envData.BuildTime,
		"IsDev":       envData.IsDev,
		"IsProd":      envData.IsProd,
		"IsTest":      envData.IsTest,
	} {
		data[key] = value
	}
	
	return data
} 