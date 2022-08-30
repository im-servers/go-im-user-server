package logic

import (
	"context"

	"go-im-user-server/rpc/internal/svc"

	"github.com/heyehang/go-im-grpc/user_server"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserLogic {
	return &CreateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateUserLogic) CreateUser(in *user_server.CreateUserReq) (*user_server.CreateUserReply, error) {
	// todo: add your logic here and delete this line

	return &user_server.CreateUserReply{}, nil
}
