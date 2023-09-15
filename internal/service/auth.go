package service

import (
	"context"
	"flynoob/bibirt-sock/api"
	"flynoob/bibirt-sock/internal/conf"

	"github.com/go-kratos/kratos/v2/metadata"
	mmd "github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

type AuthService struct {
	client api.AuthClient
}

func NewAuthService(c *conf.Server) *AuthService {
	md := metadata.New()
	md.Add("x-md-global-app-id", c.Api.AppId)
	md.Add("x-md-global-app-secret", c.Api.AppSecret)
	cc, err := grpc.Dial(
		context.Background(),
		grpc.WithEndpoint(c.Api.Addr),
		grpc.WithMiddleware(
			mmd.Client(mmd.WithConstants(md)),
		),
	)
	if err != nil {
		panic(err)
	}
	client := api.NewAuthClient(cc)

	return &AuthService{client}
}

func (s AuthService) ConnUUID(tokStr string) (string, error) {
	reply, err := s.client.ConnUUID(context.Background(), &api.ConnUUIDRequest{Token: tokStr})
	if err != nil {
		return "", err
	}

	return reply.Uuid, nil
}
