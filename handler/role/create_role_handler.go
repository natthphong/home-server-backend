package role

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/natthphong/home-server-backend/api"
)

func CreateRoleHandler(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req CreateRoleRequest
		if err := c.BodyParser(&req); err != nil {
			return api.BadRequest(c, "Invalid input")
		}
		if req.RoleCode == "" || req.AppCode == "" {
			return api.BadRequest(c, "Invalid input")
		}
		roleCode := fmt.Sprintf("%s_%s_%s", req.CompanyCode, req.AppCode, req.RoleCode)
		validate := validator.New()
		err := validate.Struct(req)
		if err := validate.Struct(req); err != nil {
			return api.ValidationErrorResponse(c, err, req)
		}
		query := `
			INSERT INTO tbl_role (role_code, parent_role_code, role_name_th, role_desc_th, role_name_en, role_desc_en, is_delete)
			VALUES ($1, $2, $3, $4, $5, $6, 'N')
			ON CONFLICT (role_code) DO UPDATE
			SET parent_role_code = EXCLUDED.parent_role_code,
				role_name_th = EXCLUDED.role_name_th,
				role_desc_th = EXCLUDED.role_desc_th,
				role_name_en = EXCLUDED.role_name_en,
				role_desc_en = EXCLUDED.role_desc_en,
				is_delete = 'N',
				update_at = now()
		`
		_, err = db.Exec(c.Context(), query, roleCode, req.ParentRoleCode, req.RoleNameTh, req.RoleDescTh, req.RoleNameEn, req.RoleDescEn)
		if err != nil {
			return api.InternalError(c, err.Error())
		}

		return api.Ok(c, "Role created successfully")
	}
}
