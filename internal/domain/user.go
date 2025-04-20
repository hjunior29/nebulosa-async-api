package domain

type User struct {
	Default
	Username       string
	HashedPassword string
}

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
