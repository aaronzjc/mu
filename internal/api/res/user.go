package res

import "github.com/aaronzjc/mu/internal/application/dto"

type UserList struct {
	List []dto.User `json:"list"`
}
