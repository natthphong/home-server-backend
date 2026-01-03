package role

type Role struct {
	RoleCode       string  `json:"roleCode"`
	ParentRoleCode *string `json:"parentRoleCode"`
	RoleNameTh     *string `json:"roleNameTh"`
	RoleNameEn     *string `json:"roleNameEn"`
	RoleDescTh     *string `json:"roleDescTh"`
	RoleDescEn     *string `json:"roleDescEn"`
}

type CreateRoleRequest struct {
	RoleCode       string  `json:"roleCode" validate:"required"`
	AppCode        string  `json:"appCode" validate:"required"`
	CompanyCode    string  `json:"companyCode" validate:"required"`
	RoleNameTh     *string `json:"roleNameTh"`
	RoleNameEn     *string `json:"roleNameEn"`
	RoleDescTh     *string `json:"roleDescTh"`
	RoleDescEn     *string `json:"roleDescEn"`
	ParentRoleCode *string `json:"parentRoleCode,omitempty"`
}

type UpdateRoleRequest struct {
	RoleNameTh   *string `json:"roleNameTh"`
	RoleNameEn   *string `json:"roleNameEn"`
	RoleDescTh   *string `json:"roleDescTh"`
	RoleDescEn   *string `json:"roleDescEn"`
	ParentRoleId *string `json:"parentRoleCode,omitempty"`
}
