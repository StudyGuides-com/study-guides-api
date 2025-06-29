package main

import (
	"os"
	"strconv"

	"golang.org/x/time/rate"
)

// parseEnvAsInt parses an environment variable as an integer with a fallback value
func parseEnvAsInt(key string, fallback int) int {
	valStr := os.Getenv(key)
	if valStr == "" {
		return fallback
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return fallback
	}
	return val
}

// parseEnvAsRate parses an environment variable as a rate limit with a fallback value
func parseEnvAsRate(key string, fallback rate.Limit) rate.Limit {
	valStr := os.Getenv(key)
	if valStr == "" {
		return fallback
	}
	val, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		return fallback
	}
	return rate.Limit(val)
}

// getPort returns the port from environment variable or defaults to 8080
func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
} 