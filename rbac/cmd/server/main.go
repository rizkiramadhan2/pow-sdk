package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	apphttp "github.com/rizkiramadhan2/rbac/app/http"
	rbacpg "github.com/rizkiramadhan2/rbac/internal/repository/postgres"
	"github.com/rizkiramadhan2/rbac/internal/usecase"
)

func main() {
	ctx := context.Background()

	databaseURL := env("DATABASE_URL", "postgres://postgres:postgres@localhost:5454/rbac?sslmode=disable")
	port := env("PORT", "10010")
	secret := os.Getenv("RBAC_SECRET")
	autoMigrate := env("RBAC_AUTO_MIGRATE", "true")

	db, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if autoMigrate == "true" {
		log.Println("running database migration")

		if err := rbacpg.Migrate(ctx, db); err != nil {
			log.Fatal(err)
		}

		log.Println("database migration completed")
	}

	store := rbacpg.NewStore(db)
	// workspace :=rbacpg.NewWorkspaceRepo(db)
	// user := rbacpg.NewUserRepo(db)
	// role := rbacpg.NewRoleRepo(db)
	// policy := rbacpg.NewPolicyRepo(db)

	server := apphttp.NewServer(apphttp.ServerConfig{
		Port:   port,
		Secret: secret,
	}, apphttp.Dependencies{
		AuthzUsecase:     usecase.NewAuthzUsecase(store),
		WorkspaceUsecase: usecase.NewWorkspaceUsecase(store),
		UserUsecase:      usecase.NewUserUsecase(store),
		RoleUsecase:      usecase.NewRoleUsecase(store),
		PolicyUsecase:    usecase.NewPolicyUsecase(store),
	})

	addr := ":" + port

	log.Println("rbac service running on", addr)

	if err := server.Run(addr); err != nil {
		log.Fatal(err)
	}
}

func env(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
