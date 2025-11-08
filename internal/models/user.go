package models

type Login struct {
	HashedPassword string
	SessionToken   string
	CSRFToken      string
	//ExpiresAt      time.Time
}

type User struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Login
}
