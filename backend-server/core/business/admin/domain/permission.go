package domain

type Permission struct {
	ID         int64  `db:"id" json:"id"`
	RoleID     int64  `db:"role_id" json:"role_id"`
	ModuleName string `db:"module_name" json:"module_name"`
	Create     bool   `db:"create_perm" json:"create"`
	Read       bool   `db:"read_perm" json:"read"`
	Update     bool   `db:"update_perm" json:"update"`
	Delete     bool   `db:"delete_perm" json:"delete"`
	Export     bool   `db:"export_perm" json:"export"`
	Import     bool   `db:"import_perm" json:"import"`
}
