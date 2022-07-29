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

	// Middlewares
	authenticatedRoute.Use(middlewares.BasicAuth(server.tokenBuilder))

	// UnAuthenticated Route
	unauthenticatedRoute.POST("/users/signup", server.createUser)
	unauthenticatedRoute.POST("/users/signin", server.loginUser)
	unauthenticatedRoute.POST("/users/renew-access", server.renewAccessToken)

	// Authenticated Route
	authenticatedRoute.GET("/users", server.listUsers)
	authenticatedRoute.GET("/users/:user", server.getUser)
	authenticatedRoute.DELETE("/users/:user", server.deleteUser)

	// Manage Table
	authenticatedRoute.POST("/users/:user/tables", server.createTable)
	authenticatedRoute.GET("/users/:user/tables", server.listTables)
	authenticatedRoute.GET("/users/:user/tables/:table", server.getTable)
	authenticatedRoute.DELETE("/users/:user/tables/:table", server.deleteTable)

	// Manage Columns
	authenticatedRoute.POST("/users/:user/tables/:table/columns", server.addColumns)
	authenticatedRoute.DELETE("/users/:user/tables/:table/columns", server.deleteColumns)

	// Manage Rows
	authenticatedRoute.POST("/users/:user/tables/:table/rows", server.insertRows)
	authenticatedRoute.DELETE("/users/:user/tables/:table/rows", server.deleteRows)
	authenticatedRoute.GET("/users/:user/tables/:table/rows", server.getRows)
	authenticatedRoute.GET("/users/:user/tables/:table/rows:row", server.getRows)
	server.router = router
}
