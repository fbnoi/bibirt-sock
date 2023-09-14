package service

import (
	"flynoob/bibirt-sock/api"

	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewAuthService)

type AuthService struct {
	client api.AuthClient

	appId     string
	appSecret string
}
