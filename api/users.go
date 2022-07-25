package api

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	repo "github.com/SirJager/tables/service/core/repo"
	"github.com/SirJager/tables/service/core/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type userAsResponse struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Fullname string `json:"fullname"`
	Public   bool   `json:"public"`
	Blocked  bool   `json:"blocked"`
	Verified bool   `json:"verified"`
	Updated  string `json:"updated"`
	Created  string `json:"created"`
}

// ------------------------------------------------------------------------------------------------------------
func removePassword(user repo.CoreUser) userAsResponse {
	return userAsResponse{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
		Fullname: user.Fullname,
		Public:   user.Public,
		Verified: user.Verified,
		Blocked:  user.Blocked,
		Created:  user.Created.String(),
		Updated:  user.Updated.String(),
	}
}

func decodeBasicAuth(basicAuth string) (string, string, error) {
	decoded, err := base64.URLEncoding.DecodeString(strings.Split(basicAuth, " ")[1])
	if err != nil {
		return "", "", err
	}
	username := strings.Split(string(decoded), ":")[0]
	password := string(decoded)[(len(username) + 1):len(string(decoded))]
	return username, password, nil
}

// ------------------------------------------------------------------------------------------------------------
type createUserRequest struct {
	Fullname string `json:"fullname" binding:"required,gte=3,lt=254"`
	Email    string `json:"email" binding:"required,email,lowercase,lt=320"`
	Username string `json:"username" binding:"required,alphanum,lowercase,gte=3,lte=60"`
	Password string `json:"password" binding:"required,gte=8,lte=320"`
}

func (server *HttpServer) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	// Hash Password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	arg := repo.AddCoreUserParams{Email: req.Email, Username: req.Username, Password: hashedPassword, Fullname: req.Fullname}
	dbuser, err := server.store.AddCoreUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "user with email or username already exists"})
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, removePassword(dbuser))
}

// ------------------------------------------------------------------------------------------------------------
type loginUserRequest struct {
	Authorization string `header:"Authorization" binding:"gte=10,lt=320"`
}

type loginUserResponse struct {
	AccessToken           string         `json:"access_token"`
	RefreshToken          string         `json:"refresh_token"`
	SessionId             uuid.UUID      `json:"session_id"`
	AccessTokenExpiresAt  time.Time      `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time      `json:"refresh_token_expires_at"`
	User                  userAsResponse `json:"user"`
}
type Username struct {
	Username string `json:"username" validate:"lowercase,alphanum,gte=3,lte=320"`
}
type Email struct {
	Email string `json:"email" validate:"lowercase,email,lte=320"`
}
type Password struct {
	Password string `json:"password" validate:"gte=8,lte=255"`
}

func (server *HttpServer) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	var err error
	if err := ctx.ShouldBindHeader(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	decodedIdentifier, decodedPass, err := decodeBasicAuth(req.Authorization)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	validate := validator.New()

	password := Password{Password: decodedPass}
	err = validate.Struct(password)
	if err != nil {
		println()
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	var dbuser repo.CoreUser
	if strings.Contains(decodedIdentifier, "@") {
		email := Email{Email: decodedIdentifier}
		err = validate.Struct(email)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}
		dbuser, err = server.store.GetCoreUserWithEmail(ctx, email.Email)
	} else {
		username := Username{Username: decodedIdentifier}
		err = validate.Struct(username)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}
		dbuser, err = server.store.GetCoreUserWithUsername(ctx, username.Username)
	}
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid credentials "})
			return
		}
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	// This will send back error if password do not match

	err = utils.VerifyPassword(password.Password, dbuser.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid credentials"})
		return
	}

	// create access token
	access_token, accessPayload, err := server.tokenBuilder.CreateToken(fmt.Sprintf("%d", dbuser.ID), server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	refresh_token, refreshPayload, err := server.tokenBuilder.CreateToken(fmt.Sprintf("%d", dbuser.ID), server.config.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	session, err := server.store.AddSession(ctx, repo.AddSessionParams{
		Expires:      refreshPayload.ExpiredAt,
		UserAgent:    ctx.Request.UserAgent(),
		ID:           refreshPayload.Id,
		ClientIp:     ctx.ClientIP(),
		RefreshToken: refresh_token,
		UserID:       dbuser.ID,
		Blocked:      false,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	loginResponse := loginUserResponse{
		SessionId:             session.ID,
		AccessToken:           access_token,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refresh_token,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User:                  removePassword(dbuser),
	}

	ctx.JSON(http.StatusOK, loginResponse)
}

// ------------------------------------------------------------------------------------------------------------
type UriParamsUser struct {
	User int64 `uri:"user" binding:"required,numeric,min=1"`
}

func (server *HttpServer) getUser(ctx *gin.Context) {
	var req UriParamsUser
	var err error
	if err = ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	dbuser, err := server.store.GetCoreUserWithUid(ctx, req.User)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, ErrorResponse{Error: "User not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, removePassword(dbuser))

}

// ------------------------------------------------------------------------------------------------------------
type UriParamsUid struct {
	Uid int64 `uri:"uid" binding:"required,min=1"`
}

func (server *HttpServer) deleteUser(ctx *gin.Context) {
	var req UriParamsUid
	var err error
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	err = server.store.RemoveCoreUserWithUid(ctx, req.Uid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, MessageResponse{Message: fmt.Sprintf("user=(%d) successfully deleted", req.Uid)})

}

// ------------------------------------------------------------------------------------------------------------
type listUsersRequest struct {
	Limit int32 `form:"limit" binding:"numeric,min=0"`
	Page  int32 `form:"page" binding:"numeric,min=0"`
}

type listUsersResponse struct {
	Users []userAsResponse
}

func (server *HttpServer) listUsers(ctx *gin.Context) {
	var req listUsersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if req.Limit == 0 && req.Page == 0 {
		//  [[ Returns All Users ]]   if req.Limit and req.Page == 0
		dbusers, err := server.store.ListCoreUsers(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			return
		}
		users := []userAsResponse{}
		for _, dbuser := range dbusers {
			users = append(users, removePassword(dbuser))
		}
		ctx.JSON(http.StatusOK, users)
		return
	} else if req.Limit > 0 && req.Page < 1 {
		//  [[ Returns Limited User ]]
		dbusers, err := server.store.ListCoreUsersWithLimit(ctx, req.Limit)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			return
		}
		users := []userAsResponse{}
		for _, dbuser := range dbusers {
			users = append(users, removePassword(dbuser))
		}
		ctx.JSON(http.StatusOK, users)
		return
	}

	// ------------------------------------------------------------------------------------------------------------
	// Returns Limited Users of request page (offset).
	arg := repo.ListCoreUsersWithLimitOffsetParams{
		Limit:  req.Limit,
		Offset: (req.Page - 1) * req.Limit,
	}
	dbusers, err := server.store.ListCoreUsersWithLimitOffset(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	users := []userAsResponse{}
	for _, dbuser := range dbusers {
		users = append(users, removePassword(dbuser))
	}

	ctx.JSON(http.StatusOK, listUsersResponse{Users: users})
}
