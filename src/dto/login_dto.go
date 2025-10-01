package dto

type LoginDTO struct {
	AccessToken string  `json:"access_token"`
	User        UserDTO `json:"user"`
}
