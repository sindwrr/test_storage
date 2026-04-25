package models

import "time"

type User struct {
    ID        int       `json:"id"`
    Username  string    `json:"username"`
    GroupID   int       `json:"group_id"`
    IsActive  bool      `json:"is_active"`
    IsLDAP    bool      `json:"is_ldap"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
