package object

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/natthphong/home-server-backend/api"
)

func CreateObjectHandler(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req CreateObjectRequest
		if err := c.BodyParser(&req); err != nil {
			return api.BadRequest(c, "Invalid input")
		}
		validate := validator.New()
		err := validate.Struct(req)
		if err := validate.Struct(req); err != nil {
			return api.ValidationErrorResponse(c, err, req)
		}
		if req.ObjectCode == "" || req.AppCode == "" || req.CompanyCode == "" {
			return api.BadRequest(c, "Invalid input")
		}
		objectCode := fmt.Sprintf("%s_%s_%s", req.CompanyCode, req.AppCode, req.ObjectCode)
		query := `
		    INSERT INTO tbl_object ( object_code, object_name, object_desc, is_delete)
		    VALUES ($1, $2, $3, 'N' )
		    ON CONFLICT (object_code) DO UPDATE
		    SET object_name = EXCLUDED.object_name,
		        object_desc = EXCLUDED.object_desc,
		        is_delete = 'N',
		        update_at = now()
		`
		_, err = db.Exec(c.Context(), query, objectCode, req.ObjectName, req.ObjectDesc)
		if err != nil {
			return api.InternalError(c, err.Error())
		}

		return api.Ok(c, "Object created successfully")
	}
}
