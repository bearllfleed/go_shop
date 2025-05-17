package response

import "github.com/bearllfleed/go_shop/model"

type LoginResponse struct {
	User  *model.User `json:"user"`
	Token string      `json:"token"`
}
