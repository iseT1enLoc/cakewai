package routes

import (
	"cakewai/cakewai.com/api/handlers"
	appconfig "cakewai/cakewai.com/component/appcfg"
	"cakewai/cakewai.com/repository"
	"cakewai/cakewai.com/usecase"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gin-gonic/gin"
)

func StartTokenCleanupTask(ctx context.Context, repo repository.RefreshTokenRepository) {
	fmt.Println("START TOKEN CLEAN UP WORK")
	// Run cleanup every hour (adjust as necessary)
	//10 ngay xoa mot lan
	ticker := time.NewTicker(time.Hour * 240)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Perform cleanup task
			if err := repo.CleanupExpiredTokens(ctx); err != nil {
				fmt.Errorf("Failed to cleanup expired tokens: %v", err)
			}
		case <-ctx.Done():
			// Stop cleanup when context is done (e.g., server shutdown)
			return
		}
	}
}

func NewRefreshTokenRoute(env *appconfig.Env, timeout time.Duration, db *mongo.Database, r *gin.RouterGroup) {
	repo := repository.NewrefreshTokenRepository(db, "refresh_token")
	sc := handlers.RefreshTokenHandler{
		RefreshTokenUsecase: usecase.NewRefreshTokenUseCase(repo, timeout),
		Env:                 env,
	}
	// Start cleanup task in the background
	go StartTokenCleanupTask(context.Background(), repo)
	r.GET("/refreshtoken", sc.RefreshTokenHandler())
	r.PUT("/revoke_token", sc.RevokeRefreshTokenHandler())
}
