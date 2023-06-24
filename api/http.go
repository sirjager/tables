package api

import (
	"database/sql"
	"fmt"

	"github.com/SirJager/tables/config"
	"github.com/SirJager/tables/middlewares"
	repo "github.com/SirJager/tables/service/core/repo"
	"github.com/SirJager/tables/service/core/tokens"
	"github.com/gin-gonic/gin"
)

func init() {
	c, _ := config.LoadConfig(".")
	if c.GinMode != gin.DebugMode {
		gin.SetMode(gin.ReleaseMode)
	}
}

type HttpServer struct {
	db           *sql.DB
	tokenBuilder tokens.Builder
	store        repo.Store
	router       *gin.Engine
	config       config.ServerConfig
}

func NewHttpServer(store repo.Store, db *sql.DB, cfg config.ServerConfig) (*HttpServer, error) {
	tokenBuilder, err := tokens.NewPasetoBuilder(cfg.TokenSecretKey)
	if err != nil {
		return nil, fmt.Errorf("could not initiate token builder : %w", err)
	}
	server := &HttpServer{
		store:        store,
		db:           db,
		tokenBuilder: tokenBuilder,
		config:       cfg,
	}
	server.setupHttpRouter()
	return server, nil
}

func (server *HttpServer) Start(address string) error {

	return server.router.Run(address)
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

func (server *HttpServer) setupHttpRouter() {
	router := gin.New()
	router.Use(middlewares.Logger())
	router.Use(middlewares.CORSMiddleware())

	authenticatedRoute := router.Group("/")
	unauthenticatedRoute := router.Group("/")
	onlyAdminCanAccessRoute := router.Group("/")

	// Middlewares
	onlyAdminCanAccessRoute.Use(middlewares.HasAdminPass())
	authenticatedRoute.Use(middlewares.BasicAuth(server.tokenBuilder))

	// UnAuthenticated Route
	unauthenticatedRoute.POST("/users/signup", server.createUser)
	unauthenticatedRoute.POST("/users/signin", server.loginUser)
	unauthenticatedRoute.POST("/users/renew-access", server.renewAccessToken)

	//* Authenticated Route

	onlyAdminCanAccessRoute.GET("/users", server.listUsers)

	authenticatedRoute.GET("/users/me", server.getUser)
	authenticatedRoute.DELETE("/users/me", server.deleteUser)

	// Manage Table
	authenticatedRoute.POST("/tables", server.createTable)
	authenticatedRoute.GET("/tables", server.listTables)
	authenticatedRoute.GET("/tables/:table", server.getTable)
	authenticatedRoute.DELETE("/tables/:table", server.deleteTable)

	// Manage Columns
	authenticatedRoute.POST("/tables/:table/columns", server.addColumns)
	authenticatedRoute.DELETE("/tables/:table/columns", server.deleteColumns)

	// Manage Rows
	authenticatedRoute.GET("/tables/:table/rows", server.getRows)
	authenticatedRoute.POST("/tables/:table/rows", server.insertRows)
	authenticatedRoute.PATCH("/tables/:table/rows", server.updateRows)
	authenticatedRoute.DELETE("/tables/:table/rows", server.deleteRows)

	server.router = router
}
