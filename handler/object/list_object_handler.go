package object

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/natthphong/home-server-backend/api"
)

func GetObjectsHandler(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		page := c.QueryInt("page", 1)
		size := c.QueryInt("size", 10)
		searchTerm := c.Query("search", "")        // Optional search parameter
		orderBy := c.Query("orderBy", "create_at") // Optional search parameter
		appCode := c.Query("appCode", "")          // Optional search parameter
		companyCode := c.Query("companyCode", "")  // Optional search parameter

		if page < 1 || size < 10 {
			return api.BadRequest(c, "Invalid pagination parameters: page must be >= 1 and size must be >= 10")
		}
		if orderBy != "create_at" || orderBy != "object_cde" || orderBy != "object_de" {

		}

		offset := (page - 1) * size
		objects := []Object{}

		query := `SELECT object_code, object_name,  object_desc 
				  FROM tbl_object 
				  WHERE is_delete = 'N'`
		countQuery := `SELECT COUNT(*) FROM tbl_object WHERE is_delete = 'N'`

		queryArgs := []interface{}{}
		countArgs := []interface{}{}
		i := 1

		if searchTerm != "" {
			likeClause := fmt.Sprintf(" AND (object_code ||  COALESCE(object_desc,'')  ) LIKE $%d", i)
			query += likeClause
			countQuery += likeClause
			queryArgs = append(queryArgs, "%"+searchTerm+"%")
			countArgs = append(countArgs, "%"+searchTerm+"%")
			i++
		}
		if companyCode != "" && appCode != "" {
			prefixObjectCode := fmt.Sprintf("%s_%s", companyCode, appCode)
			likeClause := fmt.Sprintf(" AND object_code like $%d", i)
			query += likeClause
			countQuery += likeClause
			queryArgs = append(queryArgs, "%"+prefixObjectCode)
			countArgs = append(countArgs, "%"+prefixObjectCode)
			i++
		}

		query += fmt.Sprintf(" ORDER BY create_at DESC LIMIT $%d OFFSET $%d", i, i+1)
		queryArgs = append(queryArgs, size, offset)

		rows, err := db.Query(c.Context(), query, queryArgs...)
		if err != nil {
			return api.InternalError(c, err.Error())
		}
		defer rows.Close()

		for rows.Next() {
			var object Object
			if err := rows.Scan(&object.ObjectCode, &object.ObjectName, &object.ObjectDesc); err != nil {
				return api.InternalError(c, err.Error())
			}
			objects = append(objects, object)
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
			"list":       objects,
		}

		return api.Ok(c, response)
	}
}
