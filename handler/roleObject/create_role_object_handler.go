package roleObject

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/natthphong/home-server-backend/api"
)

func CreateRoleObjectHandler(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req CreateRoleObjectRequest
		if err := c.BodyParser(&req); err != nil {
			return api.BadRequest(c, "Invalid input")
		}
		validate := validator.New()
		err := validate.Struct(req)
		if err := validate.Struct(req); err != nil {
			return api.ValidationErrorResponse(c, err, req)
		}

		query := `
			INSERT INTO tbl_role_object (role_code, object_code, is_delete)
			VALUES ($1, $2, 'N')
			ON CONFLICT (role_code, object_code) DO UPDATE
			SET is_delete = 'N',
				update_at = CURRENT_TIMESTAMP
		`
		_, err = db.Exec(c.Context(), query, req.RoleCode, req.ObjectCode)
		if err != nil {
			return api.InternalError(c, err.Error())
		}

		return api.Ok(c, "Role-object created successfully")
	}
}
