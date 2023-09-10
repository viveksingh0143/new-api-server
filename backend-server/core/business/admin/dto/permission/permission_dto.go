package permission

type PermissionDto struct {
	ID         int64  `json:"id"`
	RoleID     int64  `json:"role_id"`
	ModuleName string `json:"module_name"`
	Create     bool   `json:"create_perm"`
	Read       bool   `json:"read_perm"`
	Update     bool   `json:"update_perm"`
	Delete     bool   `json:"delete_perm"`
	Export     bool   `json:"export_perm"`
	Import     bool   `json:"import_perm"`
}
