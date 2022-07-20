package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}
type renewAccessTokenRequest struct {
	RefreshToken string `header:"RefreshToken" binding:"required,gte=250,lt=360"`
}

func (server *HttpServer) renewAccessToken(ctx *gin.Context) {
	var req renewAccessTokenRequest
	var err error
	if err := ctx.ShouldBindHeader(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	refreshPayload, err := server.tokenBuilder.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, ErrorResponse{Error: err.Error()})
		return
	}

	session, err := server.store.GetSession(ctx, refreshPayload.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid token"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	if session.IsBlocked {
		ctx.JSON(http.StatusUnauthorized, ErrorResponse{Error: "session blocked"})
		return
	}
	if fmt.Sprintf("%d", session.Uid) != refreshPayload.User {
		ctx.JSON(http.StatusUnauthorized, ErrorResponse{Error: "incorrect session user"})
		return
	}
	if session.RefreshToken != req.RefreshToken {
		ctx.JSON(http.StatusUnauthorized, ErrorResponse{Error: "mismatch session token"})
		return
	}
	if time.Now().After(session.ExpiresAt) {
		ctx.JSON(http.StatusUnauthorized, ErrorResponse{Error: "session expired"})
		return
	}

	// create access token
	access_token, accessPayload, err := server.tokenBuilder.CreateToken(
		refreshPayload.User, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	refresTokenResponse := renewAccessTokenResponse{
		AccessToken:          access_token,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}

	ctx.JSON(http.StatusOK, refresTokenResponse)
}
