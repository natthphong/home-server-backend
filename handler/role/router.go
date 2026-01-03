package role

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Register(app fiber.Router, dbPool *pgxpool.Pool, jwtSecret string) {

	roleGroup := app.Group("/role")
	roleGroup.Get("", GetRolesHandler(dbPool))
	roleGroup.Get("/under", GetRolesUnderHandler(dbPool))
	roleGroup.Post("", CreateRoleHandler(dbPool))
	roleGroup.Delete("/:roleCode", DeleteRoleHandler(dbPool))
	roleGroup.Put("/:roleCode", UpdateRoleHandler(dbPool))
}
