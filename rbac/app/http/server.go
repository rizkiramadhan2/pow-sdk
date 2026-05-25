package http

import (
	"github.com/gin-gonic/gin"

	httpHandler "github.com/rizkiramadhan2/rbac/handler/http"
	"github.com/rizkiramadhan2/rbac/internal/usecase"
)

type Server struct {
	router *gin.Engine

	authzHandler     *httpHandler.AuthzHandler
	workspaceHandler *httpHandler.WorkspaceHandler
	userHandler      *httpHandler.UserHandler
	roleHandler      *httpHandler.RoleHandler
	policyHandler    *httpHandler.PolicyHandler
}

type ServerConfig struct {
	Port   string
	Secret string
}

type Dependencies struct {
	AuthzUsecase     *usecase.AuthzUsecase
	WorkspaceUsecase *usecase.WorkspaceUsecase
	UserUsecase      *usecase.UserUsecase
	RoleUsecase      *usecase.RoleUsecase
	PolicyUsecase    *usecase.PolicyUsecase
}

func NewServer(config ServerConfig, deps Dependencies) *Server {
	router := gin.Default()

	if config.Secret != "" {
		router.Use(SecretAuthMiddleware(config.Secret))
	}

	server := &Server{
		router: router,

		authzHandler:     httpHandler.NewAuthzHandler(deps.AuthzUsecase),
		workspaceHandler: httpHandler.NewWorkspaceHandler(deps.WorkspaceUsecase),
		userHandler:      httpHandler.NewUserHandler(deps.UserUsecase),
		roleHandler:      httpHandler.NewRoleHandler(deps.RoleUsecase),
		policyHandler:    httpHandler.NewPolicyHandler(deps.PolicyUsecase),
	}

	server.registerRoutes()

	return server
}

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}
