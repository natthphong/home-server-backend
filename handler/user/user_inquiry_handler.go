package user

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/natthphong/home-server-backend/api"
)

func InquiryUsersHandler(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req InquiryRequest
		if err := c.BodyParser(&req); err != nil {
			return api.BadRequest(c, "Invalid input")
		}

		query := `
				SELECT user_id, first_name_th, first_name_en, mid_name_th, mid_name_en, last_name_th,
							last_name_en, phone, user_id_type, email, nationality, occupation, request_ref,
							birth_date, gender, tax_id, second_email, occupation_other_desc, is_active, status, account_name,
							external_id, user_details, in_active
				   FROM tbl_user
				   WHERE is_delete = 'N'
		`

		args := []interface{}{}
		i := 1
		if req.AppCode != nil && *req.AppCode != "" {
			query += fmt.Sprintf(" AND appCode = $%d", i)
			args = append(args, *req.AppCode)
			i++
		}
		if req.Status != nil && *req.Status != "" {
			query += fmt.Sprintf(" AND status = $%d", i)
			args = append(args, *req.Status)
			i++
		}
		if req.UserID != nil && *req.UserID != "" {
			query += fmt.Sprintf(" AND user_id = $%d", i)
			args = append(args, *req.UserID)
			i++
		}
		if req.ExternalID != nil && *req.ExternalID != "" {
			query += fmt.Sprintf(" AND external_id = $%d", i)
			args = append(args, *req.ExternalID)
			i++
		}
		if req.Email != nil && *req.Email != "" {
			query += fmt.Sprintf(" AND email = $%d", i)
			args = append(args, *req.Email)
			i++
		}
		if req.Phone != nil && *req.Phone != "" {
			query += fmt.Sprintf(" AND phone = $%d", i)
			args = append(args, *req.Phone)
			i++
		}
		if i == 1 {
			return api.BadRequest(c, "Invalid input")
		}

		query += " ORDER BY create_at DESC"
		rows, err := db.Query(c.Context(), query, args...)
		if err != nil {
			return api.InternalError(c, err.Error())
		}
		defer rows.Close()

		// Parse results
		users := []User{}
		for rows.Next() {
			var user User
			if err := rows.Scan(
				&user.UserID, &user.FirstNameTh, &user.FirstNameEn, &user.MidNameTh, &user.MidNameEn,
				&user.LastNameTh, &user.LastNameEn, &user.Phone, &user.UserIDType, &user.Email,
				&user.Nationality, &user.Occupation, &user.RequestRef, &user.BirthDate, &user.Gender,
				&user.TaxID, &user.SecondEmail, &user.OccupationOtherDesc, &user.IsActive,
				&user.Status, &user.AccountName, &user.ExternalID, &user.UserDetails, &user.InActive,
			); err != nil {
				return api.InternalError(c, err.Error())
			}
			users = append(users, user)
		}

		response := fiber.Map{
			"users": users,
		}

		return api.Ok(c, response)
	}
}
