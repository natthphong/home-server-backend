package object

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Register(app fiber.Router, dbPool *pgxpool.Pool, jwtSecret string) {
	//objectPermission, ok := permission["object"]
	//if !ok {
	//	objectPermission = []string{}
	//}
	objectGroup := app.Group("/object")
	objectGroup.Get("", GetObjectsHandler(dbPool))
	objectGroup.Post("", CreateObjectHandler(dbPool))
	objectGroup.Delete("/:objectCode", DeleteObjectHandler(dbPool))
}
