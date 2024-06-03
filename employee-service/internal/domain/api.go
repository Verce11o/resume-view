package domain

type CreateEmployeeRequest struct {
	FirstName    string `validate:"required" json:"first_name"`
	LastName     string `validate:"required" json:"last_name"`
	PositionName string `validate:"required" json:"position_name"`
	Salary       int    `validate:"required" json:"salary"`
}

type UpdateEmployeeRequest struct {
	FirstName  string `validate:"required" json:"first_name"`
	LastName   string `validate:"required" json:"last_name"`
	PositionID string `validate:"required" json:"position_id"`
}

type SignInEmployeeRequest struct {
	ID string `validate:"required" json:"id"`
}

type UpdatePositionRequest struct {
	Name   string `validate:"required" json:"name"`
	Salary int    `validate:"required" json:"salary"`
}

type CreatePositionRequest struct {
	Name   string `validate:"required" json:"name"`
	Salary int    `validate:"required" json:"salary"`
}
