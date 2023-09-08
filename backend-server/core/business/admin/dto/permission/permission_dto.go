package permission

type PermissionDto struct {
	ID         int64  `json:"id"`
	RoleID     int64  `json:"role_id"`
	ModuleName string `json:"module_name"`
	Create     bool   `json:"create"`
	Read       bool   `json:"read"`
	Update     bool   `json:"update"`
	Delete     bool   `json:"delete"`
	Export     bool   `json:"export"`
	Import     bool   `json:"import"`
}
