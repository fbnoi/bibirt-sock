package service

import (
	"context"
	"flynoob/bibirt-sock/api"
	"flynoob/bibirt-sock/internal/conf"

	"github.com/go-kratos/kratos/v2/metadata"
	mmd "github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

func NewAuthService(c *conf.Server) *AuthService {

	cc, err := grpc.Dial(
		context.Background(),
		grpc.WithEndpoint(c.Api.Addr),
		grpc.WithMiddleware(
			mmd.Client(),
		),
	)
	if err != nil {
		panic(err)
	}
	client := api.NewAuthClient(cc)

	return &AuthService{client, c.Api.AppId, c.Api.AppSecret}
}

func (s AuthService) ConnUUID(tokStr string) (string, error) {
	ctx := context.Background()
	ctx = metadata.AppendToClientContext(ctx, "x-md-global-app-id", s.appId)
	ctx = metadata.AppendToClientContext(ctx, "x-md-global-app-secret", s.appSecret)
	reply, err := s.client.ConnUUID(ctx, &api.ConnUUIDRequest{Token: tokStr})
	if err != nil {
		return "", err
	}

	return reply.Uuid, nil
}
