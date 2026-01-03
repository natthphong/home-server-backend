package roleObject

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/natthphong/home-server-backend/api"
)

func GetRoleObjectsHandler(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		page := c.QueryInt("page", 1)
		size := c.QueryInt("size", 10)
		roleCode := c.Query("roleCode", "")
		if page < 1 || size < 10 {
			return api.BadRequest(c, "Invalid pagination parameters: page must be >= 1 and size must be >= 10")
		}

		offset := (page - 1) * size
		roleObjects := []RoleObject{}

		query := `SELECT role_code, object_code FROM tbl_role_object WHERE is_delete = 'N'`
		countQuery := `SELECT COUNT(*) FROM tbl_role_object WHERE is_delete = 'N'`

		args := []interface{}{}
		countArgs := []interface{}{}
		i := 1

		if roleCode != "" {
			roleCodeCondition := fmt.Sprintf(" AND role_code = $%d", i)
			query += roleCodeCondition
			countQuery += roleCodeCondition
			args = append(args, roleCode)
			countArgs = append(countArgs, roleCode)
			i++
		}

		query += fmt.Sprintf(" ORDER BY create_at DESC LIMIT $%d OFFSET $%d", i, i+1)
		args = append(args, size, offset)

		rows, err := db.Query(c.Context(), query, args...)
		if err != nil {
			return api.InternalError(c, err.Error())
		}
		defer rows.Close()

		for rows.Next() {
			var roleObject RoleObject
			if err := rows.Scan(&roleObject.RoleCode, &roleObject.ObjectCode); err != nil {
				return api.InternalError(c, err.Error())
			}
			roleObjects = append(roleObjects, roleObject)
		}

		var totalCount int
		err = db.QueryRow(c.Context(), countQuery, countArgs...).Scan(&totalCount)
		if err != nil {
			return api.InternalError(c, err.Error())
		}

		totalPage := (totalCount + size - 1) / size
		
		response := fiber.Map{
			"totalCount": totalCount,
			"totalPage":  totalPage,
			"list":       roleObjects,
		}

		return api.Ok(c, response)
	}
}
