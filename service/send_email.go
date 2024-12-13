package service

import (
	"fmt"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendPasswordRecoveryEmail sends the email with a password recovery link
func SendPasswordRecoveryEmail(senderName, senderEmail, receiverName, receiverEmail string, resetToken string, reset_link string) error {
	// Sender and recipient details
	from := mail.NewEmail(senderName, senderEmail)
	to := mail.NewEmail(receiverName, receiverEmail)

	// Email subject and content
	subject := "Password recover"
	plainTextContent := fmt.Sprintf("Hi %s,\n\nWe received a request to reset your password. Please click the link below to reset your password:\n\n%s\n\nIf you did not request a password reset, please ignore this email.\n\nBest regards,\nYour App Team", receiverName, reset_link)
	htmlContent := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<body>
			<p>Hi %s,</p>
			<p>We received a request to reset your password. Please click the link below to reset your password:</p>
			<a href="%s">Reset Password</a>
			<p>If you did not request a password reset, please ignore this email.</p>
			<p>Best regards,<br> <strong>Cakewai Team</strong></p>
		</body>
		</html>
	`, receiverName, reset_link)

	// Create the email message
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	// Fetch SendGrid API Key from environment variables
	apiKey := os.Getenv("EMAIL_API")
	if apiKey == "" {
		return fmt.Errorf("SENDGRID_API_KEY not set in environment variables")
	}

	// Send email using SendGrid
	client := sendgrid.NewSendClient(apiKey)
	response, err := client.Send(message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	// Log response for debugging purposes
	if response.StatusCode >= 400 {
		fmt.Printf("Failed to send email. Status Code: %d, Body: %s\n", response.StatusCode, response.Body)
		return fmt.Errorf("email send failed with status code: %d", response.StatusCode)
	}

	fmt.Printf("Password recovery email sent successfully! Status code: %d\n", response.StatusCode)
	return nil
}
