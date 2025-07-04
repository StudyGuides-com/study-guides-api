package utils

import (
	"os"
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
}

// GetEnvironmentData returns environment information for templates
func GetEnvironmentData() EnvironmentData {
	// Environment detection - prioritize explicit ENVIRONMENT setting
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		// Only auto-detect production if explicitly on Digital Ocean App Platform
		if os.Getenv("DIGITALOCEAN_APP_PLATFORM") != "" {
			env = "production"
		} else {
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
		os.Getenv("DIGITALOCEAN_APP_DEPLOYMENT_ID") != ""

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
	}
}

// MergeWithEnvData merges environment data with existing template data
func MergeWithEnvData(data map[string]interface{}) map[string]interface{} {
	envData := GetEnvironmentData()
	
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
	} {
		data[key] = value
	}
	
	return data
} 