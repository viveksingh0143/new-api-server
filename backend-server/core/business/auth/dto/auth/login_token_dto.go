package auth

type LoginTokenDto struct {
	AccessToken  string           `json:"access_token"`
	RefreshToken string           `json:"refresh_token"`
	Name         string           `json:"name"`
	StaffID      string           `json:"staff_id"`
	Permissions  []*PermissionDto `json:"permissions"`
}
