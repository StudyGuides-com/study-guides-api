package utils

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// EnvironmentData contains environment information for templates
type EnvironmentData struct {
	Environment    string
	Version        string
	BuildTime      string
	DeploymentID   string
	AppName        string
	Region         string
	IsDev          bool
	IsProd         bool
	IsStaging      bool
	IsDigitalOcean bool
	AllDOEnvVars   map[string]string
}

// detectEnvironmentFromHostname determines environment based on hostname
func detectEnvironmentFromHostname(hostname string) string {
	hostname = strings.ToLower(hostname)
	
	// Check for common environment subdomains
	if strings.HasPrefix(hostname, "dev.") || strings.HasPrefix(hostname, "development.") {
		return "development"
	}
	if strings.HasPrefix(hostname, "test.") || strings.HasPrefix(hostname, "staging.") {
		return "staging"
	}
	
	// Production is the main domain without subdomain (api.studyguides.com)
	// or explicitly prod subdomain
	if strings.HasPrefix(hostname, "prod.") || strings.HasPrefix(hostname, "production.") {
		return "production"
	}
	
	// If it's the main domain (api.studyguides.com), it's production
	if hostname == "api.studyguides.com" {
		return "production"
	}
	
	// For any other domain that doesn't have a known subdomain, assume production
	// This handles cases where the domain might be different or has other subdomains
	return "production"
}

// GetEnvironmentData returns environment information for templates
func GetEnvironmentData(r *http.Request) EnvironmentData {
	// Environment detection - prioritize explicit ENVIRONMENT setting
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		// Check for Digital Ocean specific environment variable
		if os.Getenv("DIGITALOCEAN_APP_ENV") != "" {
			env = os.Getenv("DIGITALOCEAN_APP_ENV")
		} else if r != nil {
			// Detect from hostname if request is available
			env = detectEnvironmentFromHostname(r.Host)
		} else {
			// Default to development for safety
			env = "development"
		}
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
	
	// Digital Ocean specific information (only available on DO)
	deploymentID := os.Getenv("DIGITALOCEAN_APP_DEPLOYMENT_ID")
	appName := os.Getenv("DIGITALOCEAN_APP_NAME")
	region := os.Getenv("DIGITALOCEAN_APP_REGION")
	
	// Check if running on Digital Ocean (any DO env var present)
	isDigitalOcean := os.Getenv("DIGITALOCEAN_APP_PLATFORM") != "" || 
		os.Getenv("DIGITALOCEAN_APP_NAME") != "" ||
		os.Getenv("DIGITALOCEAN_APP_DEPLOYMENT_ID") != "" ||
		os.Getenv("DIGITALOCEAN_APP_ENV") != "" ||
		os.Getenv("DIGITALOCEAN_APP_REGION") != ""

	// Collect all Digital Ocean environment variables for debugging
	allDOEnvVars := make(map[string]string)
	for _, envVar := range os.Environ() {
		if len(envVar) > 0 {
			parts := strings.SplitN(envVar, "=", 2)
			if len(parts) == 2 && strings.HasPrefix(parts[0], "DIGITALOCEAN_") {
				allDOEnvVars[parts[0]] = parts[1]
			}
		}
	}

	// Debug: Log environment detection for troubleshooting
	log.Printf("Environment Detection Debug:")
	log.Printf("  ENVIRONMENT: %s", os.Getenv("ENVIRONMENT"))
	log.Printf("  DIGITALOCEAN_APP_ENV: %s", os.Getenv("DIGITALOCEAN_APP_ENV"))
	if r != nil {
		log.Printf("  Hostname: %s", r.Host)
		log.Printf("  Detected from hostname: %s", detectEnvironmentFromHostname(r.Host))
		log.Printf("  Request URL: %s", r.URL.String())
		log.Printf("  Request Method: %s", r.Method)
	}
	log.Printf("  DIGITALOCEAN_APP_PLATFORM: %s", os.Getenv("DIGITALOCEAN_APP_PLATFORM"))
	log.Printf("  DIGITALOCEAN_APP_NAME: %s", os.Getenv("DIGITALOCEAN_APP_NAME"))
	log.Printf("  DIGITALOCEAN_APP_DEPLOYMENT_ID: %s", os.Getenv("DIGITALOCEAN_APP_DEPLOYMENT_ID"))
	log.Printf("  DIGITALOCEAN_APP_REGION: %s", os.Getenv("DIGITALOCEAN_APP_REGION"))
	log.Printf("  Final Environment: %s", env)
	log.Printf("  Is Digital Ocean: %t", isDigitalOcean)
	log.Printf("  All DO Env Vars Count: %d", len(allDOEnvVars))

	return EnvironmentData{
		Environment:    env,
		Version:        version,
		BuildTime:      buildTime,
		DeploymentID:   deploymentID,
		AppName:        appName,
		Region:         region,
		IsDev:          env == "development",
		IsProd:         env == "production",
		IsStaging:      env == "staging",
		IsDigitalOcean: isDigitalOcean,
		AllDOEnvVars:   allDOEnvVars,
	}
}

// MergeWithEnvData merges environment data with existing template data
func MergeWithEnvData(data map[string]interface{}, r *http.Request) map[string]interface{} {
	envData := GetEnvironmentData(r)
	
	// Merge environment data with existing data
	for key, value := range map[string]interface{}{
		"Environment":    envData.Environment,
		"Version":        envData.Version,
		"BuildTime":      envData.BuildTime,
		"DeploymentID":   envData.DeploymentID,
		"AppName":        envData.AppName,
		"Region":         envData.Region,
		"IsDev":          envData.IsDev,
		"IsProd":         envData.IsProd,
		"IsStaging":      envData.IsStaging,
		"IsDigitalOcean": envData.IsDigitalOcean,
		"AllDOEnvVars":   envData.AllDOEnvVars,
	} {
		data[key] = value
	}
	
	return data
} 