package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/SirJager/tables/middlewares"
	repo "github.com/SirJager/tables/service/core/repo"
	"github.com/SirJager/tables/service/core/tokens"
	"github.com/gin-gonic/gin"
)

const (
	NO_ACCESS_TO_RESOURCE = "you do not have access to this resources"
)

type createTableRequest struct {
	Table   string        `json:"table" binding:"required,alphanum,gte=3,lte=60"`
	Columns []repo.Column `json:"columns" binding:"required"`
}

func (server *HttpServer) createTable(ctx *gin.Context) {
	var req createTableRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	authPayload := ctx.MustGet(middlewares.AuthorizationPayloadKey).(*tokens.Payload)
	UserID, err := strconv.Atoi(authPayload.User)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, ErrorResponse{Error: NO_ACCESS_TO_RESOURCE})
		return
	}

	mytable, err := server.store.CreateTableTx(ctx, repo.CreateTableTxParams{Table: req.Table, UserID: int64(UserID), Columns: req.Columns})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, mytable)
}

func (server *HttpServer) listTables(ctx *gin.Context) {
	authPayload := ctx.MustGet(middlewares.AuthorizationPayloadKey).(*tokens.Payload)
	UserID, err := strconv.Atoi(authPayload.User)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, ErrorResponse{Error: NO_ACCESS_TO_RESOURCE})
		return
	}
	res, err := server.store.GetTablesWhereUser(ctx, int64(UserID))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	var tables []repo.TableSchema
	for _, t := range res {
		table, err := repo.FormatTableEntryToTable(t)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			return
		}
		tables = append(tables, table)
	}
	ctx.JSON(http.StatusOK, tables)
}

type getTableParams struct {
	Table string `uri:"table" binding:"required,alphanum,gte=3,lte=60"`
}

func (server *HttpServer) getTable(ctx *gin.Context) {
	var req getTableParams
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	authPayload := ctx.MustGet(middlewares.AuthorizationPayloadKey).(*tokens.Payload)
	UserID, err := strconv.Atoi(authPayload.User)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, ErrorResponse{Error: NO_ACCESS_TO_RESOURCE})
		return
	}
	res, err := server.store.GetTableByUserIdAndTableName(ctx, repo.GetTableByUserIdAndTableNameParams{UserID: int64(UserID), Name: req.Table})
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, ErrorResponse{Error: fmt.Sprintf("table '%v' not found", req.Table)})
			return
		}
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	table, err := repo.FormatTableEntryToTable(res)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, table)
}

// -----------------------------------------------------------------------------------------------------

type deleteTableRequest struct {
	Table string `uri:"table" binding:"required,alphanum,gte=3,lte=60"`
}

func (server *HttpServer) deleteTable(ctx *gin.Context) {
	var req deleteTableRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	authPayload := ctx.MustGet(middlewares.AuthorizationPayloadKey).(*tokens.Payload)
	UserID, err := strconv.Atoi(authPayload.User)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, ErrorResponse{Error: NO_ACCESS_TO_RESOURCE})
		return
	}
	// Our First we will get the table by using userid and tablename
	table, err := server.store.GetTableByUserIdAndTableName(ctx, repo.GetTableByUserIdAndTableNameParams{UserID: int64(UserID), Name: req.Table})
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, ErrorResponse{Error: fmt.Sprintf("table '%v' not found", req.Table)})
			return
		}
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	err = server.store.DropTableTx(ctx, repo.Name{Value: table.Name})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, MessageResponse{Message: fmt.Sprintf("Table %v successfully deleted", table.Name)})
}

// -----------------------------------------------------------------------------------------------------
type insertRowsRequest struct {
	Rows []map[string]interface{} `json:"rows" binding:"required"`
}
type insertRowsRequestUri struct {
	Table string `uri:"table" binding:"required,alphanum,gte=3,lte=60"`
}

func (server *HttpServer) insertRows(ctx *gin.Context) {
	var req insertRowsRequest
	var uri insertRowsRequestUri
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	println(fmt.Sprintf("%v", req))
	authPayload := ctx.MustGet(middlewares.AuthorizationPayloadKey).(*tokens.Payload)
	UserID, err := strconv.Atoi(authPayload.User)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, ErrorResponse{Error: NO_ACCESS_TO_RESOURCE})
		return
	}
	// If table belongs to user or not is checked inside insert rows func
	err = server.store.InsertRows(ctx, repo.InsertRowsParams{UserID: int64(UserID), Table: uri.Table, Rows: req.Rows})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, MessageResponse{Message: fmt.Sprintf("%d rows inserted in %s", len(req.Rows), uri.Table)})
}

// -----------------------------------------------------------------------------------------------------
type updateRowsRequest struct {
	Rows []map[string]interface{} `json:"rows" binding:"required"`
}
type updateRowsRequestUri struct {
	Table string `uri:"table" binding:"required,alphanum,gte=3,lte=60"`
}

