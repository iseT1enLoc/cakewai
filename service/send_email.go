package service

import (
	"fmt"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendPasswordRecoveryEmail sends the email with a password recovery link
func SendPasswordRecoveryEmail(senderName, senderEmail, receiverName, receiverEmail string, new_password string) error {
	// Sender and recipient details
	from := mail.NewEmail(senderName, senderEmail)
	to := mail.NewEmail(receiverName, receiverEmail)

	// Email subject and content
	// Email subject and content
	subject := "Khôi phục mật khẩu - Mật khẩu mới"
	plainTextContent := fmt.Sprintf("Chào %s,\n\nChúng tôi nhận được yêu cầu đặt lại mật khẩu của bạn. Mật khẩu mới của bạn là: %s\n\nVui lòng giữ nó an toàn. Nếu bạn không yêu cầu khôi phục mật khẩu, vui lòng bỏ qua email này.\n\nTrân trọng,\nNhóm Cakewai", receiverName, new_password)

	// HTML email content with the new password
	htmlContent := fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<body>
	    <p>CỬA HÀNG BÁN BÁNH CAKEWAI</p>
		<p>Chào %s,</p>
		<p>Chúng tôi nhận được yêu cầu đặt lại mật khẩu của bạn. Mật khẩu mới của bạn là:</p>
		<p><strong>%s</strong></p>
		<a href = "http://localhost:5173">Đăng nhập lại tại đây!<a>
		<p>Trân trọng,<br> <strong>Cửa hàng Cakewai</strong></p>
	</body>
	</html>
`, receiverName, new_password)

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
