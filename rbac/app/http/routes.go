package http

func (s *Server) registerRoutes() {
	v1 := s.router.Group("/v1")

	v1.POST("/authz/check", s.authzHandler.Can)

	v1.POST("/workspaces", s.workspaceHandler.Create)
	v1.GET("/workspaces/:workspace_id", s.workspaceHandler.Get)

	v1.POST("/users", s.userHandler.Create)
	v1.GET("/workspaces/:workspace_id/users/:user_id", s.userHandler.Get)

	v1.POST("/roles", s.roleHandler.Create)
	v1.GET("/workspaces/:workspace_id/roles", s.roleHandler.List)
	v1.GET("/workspaces/:workspace_id/roles/:role_id", s.roleHandler.Get)
	v1.POST("/roles/assign", s.roleHandler.AssignToUser)
	v1.DELETE("/workspaces/:workspace_id/users/:user_id/roles/:role_id", s.roleHandler.RemoveFromUser)

	v1.POST("/policies", s.policyHandler.Create)
	v1.GET("/workspaces/:workspace_id/policies", s.policyHandler.List)
	v1.GET("/workspaces/:workspace_id/policies/:policy_id", s.policyHandler.Get)
	v1.POST("/roles/policies/attach", s.policyHandler.AttachToRole)
	v1.DELETE("/workspaces/:workspace_id/roles/:role_id/policies/:policy_id", s.policyHandler.DetachFromRole)
}