func (server *HttpServer) updateRows(ctx *gin.Context) {
	var req updateRowsRequest
	var uri updateRowsRequestUri
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	println(fmt.Sprintf("%v", req))
	authPayload := ctx.MustGet(middlewares.AuthorizationPayloadKey).(*tokens.Payload)
	UserID, err := strconv.Atoi(authPayload.User)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, ErrorResponse{Error: NO_ACCESS_TO_RESOURCE})
		return
	}
	// If table belongs to user or not is checked inside insert rows func
	err = server.store.UpdateRows(ctx, repo.UpdateRowsParams{UserID: int64(UserID), Table: uri.Table, Rows: req.Rows})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, MessageResponse{Message: fmt.Sprintf("%d rows updated in %s", len(req.Rows), uri.Table)})
}

// -----------------------------------------------------------------------------------------------------
type addColumnsParams struct {
	Columns []repo.Column `json:"columns" binding:"required"`
}

type addColumnsUri struct {
	Table string `uri:"table" binding:"required,alphanum,gte=3,lte=60"`
}

func (server *HttpServer) addColumns(ctx *gin.Context) {
	var req addColumnsParams
	var uri addColumnsUri
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	authPayload := ctx.MustGet(middlewares.AuthorizationPayloadKey).(*tokens.Payload)
	UserID, err := strconv.Atoi(authPayload.User)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, ErrorResponse{Error: NO_ACCESS_TO_RESOURCE})
		return
	}

	updatedTable, err := server.store.AddColumnTx(ctx,
		repo.AddColumnsTxParams{UserID: int64(UserID), Table: uri.Table, Columns: req.Columns})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, updatedTable)
}

// -----------------------------------------------------------------------------------------------------
type dropColumnsParams struct {
	Columns []string `json:"columns" binding:"required"`
}

type dropColumnsUri struct {
	Table string `uri:"table" binding:"required,gte=3,lte=50"`
}

func (server *HttpServer) deleteColumns(ctx *gin.Context) {
	var req dropColumnsParams
	var uri dropColumnsUri
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	authPayload := ctx.MustGet(middlewares.AuthorizationPayloadKey).(*tokens.Payload)
	UserID, err := strconv.Atoi(authPayload.User)
	if err != nil {
		ctx.JSON(http.StatusNotFound, ErrorResponse{Error: "you do not have access to this table"})
		return
	}
	updatedTable, err := server.store.DropColumnTx(ctx, repo.DropColumnsTxParams{UserID: int64(UserID), Table: uri.Table, Columns: req.Columns})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedTable)
}

// -----------------------------------------------------------------------------------------------------
type deleteRowsUriParams struct {
	Table string `uri:"table" validate:"required,alphanum,gte=3,lte=60"`
}

// valid body example: { "rows": { "id": [ 1, 2, 3 ], "name": [ "user1" ] } }
type deleteRowParams struct {
	Filters map[string]interface{} `json:"filters" binding:"required"`
}

func (server *HttpServer) deleteRows(ctx *gin.Context) {
	var req deleteRowParams
	var uri deleteRowsUriParams
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	authPayload := ctx.MustGet(middlewares.AuthorizationPayloadKey).(*tokens.Payload)
	UserID, err := strconv.Atoi(authPayload.User)
	if err != nil {
		ctx.JSON(http.StatusNotFound, ErrorResponse{Error: NO_ACCESS_TO_RESOURCE})
		return
	}

	_, err = server.store.DeleteRows(ctx, repo.DeleteRowsParams{Table: uri.Table, UserID: int64(UserID), Filters: req.Filters})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, MessageResponse{Message: "Done"})
}

// -----------------------------------------------------------------------------------------------------

type getRowsParams struct {
	Table string `uri:"table" binding:"required,gte=3,lte=50"`
}

type getRowsBodyParams struct {
	Fields  []string               `json:"fields" binding:""`
	Filters map[string]interface{} `json:"filters" binding:""`
}

func (server *HttpServer) getRows(ctx *gin.Context) {
	var uri getRowsParams
	var req getRowsBodyParams
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	authPayload := ctx.MustGet(middlewares.AuthorizationPayloadKey).(*tokens.Payload)
	UserID, err := strconv.Atoi(authPayload.User)
	if err != nil {
		ctx.JSON(http.StatusNotFound, ErrorResponse{Error: "you do not have access to this table"})
		return
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		if err.Error() != "EOF" {
			ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}
	}

	// If Body is not empty
	result, err := server.store.GetRows(ctx, repo.GetRowsParams{UserID: int64(UserID), Table: uri.Table, Fields: req.Fields, Filters: req.Filters})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, result)
}
