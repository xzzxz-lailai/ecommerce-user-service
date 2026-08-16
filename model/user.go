package model

import "time"

type User struct {
	UserID    int64     `db:"user_id"`
	Username  string    `db:"username"`
	Password  string    `db:"password"`
	Phone     string    `db:"phone"`
	Email     string    `db:"email"`
	Status    int8      `db:"status"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}
