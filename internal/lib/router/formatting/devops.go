package formatting

import (
	"encoding/json"
	"fmt"
	"strings"

	devopspb "github.com/studyguides-com/study-guides-api/api/v1/devops"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// FormatDeploymentsAsJSON formats deployments as JSON
func FormatDeploymentsAsJSON(deployments []devopspb.Deployment) (string, error) {
	jsonBytes, err := json.MarshalIndent(deployments, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

// FormatDeploymentsAsCSV formats deployments as CSV
func FormatDeploymentsAsCSV(deployments []devopspb.Deployment) (string, error) {
	if len(deployments) == 0 {
		return "Deployment ID,App ID,Status,Created At", nil
	}

	var lines []string
	lines = append(lines, "Deployment ID,App ID,Status,Created At")

	for i := range deployments {
		deployment := &deployments[i]
		createdAt := formatTimestamp(deployment.CreatedAt)
		line := fmt.Sprintf("%s,%s,%s,%s",
			deployment.DeploymentId,
			deployment.AppId,
			deployment.Status.String(),
			createdAt,
		)
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n"), nil
}

// FormatDeploymentsAsTable formats deployments as a markdown table
func FormatDeploymentsAsTable(deployments []devopspb.Deployment) (string, error) {
	if len(deployments) == 0 {
		return "| Deployment ID | App ID | Status | Created At |\n|---------------|--------|--------|------------|", nil
	}

	var lines []string
	lines = append(lines, "| Deployment ID | App ID | Status | Created At |")
	lines = append(lines, "|---------------|--------|--------|------------|")

	for i := range deployments {
		deployment := &deployments[i]
		createdAt := formatTimestamp(deployment.CreatedAt)
		line := fmt.Sprintf("| %s | %s | %s | %s |",
			deployment.DeploymentId,
			deployment.AppId,
			deployment.Status.String(),
			createdAt,
		)
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n"), nil
}

// FormatDeploymentsAsList formats deployments as a human-readable list
func FormatDeploymentsAsList(deployments []devopspb.Deployment) (string, error) {
	if len(deployments) == 0 {
		return "No deployments found.", nil
	}

	var lines []string
	lines = append(lines, fmt.Sprintf("Found %d deployment(s):", len(deployments)))
	lines = append(lines, "")

	for i := range deployments {
		deployment := &deployments[i]
		createdAt := formatTimestamp(deployment.CreatedAt)
		line := fmt.Sprintf("%d. **Deployment ID:** %s\n   **App ID:** %s\n   **Status:** %s\n   **Created:** %s",
			i+1,
			deployment.DeploymentId,
			deployment.AppId,
			deployment.Status.String(),
			createdAt,
		)
		lines = append(lines, line)
		if i < len(deployments)-1 {
			lines = append(lines, "")
		}
	}

	return strings.Join(lines, "\n"), nil
}

// FormatDeploymentAsJSON formats a single deployment as JSON
func FormatDeploymentAsJSON(deployment devopspb.Deployment) (string, error) {
	jsonBytes, err := json.MarshalIndent(deployment, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

// FormatDeploymentAsCSV formats a single deployment as CSV
func FormatDeploymentAsCSV(deployment devopspb.Deployment) (string, error) {
	createdAt := formatTimestamp(deployment.CreatedAt)
	return fmt.Sprintf("Deployment ID,App ID,Status,Created At\n%s,%s,%s,%s",
		deployment.DeploymentId,
		deployment.AppId,
		deployment.Status.String(),
		createdAt,
	), nil
}

// FormatDeploymentAsTable formats a single deployment as a markdown table
func FormatDeploymentAsTable(deployment devopspb.Deployment) (string, error) {
	createdAt := formatTimestamp(deployment.CreatedAt)
	return fmt.Sprintf("| Deployment ID | App ID | Status | Created At |\n|---------------|--------|--------|------------|\n| %s | %s | %s | %s |",
		deployment.DeploymentId,
		deployment.AppId,
		deployment.Status.String(),
		createdAt,
	), nil
}

// FormatDeploymentAsList formats a single deployment as human-readable text
func FormatDeploymentAsList(deployment devopspb.Deployment) (string, error) {
	createdAt := formatTimestamp(deployment.CreatedAt)
	return fmt.Sprintf("**Deployment Details:**\n\n**Deployment ID:** %s\n**App ID:** %s\n**Status:** %s\n**Created:** %s",
		deployment.DeploymentId,
		deployment.AppId,
		deployment.Status.String(),
		createdAt,
	), nil
}

// formatTimestamp formats a protobuf timestamp for display
func formatTimestamp(ts *timestamppb.Timestamp) string {
	if ts == nil {
		return "Unknown"
	}
	t := ts.AsTime()
	return t.Format("2006-01-02 15:04:05 UTC")
}
