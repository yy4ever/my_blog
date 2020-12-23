package models

type User struct {
	ID           int      `db:"id" json:"id"`
	Uuid         string   `db:"uuid" json:"uuid"`
	Name         string   `db:"name" json:"name"`
	PasswordHash string   `db:"password_hash" json:"-"`
	Email        string   `db:"email" json:"email"`
	RoleID       int      `db:"role_id" json:"-"`
	RoleName     string   `db:"-" json:"role"`
	Role         Role     `db:"-" json:"-"`
	CreatedAt    JSONTime `db:"created_at" json:"created_at"`
	LastLoginAt  JSONTime `db:"last_login_at" json:"last_login_at"`
	Status       string   `db:"status" json:"status"`
}

type Role struct {
	ID          int    `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	Permissions int   `db:"permissions" json:"permissions"`
}

type Follow struct {
	FollowerID int      `db:"follower_id" json:"follower_id"`
	FollowedID int      `db:"followed_id" json:"followed_id"`
	FollowAt   JSONTime `db:"follow_at" json:"follow_at"`
}
