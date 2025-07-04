# Corrected Environment Detection Fix - Complete Diff

## Problem

All environments were showing as "Dev" on Digital Ocean deployments, even though the application was running in production. Specifically, both `dev.api.studyguides.com` and `test.api.studyguides.com` were showing "🚧 Development Environment" when they should show different environments based on their subdomains.

## Root Cause

The environment detection logic was completely wrong. It was only checking for Digital Ocean environment variables and defaulting to "development" when those weren't present, instead of properly detecting the environment from the hostname/subdomain or explicit environment configuration.

## Changes Made

### 1. Added Hostname-Based Environment Detection

**File: `internal/lib/webrouter/utils/env.go`**

```diff
+ import (
+ 	"log"
+ 	"net/http"
+ 	"os"
+ 	"strings"
+ 	"time"
+ )

+ // detectEnvironmentFromHostname determines environment based on hostname
+ func detectEnvironmentFromHostname(hostname string) string {
+ 	hostname = strings.ToLower(hostname)
+
+ 	// Check for common environment subdomains
+ 	if strings.HasPrefix(hostname, "dev.") || strings.HasPrefix(hostname, "development.") {
+ 		return "development"
+ 	}
+ 	if strings.HasPrefix(hostname, "test.") || strings.HasPrefix(hostname, "staging.") {
+ 		return "staging"
+ 	}
+ 	if strings.HasPrefix(hostname, "prod.") || strings.HasPrefix(hostname, "production.") {
+ 		return "production"
+ 	}
+
+ 	// If no subdomain or unknown subdomain, default to development
+ 	return "development"
+ }

- func GetEnvironmentData() EnvironmentData {
+ func GetEnvironmentData(r *http.Request) EnvironmentData {
```

### 2. Fixed Environment Detection Logic

```diff
	// Environment detection - prioritize explicit ENVIRONMENT setting
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
-		// Auto-detect environment based on Digital Ocean variables
-		if os.Getenv("DIGITALOCEAN_APP_PLATFORM") != "" {
-			// On Digital Ocean App Platform, check for specific environment indicators
-			if os.Getenv("DIGITALOCEAN_APP_ENV") != "" {
-				env = os.Getenv("DIGITALOCEAN_APP_ENV")
-			} else {
-				// Default to production on Digital Ocean unless explicitly set otherwise
-				env = "production"
-			}
-		} else {
-			env = "development"
-		}
+		// Check for Digital Ocean specific environment variable
+		if os.Getenv("DIGITALOCEAN_APP_ENV") != "" {
+			env = os.Getenv("DIGITALOCEAN_APP_ENV")
+		} else if r != nil {
+			// Detect from hostname if request is available
+			env = detectEnvironmentFromHostname(r.Host)
+		} else {
+			// Default to development for safety
+			env = "development"
+		}
	}
```

### 3. Enhanced Debug Logging

```diff
	// Debug: Log environment detection for troubleshooting
	log.Printf("Environment Detection Debug:")
	log.Printf("  ENVIRONMENT: %s", os.Getenv("ENVIRONMENT"))
+	log.Printf("  DIGITALOCEAN_APP_ENV: %s", os.Getenv("DIGITALOCEAN_APP_ENV"))
+	if r != nil {
+		log.Printf("  Hostname: %s", r.Host)
+		log.Printf("  Detected from hostname: %s", detectEnvironmentFromHostname(r.Host))
+	}
	log.Printf("  DIGITALOCEAN_APP_PLATFORM: %s", os.Getenv("DIGITALOCEAN_APP_PLATFORM"))
	log.Printf("  DIGITALOCEAN_APP_NAME: %s", os.Getenv("DIGITALOCEAN_APP_NAME"))
	log.Printf("  DIGITALOCEAN_APP_DEPLOYMENT_ID: %s", os.Getenv("DIGITALOCEAN_APP_DEPLOYMENT_ID"))
	log.Printf("  DIGITALOCEAN_APP_REGION: %s", os.Getenv("DIGITALOCEAN_APP_REGION"))
	log.Printf("  Final Environment: %s", env)
	log.Printf("  Is Digital Ocean: %t", isDigitalOcean)
```

### 4. Updated Function Signatures

```diff
- func MergeWithEnvData(data map[string]interface{}) map[string]interface{} {
- 	envData := GetEnvironmentData()
+ func MergeWithEnvData(data map[string]interface{}, r *http.Request) map[string]interface{} {
+ 	envData := GetEnvironmentData(r)
```

### 5. Updated Route Handlers

**File: `internal/lib/webrouter/routes/home.go`**

```diff
	// Add environment data
-	data = utils.MergeWithEnvData(data)
+	data = utils.MergeWithEnvData(data, r)
```

**File: `internal/lib/webrouter/routes/notfound.go`**

```diff
	// Add environment data
-	data = utils.MergeWithEnvData(data)
+	data = utils.MergeWithEnvData(data, r)
```

## Summary of Changes

1. **Hostname-Based Detection**: Added function to detect environment from subdomain (dev, test, staging, prod)
2. **Proper Priority Order**:
   - First: Explicit `ENVIRONMENT` variable
   - Second: `DIGITALOCEAN_APP_ENV` variable
   - Third: Hostname detection
   - Last: Default to "development"
3. **Request Context**: Functions now accept HTTP request to access hostname
4. **Enhanced Debugging**: Logs show hostname detection process
5. **Updated All Handlers**: All route handlers now pass request to environment functions

## Expected Results

- **dev.api.studyguides.com**: Shows "🚧 Development Environment" badge
- **test.api.studyguides.com**: Shows "🧪 Staging Environment" badge
- **prod.api.studyguides.com**: Shows "✅ Production Environment" badge
- **api.studyguides.com**: Shows "🚧 Development Environment" badge (default)
- **Local Development**: Shows "🚧 Development Environment" badge

## Environment Detection Priority

1. **Explicit ENVIRONMENT variable** (highest priority)
2. **DIGITALOCEAN_APP_ENV variable** (Digital Ocean specific)
3. **Hostname detection** (dev/test/staging/prod subdomains)
4. **Default to development** (safety fallback)

## Testing

1. Deploy to different subdomains and verify correct environment badges
2. Check server logs for hostname detection debug information
3. Test with explicit ENVIRONMENT variable to override hostname detection
4. Verify local development still shows as "Development"
