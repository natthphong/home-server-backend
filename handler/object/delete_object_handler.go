package object

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/natthphong/home-server-backend/api"
)

func DeleteObjectHandler(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO open txn and and update roleObject 'N' after object soft close
		objectCode := c.Params("objectCode")
		if objectCode == "" {
			return api.BadRequest(c, "Invalid object code")
		}
		query := `UPDATE tbl_object SET is_delete = 'Y' WHERE object_code = $1`
		cmdTag, err := db.Exec(c.Context(), query, objectCode)
		if err != nil {
			return api.InternalError(c, err.Error())
		}

		rows := cmdTag.RowsAffected()
		if rows == 0 {

			return api.NotFound(c, "Object not found")
		}

		return api.Ok(c, "Object deleted successfully")
	}
}
