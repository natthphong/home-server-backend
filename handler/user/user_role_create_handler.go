package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/natthphong/home-server-backend/api"
)

func AssignUserRoleHandler(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req RoleUser

		if err := c.BodyParser(&req); err != nil {
			return api.BadRequest(c, "Invalid input")
		}

		// Step 1: Validate user in tbl_user
		var isDeleted string
		var userIdToken string
		var inactiveUser bool

		var existingStatus string
		userCheckQuery := `
		SELECT tu.is_delete, tu.in_active, tu.status , tuca.user_id_token FROM tbl_user tu
		left outer join tbl_user_company_app tuca on tuca.user_id = tu.user_id
 		WHERE user_id = $1
		`
		// TODO check user have permission this app in this company
		// TODO select add left join with tbl_user_company_app
		err := db.QueryRow(c.Context(), userCheckQuery, req.UserID).Scan(&isDeleted, &inactiveUser, &existingStatus, &userIdToken)
		if err != nil {
			return api.BadRequest(c, err.Error())
		}
		if userIdToken == "" {
			return api.BadRequest(c, "user have no permission in this app")
		}
		if isDeleted == "Y" {
			return api.BadRequest(c, "User already exists but is marked as deleted")
		}
		if inactiveUser {
			return api.BadRequest(c, "User already exists but is inactive")
		}
		if existingStatus == StatusReject {
			return api.BadRequest(c, "User already exists but has a rejected status")
		}

		// Step 2: Insert role assignment into tbl_user_role
		insertQuery := `
			INSERT INTO tbl_user_role (
				role_code, user_id_token, create_by
			) VALUES (
				$1, $2, $3
			)
		`
		_, err = db.Exec(c.Context(), insertQuery, req.RoleCode, userIdToken, "system")
		if err != nil {
			return api.InternalError(c, err.Error())
		}

		return api.Ok(c, fiber.Map{"message": "User role assigned successfully"})
	}
}
