package role

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/natthphong/home-server-backend/api"
)

func UpdateRoleHandler(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleCode := c.Params("roleCode")
		var req UpdateRoleRequest
		if err := c.BodyParser(&req); err != nil {
			return api.BadRequest(c, "Invalid input")
		}

		query := `
           UPDATE tbl_role
           SET role_name_th = COALESCE($1, role_name_th),
               role_name_en = COALESCE($2, role_name_en),
               role_desc_th = COALESCE($3, role_desc_th),
               role_desc_en = COALESCE($4, role_desc_en),
               parent_role_id = COALESCE($5, parent_role_id),
               update_at = NOW()
           WHERE role_code = $6
       `
		_, err := db.Exec(c.Context(), query, req.RoleNameTh, req.RoleNameEn, req.RoleDescTh, req.RoleDescEn, req.ParentRoleId, roleCode)
		if err != nil {
			return api.InternalError(c, err.Error())
		}

		return api.Ok(c, "Role updated successfully")
	}
}
