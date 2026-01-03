package role

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/natthphong/home-server-backend/api"
)

func GetRolesUnderHandler(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleCode := c.Query("roleCode")

		query := `
		WITH RECURSIVE role_hierarchy AS (
			SELECT
				 role_code, parent_role_code, role_name_th, role_name_en, role_desc_th, role_desc_en, is_delete
			FROM
				tbl_role
			WHERE
				role_code =  $1

			UNION ALL

			SELECT
				 r.role_code, r.parent_role_code, r.role_name_th, r.role_name_en, r.role_desc_th, r.role_desc_en, r.is_delete
			FROM
				tbl_role r
			JOIN
				role_hierarchy rh ON r.parent_role_code = rh.role_code
		)
			SELECT  role_code, parent_role_code, role_name_th, role_name_en, role_desc_th, role_desc_en
			FROM role_hierarchy
			WHERE is_delete = 'N';
       `
		roles := []Role{}
		rows, err := db.Query(c.Context(), query, roleCode)
		if err != nil {
			return api.InternalError(c, err.Error())
		}
		defer rows.Close()

		for rows.Next() {
			var role Role
			if err := rows.Scan(&role.RoleCode, &role.ParentRoleCode, &role.RoleNameTh, &role.RoleNameEn, &role.RoleDescTh, &role.RoleDescEn); err != nil {
				return api.InternalError(c, err.Error())
			}
			roles = append(roles, role)
		}

		return api.Ok(c, roles)
	}
}
