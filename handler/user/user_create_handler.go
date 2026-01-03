package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/natthphong/home-server-backend/api"
	"golang.org/x/crypto/bcrypt"
)

func CreateUserHandler(dbPool *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			req CreateUserRequest
			//roleCode string
		)
		if err := c.BodyParser(&req); err != nil {
			return api.BadRequest(c, "Invalid input")
		}
		tx, err := dbPool.Begin(c.Context())
		if err != nil {
			return api.InternalError(c, err.Error())
		}
		defer func() {
			_ = tx.Rollback(c.Context())
		}()
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return api.InternalError(c, "Error hashing password")
		}
		// TODO
		status := StatusWaitApprove
		req.Password = string(hashedPassword)
		// 1) Insert tbl_user (ตามโครงสร้างใหม่)
		insertUserQuery := `
			INSERT INTO tbl_user (
				user_id, first_name_th, first_name_en, mid_name_th, mid_name_en, last_name_th, last_name_en, phone,
				user_id_type, email, nationality, occupation, request_ref, birth_date, gender, tax_id, second_email,
				occupation_other_desc, is_active, password, status, external_id, in_active,
				create_by, create_at, is_delete
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8,
				$9, $10, $11, $12, $13, $14, $15, $16, $17,
				$18, $19, $20, $21, $22, $23,
				$24, now(), $25
			)
		`

		_, err = tx.Exec(c.Context(), insertUserQuery,
			req.UserID,
			req.FirstNameTh, req.FirstNameEn,
			req.MidNameTh, req.MidNameEn,
			req.LastNameTh, req.LastNameEn,
			req.Phone,
			req.UserIDType,
			req.Email,
			req.Nationality,
			req.Occupation,
			req.RequestRef,
			req.BirthDate,
			req.Gender,
			req.TaxID,
			req.SecondEmail,
			req.OccupationOtherDesc,
			"Y",
			req.Password,
			status,
			req.ExternalID,
			"N",
			"system",
			"N",
		)
		if err != nil {
			return api.InternalError(c, err.Error())
		}

		userIDToken := uuid.NewString()
		insertUserCompanyAppQuery := `
			INSERT INTO tbl_user_company_app (
				user_id_token, in_active, user_id, app_code, company_code, branch_code,
				create_by, create_at, is_delete, user_active_time
			) VALUES (
				$1, $2, $3, $4, $5, $6,
				$7, now(), $8, now()
			)
		`
		_, err = tx.Exec(c.Context(), insertUserCompanyAppQuery,
			userIDToken,
			"N",
			req.UserID,
			req.AppCode,
			req.CompanyCode,
			req.BranchCode,
			"system",
			"N",
		)
		if err != nil {
			return api.InternalError(c, err.Error())
		}

		if err := tx.Commit(c.Context()); err != nil {
			return api.InternalError(c, err.Error())
		}
		return api.Ok(c, fiber.Map{"message": "User linked successfully", "status": status})
	}
}
