package data

import (
	"context"
	"database/sql"
	"time"
)

type Permissions []string

// Add func to check permissions silce contian a 
// specfic permissions check
func (p Permissions) Include(code string) bool{
	for i := range p{
		if code == p[i]{
			return true
		}
	}
	return false
}

type PermissionModel struct{
	DB *sql.DB
}

// return all the specfic permissions for the user
func (m PermissionModel) GetAllForUser(userID int64) (Permissions, error) {
	query := `
	SELECT permissions.code
	FROM permissions
	INNER JOIN users_permissions ON user_permissions.permissions_id = permissions.id
	INNER JOIN users ON users_permissions.user_id = users.id
	WHERE users.id = $1
	`
}
