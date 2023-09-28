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

func (p *Permission) GetAllActivePermissions() []string {
	result := make([]string, 0)
	if p.Read {
		result = append(result, "READ")
	}
	if p.Create {
		result = append(result, "CREATE")
	}
	if p.Update {
		result = append(result, "UPDATE")
	}
	if p.Delete {
		result = append(result, "DELETE")
	}
	if p.Export {
		result = append(result, "EXPORT")
	}
	if p.Import {
		result = append(result, "IMPORT")
	}
	return result
}
