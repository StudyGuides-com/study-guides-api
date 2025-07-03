package handlers

import (
	"context"
	"fmt"

	"github.com/studyguides-com/study-guides-api/internal/lib/router/formatting"
	"github.com/studyguides-com/study-guides-api/internal/store"
)

// HandleDeploy handles deployment requests
func HandleDeploy(ctx context.Context, store store.Store, params map[string]string) (string, error) {
	appID, ok := params["appId"]
	if !ok || appID == "" {
		return "", fmt.Errorf("appId is required")
	}

	force := false
	if forceStr, ok := params["force"]; ok && forceStr == "true" {
		force = true
	}

	deploymentID, err := store.DevopsStore().Deploy(ctx, appID)
	if err != nil {
		return "", fmt.Errorf("failed to deploy app %s: %w", appID, err)
	}

	if force {
		return fmt.Sprintf("✅ Forced deployment initiated for app %s\nDeployment ID: %s", appID, deploymentID), nil
	}
	return fmt.Sprintf("✅ Deployment initiated for app %s\nDeployment ID: %s", appID, deploymentID), nil
}

// HandleRollback handles rollback requests
func HandleRollback(ctx context.Context, store store.Store, params map[string]string) (string, error) {
	appID, ok := params["appId"]
	if !ok || appID == "" {
		return "", fmt.Errorf("appId is required")
	}

	deploymentID, ok := params["deploymentId"]
	if !ok || deploymentID == "" {
		return "", fmt.Errorf("deploymentId is required")
	}

	newDeploymentID, err := store.DevopsStore().Rollback(ctx, appID, deploymentID)
	if err != nil {
		return "", fmt.Errorf("failed to rollback app %s to deployment %s: %w", appID, deploymentID, err)
	}

	return fmt.Sprintf("✅ Rollback initiated for app %s\nRolling back to deployment: %s\nNew deployment ID: %s", appID, deploymentID, newDeploymentID), nil
}

// HandleListDeployments handles listing deployments
func HandleListDeployments(ctx context.Context, store store.Store, params map[string]string) (string, error) {
	appID, ok := params["appId"]
	if !ok || appID == "" {
		return "", fmt.Errorf("appId is required")
	}

	format := "list"
	if f, ok := params["format"]; ok {
		format = f
	}

	deployments, err := store.DevopsStore().ListDeployments(ctx, appID)
	if err != nil {
		return "", fmt.Errorf("failed to list deployments for app %s: %w", appID, err)
	}

	if len(deployments) == 0 {
		return fmt.Sprintf("No deployments found for app %s", appID), nil
	}

	switch format {
	case "json":
		return formatting.FormatDeploymentsAsJSON(deployments)
	case "csv":
		return formatting.FormatDeploymentsAsCSV(deployments)
	case "table":
		return formatting.FormatDeploymentsAsTable(deployments)
	default:
		return formatting.FormatDeploymentsAsList(deployments)
	}
}

// HandleGetDeploymentStatus handles getting deployment status
func HandleGetDeploymentStatus(ctx context.Context, store store.Store, params map[string]string) (string, error) {
	appID, ok := params["appId"]
	if !ok || appID == "" {
		return "", fmt.Errorf("appId is required")
	}

	deploymentID, ok := params["deploymentId"]
	if !ok || deploymentID == "" {
		return "", fmt.Errorf("deploymentId is required")
	}

	format := "list"
	if f, ok := params["format"]; ok {
		format = f
	}

	deployment, err := store.DevopsStore().GetDeploymentStatus(ctx, appID, deploymentID)
	if err != nil {
		return "", fmt.Errorf("failed to get deployment status for app %s deployment %s: %w", appID, deploymentID, err)
	}

	switch format {
	case "json":
		return formatting.FormatDeploymentAsJSON(deployment)
	case "csv":
		return formatting.FormatDeploymentAsCSV(deployment)
	case "table":
		return formatting.FormatDeploymentAsTable(deployment)
	default:
		return formatting.FormatDeploymentAsList(deployment)
	}
}
