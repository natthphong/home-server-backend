package user

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/natthphong/home-server-backend/api"
)

func ListUsersHandler(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		//appCode := c.Query("appCode")
		//companyCode := c.Query("companyCode")
		// TODO
		//branchCode := c.Query("branchCode")
		status := c.Query("status")
		startDate := c.Query("start_date", "")
		endDate := c.Query("end_date", "")
		page := c.QueryInt("page", 1)
		size := c.QueryInt("size", 10)

		if page < 1 || size < 10 {
			return api.BadRequest(c, "Invalid pagination parameters: page must be >= 1 and size >= 10")
		}

		offset := (page - 1) * size
		users := []User{}
		args := []interface{}{}
		countArgs := []interface{}{}

		query := `SELECT user_id, first_name_th, first_name_en, mid_name_th, mid_name_en, last_name_th,
							last_name_en, phone, user_id_type, email, nationality, occupation, request_ref,
							birth_date, gender, tax_id, second_email, occupation_other_desc, is_active, status, account_name,
							external_id, user_details, in_active
				   FROM tbl_user
				   WHERE is_delete = 'N'`
		countQuery := `SELECT COUNT(*) FROM tbl_user WHERE is_delete = 'N'`

		// Add filters for appCode and status if provided
		i := 1
		//if appCode != "" {
		//	query += fmt.Sprintf(" AND appCode = $%d", i)
		//	countQuery += fmt.Sprintf(" AND appCode = $%d", i)
		//	args = append(args, appCode)
		//	countArgs = append(countArgs, appCode)
		//	i++
		//}
		if status != "" {
			query += fmt.Sprintf(" AND status = $%d", i)
			countQuery += fmt.Sprintf(" AND status = $%d", i)
			args = append(args, status)
			countArgs = append(countArgs, status)
			i++
		}

		// Add date range filter
		if startDate != "" && endDate != "" {
			query += fmt.Sprintf(" AND DATE(create_at) BETWEEN $%d AND $%d", i, i+1)
			countQuery += fmt.Sprintf(" AND DATE(create_at) BETWEEN $%d AND $%d", i, i+1)
			args = append(args, startDate, endDate)
			countArgs = append(countArgs, startDate, endDate)
			i += 2
		}

		// Add ORDER BY and pagination
		query += fmt.Sprintf(" ORDER BY create_at DESC LIMIT $%d OFFSET $%d", i, i+1)
		args = append(args, size, offset)

		// Execute the query to get user data
		rows, err := db.Query(c.Context(), query, args...)
		if err != nil {
			return api.InternalError(c, err.Error())
		}
		defer rows.Close()

		// Scan user data into the users slice
		for rows.Next() {
			var user User
			if err := rows.Scan(
				&user.UserID, &user.FirstNameTh, &user.FirstNameEn, &user.MidNameTh, &user.MidNameEn,
				&user.LastNameTh, &user.LastNameEn, &user.Phone, &user.UserIDType, &user.Email,
				&user.Nationality, &user.Occupation, &user.RequestRef, &user.BirthDate, &user.Gender,
				&user.TaxID, &user.SecondEmail, &user.OccupationOtherDesc, &user.IsActive, &user.Status, &user.AccountName, &user.ExternalID, &user.UserDetails, &user.InActive,
			); err != nil {
				return api.InternalError(c, err.Error())
			}
			users = append(users, user)
		}

		// Get total count
		var totalCount int
		if len(countArgs) == 0 {
			err = db.QueryRow(c.Context(), countQuery).Scan(&totalCount)
		} else {
			err = db.QueryRow(c.Context(), countQuery, countArgs...).Scan(&totalCount)
		}
		if err != nil {
			return api.InternalError(c, err.Error())
		}

		// Calculate total pages
		totalPage := (totalCount + size - 1) / size

		response := UserListResponse{
			Page:       page,
			Size:       size,
			TotalCount: totalCount,
			TotalPage:  totalPage,
			Users:      users,
		}

		return api.Ok(c, response)
	}
}
