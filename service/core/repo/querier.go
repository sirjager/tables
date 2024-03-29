// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0

package core_repo

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	// -------------------------- ADD ONE TO -> _SESSIONS --------------------------
	CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error)
	// -------------------------- ADD _TABLES <-> _TABLES --------------------------
	CreateTable(ctx context.Context, arg CreateTableParams) (Table, error)
	// ------------------------------ ADD ONE _USERS <-> _USER  ------------------------------
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteTableWhereName(ctx context.Context, name string) error
	// -------------------------- REMOVE _TABLES <-> _TABLES --------------------------
	DeleteTablesWhereUser(ctx context.Context, userID int64) error
	// ------------------------------ REMOVE ONE _USERS -> nil  ------------------------------
	DeleteUser(ctx context.Context, id int64) error
	// ------------------------------ GET MULTIPLE _TABLES <== [_TABLES] ------------------------------
	GetAllTables(ctx context.Context) ([]Table, error)
	// ------------------------------ GET MULTIPLE _USERS <== [_USERS] ------------------------------
	GetAllUsers(ctx context.Context) ([]User, error)
	// -------------------------- GET ONE FROM <- _SESSIONS --------------------------
	GetSession(ctx context.Context, id uuid.UUID) (Session, error)
	GetSomeTables(ctx context.Context, arg GetSomeTablesParams) ([]Table, error)
	GetSomeTablesWhereUser(ctx context.Context, arg GetSomeTablesWhereUserParams) ([]Table, error)
	GetSomeUsers(ctx context.Context, arg GetSomeUsersParams) ([]User, error)
	GetTableByUserIdAndTableName(ctx context.Context, arg GetTableByUserIdAndTableNameParams) (Table, error)
	// -------------------------- GET ONE _TABLES <- _TABLES --------------------------
	GetTableWhereName(ctx context.Context, name string) (Table, error)
	// --------------------- GET MULTIPLE _TABLES OF _USERS.user_id <== [_TABLES] ---------------------
	GetTablesWhereUser(ctx context.Context, userID int64) ([]Table, error)
	// ------------------------------ GET ONE _USERS <== _USER  ------------------------------
	GetUser(ctx context.Context, id int64) (User, error)
	GetUserWhereEmail(ctx context.Context, email string) (User, error)
	GetUserWhereUsername(ctx context.Context, username string) (User, error)
	// -------------------------- UPDATE _TABLES <-> _TABLES --------------------------
	UpdateTableColumns(ctx context.Context, arg UpdateTableColumnsParams) (Table, error)
	UpdateUserBlocked(ctx context.Context, arg UpdateUserBlockedParams) (User, error)
	// ------------------------------ UPDATE ONE _USERS <-> _USERS  ------------------------------
	UpdateUserFullName(ctx context.Context, arg UpdateUserFullNameParams) (User, error)
	UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) (User, error)
	UpdateUserPublic(ctx context.Context, arg UpdateUserPublicParams) (User, error)
	UpdateUserUsername(ctx context.Context, arg UpdateUserUsernameParams) (User, error)
	UpdateUserVerified(ctx context.Context, arg UpdateUserVerifiedParams) (User, error)
}

var _ Querier = (*Queries)(nil)
