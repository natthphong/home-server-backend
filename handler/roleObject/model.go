package roleObject

// RoleObject represents the association between a role and an object
type RoleObject struct {
	RoleCode   string `json:"roleCode"`
	ObjectCode string `json:"objectCode"`
	IsDeleted  string `json:"isDeleted,omitempty"`
}

// CreateRoleObjectRequest represents the request body for creating a role-object association
type CreateRoleObjectRequest struct {
	RoleCode   string `json:"roleCode" validate:"required"`
	ObjectCode string `json:"objectCode" validate:"required"`
}

// DeleteRoleObjectRequest represents the request body for deleting a role-object association
type DeleteRoleObjectRequest struct {
	RoleCode   string `json:"roleCode" validate:"required"`
	ObjectCode string `json:"objectCode" validate:"required"`
}
