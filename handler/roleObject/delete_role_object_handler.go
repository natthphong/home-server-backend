package roleObject

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/natthphong/home-server-backend/api"
)

func DeleteRoleObjectHandler(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req DeleteRoleObjectRequest
		if err := c.BodyParser(&req); err != nil {
			return api.BadRequest(c, "Invalid input")
		}

		query := `UPDATE tbl_role_object SET is_delete = 'Y' WHERE role_code = $1 AND object_code = $2`
		_, err := db.Exec(c.Context(), query, req.RoleCode, req.ObjectCode)
		if err != nil {
			return api.InternalError(c, err.Error())
		}

		return api.Ok(c, "Role-object deleted successfully")
	}
}
