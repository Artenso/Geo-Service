package converter

import "github.com/Artenso/Geo-Service/internal/model"

func RequestAuthToUser(input *model.RequestAuth) *model.User {
	return &model.User{
		Name: input.Name,
		Pass: []byte(input.Pass),
	}
}
