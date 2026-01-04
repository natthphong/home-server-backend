package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Register(router fiber.Router, dbPool *pgxpool.Pool) {

	userGroup := router.Group("/user")
	userGroup.Post("", CreateUserHandler(dbPool))
	userGroup.Get("", ListUsersHandler(dbPool))
	// TODO testing
	userGroup.Post("/inquiry", InquiryUsersHandler(dbPool))
	userGroup.Put("", UpdateUserHandler(dbPool))
	userGroup.Post("/role", AssignUserRoleHandler(dbPool))
	userGroup.Get("/role", GetUserRolesHandler(dbPool))
	userGroup.Post("/approve", ApproveUserHandler(dbPool))

}
