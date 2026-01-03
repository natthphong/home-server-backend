package role

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/natthphong/home-server-backend/api"
)

func GetRolesHandler(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		page := c.QueryInt("page", 1)
		size := c.QueryInt("size", 10)
		search := c.Query("search", "")           // Optional parameter
		appCode := c.Query("appCode", "")         // Optional search parameter
		companyCode := c.Query("companyCode", "") // Optional search parameter
		if page < 1 || size < 10 {
			return api.BadRequest(c, "Invalid pagination parameters: page must be >= 1 and size must be >= 10")
		}
		offset := (page - 1) * size
		roles := []Role{}

		// Build the base queries
		query := `SELECT  role_code, role_name_th, role_name_en, role_desc_th, role_desc_en, parent_role_code
				  FROM tbl_role
				  WHERE is_delete = 'N'`
		countQuery := `SELECT COUNT(*) FROM tbl_role WHERE is_delete = 'N'`

		queryArgs := []interface{}{}
		countArgs := []interface{}{}
		argIndex := 1

		if search != "" {
			likeClause := fmt.Sprintf(" AND (role_code || role_name_th || role_name_en || COALESCE(role_desc_th,'') ||  COALESCE(role_desc_en,'')  ) LIKE $%d", argIndex)
			query += likeClause
			countQuery += likeClause
			queryArgs = append(queryArgs, "%"+search+"%")
			countArgs = append(countArgs, "%"+search+"%")
			argIndex++
		}
		if companyCode != "" && appCode != "" {
			prefixObjectCode := fmt.Sprintf("%s_%s", companyCode, appCode)
			likeClause := fmt.Sprintf(" AND role_code like $%d", argIndex)
			query += likeClause
			countQuery += likeClause
			queryArgs = append(queryArgs, "%"+prefixObjectCode)
			countArgs = append(countArgs, "%"+prefixObjectCode)
			argIndex++
		}

		query += fmt.Sprintf(" ORDER BY create_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
		queryArgs = append(queryArgs, size, offset)

		rows, err := db.Query(c.Context(), query, queryArgs...)
		if err != nil {
			return api.InternalError(c, err.Error())
		}
		defer rows.Close()

		for rows.Next() {
			var role Role
			if err := rows.Scan(&role.RoleCode, &role.RoleNameTh, &role.RoleNameEn, &role.RoleDescTh, &role.RoleDescEn, &role.ParentRoleCode); err != nil {
				return api.InternalError(c, err.Error())
			}
			roles = append(roles, role)
		}

		var totalCount int
		err = db.QueryRow(c.Context(), countQuery, countArgs...).Scan(&totalCount)
		if err != nil {
			return api.InternalError(c, err.Error())
		}

		totalPage := (totalCount + size - 1) / size

		// Response with pagination details
		response := fiber.Map{
			"totalCount": totalCount,
			"totalPage":  totalPage,
			"list":       roles,
		}

		return api.Ok(c, response)
	}
}
