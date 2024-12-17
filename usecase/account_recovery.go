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
	username, newpass, err := a.user_repo.HandleForgotPassword(ctx, email)
	if err != nil {
		return errors.New("Fail to reset password")
	}
	sender := os.Getenv("SENDER")
	fmt.Println(sender)
	err = service.SendPasswordRecoveryEmail("Cakewai", sender, *username, email, *newpass)
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
