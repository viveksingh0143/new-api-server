package auth

type RegisterUserDto struct {
	Name     string `json:"name" validate:"required"`
	Username string `json:"username" validate:"required"`
	EMail    string `json:"email" validate:"required"`
	StaffID  string `json:"staff_id" validate:"required"`
	Password string `json:"password" validate:"required"`
}
