package models

type Login struct {
	Password     string `json:"-"`
	SessionToken string `json:"-"`
	CSRFToken    string `json:"-"`
}

type User struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	IsAdmin  bool   `json:"is_admin"`
	Login
}
