package model

type User struct {
	Id       int64
	Username string
	Email    string
	Password string
}

func (User) TableName() string {
	return "user_t"
}
