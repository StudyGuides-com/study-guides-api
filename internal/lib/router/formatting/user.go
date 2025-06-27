package formatting

import (
	"encoding/json"
	"fmt"

	sharedpb "github.com/studyguides-com/study-guides-api/api/v1/shared"
)

// UserAsFormatted formats a single user according to the specified format
func UserAsFormatted(user *sharedpb.User, format FormatType) string {
	switch format {
	case FormatJSON:
		return UserAsJSON(user)
	case FormatCSV:
		return UserAsCSV(user)
	case FormatTable:
		return UserAsTable(user)
	case FormatList:
		fallthrough
	default:
		return UserAsDetailedText(user)
	}
}

// UserAsJSON formats a single user as JSON
func UserAsJSON(user *sharedpb.User) string {
	// Create a comprehensive structure for JSON output
	type UserDetailOutput struct {
		ID            string `json:"id"`
		Name          string `json:"name,omitempty"`
		GamerTag      string `json:"gamer_tag,omitempty"`
		Email         string `json:"email,omitempty"`
		EmailVerified string `json:"email_verified,omitempty"`
		Image         string `json:"image,omitempty"`
		ContentTagID  string `json:"content_tag_id,omitempty"`
	}

	output := UserDetailOutput{
		ID: user.Id,
	}

	if user.Name != nil && *user.Name != "" {
		output.Name = *user.Name
	}
	if user.GamerTag != nil && *user.GamerTag != "" {
		output.GamerTag = *user.GamerTag
	}
	if user.Email != nil && *user.Email != "" {
		output.Email = *user.Email
	}
	if user.EmailVerified != nil {
		output.EmailVerified = user.EmailVerified.AsTime().Format("2006-01-02 15:04:05")
	}
	if user.Image != nil && *user.Image != "" {
		output.Image = *user.Image
	}
	if user.ContentTagId != nil && *user.ContentTagId != "" {
		output.ContentTagID = *user.ContentTagId
	}

	jsonBytes, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error formatting as JSON: %v", err)
	}
	return string(jsonBytes)
}

// UserAsCSV formats a single user as CSV
func UserAsCSV(user *sharedpb.User) string {
	name := ""
	if user.Name != nil && *user.Name != "" {
		name = *user.Name
	}

	gamerTag := ""
	if user.GamerTag != nil && *user.GamerTag != "" {
		gamerTag = *user.GamerTag
	}

	email := ""
	if user.Email != nil && *user.Email != "" {
		email = *user.Email
	}

	emailVerified := ""
	if user.EmailVerified != nil {
		emailVerified = user.EmailVerified.AsTime().Format("2006-01-02 15:04:05")
	}

	image := ""
	if user.Image != nil && *user.Image != "" {
		image = *user.Image
	}

	contentTagID := ""
	if user.ContentTagId != nil && *user.ContentTagId != "" {
		contentTagID = *user.ContentTagId
	}

	header := "id,name,gamer_tag,email,email_verified,image,content_tag_id\n"
	row := fmt.Sprintf("\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\"\n",
		user.Id, name, gamerTag, email, emailVerified, image, contentTagID)

	return header + row
}

// UserAsTable formats a single user as a table
func UserAsTable(user *sharedpb.User) string {
	name := ""
	if user.Name != nil && *user.Name != "" {
		name = *user.Name
	}

	gamerTag := ""
	if user.GamerTag != nil && *user.GamerTag != "" {
		gamerTag = *user.GamerTag
	}

	email := ""
	if user.Email != nil && *user.Email != "" {
		email = *user.Email
	}

	emailVerified := ""
	if user.EmailVerified != nil {
		emailVerified = user.EmailVerified.AsTime().Format("2006-01-02 15:04:05")
	}

	image := ""
	if user.Image != nil && *user.Image != "" {
		image = *user.Image
	}

	contentTagID := ""
	if user.ContentTagId != nil && *user.ContentTagId != "" {
		contentTagID = *user.ContentTagId
	}

	response := "| Field | Value |\n"
	response += "|-------|-------|\n"
	response += fmt.Sprintf("| ID | %s |\n", user.Id)
	if name != "" {
		response += fmt.Sprintf("| Name | %s |\n", name)
	}
	if gamerTag != "" {
		response += fmt.Sprintf("| Gamer Tag | %s |\n", gamerTag)
	}
	if email != "" {
		response += fmt.Sprintf("| Email | %s |\n", email)
	}
	if emailVerified != "" {
		response += fmt.Sprintf("| Email Verified | %s |\n", emailVerified)
	}
	if image != "" {
		response += fmt.Sprintf("| Image | %s |\n", image)
	}
	if contentTagID != "" {
		response += fmt.Sprintf("| Content Tag ID | %s |\n", contentTagID)
	}

	return response
}

// UserAsDetailedText formats a single user with comprehensive details
func UserAsDetailedText(user *sharedpb.User) string {
	var response string
	response += fmt.Sprintf("**User Details for ID: %s**\n\n", user.Id)

	if user.Name != nil && *user.Name != "" {
		response += fmt.Sprintf("**Name:** %s\n", *user.Name)
	}

	if user.GamerTag != nil && *user.GamerTag != "" {
		response += fmt.Sprintf("**Gamer Tag:** %s\n", *user.GamerTag)
	}

	if user.Email != nil && *user.Email != "" {
		response += fmt.Sprintf("**Email:** %s\n", *user.Email)
	}

	if user.EmailVerified != nil {
		response += fmt.Sprintf("**Email Verified:** %s\n", user.EmailVerified.AsTime().Format("2006-01-02 15:04:05"))
	}

	if user.Image != nil && *user.Image != "" {
		response += fmt.Sprintf("**Image:** %s\n", *user.Image)
	}

	if user.ContentTagId != nil && *user.ContentTagId != "" {
		response += fmt.Sprintf("**Content Tag ID:** %s\n", *user.ContentTagId)
	}

	return response
}
