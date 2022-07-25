package api

import (
	"fmt"
	"net/http"

	repo "github.com/SirJager/tables/service/core/repo"
	"github.com/gin-gonic/gin"
)

type createTableRequest struct {
	TableName string        `json:"table" binding:"required,gte=3,lte=60"`
	Uid       int64         `json:"uid" binding:"required,numeric,min=1"`
	Columns   []repo.Column `json:"columns" binding:"required"`
}

func (server *HttpServer) createTable(ctx *gin.Context) {
	var req createTableRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	mytable, err := server.store.CreateTableTx(ctx, repo.CreateTableTxParams{TableName: req.TableName, Uid: req.Uid, Columns: req.Columns})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, mytable)
}

type listTableParams struct {
	User int64 `uri:"user" binding:"required,min=1"`
}

func (server *HttpServer) listTables(ctx *gin.Context) {
	var req listTableParams
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	res, err := server.store.ListCoreTablesWithUid(ctx, req.User)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	var tables []repo.RealTable
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
	User  int32  `uri:"user" binding:"required,numeric,min=1"`
	Table string `uri:"table" binding:"required,gte=3"`
}

func (server *HttpServer) getTable(ctx *gin.Context) {
	var req getTableParams
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	res, err := server.store.GetCoreTableWithName(ctx, req.Table)
	if err != nil {
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
	Table string `uri:"table" binding:"required,gte=3,lte=50"`
	User  int64  `uri:"user" binding:"required,numeric,min=1"`
}

func (server *HttpServer) deleteTable(ctx *gin.Context) {
	var req deleteTableRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	err := server.store.DropTableTx(ctx, repo.RemoveCoreTableWithUidAndNameParams{UserID: req.User, Name: req.Table})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return

	}
	ctx.JSON(http.StatusOK, MessageResponse{Message: fmt.Sprintf("Table '%s' deleted", req.Table)})
}

// -----------------------------------------------------------------------------------------------------
type insertRowsRequest struct {
	Rows [][]repo.KeyValueParams `json:"rows" binding:"required"`
}
type insertRowsRequestUri struct {
	Table string `uri:"table" binding:"required,alphanum,gte=1,lte=50"`
	User  int32  `uri:"user" binding:"required,numeric,min=1"`
}

type getRowsParams struct {
	User  int32  `uri:"user" binding:"required,numeric,min=1"`
	Table string `uri:"table" binding:"required,gte=3,lte=50"`
}

func (server *HttpServer) getRows(ctx *gin.Context) {
	var req getRowsParams
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	res, err := server.store.GetRows(ctx, repo.GetRowsParams{Uid: req.User, Tablename: req.Table})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, res)
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
	err := server.store.InsertRows(ctx, repo.InsertRowsParams{Uid: uri.User, Tablename: uri.Table, Rows: req.Rows})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, MessageResponse{Message: fmt.Sprintf("%d rows inserted in %s", len(req.Rows), uri.Table)})
}
