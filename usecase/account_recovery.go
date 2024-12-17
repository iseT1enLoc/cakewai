package usecase

import (
	appconfig "cakewai/cakewai.com/component/appcfg"
	"cakewai/cakewai.com/domain"
	tokenutil "cakewai/cakewai.com/internals/token_utils"
	"cakewai/cakewai.com/repository"
	"cakewai/cakewai.com/service"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"golang.org/x/crypto/bcrypt"
)

type accountRecovery struct {
	user_repo repository.UserRepository
	timeout   time.Duration
}

// ResetPasswordProcessing implements domain.AccountRecovery.
// func (a *accountRecovery) ResetPasswordProcessing(ctx context.Context, env *appconfig.Env, email string, new_password string, token string) error {
// 	fmt.Println("enter reset password processing")
// 	// Parse and validate the token
// 	parsedTokentoken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
// 		// Validate the algorithm
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 		}
// 		return env.ACCESS_SECRET, nil
// 	})
// 	fmt.Println("enter line 36")
// 	if err != nil {
// 		fmt.Println(err)
// 		return fmt.Errorf("failed to parse token: %w", err)
// 	}
// 	fmt.Println("enter line 40")
// 	// Check if the token is valid
// 	if !parsedTokentoken.Valid {
// 		return fmt.Errorf("invalid token")
// 	}
// 	fmt.Println("Enter line 45")
// 	// Extract claims (if needed)
// 	claims, ok := parsedTokentoken.Claims.(domain.JwtCustomClaims)
// 	if !ok {
// 		return fmt.Errorf("failed to extract claims")
// 	}
// 	fmt.Println("Enter line 51")
// 	// Extract data from claims
// 	userID, ok := claims.ID.(string)
// 	if !ok {
// 		return fmt.Errorf("user_id claim missing or invalid")
// 	}

//		user, _ := a.user_repo.GetUserById(ctx, userID)
//		encryptedPassword, _ := bcrypt.GenerateFromPassword(
//			[]byte(new_password),
//			bcrypt.DefaultCost,
//		)
//		user.Password = string(encryptedPassword)
//		err = a.user_repo.UpdateUser(ctx, user)
//		// Token validated; proceed with password reset flow
//		// Example: Allow user to reset the password (implement logic as needed)
//		if err != nil {
//			return fmt.Errorf("user_id claim missing or invalid")
//		}
//		return nil
//	}
//

// ChangesPassword implements domain.AccountRecovery.
func (a *accountRecovery) ChangesPassword(ctx context.Context, env *appconfig.Env, userId primitive.ObjectID, currentPassword string, newPassword string) error {
	err := a.user_repo.UpdateUserPassword(ctx, userId, currentPassword, newPassword)
	return err
}

func (a *accountRecovery) ResetPasswordProcessing(ctx context.Context, env *appconfig.Env, new_password string, token string) error {

	user_id, _, err := tokenutil.ExtractIDAndRole(token, env.ACCESS_SECRET)

	if err != nil {
		fmt.Println(err)
		fmt.Println(err)
		return fmt.Errorf("Error while handling token")
	}
	fmt.Println(user_id)

	user, err := a.user_repo.GetUserById(ctx, user_id)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}
	info := new_password
	fmt.Println(info)

	// Encrypt the new password
	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(new_password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return fmt.Errorf("failed to encrypt password: %w", err)
	}
	fmt.Println(new_password)
	// Update the user's password
	user.Password = string(encryptedPassword)
	err = a.user_repo.UpdateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Optionally, invalidate the token or log out the user if needed
	// You could set a flag or update any related session/token info in your system

	fmt.Println("Password reset successfully for user:", user_id)
	return nil
}

// ResetPasswordRequest implements domain.AccountRecovery.
func (a *accountRecovery) ResetPasswordRequest(ctx context.Context, env *appconfig.Env, email string) error {
	user, err := a.user_repo.GetUserByEmail(ctx, email)
	if err != nil {
		return errors.New("user not found")
	}

	token, _, err := tokenutil.CreateAccessToken(user.Id, env.ACCESS_SECRET, false, 100)
	if err != nil {
		log.Println(err)
		return err
	}

	resetLink := "http://localhost:3000/new_password.html?token=" + token
	sender := os.Getenv("SENDER")
	fmt.Println(sender)
	err = service.SendPasswordRecoveryEmail("Cakewai", sender, user.Name, user.Email, token, resetLink)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func NewAccountRecovery(user_repo repository.UserRepository, timeout time.Duration) domain.AccountRecovery {
	return &accountRecovery{
		user_repo: user_repo,
		timeout:   timeout,
	}
}
