package handlers

import (
	"context"
	"fmt"

	"github.com/studyguides-com/study-guides-api/internal/lib/router/formatting"
	"github.com/studyguides-com/study-guides-api/internal/lib/tools"
	"github.com/studyguides-com/study-guides-api/internal/store"
)

// HandleDeploy handles deployment requests
func HandleDeploy(ctx context.Context, store store.Store, params map[string]string) (string, error) {
	appID, ok := params["appId"]
	if !ok || appID == "" {
		return formatting.NewSingleResponse(nil, "appId is required").ToJSON(), nil
	}

	// Resolve app name to app ID if it's a friendly name
	resolvedAppID := tools.ResolveAppName(appID)
	originalAppID := appID
	appID = resolvedAppID

	force := false
	if forceStr, ok := params["force"]; ok && forceStr == "true" {
		force = true
	}

	deploymentID, err := store.DevopsStore().Deploy(ctx, appID)
	if err != nil {
		return formatting.NewSingleResponse(nil, fmt.Sprintf("failed to deploy app %s: %v", appID, err)).ToJSON(), nil
	}

	// Show the friendly name in the response if it was resolved
	displayName := originalAppID
	if originalAppID != appID {
		displayName = fmt.Sprintf("%s (%s)", originalAppID, appID)
	}

	// Create deployment data structure
	deploymentData := map[string]interface{}{
		"deploymentId": deploymentID,
		"appId":        appID,
		"appName":      originalAppID,
		"force":        force,
	}

	msg := ""
	if force {
		msg = fmt.Sprintf("✅ Forced deployment initiated for app %s\nDeployment ID: %s", displayName, deploymentID)
	} else {
		msg = fmt.Sprintf("✅ Deployment initiated for app %s\nDeployment ID: %s", displayName, deploymentID)
	}
	return formatting.NewSingleResponse(deploymentData, msg).ToJSON(), nil
}

// HandleRollback handles rollback requests
func HandleRollback(ctx context.Context, store store.Store, params map[string]string) (string, error) {
	appID, ok := params["appId"]
	if !ok || appID == "" {
		return formatting.NewSingleResponse(nil, "appId is required").ToJSON(), nil
	}

	// Resolve app name to app ID if it's a friendly name
	resolvedAppID := tools.ResolveAppName(appID)
	originalAppID := appID
	appID = resolvedAppID

	deploymentID, ok := params["deploymentId"]
	if !ok || deploymentID == "" {
		return formatting.NewSingleResponse(nil, "deploymentId is required").ToJSON(), nil
	}

	newDeploymentID, err := store.DevopsStore().Rollback(ctx, appID, deploymentID)
	if err != nil {
		return formatting.NewSingleResponse(nil, fmt.Sprintf("failed to rollback app %s to deployment %s: %v", appID, deploymentID, err)).ToJSON(), nil
	}

	// Show the friendly name in the response if it was resolved
	displayName := originalAppID
	if originalAppID != appID {
		displayName = fmt.Sprintf("%s (%s)", originalAppID, appID)
	}

	// Create rollback data structure
	rollbackData := map[string]interface{}{
		"newDeploymentId": newDeploymentID,
		"originalDeploymentId": deploymentID,
		"appId":        appID,
		"appName":      originalAppID,
	}

	msg := fmt.Sprintf("✅ Rollback initiated for app %s\nRolling back to deployment: %s\nNew deployment ID: %s", displayName, deploymentID, newDeploymentID)
	return formatting.NewSingleResponse(rollbackData, msg).ToJSON(), nil
}

// HandleListDeployments handles listing deployments
func HandleListDeployments(ctx context.Context, store store.Store, params map[string]string) (string, error) {
	appID, ok := params["appId"]
	if !ok || appID == "" {
		return formatting.NewSingleResponse(nil, "appId is required").ToJSON(), nil
	}

	// Resolve app name to app ID if it's a friendly name
	resolvedAppID := tools.ResolveAppName(appID)
	originalAppID := appID
	appID = resolvedAppID

	format := "list"
	if f, ok := params["format"]; ok {
		format = f
	}

	deployments, err := store.DevopsStore().ListDeployments(ctx, appID)
	if err != nil {
		return formatting.NewSingleResponse(nil, fmt.Sprintf("failed to list deployments for app %s: %v", appID, err)).ToJSON(), nil
	}

	if len(deployments) == 0 {
		// Show the friendly name in the response if it was resolved
		displayName := originalAppID
		if originalAppID != appID {
			displayName = fmt.Sprintf("%s (%s)", originalAppID, appID)
		}
		return formatting.NewSingleResponse(nil, fmt.Sprintf("No deployments found for app %s", displayName)).ToJSON(), nil
	}

	// Create deployments data structure
	deploymentsData := map[string]interface{}{
		"deployments": deployments,
		"count":       len(deployments),
		"appId":       appID,
		"appName":     originalAppID,
		"format":      format,
	}

	switch format {
	case "json":
		return formatting.NewSingleResponse(deploymentsData, "Deployments as JSON").ToJSON(), nil
	case "csv":
		return formatting.NewSingleResponse(deploymentsData, "Deployments as CSV").ToJSON(), nil
	default:
		return formatting.NewSingleResponse(deploymentsData, "Deployments as list").ToJSON(), nil
	}
}

// HandleGetDeploymentStatus handles getting deployment status
func HandleGetDeploymentStatus(ctx context.Context, store store.Store, params map[string]string) (string, error) {
	appID, ok := params["appId"]
	if !ok || appID == "" {
		return formatting.NewSingleResponse(nil, "appId is required").ToJSON(), nil
	}

	// Resolve app name to app ID if it's a friendly name
	resolvedAppID := tools.ResolveAppName(appID)
	appID = resolvedAppID

	deploymentID, ok := params["deploymentId"]
	if !ok || deploymentID == "" {
		return formatting.NewSingleResponse(nil, "deploymentId is required").ToJSON(), nil
	}

	format := "list"
	if f, ok := params["format"]; ok {
		format = f
	}

	deployment, err := store.DevopsStore().GetDeploymentStatus(ctx, appID, deploymentID)
	if err != nil {
		return formatting.NewSingleResponse(nil, fmt.Sprintf("failed to get deployment status for app %s deployment %s: %v", appID, deploymentID, err)).ToJSON(), nil
	}

	// Create deployment status data structure
	deploymentStatusData := map[string]interface{}{
		"deployment": deployment,
		"appId":      appID,
		"format":     format,
	}

	switch format {
	case "json":
		return formatting.NewSingleResponse(deploymentStatusData, "Deployment as JSON").ToJSON(), nil
	case "csv":
		return formatting.NewSingleResponse(deploymentStatusData, "Deployment as CSV").ToJSON(), nil
	default:
		return formatting.NewSingleResponse(deploymentStatusData, "Deployment as list").ToJSON(), nil
	}
}
