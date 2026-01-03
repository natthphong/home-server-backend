package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/natthphong/home-server-backend/api"
)

func GetUserRolesHandler(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Query("userId")
		//appCode := c.Query("appCode")
		page := c.QueryInt("page", 1)
		size := c.QueryInt("size", 10)

		if userID == "" {
			return api.BadRequest(c, "UserID are required")
		}
		if page < 1 || size < 10 {
			return api.BadRequest(c, "Invalid pagination parameters: page must be >= 1 and size must be >= 10")
		}

		offset := (page - 1) * size
		roleAssignments := []UserRole{} // Assuming UserRole struct to store joined data
		args := []interface{}{userID, size, offset}

		// Main query with LEFT JOINs for roles
		query := `
			SELECT u.user_id, u.first_name_th, u.last_name_th, ur.role_code, r.role_name_th, r.role_name_en
			FROM tbl_user u
			left outer join tbl_user_company_app tuca on tuca.user_id = u.user_id
			LEFT JOIN tbl_user_role ur ON tuca.user_id_token = ur.user_id_token
			LEFT JOIN tbl_role r ON ur.role_code = r.role_code
			WHERE u.user_id = $1  AND u.is_delete = 'N'
			ORDER BY u.create_at DESC
			LIMIT $2 OFFSET $3
		`

		// Count query for pagination
		countQuery := `
			SELECT COUNT(*)
			FROM tbl_user u
			left outer join tbl_user_company_app tuca on tuca.user_id = u.user_id
			LEFT JOIN tbl_user_role ur ON tuca.user_id_token = ur.user_id_token
			WHERE u.user_id = $1 AND u.is_delete = 'N'
		`

		// Execute the main query
		rows, err := db.Query(c.Context(), query, args...)
		if err != nil {
			return api.InternalError(c, err.Error())
		}
		defer rows.Close()

		// Scan each row into the roleAssignments slice
		for rows.Next() {
			var roleAssignment UserRole
			if err := rows.Scan(
				&roleAssignment.UserID,
				&roleAssignment.FirstNameTh,
				&roleAssignment.LastNameTh,
				&roleAssignment.RoleCode,
				&roleAssignment.RoleNameTh,
				&roleAssignment.RoleNameEn,
			); err != nil {
				return api.InternalError(c, err.Error())
			}
			roleAssignments = append(roleAssignments, roleAssignment)
		}

		// Execute the count query
		var totalCount int
		err = db.QueryRow(c.Context(), countQuery, userID).Scan(&totalCount)
		if err != nil {
			return api.InternalError(c, err.Error())
		}

		// Calculate total pages
		totalPage := (totalCount + size - 1) / size

		// Create response format
		response := fiber.Map{
			"totalCount": totalCount,
			"totalPage":  totalPage,
			"page":       page,
			"size":       size,
			"list":       roleAssignments,
		}

		return api.Ok(c, response)
	}
}
