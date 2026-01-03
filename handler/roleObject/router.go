package roleObject

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Register(app fiber.Router, dbPool *pgxpool.Pool, jwtSecret string) {

	roleObjectGroup := app.Group("/role-object")
	roleObjectGroup.Get("", GetRoleObjectsHandler(dbPool))
	roleObjectGroup.Post("", CreateRoleObjectHandler(dbPool))
	roleObjectGroup.Delete("", DeleteRoleObjectHandler(dbPool))
}
