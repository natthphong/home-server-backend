package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/natthphong/home-server-backend/api"
)

func ApproveUserHandler(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req ApproveUserRequest
		if err := c.BodyParser(&req); err != nil {
			return api.BadRequest(c, "Invalid input")
		}

		if req.UserID == "" {
			return api.BadRequest(c, "UserID is required")
		}
		if req.Status == StatusWaitApprove {
			return api.BadRequest(c, "status cannot WAIT_APPROVE")
		}
		if req.Status != StatusReject && req.Status != StatusSuccess {
			return api.BadRequest(c, "must be status REJECT or SUCCESS ")
		}

		updateQuery := `
			UPDATE tbl_user
			SET status = $1,
				update_at = CURRENT_TIMESTAMP
			WHERE user_id = $2
		`

		_, err := db.Exec(c.Context(), updateQuery, req.Status, req.UserID)
		if err != nil {
			return api.InternalError(c, err.Error())
		}
		if req.Status == StatusReject {
			//TODO
		}

		return api.Ok(c, fiber.Map{"message": "success", "status": req.Status})
	}
}
