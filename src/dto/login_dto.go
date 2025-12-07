package dto

type LoginDto struct {
	AccessToken string  `json:"access_token"`
	User        UserDTO `json:"user"`
}
