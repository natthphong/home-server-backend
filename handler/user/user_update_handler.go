package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/natthphong/home-server-backend/api"
)

func UpdateUserHandler(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req UpdateUserRequest
		if err := c.BodyParser(&req); err != nil {
			return api.BadRequest(c, "Invalid input")
		}

		if req.UserID == "" {
			return api.BadRequest(c, "UserID is required")
		}

		// Query to update user with COALESCE for optional fields
		updateQuery := `
			UPDATE tbl_user
			SET first_name_th = COALESCE($1, first_name_th),
				first_name_en = COALESCE($2, first_name_en),
				mid_name_th = COALESCE($3, mid_name_th),
				mid_name_en = COALESCE($4, mid_name_en),
				last_name_th = COALESCE($5, last_name_th),
				last_name_en = COALESCE($6, last_name_en),
				phone = COALESCE($7, phone),
				email = COALESCE($8, email),
				nationality = COALESCE($9, nationality),
				occupation = COALESCE($10, occupation),
				request_ref = COALESCE($11, request_ref),
				birth_date = COALESCE($12, birth_date),
				gender = COALESCE($13, gender),
				tax_id = COALESCE($14, tax_id),
				second_email = COALESCE($15, second_email),
				occupation_other_desc = COALESCE($16, occupation_other_desc),
				account_name = COALESCE($17, account_name),
				external_id = COALESCE($18, external_id),
				user_details = COALESCE($19, user_details),
				update_at = CURRENT_TIMESTAMP
			WHERE user_id = $20 AND is_delete = 'N'
		`

		// Execute the query
		_, err := db.Exec(c.Context(), updateQuery,
			req.FirstNameTh, req.FirstNameEn, req.MidNameTh, req.MidNameEn, req.LastNameTh,
			req.LastNameEn, req.Phone, req.Email, req.Nationality, req.Occupation,
			req.RequestRef, req.BirthDate, req.Gender, req.TaxID, req.SecondEmail,
			req.OccupationOtherDesc, req.AccountName, req.ExternalID, req.UserDetails, req.UserID,
		)

		if err != nil {

			return api.InternalError(c, err.Error())
		}

		return api.Ok(c, "User updated successfully")
	}
}
