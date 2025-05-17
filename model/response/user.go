package response

import "github.com/bearllflee/go_shop/model"

type LoginResponse struct {
	User  *model.User `json:"user"`
	Token string      `json:"token"`
}
