package role

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/natthphong/home-server-backend/api"
)

func DeleteRoleHandler(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleCode := c.Params("roleCode")
		// TODO open txn and and update parentRoleCode = null
		query := `UPDATE tbl_role SET is_delete = 'Y' WHERE role_code = $1`
		_, err := db.Exec(c.Context(), query, roleCode)
		if err != nil {
			return api.InternalError(c, err.Error())
		}

		return api.Ok(c, "Role deleted successfully")
	}
}
