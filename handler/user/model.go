package user

import "time"

// Status constants
const (
	StatusWaitApprove = "WAIT_APPROVE"
	StatusActive      = "ACTIVE"
	StatusReject      = "REJECT"
)

type UserRole struct {
	UserID      string `json:"userId"`
	FirstNameTh string `json:"firstNameTh"`
	LastNameTh  string `json:"lastNameTh"`
	RoleCode    string `json:"roleCode"`
	RoleNameTh  string `json:"roleNameTh"`
	RoleNameEn  string `json:"roleNameEn"`
}

type RoleUser struct {
	UserID      string  `json:"userId"`
	AppCode     string  `json:"appCode"`
	CompanyCode string  `json:"companyCode"`
	BranchCode  *string `json:"branchCode,omitempty"`
	RoleCode    string  `json:"roleCode"`
}

// User represents the user data structure corresponding to the tbl_user table.
type User struct {
	UserID              string       `json:"userId"`
	FirstNameTh         *string      `json:"firstNameTh"`
	FirstNameEn         *string      `json:"firstNameEn"`
	MidNameTh           *string      `json:"midNameTh"`
	MidNameEn           *string      `json:"midNameEn"`
	LastNameTh          *string      `json:"lastNameTh"`
	LastNameEn          *string      `json:"lastNameEn"`
	Phone               *string      `json:"phone"`
	UserIDType          *string      `json:"userIdType"`
	Email               *string      `json:"email"`
	Nationality         *string      `json:"nationality"`
	Occupation          *string      `json:"occupation"`
	RequestRef          *string      `json:"requestRef"`
	BirthDate           *time.Time   `json:"birthDate"`
	Gender              *string      `json:"gender"`
	TaxID               *string      `json:"taxId"`
	SecondEmail         *string      `json:"secondEmail"`
	OccupationOtherDesc *string      `json:"occupationOtherDesc"`
	IsActive            *string      `json:"isActive"`
	Password            string       `json:"password"`
	Status              string       `json:"status"`
	AccountName         *string      `json:"accountName"`
	ExternalID          *string      `json:"externalId"`
	UserDetails         *interface{} `json:"userDetails,omitempty"`
	InActive            string       `json:"inActive"`
}

// CreateUserRequest represents the request body for linking a user to an app.
type CreateUserRequest struct {
	UserID      string  `json:"userId"`
	Password    string  `json:"password"`
	AppCode     string  `json:"appCode"`
	RoleCode    string  `json:"roleCode,omitempty"`
	CompanyCode string  `json:"companyCode"`
	BranchCode  *string `json:"branchCode,omitempty"`
	ExternalID  *string `json:"externalId,omitempty"`

	FirstNameTh         *string      `json:"firstNameTh"`
	FirstNameEn         *string      `json:"firstNameEn"`
	MidNameTh           *string      `json:"midNameTh"`
	MidNameEn           *string      `json:"midNameEn"`
	LastNameTh          *string      `json:"lastNameTh"`
	LastNameEn          *string      `json:"lastNameEn"`
	Phone               *string      `json:"phone"`
	UserIDType          *string      `json:"userIdType"`
	Email               *string      `json:"email"`
	Nationality         *string      `json:"nationality"`
	Occupation          *string      `json:"occupation"`
	RequestRef          *string      `json:"requestRef"`
	BirthDate           *string      `json:"birthDate"` // yyyy-mm-dd
	Gender              *string      `json:"gender"`
	TaxID               *string      `json:"taxId"`
	SecondEmail         *string      `json:"secondEmail"`
	OccupationOtherDesc *string      `json:"occupationOtherDesc"`
	IsActive            *string      `json:"isActive"`
	Status              string       `json:"status"`
	AccountName         *string      `json:"accountName"`
	UserDetails         *interface{} `json:"userDetails,omitempty"`
	InActive            string       `json:"inActive"`
}

// UpdateUserRequest represents the request body for updating user documents.
type UpdateUserRequest struct {
	UserID              string                    `json:"userId"`
	FirstNameTh         *string                   `json:"firstNameTh"`
	FirstNameEn         *string                   `json:"firstNameEn"`
	MidNameTh           *string                   `json:"midNameTh"`
	MidNameEn           *string                   `json:"midNameEn"`
	LastNameTh          *string                   `json:"lastNameTh"`
	LastNameEn          *string                   `json:"lastNameEn"`
	Phone               *string                   `json:"phone"`
	Email               *string                   `json:"email"`
	Nationality         *string                   `json:"nationality"`
	Occupation          *string                   `json:"occupation"`
	RequestRef          *string                   `json:"requestRef"`
	BirthDate           *string                   `json:"birthDate"` //(ISO 8601 format: "YYYY-MM-DD")
	Gender              *string                   `json:"gender"`
	TaxID               *string                   `json:"taxId"`
	SecondEmail         *string                   `json:"secondEmail"`
	OccupationOtherDesc *string                   `json:"occupationOtherDesc"`
	AccountName         *string                   `json:"accountName"`
	ExternalID          *string                   `json:"externalId"`  //
	UserDetails         *[]map[string]interface{} `json:"userDetails"` //  (JSON object as string)
	InActive            *bool                     `json:"inActive"`
}

// ApproveUserRequest represents the request body for approving a user.
type ApproveUserRequest struct {
	UserID string `json:"userId"`
	Status string `json:"status"`
}

// UserListResponse represents the paginated response for listing users.
type UserListResponse struct {
	Page       int    `json:"page"`
	Size       int    `json:"size"`
	TotalCount int    `json:"totalCount"`
	TotalPage  int    `json:"totalPage"`
	Users      []User `json:"users"`
}
type InquiryRequest struct {
	Status     *string `json:"status"`
	AppCode    *string `json:"appCode"`
	UserID     *string `json:"userId"`
	ExternalID *string `json:"externalId"`
	Email      *string `json:"email"`
	Phone      *string `json:"phone"`
}

type RoleWithUsers struct {
	RoleID              int     `json:"roleId"`
	RoleCode            string  `json:"roleCode"`
	RoleNameTh          string  `json:"roleNameTh"`
	RoleNameEn          *string `json:"roleNameEn"`
	UserID              *string `json:"userId"`
	AppCode             *string `json:"appCode"`
	FirstNameTh         *string `json:"firstNameTh"`
	FirstNameEn         *string `json:"firstNameEn"`
	MidNameTh           *string `json:"midNameTh"`
	MidNameEn           *string `json:"midNameEn"`
	LastNameTh          *string `json:"lastNameTh"`
	LastNameEn          *string `json:"lastNameEn"`
	IsInactive          *bool   `json:"isInactive"`
	Nationality         *string `json:"nationality"`
	Occupation          *string `json:"occupation"`
	OccupationOtherDesc *string `json:"occupationOtherDesc"`
	CompanyCode         *string `json:"companyCode"`
	BranchCode          *string `json:"branchCode"`
}
