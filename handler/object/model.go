package object

// Object represents the object entity
type Object struct {
	ObjectCode string  `json:"objectCode"`
	ObjectName string  `json:"objectName"`
	ObjectDesc *string `json:"objectDesc"`
	IsDeleted  *string `json:"isDeleted,omitempty"`
	CreateAt   *string `json:"createAt,omitempty"`
	CreateBy   *string `json:"createBy,omitempty"`
	UpdateBy   *string `json:"updateBy,omitempty"`
	UpdateAt   *string `json:"updateAt,omitempty"`
}

type CreateObjectRequest struct {
	ObjectCode  string `json:"objectCode" validate:"required,min=3,max=30"`
	ObjectName  string `json:"objectName" validate:"required,min=3,max=100"`
	ObjectDesc  string `json:"objectDesc" validate:"max=255"`
	AppCode     string `json:"appCode" validate:"required"`
	CompanyCode string `json:"companyCode" validate:"required"`
}

type DeleteObjectRequest struct {
	ObjectCode string `json:"objectCode" validate:"required"`
}
