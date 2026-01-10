package auth

import (
	"context"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/natthphong/home-server-backend/api"
	"github.com/natthphong/home-server-backend/handler/user"
	"github.com/natthphong/home-server-backend/internal/logz"
	"golang.org/x/crypto/bcrypt"
)

func GenerateJWTForUser(
	db *pgxpool.Pool,
	userID, password, appCode, companyCode string,
	jwtSecret string,
	accessTokenDuration, refreshTokenDuration time.Duration,
	refreshTokenFlag bool,
) (map[string]interface{}, error) {
	logger := logz.NewLogger()
	var userDto User
	var userIdToken string
	query := `
			SELECT
			  u.user_id,
			  u.first_name_th, u.first_name_en, u.mid_name_th, u.mid_name_en, u.last_name_th, u.last_name_en,
			  u.phone, u.user_id_type, u.email, u.nationality, u.occupation, u.request_ref,
			  u.birth_date, u.gender, u.tax_id, u.second_email, u.occupation_other_desc,
			  u.is_active, u.status, u.account_name,
			  u.external_id, u.user_details, u.in_active,
			  tuca.user_id_token,
			  tuca.branch_code, tuca.app_code, tuca.company_code, tuca.user_active_time,u.password
			FROM tbl_user u
			JOIN tbl_user_company_app tuca
			  ON tuca.user_id = u.user_id
			 AND tuca.app_code = $2
			 AND tuca.company_code = $3
			 AND tuca.is_delete = 'N'
			WHERE u.is_delete = 'N'
			  AND u.user_id = $1
			LIMIT 1
		`
	err := db.QueryRow(context.Background(), query, userID, appCode, companyCode).Scan(
		&userDto.UserID,
		&userDto.FirstNameTh, &userDto.FirstNameEn, &userDto.MidNameTh, &userDto.MidNameEn,
		&userDto.LastNameTh, &userDto.LastNameEn,
		&userDto.Phone, &userDto.UserIDType, &userDto.Email, &userDto.Nationality,
		&userDto.Occupation, &userDto.RequestRef,
		&userDto.BirthDate, &userDto.Gender, &userDto.TaxID, &userDto.SecondEmail,
		&userDto.OccupationOtherDesc,
		&userDto.IsActive, &userDto.Status, &userDto.AccountName,
		&userDto.ExternalID, &userDto.UserDetails, &userDto.InActive,
		&userIdToken,
		&userDto.BranchCode, &userDto.AppCode, &userDto.CompanyCode, &userDto.UserActiveTime, &userDto.Password,
	)
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("userDto not found")
	}

	if userDto.InActive == "Y" {
		return nil, errors.New("userDto is inactive")
	}
	if userDto.Status != user.StatusActive {
		return nil, errors.New("userDto status not Success")
	}
	if !refreshTokenFlag {
		if err := bcrypt.CompareHashAndPassword([]byte(userDto.Password), []byte(password)); err != nil {
			return nil, errors.New("invalid password")
		}
	}

	// 2. Retrieve roles and objects for the userDto
	rolesQuery := `
		SELECT tr.role_code, tr.role_name_th, tr.role_name_en, tro.object_code
		FROM tbl_user_role tur
		LEFT JOIN tbl_role tr ON tur.role_code = tr.role_code
		LEFT JOIN tbl_role_object tro ON tr.role_code = tro.role_code
		LEFT JOIN tbl_object toj on tro.object_code = toj.object_code
		WHERE tur.user_id_token = $1
		and tur.is_delete = 'N' and tr.is_delete='N' and tro.is_delete='N' and toj.is_delete='N'
		ORDER BY tr.role_code;
	`
	rows, err := db.Query(context.Background(), rolesQuery, userIdToken)
	if err != nil {
		return nil, errors.New("Failed to retrieve roles")
	}
	defer rows.Close()

	roleMap := make(map[string]*Role)
	for rows.Next() {
		var objectCode *string
		var roleCode, roleNameTh, roleNameEn string
		if err := rows.Scan(&roleCode, &roleNameTh, &roleNameEn, &objectCode); err != nil {
			return nil, errors.New("Failed to scan roles")
		}

		if _, exists := roleMap[roleCode]; !exists {
			roleMap[roleCode] = &Role{
				RoleCode:   roleCode,
				RoleNameTh: roleNameTh,
				RoleNameEn: roleNameEn,
				Objects:    []string{},
			}
		}
		if objectCode != nil {
			roleMap[roleCode].Objects = append(roleMap[roleCode].Objects, *objectCode)
		}

	}

	// Convert roleMap to a slice of roles
	roles := make([]Role, 0, len(roleMap))
	for _, role := range roleMap {
		roles = append(roles, *role)
	}

	// 3. Generate JWT tokens
	// Access Token
	accessTokenClaims := jwt.MapClaims{
		"userIdToken": userIdToken,
		"userId":      userDto.UserID,
		"firstNameTh": userDto.FirstNameTh,
		"lastNameTh":  userDto.LastNameTh,
		"appCode":     userDto.AppCode,
		"companyCode": userDto.CompanyCode,
		"accountName": userDto.AccountName,
		"status":      userDto.Status,
		"roles":       roles, // Add roles to the JWT
		"exp":         time.Now().Add(accessTokenDuration).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return nil, errors.New("Failed to generate access token")
	}

	// Refresh Token
	refreshTokenClaims := jwt.MapClaims{
		"userIdToken": userIdToken,
		"userId":      userDto.UserID,
		"appCode":     userDto.AppCode,
		"companyCode": userDto.CompanyCode,
		"exp":         time.Now().Add(refreshTokenDuration).Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return nil, errors.New("Failed to generate refresh token")
	}

	// Response
	response := map[string]interface{}{
		"accessToken":  accessTokenString,
		"refreshToken": refreshTokenString,
		"jwtBody":      accessTokenClaims,
		"userIdToken":  userIdToken,
	}
	return response, nil
}

func LoginHandler(db *pgxpool.Pool, jwtSecret string, accessTokenDuration, refreshTokenDuration time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return api.BadRequest(c, "Invalid input")
		}

		// Call GenerateJWTForUser
		response, err := GenerateJWTForUser(db, req.Username, req.Password, req.AppCode, req.CompanyCode, jwtSecret, accessTokenDuration, refreshTokenDuration, false)
		if err != nil {
			return api.InternalError(c, err.Error())
		}

		return api.Ok(c, response)
	}
}
