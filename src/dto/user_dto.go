package dto

import "room-planner/model"

type UserDTO struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func FromUserModel(user *model.User) *UserDTO {
	return &UserDTO{
		Id:       user.Id,
		Username: user.Username,
		Email:    user.Email,
	}
}
